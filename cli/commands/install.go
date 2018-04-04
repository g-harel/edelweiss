package commands

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/g-harel/edelweiss/cli"
	"github.com/g-harel/edelweiss/cli/resources"
	"github.com/g-harel/edelweiss/client"
	"github.com/spf13/cobra"
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
	log.Progress("Installing registry in the cluster")

	_, err := regexp.Compile("(?i)already\\s*exists")
	log.Fatal(err, "Could not compile regular expression")

	log.Progress("Applying registry resources to cluster")
	c, err := client.New()
	log.Fatal(err, "Could not connect")
	err = c.
		Namespace("kube-system").
		Apply(resources.Registry)
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

	// checking if cluster is running on minikube
	out, err := cli.Run(KUBECTL, "get", "nodes",
		"--output=jsonpath={$.items[?(@.spec.externalID==\"minikube\")].status.addresses[?(@.type==\"InternalIP\")].address}",
	)
	log.Fatal(err, "Could not query cluster's nodes: %v", out)
	isMinikube := out != ""

	if isMinikube {
		log.Progress("Fetching registry's adress")
		host = out

		log.Progress("Fetching registry's port")
		out, err = cli.Run(KUBECTL, "get", "svc",
			"--all-namespaces",
			"--selector=role=registry",
			"--output=jsonpath={.items[0].spec.ports[0].nodePort}",
		)
		if out == "" {
			err = fmt.Errorf("Could not find service")
		}
		log.Fatal(err, "Could not get service's port: %v", out)
		port = out
	} else {
		log.Progress("Fetching registry's adress")
		host = waitForResource("Registry LoadBalancer", ".+", func() (string, error) {
			return cli.Run(KUBECTL, "get", "svc",
				"--all-namespaces",
				"--selector=role=registry",
				"--output=jsonpath={.items[0].status.loadBalancer.ingress[0].ip}",
			)
		})

		log.Progress("Fetching registry's port")
		out, err = cli.Run(KUBECTL, "get", "svc",
			"--all-namespaces",
			"--selector=role=registry",
			"--output=jsonpath={.items[0].spec.ports[0].port}",
		)
		if out == "" {
			err = fmt.Errorf("Could not find service")
		}
		log.Fatal(err, "Could not get service's port: %v", out)
		port = out
	}

	name := "registry-proxy"
	cli.Run(DOCKER, "rm", "-f", name)
	out, err = cli.Run(DOCKER, "run",
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
