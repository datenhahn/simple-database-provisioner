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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneratePassword(t *testing.T) {

	password := GeneratePassword(10)
	t.Log("Password: " + password)

	assert.Equal(t, 10, len(password))
	assert.Regexp(t, "^[ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789]{10}$", password)

	password2 := GeneratePassword(20)
	t.Log("Password2: " + password2)

	assert.Equal(t, 20, len(password2))
	assert.Regexp(t, "^[ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789]{20}$", password2)
}

func TestMd5(t *testing.T) {

	inputString := "This is my md5 input"

	md5sum := Md5(inputString)
	assert.Equal(t, "3e8002928793cb7245bfd4ec09ad23a0", md5sum)
}

func TestMd5Short(t *testing.T) {

	inputString := "This is my md5 input"

	md5sum := Md5Short(inputString)
	assert.Equal(t, "3e800292", md5sum)
}
