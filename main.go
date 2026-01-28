package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)


func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		namespace, err := createNamespace(ctx)
		if err != nil {
			return err
		}

		/////////////////////////////
		// BACKEND
		/////////////////////////////

		backendDeployment, err := createBackendDeployment(ctx, namespace)
		if err != nil {
			return err
		}

		backendService, err := createBackendService(ctx, backendDeployment)
		if err != nil {
			return err
		}

		/////////////////////////////
		// API
		/////////////////////////////

		apiDeployment, err := createApiDeployment(ctx, namespace, backendService)
		if err != nil {
			return err
		}

		apiService, err := createApiService(ctx, apiDeployment)
		if err != nil {
			return err
		}

		/////////////////////////////
		// RAFT
		/////////////////////////////

		raftHeadlessServiceName := "k8s-demo-raft-headless"
		raftDeployment, err := createRaftDeployment(ctx, namespace, raftHeadlessServiceName)
		if err != nil {
			return err
		}

		raftService, err := createRaftService(ctx, raftDeployment)
		if err != nil {
			return err
		}

		raftHeadlessService, err := createRaftHeadlessService(ctx, raftDeployment, raftHeadlessServiceName)
		if err != nil {
			return err
		}

		/////////////////////////////
		// INGRESS
		/////////////////////////////

		ingress, err := createIngress(ctx, apiService, raftService)
		if err != nil {
			return err
		}

		ctx.Export("namespace", namespace.Metadata.Name())
		ctx.Export("raftDeployment", raftDeployment.Deployment.Metadata.Name())
		ctx.Export("raftService", raftService.Service.Metadata.Name())
		ctx.Export("raftHeadlessService", raftHeadlessService.Service.Metadata.Name())
		ctx.Export("backendDeployment", backendDeployment.Deployment.Metadata.Name())
		ctx.Export("backendService", backendService.Service.Metadata.Name())
		ctx.Export("apiDeployment", apiDeployment.Deployment.Metadata.Name())
		ctx.Export("apiService", apiService.Service.Metadata.Name())
		ctx.Export("ingress", ingress.Metadata.Name())

		return nil
	})
}
