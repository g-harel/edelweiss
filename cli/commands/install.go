package commands

import (
	"context"
	"fmt"
	"strconv"

	types "github.com/docker/docker/api/types"
	container "github.com/docker/docker/api/types/container"
	nat "github.com/docker/go-connections/nat"
	resources "github.com/g-harel/edelweiss/cli/resources"
	kubeclient "github.com/g-harel/edelweiss/clients/kubernetes"
	mobyclient "github.com/g-harel/edelweiss/clients/moby"
	cobra "github.com/spf13/cobra"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dependencies in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := kubeclient.New()
		clilog.Fatal(err, "Could not create client")
		client = client.Namespace("kube-system")

		clilog.Progress("Applying registry resources to cluster")

		err = client.Apply(resources.Registry)
		clilog.Fatal(err, "Could not create resource")

		d, err := client.Deployments().Watch(metav1.ListOptions{
			LabelSelector: "role=" + resources.Registry.Deployments[0].ObjectMeta.Labels["role"],
			Limit:         1,
		})
		clilog.Fatal(err, "Could not watch deployments")

		for event := range d.ResultChan() {
			deployment, ok := event.Object.(*appsv1beta1.Deployment)
			if !ok {
				clilog.Fatal(fmt.Errorf("Type assertion failed"), "Could not read watched deployment")
			}
			if deployment.Status.ReadyReplicas == *deployment.Spec.Replicas {
				break
			}
		}

		isMinikube, err := client.IsMinikube()
		clilog.Fatal(err, "Could not check if cluster running on minikube")

		service, err := client.Services().Get("registry", metav1.GetOptions{})
		clilog.Fatal(err, "Could not get registry service")

		clilog.Progress("Fetching registry's adress")

		var port string
		var host string

		if isMinikube {
			nodes, err := client.Nodes().List(metav1.ListOptions{})
			clilog.Fatal(err, "Could not get cluster nodes")

			host = nodes.Items[0].Status.Addresses[0].Address
			port = strconv.Itoa(int(service.Spec.Ports[0].NodePort))
		} else {
			s, err := client.Services().Watch(metav1.ListOptions{
				LabelSelector: "role=" + resources.Registry.Services[0].ObjectMeta.Labels["role"],
				Limit:         1,
			})
			clilog.Fatal(err, "Could not watch services")

			for event := range s.ResultChan() {
				var ok bool
				service, ok = event.Object.(*apicorev1.Service)
				if !ok {
					clilog.Fatal(fmt.Errorf("Type assertion failed"), "Could not read watched service")
				}
				if len(service.Status.LoadBalancer.Ingress) > 0 {
					break
				}
			}

			service, err = client.Services().Update(service)
			clilog.Fatal(err, "Could not update service status")

			host = service.Status.LoadBalancer.Ingress[0].IP
			port = strconv.Itoa(int(service.Spec.Ports[0].Port))
		}

		clilog.Progress("Setting up registry proxy")

		err = runRegistryContainer("registry-proxy", host, port)
		clilog.Fatal(err, "Could not create docker proxy")

		clilog.Success("Install complete")
	},
}

func runRegistryContainer(name, host, port string) error {
	client, err := mobyclient.New()
	if err != nil {
		return err
	}

	// error ignored because removal is optional and container often does not exist
	_ = client.ContainerRemoveByName(name)

	p, err := nat.NewPort("tcp", port)
	if err != nil {
		return err
	}

	body, err := client.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: "demandbase/docker-tcp-proxy",
			ExposedPorts: nat.PortSet{
				p: struct{}{},
			},
			Env: []string{
				"BACKEND_HOST=" + host,
				"BACKEND_PORT=" + port,
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				p: []nat.PortBinding{{
					HostIP:   "localhost",
					HostPort: "5000",
				}},
			},
		},
		nil,
		name,
	)
	if err != nil {
		return err
	}

	return client.ContainerStart(context.Background(), body.ID, types.ContainerStartOptions{})
}
