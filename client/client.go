package client

import (
	"fmt"
	"path/filepath"

	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// connect to cluster
// check if running on minikube cluster
// apply registry resources

// select pod by name (tiller, registry) + get status
// get service info

func A(g *Group) {
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

	fmt.Println(apply(clientset, g, "kube-system"))
	fmt.Println(isMinikube(clientset))
	fmt.Println(getPodByRole(clientset, "kube-system", "registry"))
}

func apply(cs *kubernetes.Clientset, g *Group, ns string) error {
	if len(g.Deployments) > 0 {
		dc := cs.AppsV1beta1().Deployments(ns)
		for _, d := range g.Deployments {
			_, err := dc.Create(d)
			if errors.IsAlreadyExists(err) {
				_, err = dc.Update(d)
				continue
			}
			if err != nil {
				return err
			}
		}
	}

	if len(g.Services) > 0 {
		sc := cs.CoreV1().Services(ns)
		for _, s := range g.Services {
			_, err := sc.Create(s)
			if errors.IsAlreadyExists(err) {
				_, err = sc.Update(s)
				continue
			}
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isMinikube(cs *kubernetes.Clientset) (bool, error) {
	nodeClient := cs.CoreV1().Nodes()
	l, err := nodeClient.List(metav1.ListOptions{})
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
	l, err := podClient.List(metav1.ListOptions{
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
