package main

import (
	"context"
	"fmt"

	//appsv1 "k8s.io/api/apps/v1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func listCronJobs(clientset *kubernetes.Clientset, namespace string) (*batchv1.CronJobList, error) {
	return clientset.BatchV1().CronJobs(namespace).List(context.TODO(), metav1.ListOptions{})
}

func listDeployments(clientset *kubernetes.Clientset, namespace string) (*appsv1.DeploymentList, error) {
	return clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
}

func int32Ptr(i int32) *int32 { return &i }

func createDeployment(clientset *kubernetes.Clientset) {

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		//panic(err)
		fmt.Println("Seem deployment already exists")
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

}

func main() {
	// Load Kubernetes configuration
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/ban//.kube/config")
	if err != nil {
		panic(err)
	}

	// Create Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// List CronJobs
	cronJobs, err := listCronJobs(clientset, "default")
	if err != nil {
		panic(err)
	}

	// Print CronJobs
	fmt.Printf("CronJobs in the default namespace:\n")
	for _, cronJob := range cronJobs.Items {
		fmt.Printf("- Name: %s\n", cronJob.Name)
		fmt.Printf("  Schedule: %s\n", cronJob.Spec.Schedule)
		fmt.Printf("  Suspend: %v\n", *cronJob.Spec.Suspend)
		fmt.Printf("  Last Schedule Time: %v\n", cronJob.Status.LastScheduleTime)
		fmt.Println()
	}

	deployments, err := listDeployments(clientset, "default")
	if err != nil {
		panic(err)
	}

	//fmt.Println(clientset.AppsV1().Deployments("default").List(context.TODO(), metav1.ListOptions{}))

	// Print Deployments
	fmt.Printf("Deployments in the default namespace:\n")
	for _, deployment := range deployments.Items {
		fmt.Printf("- Name: %s\n", deployment.Name)
		fmt.Printf("  Replicas: %d/%d\n", deployment.Status.ReadyReplicas, *deployment.Spec.Replicas)
		fmt.Printf("  Strategy: %s\n", deployment.Spec.Strategy.Type)
		fmt.Printf("  Created: %v\n", deployment.CreationTimestamp)
		if len(deployment.Spec.Template.Spec.Containers) > 0 {
			fmt.Printf("  Image: %s\n", deployment.Spec.Template.Spec.Containers[0].Image)
		}
		fmt.Println()
	}

	createDeployment(clientset)

	// deployment interface with create delete method
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(context.TODO(), "demo-deployment", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment.")

}
