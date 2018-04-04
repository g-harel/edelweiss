package client

import (
	"fmt"
	"path/filepath"

	apicorev1 "k8s.io/api/core/v1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	homedir "k8s.io/client-go/util/homedir"
)

// connect to cluster
// check if running on minikube cluster
// apply registry resources

// select pod by name (tiller, registry) + get status
// get service info

type Client struct {
	clientset  *kubernetes.Clientset
	namespace  string
	isMinikube *bool
}

func New() (*Client, error) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		clientset: clientset,
		namespace: "default",
	}, nil
}

func (c *Client) Namespace(namespace string) *Client {
	t := Client{
		clientset:  c.clientset,
		isMinikube: c.isMinikube,
		namespace:  namespace,
	}
	return &t
}

func (c *Client) Apply(g *SpecGroup) error {
	if len(g.Deployments) > 0 {
		dc := c.clientset.AppsV1beta1().Deployments(c.namespace)
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
		sc := c.clientset.CoreV1().Services(c.namespace)
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

func (c *Client) IsMinikube() (bool, error) {
	if c.isMinikube != nil {
		return *c.isMinikube, nil
	}

	nodeClient := c.clientset.CoreV1().Nodes()
	l, err := nodeClient.List(metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	t := true
	if len(l.Items) != 1 || l.Items[0].GetName() != "minikube" {
		t = false
	}
	c.isMinikube = &t

	return t, nil
}

func (c *Client) GetPodByRole(role string) (*apicorev1.Pod, error) {
	podClient := c.clientset.CoreV1().Pods(c.namespace)
	l, err := podClient.List(metav1.ListOptions{
		LabelSelector: "role=" + role,
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
