package kubernetes

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func LookupService(ctx context.Context, context *string, namespace, service string) ([]*v1.Deployment, error) {
	var err error

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	if context != nil {
		configOverrides.CurrentContext = *context
	}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	svc, err := clientset.CoreV1().Services(namespace).Get(ctx, service, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	deps, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	matchingDep := []*v1.Deployment{}
	for idx, dep := range deps.Items {
		match := true
		for k, v := range svc.Spec.Selector {
			if val, ok := dep.Spec.Template.ObjectMeta.Labels[k]; ok {
				if v != val {
					match = false
					break
				}
			} else {
				match = false
				break
			}
		}
		if match {
			matchingDep = append(matchingDep, &deps.Items[idx])
		}
	}

	return matchingDep, nil
}
