package commands

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	cli "github.com/g-harel/edelweiss/cli"
	resources "github.com/g-harel/edelweiss/cli/resources"
	client "github.com/g-harel/edelweiss/client"
	cobra "github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dependencies in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		log.Progress("Checking that kubectl points to a running cluster")
		exited := false
		go func() {
			time.Sleep(time.Second * 3)
			if !exited {
				log.Fatal(fmt.Errorf("kubectl cluster-info timeout"), "Could not connect to cluster")
			}
		}()
		out, err := cli.Run(KUBECTL, "cluster-info")
		exited = true
		log.Fatal(err, "Could not connect to cluster: %v", out)

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
	c, err := client.New()
	log.Fatal(err, "Could not connect")
	c = c.Namespace("kube-system")

	log.Progress("Applying registry resources to cluster")

	err = c.Apply(resources.Registry)
	log.Fatal(err, "Could not create resource")

	waitForResource("Registry pod", "Running", func() (string, error) {
		return cli.Run(KUBECTL, "get", "pods",
			"--all-namespaces",
			"--selector=role=registry",
			"--output=jsonpath={.items[0].status.phase}",
		)
	})

	log.Progress("Setting up registry proxy")
	var port string
	var host string

	isMinikube, err := c.IsMinikube()
	log.Fatal(err, "Could not check if cluster running on minikube")

	service, err := c.Services().Get("registry", metav1.GetOptions{})
	log.Fatal(err, "Could not get registry service")

	if isMinikube {
		log.Progress("Fetching registry's adress")

		nodes, err := c.Nodes().List(metav1.ListOptions{})
		log.Fatal(err, "Could not get cluster nodes")

		host = nodes.Items[0].Status.Addresses[0].Address
		port = strconv.Itoa(int(service.Spec.Ports[0].NodePort))
	} else {
		log.Progress("Fetching registry's adress")

		waitForResource("Registry LoadBalancer", ".+", func() (string, error) {
			return cli.Run(KUBECTL, "get", "svc",
				"--all-namespaces",
				"--selector=role=registry",
				"--output=jsonpath={.items[0].status.loadBalancer.ingress[0].ip}",
			)
		})

		service, err := c.Services().Update(service)
		log.Fatal(err, "Could not update service status")

		log.Progress("Fetching registry's port")
		host = service.Status.LoadBalancer.Ingress[0].IP
		port = strconv.Itoa(int(service.Spec.Ports[0].Port))
	}

	name := "registry-proxy"
	cli.Run(DOCKER, "rm", "-f", name)
	out, err := cli.Run(DOCKER, "run",
		"-d", "--name", name,
		"-p", "5000:"+port,
		"-e", "BACKEND_HOST="+host,
		"-e", "BACKEND_PORT="+port,
		"demandbase/docker-tcp-proxy",
	)
	log.Fatal(err, "Could not create docker proxy: %v", out)
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
