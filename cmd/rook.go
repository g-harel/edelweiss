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

		out, err := run("helm", "repo", "list")
		fatal(err)("Could not query helm repos")
		found := strings.Index(out, repoName) > 0
		if !found {
			_, err := run("helm", "repo", "add", repoName, "https://charts.rook.io/master")
			fatal(err)("Could not add rook repo")
		}

		out, err = run("helm", "init", "--upgrade")
		fatal(err)("Could not init helm in cluster")

		for {
			out, err = run("kubectl", "get", "pods",
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

		out, err = run("helm", "install", repoName+"/rook",
			"--name", "rook",
			"--namespace", "kube-system",
			"--version", "v0.7.0-27.gbfc8ec6",
			"--set", "rbacEnable=false",
		)
		fatal(err)("Could not install helm to cluster")

		color.Green("\n✓ hia\n\n")
	},
}
