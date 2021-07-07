package reconciler

import (
	"fmt"
	"os"

	meridiov1alpha1 "github.com/nordix/meridio-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"

	"k8s.io/apimachinery/pkg/util/yaml"
)

const (
	Registry     = "registry.nordix.org"
	Organization = "cloud-native/meridio"
	Tag          = "latest"
	PullPolicy   = "IfNotPresent"
)

func GetReadinessProbe(cr *meridiov1alpha1.Trench) *corev1.Probe {
	// if readiness probe is set in the cr do something
	// else use the default readiness probe
	return &corev1.Probe{
		Handler: corev1.Handler{
			Exec: &corev1.ExecAction{
				Command: []string{"/bin/grpc_health_probe", "-addr=:8000", "-connect-timeout=100ms", "-rpc-timeout=150ms"},
			},
		},
		InitialDelaySeconds: 0,
		PeriodSeconds:       10,
		TimeoutSeconds:      3,
		FailureThreshold:    5,
	}
}

func GetLivenessProbe(cr *meridiov1alpha1.Trench) *corev1.Probe {
	// if liveness probe is set in the cr do something
	// else use the default liveness probe
	return &corev1.Probe{
		Handler: corev1.Handler{
			Exec: &corev1.ExecAction{
				Command: []string{"/bin/grpc_health_probe", "-addr=:8000", "-connect-timeout=100ms", "-rpc-timeout=150ms"},
			},
		},
		InitialDelaySeconds: 0,
		PeriodSeconds:       10,
		TimeoutSeconds:      3,
		FailureThreshold:    5,
	}
}

func getDeploymentModel(f string) (*appsv1.Deployment, error) {
	data, err := os.Open(f)
	if err != nil {
		return nil, fmt.Errorf("open %s error: %s", f, err)
	}
	deployment := &appsv1.Deployment{}
	err = yaml.NewYAMLOrJSONDecoder(data, 4096).Decode(deployment)
	if err != nil {
		return nil, fmt.Errorf("decode %s error: %s", f, err)
	}
	return deployment, nil
}

func getServiceModel(f string) (*corev1.Service, error) {
	data, err := os.Open(f)
	if err != nil {
		return nil, fmt.Errorf("open %s error: %s", f, err)
	}
	service := &corev1.Service{}
	err = yaml.NewYAMLOrJSONDecoder(data, 4096).Decode(service)
	if err != nil {
		return nil, fmt.Errorf("decode %s error: %s", f, err)
	}
	return service, nil
}

func getRoleModel(f string) (*rbacv1.Role, error) {
	data, err := os.Open(f)
	if err != nil {
		return nil, fmt.Errorf("open %s error: %s", f, err)
	}
	role := &rbacv1.Role{}
	err = yaml.NewYAMLOrJSONDecoder(data, 4096).Decode(role)
	if err != nil {
		return nil, fmt.Errorf("decode %s error: %s", f, err)
	}
	return role, nil
}

func getRoleBindingModel(f string) (*rbacv1.RoleBinding, error) {
	data, err := os.Open(f)
	if err != nil {
		return nil, fmt.Errorf("open %s error: %s", f, err)
	}
	rb := &rbacv1.RoleBinding{}
	err = yaml.NewYAMLOrJSONDecoder(data, 4096).Decode(rb)
	if err != nil {
		return nil, fmt.Errorf("decode %s error: %s", f, err)
	}
	return rb, nil
}
