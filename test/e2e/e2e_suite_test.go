/*
Copyright (c) 2022 Nordix Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e_test

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/nordix/meridio/pkg/ipam/prefix"
	"github.com/nordix/meridio/test/e2e/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	namespace string
	script    string

	streamA1Name string

	lbfeDeploymentName string

	clientset *kubernetes.Clientset
)

const (
	timeout  = time.Minute * 3
	interval = time.Second * 2

	targetDeploymentName = "target-a"
	numberOfTargets      = 4
)

func init() {
	flag.StringVar(&namespace, "namespace", "red", "the namespace where expects operator to exist")
	flag.StringVar(&script, "script", "./data/kind/test.sh", "path + script used by the e2e tests")
	flag.StringVar(&streamA1Name, "stream-a-1-name", "stream-a-1", "Name of stream-a-1 (see e2e documentation diagram)")
	flag.StringVar(&lbfeDeploymentName, "lb-fe-deployment-name", "lb-fe-attractor-a-1", "Name of load-balancer deployment in trench-a")
}

func TestE2e(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2e Suite")
}

var _ = BeforeSuite(func() {
	prefix.GetFamily("")
	var err error
	clientset, err = utils.GetClientSet()
	Expect(err).ToNot(HaveOccurred())

	cmd := exec.Command(script, "init")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	Expect(stderr.String()).To(BeEmpty())
	Expect(err).ToNot(HaveOccurred())

	listOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", lbfeDeploymentName),
	}
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), listOptions)
	Expect(err).NotTo(HaveOccurred())
	Eventually(func() bool {
		for _, pod := range pods.Items {
			nfqlbOutput, err := utils.PodExec(&pod, "load-balancer", []string{"nfqlb", "show", fmt.Sprintf("--shm=tshm-%v", streamA1Name)})
			Expect(err).NotTo(HaveOccurred())
			if utils.ParseNFQLB(nfqlbOutput) != numberOfTargets {
				return false
			}
		}
		return true
	}, timeout, interval).Should(BeTrue())
})

var _ = AfterSuite(func() {
	cmd := exec.Command(script, "end")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	Expect(stderr.String()).To(BeEmpty())
	Expect(err).ToNot(HaveOccurred())
})
