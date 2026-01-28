package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	networkingv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/networking/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type DeplymentInfo struct {
	Deployment *appsv1.Deployment
	Name       string
	Port       int
	Ports      []int
}

type ServiceInfo struct {
	Service *corev1.Service
	Name    string
	Port    int
}

func createNamespace(ctx *pulumi.Context) (*corev1.Namespace, error) {
	name := "k8s-demo-cdk"

	appLabels := pulumi.StringMap{
		"app.kubernetes.io/name":     pulumi.String(name),
		"app.kubernetes.io/instance": pulumi.String(ctx.Stack()),
	}

	namespace, err := corev1.NewNamespace(ctx, name, &corev1.NamespaceArgs{
		ApiVersion: pulumi.String("v1"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:   pulumi.String(name),
			Labels: appLabels,
		},
	})
	if err != nil {
		return nil, err
	}
	return namespace, nil
}

func createIngress(ctx *pulumi.Context, service *ServiceInfo, raftService *ServiceInfo) (*networkingv1.Ingress, error) {
	name := service.Name
	ingressClass := "nginx"
	hostName := "demok.cdk.here"

	ingress, error := networkingv1.NewIngress(ctx, name, &networkingv1.IngressArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(name),
			Labels:    service.Service.Metadata.Labels(),
			Namespace: service.Service.Metadata.Namespace(),
			Annotations: pulumi.StringMap{
				"kubernetes.io/ingress.class": pulumi.String(ingressClass),
			},
		},
		Spec: &networkingv1.IngressSpecArgs{
			IngressClassName: pulumi.String(ingressClass),
			Rules: networkingv1.IngressRuleArray{
				&networkingv1.IngressRuleArgs{
					Host: pulumi.String(hostName),
					Http: &networkingv1.HTTPIngressRuleValueArgs{
						Paths: networkingv1.HTTPIngressPathArray{
							&networkingv1.HTTPIngressPathArgs{
								Path:     pulumi.String("/test"),
								PathType: pulumi.String("Prefix"),
								Backend: networkingv1.IngressBackendArgs{
									Service: &networkingv1.IngressServiceBackendArgs{
										Name: pulumi.String(name),
										Port: &networkingv1.ServiceBackendPortArgs{
											Number: pulumi.Int(service.Port),
										},
									},
								},
							},
							&networkingv1.HTTPIngressPathArgs{
								Path:     pulumi.String("/raft"),
								PathType: pulumi.String("Prefix"),
								Backend: networkingv1.IngressBackendArgs{
									Service: &networkingv1.IngressServiceBackendArgs{
										Name: pulumi.String(raftService.Name),
										Port: &networkingv1.ServiceBackendPortArgs{
											Name: pulumi.String("raft"),
										},
									},
								},
							},
							&networkingv1.HTTPIngressPathArgs{
								Path:     pulumi.String("/"),
								PathType: pulumi.String("Prefix"),
								Backend: networkingv1.IngressBackendArgs{
									Service: &networkingv1.IngressServiceBackendArgs{
										Name: pulumi.String(raftService.Name),
										Port: &networkingv1.ServiceBackendPortArgs{
											Name: pulumi.String("http"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})
	if error != nil {
		return nil, error
	}
	return ingress, nil
}
