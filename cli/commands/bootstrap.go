package commands

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/g-harel/edelweiss/cli/resources"
	"github.com/g-harel/edelweiss/client"
	"github.com/spf13/cobra"
)

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Install dependencies in the cluster",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkFatal(DOCKER, HELM, KUBECTL)
	},
	Run: func(cmd *cobra.Command, args []string) {
		log("Checking that kubectl points to a running cluster")
		exited := false
		go func() {
			time.Sleep(time.Second * 3)
			if !exited {
				fatal(fmt.Errorf("kubectl cluster-info timeout"))("Could not connect to cluster")
			}
		}()
		out, err := run(KUBECTL, "cluster-info")
		exited = true
		fatal(err)("Could not connect to cluster: %v", out)

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
					bootstrapRook()
				case "registry":
					bootstrapRegistry()
				}
			}
		}

		color.Green("\n✓ Bootstrap complete\n\n")
	},
}

func bootstrapRook() {
	repoName := "rook-master"

	log("Making sure rook repo is registered with helm")
	out, err := run(HELM, "repo", "list")
	fatal(err)("Could not query helm repos")
	if strings.Index(out, repoName) < 0 {
		out, err := run(HELM, "repo", "add", repoName, "https://charts.rook.io/master")
		fatal(err)("Could not add rook repo: %v", out)
	}

	log("Initializing helm in the cluster")
	out, err = run(HELM, "init", "--upgrade")
	fatal(err)("Could not init helm in cluster")

	waitForResource("Tiller pod", "Running", func() (string, error) {
		return run(KUBECTL, "get", "pods",
			"--all-namespaces",
			"--selector=name=tiller",
			"--output=jsonpath={.items[0].status.phase}",
		)
	})

	log("Installing rook in the cluster")
	out, err = run(HELM, "install", repoName+"/rook",
		"--name", "rook",
		"--namespace", "kube-system",
		"--version", "v0.7.0-27.gbfc8ec6",
		"--set", "rbacEnable=false",
	)
	fatal(err)("Could not install helm to cluster: %v", out)
}

func bootstrapRegistry() {
	log("Installing registry in the cluster")

	_, err := regexp.Compile("(?i)already\\s*exists")
	fatal(err)("Could not compile regular expression")

	log("Applying registry resources to cluster")
	// out, err := run(KUBECTL, "apply", "-f", "./cli/resources/registry.yaml")
	// if p.MatchString(out) {
	// 	err = nil
	// }
	out := ""
	client.A(resources.Registry)
	fatal(err)("Could not create resource: %v", out)

	waitForResource("Registry pod", "Running", func() (string, error) {
		return run(KUBECTL, "get", "pods",
			"--all-namespaces",
			"--selector=role=registry",
			"--output=jsonpath={.items[0].status.phase}",
		)
	})

	log("Setting up registry proxy")
	var port string
	var host string

	// checking if cluster is running on minikube
	out, err = run(KUBECTL, "get", "nodes",
		"--output=jsonpath={$.items[?(@.spec.externalID==\"minikube\")].status.addresses[?(@.type==\"InternalIP\")].address}",
	)
	fatal(err)("Could not query cluster's nodes: %v", out)
	isMinikube := out != ""

	if isMinikube {
		log("Fetching registry's adress")
		host = out

		log("Fetching registry's port")
		out, err = run(KUBECTL, "get", "svc",
			"--all-namespaces",
			"--selector=role=registry",
			"--output=jsonpath={.items[0].spec.ports[0].nodePort}",
		)
		if out == "" {
			err = fmt.Errorf("Could not find service")
		}
		fatal(err)("Could not get service's port: %v", out)
		port = out
	} else {
		log("Fetching registry's adress")
		host = waitForResource("Registry LoadBalancer", ".+", func() (string, error) {
			return run(KUBECTL, "get", "svc",
				"--all-namespaces",
				"--selector=role=registry",
				"--output=jsonpath={.items[0].status.loadBalancer.ingress[0].ip}",
			)
		})

		log("Fetching registry's port")
		out, err = run(KUBECTL, "get", "svc",
			"--all-namespaces",
			"--selector=role=registry",
			"--output=jsonpath={.items[0].spec.ports[0].port}",
		)
		if out == "" {
			err = fmt.Errorf("Could not find service")
		}
		fatal(err)("Could not get service's port: %v", out)
		port = out
	}

	name := "registry-proxy"
	run(DOCKER, "rm", "-f", name)
	out, err = run(DOCKER, "run",
		"-d", "--name", name,
		"-p", "5000:"+port,
		"-e", "BACKEND_HOST="+host,
		"-e", "BACKEND_PORT="+port,
		"demandbase/docker-tcp-proxy",
	)
	fatal(err)("Could not create docker proxy: %v", out)
}

func waitForResource(displayName, expected string, runner func() (string, error)) string {
	time.Sleep(time.Second)
	p := regexp.MustCompile(expected)
	var out string
	var err error
	for {
		out, err = runner()
		fatal(err)("Cannot get %v: %v", displayName, out)
		if p.MatchString(out) {
			break
		}
		log("Waiting for %v (%v)", displayName, out)
		time.Sleep(time.Second * 3)
	}
	time.Sleep(time.Second)
	return out
}
