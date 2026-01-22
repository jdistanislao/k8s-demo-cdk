package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	networkingv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/networking/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createNamespace(ctx *pulumi.Context) (*corev1.Namespace, error) {
	name := "k8s-demo-cdk"
	
	appLabels := pulumi.StringMap{
			"app.kubernetes.io/name": pulumi.String(name),
			"app.kubernetes.io/instance": pulumi.String(ctx.Stack()),
		}

	namespace, err := corev1.NewNamespace(ctx, name, &corev1.NamespaceArgs{
		ApiVersion: pulumi.String("v1"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(name),
			Labels: appLabels,
		},

	})
	if err != nil {
		return nil, err
	}
	return namespace, nil
}

func createBackendDeployment(ctx *pulumi.Context, namespace *corev1.Namespace) (*appsv1.Deployment, error) {
	appName := "k8s-demo-backend"
	appVersion := "0.1"
	image := appName+":"+appVersion

	appSelectorLabels := pulumi.StringMap{
		"app.kubernetes.io/name": pulumi.String(appName),
		"app.kubernetes.io/instance": pulumi.String(ctx.Stack()),
	}

	appLabels := pulumi.StringMap{
		"app.kubernetes.io/name": pulumi.String(appName),
		"app.kubernetes.io/instance": pulumi.String(ctx.Stack()),
		"app.kubernetes.io/version": pulumi.String(appVersion),
		"app.kubernetes.io/component": pulumi.String("api"),
		"app.kubernetes.io/managed-by": pulumi.String("Pulumi-cdk"),
	}

	deployment, err := appsv1.NewDeployment(ctx, appName, &appsv1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String(appName),
				Labels: appLabels,
				Namespace: namespace.Metadata.Name(),
			},
			Spec: &appsv1.DeploymentSpecArgs{
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: appSelectorLabels,
				},
				Replicas: pulumi.Int(3),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: appLabels,
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Name:  pulumi.String(appName),
								Image: pulumi.String(image),
								ImagePullPolicy: pulumi.String("IfNotPresent"),
								Ports: corev1.ContainerPortArray{
									&corev1.ContainerPortArgs{
										ContainerPort: pulumi.Int(8080),
										Name: pulumi.String("http"),
									},
								},
								LivenessProbe: &corev1.ProbeArgs{
									HttpGet: &corev1.HTTPGetActionArgs{
										Path: pulumi.String("/live"),
										Port: pulumi.String("http"),
									},
								},
								ReadinessProbe: &corev1.ProbeArgs{
									HttpGet: &corev1.HTTPGetActionArgs{
										Path: pulumi.String("/ready"),
										Port: pulumi.String("http"),
									},
								},
							},
						},
					},
				},
			},
		})
		if err != nil {
			return nil, err
		}
		return deployment, nil
}

func createApiDeployment(ctx *pulumi.Context, namespace *corev1.Namespace) (*appsv1.Deployment, error) {
	appName := "k8s-demo-api"
	appVersion := "0.3"
	image := appName+":"+appVersion

	appSelectorLabels := pulumi.StringMap{
		"app.kubernetes.io/name": pulumi.String(appName),
		"app.kubernetes.io/instance": pulumi.String(ctx.Stack()),
	}

	appLabels := pulumi.StringMap{
		"app.kubernetes.io/name": pulumi.String(appName),
		"app.kubernetes.io/instance": pulumi.String(ctx.Stack()),
		"app.kubernetes.io/version": pulumi.String(appVersion),
		"app.kubernetes.io/component": pulumi.String("api"),
		"app.kubernetes.io/managed-by": pulumi.String("Pulumi-cdk"),
	}

	deployment, err := appsv1.NewDeployment(ctx, appName, &appsv1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String(appName),
				Labels: appLabels,
				Namespace: namespace.Metadata.Name(),
			},
			Spec: &appsv1.DeploymentSpecArgs{
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: appSelectorLabels,
				},
				Replicas: pulumi.Int(3),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: appLabels,
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Name:  pulumi.String(appName),
								Image: pulumi.String(image),
								ImagePullPolicy: pulumi.String("IfNotPresent"),
								Ports: corev1.ContainerPortArray{
									&corev1.ContainerPortArgs{
										ContainerPort: pulumi.Int(8080),
										Name: pulumi.String("http"),
									},
								},
								LivenessProbe: &corev1.ProbeArgs{
									HttpGet: &corev1.HTTPGetActionArgs{
										Path: pulumi.String("/live"),
										Port: pulumi.String("http"),
									},
								},
								ReadinessProbe: &corev1.ProbeArgs{
									HttpGet: &corev1.HTTPGetActionArgs{
										Path: pulumi.String("/ready"),
										Port: pulumi.String("http"),
									},
								},
							},
						},
					},
				},
			},
		})
		if err != nil {
			return nil, err
		}
		return deployment, nil
}

func createApiService(ctx *pulumi.Context, deployment *appsv1.Deployment) (*corev1.Service, error) {
	name := "k8s-demo-api"

	service, error := corev1.NewService(ctx, name, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(name),
			Labels: deployment.Metadata.Labels(),
			Namespace: deployment.Metadata.Namespace(),
		},
		Spec: &corev1.ServiceSpecArgs{
			Type: corev1.ServiceSpecTypeClusterIP,
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Name: pulumi.String(name),
					Port: pulumi.Int(8080),
					TargetPort: pulumi.String("http"),
				},
			},
			Selector: deployment.Spec.Selector().MatchLabels(),
		},
		
	})
	if error != nil {
		return nil, error
	}
	return service, nil
}

func createApiIngress(ctx *pulumi.Context, service *corev1.Service) (*networkingv1.Ingress, error) {
	name := "k8s-demo-api"
	ingressClass := "nginx"
	hostName := "demok.cdk.here"

	ingress, error := networkingv1.NewIngress(ctx, name, &networkingv1.IngressArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(name),
			Labels: service.Metadata.Labels(),
			Namespace: service.Metadata.Namespace(),
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
								Backend: networkingv1.IngressBackendArgs{
									Service: &networkingv1.IngressServiceBackendArgs{
										Name: pulumi.String(name),
										Port: &networkingv1.ServiceBackendPortArgs{
											// Name: pulumi.String(name),
											Number: pulumi.Int(8080),
										},
									},
								},
								Path: pulumi.String("/test"),
								PathType: pulumi.String("Prefix"),
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


func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		namespace, err := createNamespace(ctx)
		if err != nil {
			return err
		}

		backendDeployment, err := createBackendDeployment(ctx, namespace)
		if err != nil {
			return err
		}

		apiDeployment, err := createApiDeployment(ctx, namespace)
		if err != nil {
			return err
		}

		apiService, err := createApiService(ctx, apiDeployment)
		if err != nil {
			return err
		}

		apiIngress, err := createApiIngress(ctx, apiService)
		if err != nil {
			return err
		}

		ctx.Export("namespace", namespace.Metadata.Name())
		ctx.Export("backendDeployment", backendDeployment.Metadata.Name())
		ctx.Export("apiDeployment", apiDeployment.Metadata.Name())
		ctx.Export("apiService", apiService.Metadata.Name())
		ctx.Export("apiIngress", apiIngress.Metadata.Name())

		return nil
	})
}
