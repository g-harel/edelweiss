package client

import (
	typedappsv1beta1 "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func (c *Client) Services() typedcorev1.ServiceInterface {
	return c.CoreV1().Services(c.namespace)
}

func (c *Client) Deployments() typedappsv1beta1.DeploymentInterface {
	return c.AppsV1beta1().Deployments(c.namespace)
}

func (c *Client) Nodes() typedcorev1.NodeInterface {
	return c.CoreV1().Nodes()
}
