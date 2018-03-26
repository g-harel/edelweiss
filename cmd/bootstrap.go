package cmd

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Install dependencies in the cluster",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkFatal(HELM, KUBECTL)
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
	found := strings.Index(out, repoName) > 0
	if !found {
		out, err := run(HELM, "repo", "add", repoName, "https://charts.rook.io/master")
		fatal(err)("Could not add rook repo: %v", out)
	}

	log("Initializing helm in the cluster")
	out, err = run(HELM, "init", "--upgrade")
	fatal(err)("Could not init helm in cluster")

	waitForResource("Tiller pod", func() (string, error) {
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

	p, err := regexp.Compile("(?i)already\\s*exists")
	fatal(err)("Could not compile regular expression")

	out, err := run(KUBECTL, "create", "-f", "./resources/registry.yaml")
	if p.MatchString(out) {
		err = nil
	}
	fatal(err)("Could not create resource: %v", out)

	waitForResource("Registry pod", func() (string, error) {
		return run(KUBECTL, "get", "pods",
			"--all-namespaces",
			"--selector=role=registry",
			"--output=jsonpath={.items[0].status.phase}",
		)
	})
}

func waitForResource(displayName string, runner func() (string, error)) {
	time.Sleep(time.Second)
	for {
		out, err := runner()
		fatal(err)("Cannot get %v: %v", displayName, out)
		if out == "Running" {
			break
		}
		log("Waiting for %v (%v)", displayName, out)
		time.Sleep(time.Second * 3)
	}
	time.Sleep(time.Second)
}
