package resources

import (
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Registry resource group deploys a container registry to the cluster and
// exposes it through a LoadBalancer service.
var Registry = &Group{
	Services: []apicorev1.Service{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "registry",
				Namespace: "kube-system",
				Labels: map[string]string{
					"role": "registry",
				},
			},
			Spec: apicorev1.ServiceSpec{
				Ports: []apicorev1.ServicePort{
					{
						Port: 5000,
						TargetPort: intstr.IntOrString{
							Type:   0,
							IntVal: 5000,
						},
					},
				},
				Selector: map[string]string{
					"role": "registry",
				},
			},
		},
	},
	Deployments: []appsv1beta1.Deployment{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "registry",
				Namespace: "kube-system",
				Labels: map[string]string{
					"role": "registry",
				},
			},
			Spec: appsv1beta1.DeploymentSpec{
				Replicas: int32Ptr(1),
				Template: apicorev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"role": "registry",
						},
					},
					Spec: apicorev1.PodSpec{
						Containers: []apicorev1.Container{
							{
								Name:  "registry",
								Image: "registry:2.6.1",
								Ports: []apicorev1.ContainerPort{
									{
										Protocol:      apicorev1.ProtocolTCP,
										ContainerPort: 80,
									},
								},
							},
						},
					},
				},
			},
		},
	},
}
