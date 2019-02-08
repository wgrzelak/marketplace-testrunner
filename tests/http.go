// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests // import "github.com/GoogleCloudPlatform/marketplace-testrunner/tests"

import (
	"fmt"
	"net/http"

	"io/ioutil"

	"github.com/GoogleCloudPlatform/marketplace-testrunner/asserts"
	"github.com/GoogleCloudPlatform/marketplace-testrunner/specs"
)

// RunHttpTest executes an HttpTest rule, returning an empty string if
// the test passes, otherwise the error message.
func RunHttpTest(test *specs.HttpTest, client *http.Client, expandCommand func(string) string) string {
	if msg := validate(test); msg != "" {
		return asserts.MessageWithContext(msg, "Malformed HttpTest")
	}
	method := http.MethodGet
	if test.Method != nil {
		method = expandCommand(*test.Method)
	}
	req, err := http.NewRequest(method, expandCommand(test.Url), nil)
	if err != nil {
		return fmt.Sprintf("HTTP request creation error: %v", err)
	}
	for k, v := range test.Headers {
		req.Header.Add(k, expandCommand(v))
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("HTTP request error: %v", err)
	}
	return doExpect(&test.Expect, res, expandCommand)
}

func doExpect(expect *specs.HttpExpect, response *http.Response, expandCommand func(string) string) string {
	if expect.StatusCode != nil {
		if msg := asserts.DoAssert(response.StatusCode, expect.StatusCode, expandCommand); msg != "" {
			return asserts.MessageWithContext(msg, "Unexpected response status code")
		}
	}
	if expect.StatusText != nil {
		if msg := asserts.DoAssert(response.Status, expect.StatusText, expandCommand); msg != "" {
			return asserts.MessageWithContext(msg, "Unexpected response status text")
		}
	}
	if expect.BodyText != nil {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "Unexpected error reading the body"
		}
		if msg := asserts.DoAssert(string(body), expect.BodyText, expandCommand); msg != "" {
			return asserts.MessageWithContext(msg, "Unexpected body text content")
		}
	}
	return ""
}

func validate(test *specs.HttpTest) string {
	if len(test.Url) <= 0 {
		return "Url is required"
	}
	return ""
}
