package common

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"

	meridiov1alpha1 "github.com/nordix/meridio-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	SAName                 = "meridio-sa"
	ResourceNamePrefixEnv  = "RESOURCE_NAME_PREFIX"
	ImagePullSecretEnv     = "IMAGE_PULL_SECRET"
	IdentifierRangeSizeEnv = "IDENTIFIER_RANGE_SIZE"

	Registry        = "registry.nordix.org"
	Organization    = "cloud-native/meridio"
	OrganizationNsm = "cloud-native/nsm"
	Tag             = "latest"

	BusyboxImage = "busybox"
	BusyboxTag   = "1.29"

	SubnetPoolIpv4          = "172.16.0.0/16"
	SubnetPoolIpv6          = "fd00::/48"
	ConduitPrefixLengthIpv4 = "20"
	ConduitPrefixLengthIpv6 = "56"
	NodePrefixLengthIpv4    = "24"
	NodePrefixLengthIpv6    = "64"

	NspPort        = 7778
	NspTargetPort  = 7778
	IpamPort       = 7777
	IpamTargetPort = 7777
	VlanNsePort    = 5003

	IpamSvcName = "ipam-service"
	NspSvcName  = "nsp-service"
	LBName      = "lb-fe"
	ProxyName   = "proxy"
	IpamName    = "ipam"
	NseName     = "nse-vlan"
	NspName     = "nsp"
	RlName      = "meridio-configuration-role"
	RBName      = "meridio-configuration-role-binding"
	CMName      = "meridio-configuration"

	NetworkServiceName = "external-vlan"

	IdentifierStart            = 1
	DefaultIdentifierRangeSize = 100
)

func ServiceAccountName(trench *meridiov1alpha1.Trench) string {
	return getSuffixedName(SAName, trench.ObjectMeta.Name)
}

func IPAMServiceName(trench *meridiov1alpha1.Trench) string {
	return getSuffixedName(IpamSvcName, trench.ObjectMeta.Name)
}

func NSPServiceName(trench *meridiov1alpha1.Trench) string {
	return getSuffixedName(NspSvcName, trench.ObjectMeta.Name)
}

func LbFeDeploymentName(attractor *meridiov1alpha1.Attractor) string {
	return getSuffixedName(LBName, attractor.ObjectMeta.Name)
}

func ProxyDeploymentName(trench *meridiov1alpha1.Trench) string {
	return getSuffixedName(ProxyName, trench.ObjectMeta.Name)
}

func IPAMStatefulSetName(trench *meridiov1alpha1.Trench) string {
	return getSuffixedName(IpamName, trench.ObjectMeta.Name)
}

func NSEDeploymentName(attractor *meridiov1alpha1.Attractor) string {
	return getSuffixedName(NseName, attractor.ObjectMeta.Name)
}

func NSPStatefulSetName(trench *meridiov1alpha1.Trench) string {
	return getSuffixedName(NspName, trench.ObjectMeta.Name)
}

func RoleName(trench *meridiov1alpha1.Trench) string {
	return getSuffixedName(RlName, trench.ObjectMeta.Name)
}

func RoleBindingName(trench *meridiov1alpha1.Trench) string {
	return getSuffixedName(RBName, trench.ObjectMeta.Name)
}

func ConfigMapName(trench *meridiov1alpha1.Trench) string {
	return getSuffixedName(CMName, trench.ObjectMeta.Name)
}

func NSPServiceWithPort(cr *meridiov1alpha1.Trench) string {
	return fmt.Sprintf("%s:%d", NSPServiceName(cr), NspTargetPort)
}

func IPAMServiceWithPort(cr *meridiov1alpha1.Trench) string {
	return fmt.Sprintf("%s:%d", IPAMServiceName(cr), IpamTargetPort)
}

func getSuffixedName(resourceName, suffix string) string {
	return fmt.Sprintf("%s%s-%s", getResourceNamePrefix(), resourceName, suffix)
}

func ProxyNtwkSvcNsName(cr *meridiov1alpha1.Conduit) string {
	return strings.Join([]string{ProxyName, cr.ObjectMeta.Name, cr.ObjectMeta.Labels["trench"], cr.ObjectMeta.Namespace}, ".")
}

func LoadBalancerNsName(conduit, trench, namespace string) string {
	return strings.Join([]string{conduit, trench, namespace}, ".")
}

func VlanNtwkSvcName(cr *meridiov1alpha1.Trench) string {
	return strings.Join([]string{NetworkServiceName, cr.ObjectMeta.Name, cr.ObjectMeta.Namespace}, ".")
}

func getResourceNamePrefix() string {
	return os.Getenv(ResourceNamePrefixEnv)
}

func GetIdentifierRangeSize() int {
	IdentifierRangeSize, err := strconv.Atoi(os.Getenv(IdentifierRangeSizeEnv))
	if err != nil {
		return DefaultIdentifierRangeSize
	}
	return IdentifierRangeSize
}

func GetImagePullSecrets() []corev1.LocalObjectReference {
	secstr := os.Getenv(ImagePullSecretEnv)
	secs := strings.Split(secstr, ",")
	var pullSecs []corev1.LocalObjectReference
	for _, sec := range secs {
		pullSecs = append(pullSecs, corev1.LocalObjectReference{
			Name: strings.TrimSpace(sec),
		})
	}
	return pullSecs
}

func NsName(meta metav1.ObjectMeta) string {
	return fmt.Sprintf("%s/%s", meta.Namespace, meta.Name)
}

func GetExternalInterfaceName(attractor *meridiov1alpha1.Attractor) string {
	return fmt.Sprintf("ext-vlan.%d", *attractor.Spec.Interface.NSMVlan.VlanID)
}
