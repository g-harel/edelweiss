package cmd

import (
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rookCmd = &cobra.Command{
	Use:   "rook",
	Short: "Deploy rook to the cluster",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkFatal(HELM, KUBECTL)
	},
	Run: func(cmd *cobra.Command, args []string) {
		repoName := "rook-master"

		// check that rook repo is registered with helm
		out, err := run(HELM, "repo", "list")
		fatal(err)("Could not query helm repos")
		found := strings.Index(out, repoName) > 0
		if !found {
			_, err := run(HELM, "repo", "add", repoName, "https://charts.rook.io/master")
			fatal(err)("Could not add rook repo")
		}

		// check that kubectl points to a running cluster.
		out, err = run(KUBECTL, "cluster-info")
		fatal(err)("Could not connect to cluster")

		// initalize helm in the cluster
		out, err = run(HELM, "init", "--upgrade")
		fatal(err)("Could not init helm in cluster")
		time.Sleep(time.Second)
		for {
			out, err = run(KUBECTL, "get", "pods",
				"--all-namespaces",
				"--selector=name=tiller",
				"--output=jsonpath={.items[0].status.phase}",
			)
			fatal(err)("Cannot get Tiller pod")
			if out == "Running" {
				break
			}
			color.White("Waiting for Tiller pod, status: %v", out)
			time.Sleep(time.Second * 3)
		}
		time.Sleep(time.Second)

		// install rook in the cluster with helm
		out, err = run(HELM, "install", repoName+"/rook",
			"--name", "rook",
			"--namespace", "kube-system",
			"--version", "v0.7.0-27.gbfc8ec6",
			"--set", "rbacEnable=false",
		)
		fatal(err)("Could not install helm to cluster")

		color.Green("\n✓ hia\n\n")
	},
}
