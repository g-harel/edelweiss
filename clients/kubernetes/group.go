package client

import (
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apicorev1 "k8s.io/api/core/v1"
)

// SpecGroup is a collection of specs that are logically linked.
type SpecGroup struct {
	Deployments []*appsv1beta1.Deployment
	Services    []*apicorev1.Service
}
