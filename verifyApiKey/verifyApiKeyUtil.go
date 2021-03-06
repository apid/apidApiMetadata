// Copyright 2017 Google Inc.
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

package verifyApiKey

import (
	"regexp"
	"strings"
)

/*
 * Check for the base path (API_Product) match with the path
 * received in the Request, via the customized regex, where
 * "**" gets de-normalized as ".*" and "*" as everything till
 * the next "/".
 */
func validatePath(fs []string, requestBase string) bool {

	for _, a := range fs {
		str1 := strings.Replace(a, "**", "(.*)", -1)
		str2 := strings.Replace(a, "*", "([^/]+)", -1)
		if a != str1 {
			reg, _ := regexp.Compile(str1)
			res := reg.MatchString(requestBase)
			if res == true {
				return true
			}
		} else if a != str2 {
			reg, _ := regexp.Compile(str2)
			res := reg.MatchString(requestBase)
			if res == true {
				return true
			}
		} else if requestBase == a {
			return true
		}

		/*
		 * FIXME: SINGLE_FORWARD_SLASH_PATTERN not supported yet
		 */
	}
	/* if the i/p resource is empty, no checks need to be made */
	return len(fs) == 0
}
