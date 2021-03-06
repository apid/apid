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

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

/*
 * Program takes glide.lock as command line arg and dumps to
 * stdout the go build command with all the plugins based on
 * glide.lock. If error encountered the stderr will have the
 * error details
 */
func main() {
	var buildstr string = "go build -ldflags \""
	var key, value, str string
	glideLock := os.Args[1]
	file, err := os.Open(glideLock)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open glide.lock")
		return
	}
	defer file.Close()
	m := make(map[string]string)
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read glide.lock")
		return
	}
	r1 := regexp.MustCompile("- name: github.com/apid/(.*)")
	r2 := regexp.MustCompile("^\\s*version: (.*)")
	for scanner.Scan() {
		line := scanner.Text()
		status1 := r1.FindStringSubmatch(line)
		if status1 != nil {
			key = status1[1]
		}
		status2 := r2.FindStringSubmatch(line)
		if status2 != nil {
			value = status2[1]
		}
		if key != "" && value != "" {
			if key == "apid-core" {
				key = "apidCore"
			}
			m[key] = value
			key = ""
			value = ""
		}
	}
	for k, v := range m {
		str = " -X main." + k + "=" + v
		buildstr += str
	}
	if len(m) == 0 {
		fmt.Fprintf(os.Stderr, "Is glide.lock corrupted?\n")
	} else {
		buildstr += "\""
		fmt.Fprintf(os.Stdout, "%s", buildstr)
	}

}
