package haproxy_test

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Deploy workload", func() {

	It("exposes routes via haproxy", func() {

		appUrl := fmt.Sprintf("http://%s:%s", workerAddress, nodePort)

		timeout := time.Duration(5 * time.Second)
		httpClient := http.Client{
			Timeout: timeout,
		}

		_, err := httpClient.Get(appUrl)
		Expect(err).To(HaveOccurred())

		deployNginx := runner.RunKubectlCommand("create", "-f", nginxSpec)
		Eventually(deployNginx, "60s").Should(gexec.Exit(0))
		rolloutWatch := runner.RunKubectlCommand("rollout", "status", "deployment/nginx", "-w")
		Eventually(rolloutWatch, "120s").Should(gexec.Exit(0))

		Eventually(func() string {
			result, err := httpClient.Get(appUrl)
			if err != nil {
				return err.Error()
			}
			return result.Status
		}, "120s", "5s").Should(Equal("200 OK"))
	})

	AfterEach(func() {
		session := runner.RunKubectlCommand("delete", "-f", nginxSpec)
		session.Wait("60s")
	})

})
