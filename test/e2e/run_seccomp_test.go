package integration

import (
	. "github.com/containers/podman/v5/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Podman run", func() {

	It("podman run --seccomp-policy default", func() {
		session := podmanTest.Podman([]string{"run", "-q", "--seccomp-policy", "default", alpineSeccomp, "ls"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(ExitCleanly())
	})

	It("podman run --seccomp-policy ''", func() {
		// Empty string is interpreted as "default".
		session := podmanTest.Podman([]string{"run", "-q", "--seccomp-policy", "", alpineSeccomp, "ls"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(ExitCleanly())
	})

	It("podman run --seccomp-policy invalid", func() {
		session := podmanTest.Podman([]string{"run", "--seccomp-policy", "invalid", alpineSeccomp, "ls"})
		session.WaitWithDefaultTimeout()
		Expect(session).To(ExitWithError(125, `invalid seccomp policy "invalid": valid policies are ["default" "image"]`))
	})

	It("podman run --seccomp-policy image (block all syscalls)", func() {
		session := podmanTest.Podman([]string{"run", "--seccomp-policy", "image", alpineSeccomp, "ls"})
		session.WaitWithDefaultTimeout()
		// TODO: we're getting a "cannot start a container that has
		//       stopped" error which seems surprising.  Investigate
		//       why that is so.
		if podmanTest.OCIRuntime == "runc" {
			// TODO: worse than that. With runc, we get two alternating failures:
			//   126 + cannot start a container that has stopped
			//   127 + failed to connect to container's attach socket ... ENOENT
			Expect(session.ExitCode()).To(BeNumerically(">=", 126), "Exit status using runc")
		} else {
			expect := "OCI runtime error: crun: read from the init process"
			if IsRemote() {
				expect = "for attach: crun: read from the init process: OCI runtime error"
			}
			Expect(session).To(ExitWithError(126, expect))
		}
	})

	It("podman run --seccomp-policy image (bogus profile)", func() {
		session := podmanTest.Podman([]string{"run", "--seccomp-policy", "image", alpineBogusSeccomp, "ls"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(ExitWithError(125, "loading seccomp profile failed: decoding seccomp profile failed: invalid character 'B' looking for beginning of value"))
	})
})
