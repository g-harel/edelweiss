package client

import (
	typedappsv1beta1 "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// Services is a shortcut to the service client which uses the configured namespace.
func (c *Client) Services() typedcorev1.ServiceInterface {
	return c.CoreV1().Services(c.namespace)
}

// Deployments is a shortcut to the deployment client which uses the configured namespace.
func (c *Client) Deployments() typedappsv1beta1.DeploymentInterface {
	return c.AppsV1beta1().Deployments(c.namespace)
}

// Nodes is a shortcut to the node client.
func (c *Client) Nodes() typedcorev1.NodeInterface {
	return c.CoreV1().Nodes()
}
