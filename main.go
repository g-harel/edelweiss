package main

import (
	"fmt"
	"path/filepath"

	"github.com/g-harel/edelweiss/cli/commands"
	apicorev1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// connect to cluster
// check if running on minikube cluster

// select pod by name (tiller, registry) + get status
// get service info
// apply registry resources

func main() {
	var customconfig string

	var kubeconfig string
	if customconfig == "" {
		kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	} else {
		kubeconfig = customconfig
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	commands.Execute()
	return
	fmt.Println(isMinikube(clientset))
	fmt.Println(getPodByRole(clientset, "kube-system", "registry"))

	return
}

func isMinikube(cs *kubernetes.Clientset) (bool, error) {
	nodeClient := cs.CoreV1().Nodes()
	l, err := nodeClient.List(meta_v1.ListOptions{})
	if err != nil {
		return false, err
	}
	if len(l.Items) != 1 {
		return false, nil
	}
	if l.Items[0].GetName() != "minikube" {
		return false, nil
	}
	return true, nil
}

func getPodByRole(cs *kubernetes.Clientset, ns, rl string) (*apicorev1.Pod, error) {
	podClient := cs.CoreV1().Pods(ns)
	l, err := podClient.List(meta_v1.ListOptions{
		LabelSelector: "role=" + rl,
		Limit:         1,
	})
	if err != nil {
		panic(err)
	}
	if len(l.Items) < 1 {
		return nil, fmt.Errorf("Pods not found")
	}
	return &l.Items[0], nil
}
