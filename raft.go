package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createRaftDeployment(ctx *pulumi.Context, namespace *corev1.Namespace, raftHeadlessServiceName string) (*DeplymentInfo, error) {
	appName := "k8s-demo-raft"
	appVersion := "latest"
	appPortHttp := 8080
	appPortRaft := 6666
	image := appName + ":" + appVersion

	appSelectorLabels := pulumi.StringMap{
		"app.kubernetes.io/name":     pulumi.String(appName),
		"app.kubernetes.io/instance": pulumi.String(ctx.Stack()),
	}

	appLabels := pulumi.StringMap{
		"app.kubernetes.io/name":       pulumi.String(appName),
		"app.kubernetes.io/instance":   pulumi.String(ctx.Stack()),
		"app.kubernetes.io/version":    pulumi.String(appVersion),
		"app.kubernetes.io/component":  pulumi.String("raft"),
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
									ContainerPort: pulumi.Int(appPortHttp),
									Name:          pulumi.String("http"),
								},
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(appPortRaft),
									Name:          pulumi.String("raft"),
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
							Env: corev1.EnvVarArray{
								corev1.EnvVarArgs{
									Name:  pulumi.String("DEMOK_NAMESPACE"),
									Value: namespace.Metadata.Name(),
								},
								corev1.EnvVarArgs{
									Name:  pulumi.String("DEMOK_RAFT_HEADLESS_SERVICE_NAME"),
									Value: pulumi.String(raftHeadlessServiceName),
								},
							},
							// Lifecycle: &corev1.LifecycleArgs{
							// 	PreStop: &corev1.LifecycleHandlerArgs{
							// 		Exec: &corev1.ExecActionArgs{
							// 			Command: pulumi.ToStringArray([]string{"sh", "-c", "sleep", "30"}),
							// 		},
							// 	},
							// },
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &DeplymentInfo{Deployment: deployment, Name: appName, Ports: []int{appPortHttp, appPortRaft}}, nil
}

func createRaftService(ctx *pulumi.Context, deployment *DeplymentInfo) (*ServiceInfo, error) {
	name := deployment.Name

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
					Name:       pulumi.String(deployment.Name + "-http"),
					Port:       pulumi.Int(deployment.Ports[0]),
					TargetPort: pulumi.String("http"),
				},
				&corev1.ServicePortArgs{
					Name:       pulumi.String(deployment.Name + "-raft"),
					Port:       pulumi.Int(deployment.Ports[1]),
					TargetPort: pulumi.String("raft"),
				},
			},
			Selector: deployment.Deployment.Spec.Selector().MatchLabels(),
		},
	})
	if error != nil {
		return nil, error
	}
	return &ServiceInfo{Service: service, Name: name}, nil
}

func createRaftHeadlessService(ctx *pulumi.Context, deployment *DeplymentInfo, serviceName string) (*ServiceInfo, error) {

	service, error := corev1.NewService(ctx, serviceName, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(serviceName),
			Labels:    deployment.Deployment.Metadata.Labels(),
			Namespace: deployment.Deployment.Metadata.Namespace(),
		},
		Spec: &corev1.ServiceSpecArgs{
			ClusterIP: pulumi.String("None"),
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Name:       pulumi.String(deployment.Name),
					Port:       pulumi.Int(deployment.Ports[1]),
					TargetPort: pulumi.String("raft"),
				},
			},
			Selector: deployment.Deployment.Spec.Selector().MatchLabels(),
		},
	})
	if error != nil {
		return nil, error
	}
	return &ServiceInfo{Service: service, Name: serviceName, Port: deployment.Ports[1]}, nil
}
