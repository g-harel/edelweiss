package main

import (
	"fmt"

	"github.com/g-harel/edelweiss/cli/commands"
	"github.com/g-harel/edelweiss/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	c, _ := client.New()
	i, _ := c.IsMinikube()
	fmt.Println(i)

	p, _ := c.Namespace("kube-system").GetPodByRole("registry")
	fmt.Println(p.Status.Phase)
	fmt.Println(p.GetName())

	s, _ := c.CoreV1().Services("kube-system").Get("registry", metav1.GetOptions{})
	fmt.Println(s.Spec.Ports[0].NodePort)

	// return
	commands.Execute()
}
