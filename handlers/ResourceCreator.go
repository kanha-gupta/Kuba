package handlers

import (
	"context"
	"fmt"
	"io/ioutil"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

func YamlResourceCreator(clientset *kubernetes.Clientset, namespace string, filePath string) error {
	yamlContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	decoder := serializer.NewCodecFactory(scheme.Scheme).UniversalDeserializer()

	obj, _, err := decoder.Decode(yamlContent, nil, nil)
	if err != nil {
		return err
	}

	typedObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return err
	}

	switch kind := typedObj["kind"]; kind {
	case "Deployment":
		deployment := &appsv1.Deployment{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(typedObj, deployment)
		if err != nil {
			return err
		}

		_, err = clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
		if err != nil {
			return err
		}

		fmt.Println("Deployment created:", deployment.GetName())

	case "Service":
		service := &corev1.Service{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(typedObj, service)
		if err != nil {
			return err
		}

		_, err = clientset.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
		if err != nil {
			return err
		}

		fmt.Println("Service created:", service.GetName())

	case "StatefulSet":
		statefulSet := &appsv1.StatefulSet{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(typedObj, statefulSet)
		if err != nil {
			return err
		}

		_, err = clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), statefulSet, metav1.CreateOptions{})
		if err != nil {
			return err
		}

		fmt.Println("StatefulSet created:", statefulSet.GetName())

	case "DaemonSet":
		daemonSet := &appsv1.DaemonSet{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(typedObj, daemonSet)
		if err != nil {
			return err
		}

		_, err = clientset.AppsV1().DaemonSets(namespace).Create(context.TODO(), daemonSet, metav1.CreateOptions{})
		if err != nil {
			return err
		}

		fmt.Println("DaemonSet created:", daemonSet.GetName())

	case "Job":
		job := &batchv1.Job{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(typedObj, job)
		if err != nil {
			return err
		}

		_, err = clientset.BatchV1().Jobs(namespace).Create(context.TODO(), job, metav1.CreateOptions{})
		if err != nil {
			return err
		}

		fmt.Println("Job created:", job.GetName())

	case "CronJob":
		cronJob := &batchv1beta1.CronJob{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(typedObj, cronJob)
		if err != nil {
			return err
		}

		_, err = clientset.BatchV1beta1().CronJobs(namespace).Create(context.TODO(), cronJob, metav1.CreateOptions{})
		if err != nil {
			return err
		}

		fmt.Println("CronJob created:", cronJob.GetName())

	case "Namespace":
		namespaceObj := &corev1.Namespace{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(typedObj, namespaceObj)
		if err != nil {
			return err
		}

		_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), namespaceObj, metav1.CreateOptions{})
		if err != nil {
			return err
		}

		fmt.Println("Namespace created:", namespaceObj.GetName())

	default:
		return fmt.Errorf("unsupported kind: %s", kind)
	}

	return nil
}
