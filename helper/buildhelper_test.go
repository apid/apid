/*
Copyright 2016 The Transicator Authors

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
package main

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
	"strings"
	"time"
)

// Make sure all plugins are accounted for
var apidpluginlist = []string{"apidVerifyApiKey", "apidGatewayConfDeploy", "apidApigeeSync", "apidAnalytics", "goscaffold", "apidCore"}

var _ = Describe("glide contains all apid plugins", func() {
	It("glide contains all apid plugins", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		cmd, err := exec.CommandContext(ctx, "./buildhelper", "../glide.lock").Output()
		Expect(err).NotTo(HaveOccurred())
		for _, v := range apidpluginlist {
			Ω(strings.Contains(string(cmd), v)).Should(BeTrue())
		}
	})

})
