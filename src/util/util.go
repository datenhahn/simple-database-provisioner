/*
 * Copyright (c) 2019 Ecodia GmbH & Co. KG <opensource@ecodia.de>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"crypto/md5"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// GeneratePassword generates a password from the folowing runes in a given length
// runes = ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789
func GeneratePassword(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()
	return str
}

// Md5 calculates the md5 sum hash of the supplied string
func Md5(input string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}

// Md5Short returns the last 8 characters of the md5 hash
// calculated for a string
func Md5Short(input string) string {
	return Md5(input)[:8]
}

// CreateTempFile returns the filename of a temporary file.
// Other than the ioutil.TempFile function, this function
// is not safe against race conditions, but it can be used
// to create filenames for tests which require a filename
// as input. The likeliness of collisions is small enough
// for that usecase.
func CreateTempFile() string {
	file, err := ioutil.TempFile("", "sdp-")
	defer os.Remove(file.Name())

	if err != nil {
		logrus.Panicf("Error creating tempfile: %v", err)
	}

	return file.Name()
}

// GetKubectlContext uses the kubectl commandline tool to return
// the current kubernetes context. Should only be used in test
// contexts.
func GetKubectlContext() string {

	out, err := exec.Command("kubectl", "config", "current-context").Output()
	if err != nil {
		logrus.Panic(err)
	}
	return strings.TrimSpace(string(out))
}

// PanicIfNotMinikube panics if the current kubectl context is not minikube.
func PanicIfNotMinikube() {

	cluster := GetKubectlContext()

	if cluster != "minikube" {
		logrus.Panicf("MUST RUN ON MINIKUBE, current cluster is: %s", cluster)
	}
}
