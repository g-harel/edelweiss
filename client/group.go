package client

import (
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apicorev1 "k8s.io/api/core/v1"
)

// Group is a collection of specs that are logically linked.
type Group struct {
	Deployments []*appsv1beta1.Deployment
	Services    []*apicorev1.Service
}
