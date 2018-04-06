package commands

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	types "github.com/docker/docker/api/types"
	container "github.com/docker/docker/api/types/container"
	filters "github.com/docker/docker/api/types/filters"
	mobyClient "github.com/docker/docker/client"
	nat "github.com/docker/go-connections/nat"
	cli "github.com/g-harel/edelweiss/cli"
	resources "github.com/g-harel/edelweiss/cli/resources"
	kubeClient "github.com/g-harel/edelweiss/client"
	cobra "github.com/spf13/cobra"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dependencies in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			args = []string{"registry", "rook"}
		}
		specs := make(map[string]bool)
		for _, s := range args {
			specs[s] = true
		}
		for s, ok := range specs {
			if ok {
				switch s {
				case "rook":
					installRook()
				case "registry":
					installRegistry()
				}
			}
		}

		log.Success("Install complete")
	},
}

// TODO refactor
// https://akomljen.com/rook-cloud-native-on-premises-persistent-storage-for-kubernetes-on-kubernetes/
func installRook() {
	repoName := "rook-master"

	log.Progress("Making sure rook repo is registered with helm")
	out, err := cli.Run(HELM, "repo", "list")
	log.Fatal(err, "Could not query helm repos")
	if strings.Index(out, repoName) < 0 {
		out, err := cli.Run(HELM, "repo", "add", repoName, "https://charts.rook.io/master")
		log.Fatal(err, "Could not add rook repo: %v", out)
	}

	log.Progress("Initializing helm in the cluster")
	out, err = cli.Run(HELM, "init", "--upgrade")
	log.Fatal(err, "Could not init helm in cluster")

	waitForResource("Tiller pod", "Running", func() (string, error) {
		return cli.Run(KUBECTL, "get", "pods",
			"--all-namespaces",
			"--selector=name=tiller",
			"--output=jsonpath={.items[0].status.phase}",
		)
	})

	log.Progress("Installing rook in the cluster")
	out, err = cli.Run(HELM, "install", repoName+"/rook",
		"--name", "rook",
		"--namespace", "kube-system",
		"--version", "v0.7.0-27.gbfc8ec6",
		"--set", "rbacEnable=false",
	)
	log.Fatal(err, "Could not install helm to cluster: %v", out)
}

func installRegistry() {
	client, err := kubeClient.New()
	log.Fatal(err, "Could not connect")
	client = client.Namespace("kube-system")

	log.Progress("Applying registry resources to cluster")

	err = client.Apply(resources.Registry)
	log.Fatal(err, "Could not create resource")

	d, err := client.Deployments().Watch(metav1.ListOptions{
		LabelSelector: "role=" + resources.Registry.Deployments[0].ObjectMeta.Labels["role"],
		Limit:         1,
	})
	log.Fatal(err, "Could not watch deployments")

	for event := range d.ResultChan() {
		deployment, ok := event.Object.(*appsv1beta1.Deployment)
		if !ok {
			log.Fatal(fmt.Errorf("Type assertion failed"), "Could not read watched deployment")
		}
		if deployment.Status.ReadyReplicas == *deployment.Spec.Replicas {
			break
		}
	}

	isMinikube, err := client.IsMinikube()
	log.Fatal(err, "Could not check if cluster running on minikube")

	service, err := client.Services().Get("registry", metav1.GetOptions{})
	log.Fatal(err, "Could not get registry service")

	log.Progress("Fetching registry's adress")

	var port string
	var host string

	if isMinikube {
		nodes, err := client.Nodes().List(metav1.ListOptions{})
		log.Fatal(err, "Could not get cluster nodes")

		host = nodes.Items[0].Status.Addresses[0].Address
		port = strconv.Itoa(int(service.Spec.Ports[0].NodePort))
	} else {
		s, err := client.Services().Watch(metav1.ListOptions{
			LabelSelector: "role=" + resources.Registry.Services[0].ObjectMeta.Labels["role"],
			Limit:         1,
		})
		log.Fatal(err, "Could not watch services")

		for event := range s.ResultChan() {
			var ok bool
			service, ok = event.Object.(*apicorev1.Service)
			if !ok {
				log.Fatal(fmt.Errorf("Type assertion failed"), "Could not read watched service")
			}
			if len(service.Status.LoadBalancer.Ingress) > 0 {
				break
			}
		}

		service, err = client.Services().Update(service)
		log.Fatal(err, "Could not update service status")

		host = service.Status.LoadBalancer.Ingress[0].IP
		port = strconv.Itoa(int(service.Spec.Ports[0].Port))
	}

	log.Progress("Setting up registry proxy")

	err = runRegistryContainer("registry-proxy", host, port)
	log.Fatal(err, "Could not create docker proxy")
}

func waitForResource(displayName, expected string, runner func() (string, error)) string {
	time.Sleep(time.Second)
	p := regexp.MustCompile(expected)
	var out string
	var err error
	for {
		out, err = runner()
		log.Fatal(err, "Cannot get %v: %v", displayName, out)
		if p.MatchString(out) {
			break
		}
		log.Progress("Waiting for %v (%v)", displayName, out)
		time.Sleep(time.Second * 3)
	}
	time.Sleep(time.Second)
	return out
}

func runRegistryContainer(name, host, port string) error {
	client, err := mobyClient.NewEnvClient()
	if err != nil {
		return err
	}

	containers, err := client.ContainerList(
		context.Background(),
		types.ContainerListOptions{
			Limit: 1,
			Filters: filters.NewArgs(filters.KeyValuePair{
				Key:   "name",
				Value: name,
			}),
		},
	)
	if err != nil {
		return err
	}

	if len(containers) > 0 {
		// error ignored because removal is optional and container often does not exist
		_ = client.ContainerRemove(
			context.Background(),
			containers[0].ID,
			types.ContainerRemoveOptions{Force: true},
		)
	}

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

	err = client.ContainerStart(context.Background(), body.ID, types.ContainerStartOptions{})

	return err
}
