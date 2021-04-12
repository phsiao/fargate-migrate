package kubernetes

import (
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NormalizeEnvironment(ctx context.Context, context *string, namespace string, containers []corev1.Container) ([]corev1.Container, error) {

	clientset, err := getClient(context)
	if err != nil {
		return nil, err
	}

	// expand and filter environment variables
	for idx, container := range containers {
		env := []corev1.EnvVar{}
		for _, e := range container.Env {
			if e.Value != "" {
				env = append(env, e)
				continue
			}
			if e.ValueFrom != nil {
				if e.ValueFrom.ConfigMapKeyRef != nil {
					cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(ctx, e.ValueFrom.ConfigMapKeyRef.Name, metav1.GetOptions{})
					if err != nil {
						return nil, err
					}
					if val, ok := cm.Data[e.ValueFrom.ConfigMapKeyRef.Key]; ok {
						env = append(env, corev1.EnvVar{
							Name:  e.Name,
							Value: val,
						})
					}
					continue
				}
				if e.ValueFrom.SecretKeyRef != nil {
					sec, err := clientset.CoreV1().Secrets(namespace).Get(ctx, e.ValueFrom.SecretKeyRef.Name, metav1.GetOptions{})
					if err != nil {
						return nil, err
					}
					if val, ok := sec.Data[e.ValueFrom.SecretKeyRef.Key]; ok {
						env = append(env, corev1.EnvVar{
							Name:  e.Name,
							Value: string(val),
						})
					}
					continue
				}
			}
		}

		for _, e := range container.EnvFrom {
			if e.ConfigMapRef != nil {
				cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(ctx, e.ConfigMapRef.Name, metav1.GetOptions{})
				if err != nil {
					return nil, err
				}
				for k, v := range cm.Data {
					if e.Prefix != "" {
						if strings.HasPrefix(k, e.Prefix+"_") {
							k = k[len(e.Prefix)+1:]
						}
					}
					env = append(env, corev1.EnvVar{
						Name:  k,
						Value: v,
					})
				}
				continue
			}
			if e.SecretRef != nil {
				sec, err := clientset.CoreV1().Secrets(namespace).Get(ctx, e.SecretRef.Name, metav1.GetOptions{})
				if err != nil {
					return nil, err
				}
				for k, v := range sec.Data {
					if e.Prefix != "" {
						if strings.HasPrefix(k, e.Prefix+"_") {
							k = k[len(e.Prefix)+1:]
						}
					}
					env = append(env, corev1.EnvVar{
						Name:  k,
						Value: string(v),
					})
				}
				continue
			}
		}

		containers[idx].Env = env
	}

	return containers, nil
}
