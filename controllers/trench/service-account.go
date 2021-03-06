package trench

import (
	meridiov1alpha1 "github.com/nordix/meridio-operator/api/v1alpha1"
	common "github.com/nordix/meridio-operator/controllers/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ServiceAccount struct {
	trench *meridiov1alpha1.Trench
	model  *corev1.ServiceAccount
	exec   *common.Executor
}

func NewServiceAccount(e *common.Executor, t *meridiov1alpha1.Trench) (*ServiceAccount, error) {
	l := &ServiceAccount{
		trench: t.DeepCopy(),
		exec:   e,
	}

	// get model
	if err := l.getModel(); err != nil {
		return nil, err
	}
	return l, nil
}

func (sa *ServiceAccount) getSelector() client.ObjectKey {
	return client.ObjectKey{
		Namespace: sa.trench.ObjectMeta.Namespace,
		Name:      common.ServiceAccountName(sa.trench),
	}
}

func (i *ServiceAccount) getModel() error {
	model, err := common.GetServiceAccountModel("deployment/service-account.yaml")
	if err != nil {
		return err
	}
	i.model = model
	return nil
}

func (sa *ServiceAccount) insertParameters(init *corev1.ServiceAccount) *corev1.ServiceAccount {
	ret := init.DeepCopy()
	ret.ObjectMeta.Name = common.ServiceAccountName(sa.trench)
	ret.ObjectMeta.Namespace = sa.trench.ObjectMeta.Namespace
	return ret
}

func (sa *ServiceAccount) getCurrentStatus() (*corev1.ServiceAccount, error) {
	currentState := &corev1.ServiceAccount{}
	selector := sa.getSelector()
	err := sa.exec.GetObject(selector, currentState)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return currentState, nil
}

func (sa *ServiceAccount) getDesiredStatus() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      common.ServiceAccountName(sa.trench),
			Namespace: sa.trench.ObjectMeta.Namespace,
		},
	}
}

func (sa *ServiceAccount) getReconciledDesiredStatus(current *corev1.ServiceAccount) *corev1.ServiceAccount {
	return sa.insertParameters(current)
}

func (sa *ServiceAccount) getAction() error {
	cs, err := sa.getCurrentStatus()
	if err != nil {
		return err
	}
	if cs == nil {
		ds := sa.getDesiredStatus()
		sa.exec.AddCreateAction(ds)
	} else {
		ds := sa.getReconciledDesiredStatus(cs)
		if !equality.Semantic.DeepEqual(ds, cs) {
			sa.exec.AddUpdateAction(ds)
		}
	}
	return nil
}
