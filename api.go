package main

import (
	"strconv"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createApiDeployment(ctx *pulumi.Context, namespace *corev1.Namespace, backendService *ServiceInfo) (*DeplymentInfo, error) {
	appName := "k8s-demo-api"
	appVersion := "latest"
	appPort := 8080
	image := appName + ":" + appVersion

	appSelectorLabels := pulumi.StringMap{
		"app.kubernetes.io/name":     pulumi.String(appName),
		"app.kubernetes.io/instance": pulumi.String(ctx.Stack()),
	}

	appLabels := pulumi.StringMap{
		"app.kubernetes.io/name":       pulumi.String(appName),
		"app.kubernetes.io/instance":   pulumi.String(ctx.Stack()),
		"app.kubernetes.io/version":    pulumi.String(appVersion),
		"app.kubernetes.io/component":  pulumi.String("api"),
		"app.kubernetes.io/managed-by": pulumi.String("Pulumi-cdk"),
	}

	deployment, err := appsv1.NewDeployment(ctx, appName, &appsv1.DeploymentArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(appName),
			Labels:    appLabels,
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
							Name:            pulumi.String(appName),
							Image:           pulumi.String(image),
							ImagePullPolicy: pulumi.String("IfNotPresent"),
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(appPort),
									Name:          pulumi.String("http"),
								},
							},
							LivenessProbe: &corev1.ProbeArgs{
								HttpGet: &corev1.HTTPGetActionArgs{
									Path: pulumi.String("/actuator/health/liveness"),
									Port: pulumi.String("http"),
								},
							},
							ReadinessProbe: &corev1.ProbeArgs{
								HttpGet: &corev1.HTTPGetActionArgs{
									Path: pulumi.String("/actuator/health/readiness"),
									Port: pulumi.String("http"),
								},
							},
							Env: corev1.EnvVarArray{
								corev1.EnvVarArgs{
									Name:  pulumi.String("DEMOK_BACKEND_SERVICE_NAME"),
									Value: pulumi.String(backendService.Name),
								},
								corev1.EnvVarArgs{
									Name:  pulumi.String("DEMOK_BACKEND_SERVICE_PORT"),
									Value: pulumi.String(strconv.Itoa(backendService.Port)),
								},
							},
							Lifecycle: &corev1.LifecycleArgs{
								PreStop: &corev1.LifecycleHandlerArgs{
									Exec: &corev1.ExecActionArgs{
										Command: pulumi.ToStringArray([]string{"sh", "-c", "sleep", "30"}),
									},
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
	return &DeplymentInfo{Deployment: deployment, Name: appName, Port: appPort}, nil
}

func createApiService(ctx *pulumi.Context, deployment *DeplymentInfo) (*ServiceInfo, error) {
	name := deployment.Name
	port := deployment.Port

	service, error := corev1.NewService(ctx, name, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(name),
			Labels:    deployment.Deployment.Metadata.Labels(),
			Namespace: deployment.Deployment.Metadata.Namespace(),
		},
		Spec: &corev1.ServiceSpecArgs{
			Type: corev1.ServiceSpecTypeClusterIP,
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Name:       pulumi.String(name),
					Port:       pulumi.Int(port),
					TargetPort: pulumi.String("http"),
				},
			},
			Selector: deployment.Deployment.Spec.Selector().MatchLabels(),
		},
	})
	if error != nil {
		return nil, error
	}
	return &ServiceInfo{Service: service, Name: name, Port: port}, nil
}