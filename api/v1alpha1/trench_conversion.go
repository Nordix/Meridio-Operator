package v1alpha1

import (
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	meridiov1beta1 "github.com/nordix/meridio-operator/api/v1beta1"
)

// convertTo converts this Trench to Hub version (v1beta1)
func (src *Trench) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*meridiov1beta1.Trench)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta
	// Spec
	dst.Spec.Family = src.Spec.IPFamily
	// Status
	//...
	return nil
}

// convertTo converts from Hub version (v1beta1) to this version
func (dst *Trench) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*meridiov1beta1.Trench)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta
	// Spec
	dst.Spec.IPFamily = src.Spec.Family
	// Status
	//...
	return nil
}
