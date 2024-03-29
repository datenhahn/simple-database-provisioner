// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// CustomResourceDefinitionManager is an autogenerated mock type for the CustomResourceDefinitionManager type
type CustomResourceDefinitionManager struct {
	mock.Mock
}

// InstallCustomResourceDefinition provides a mock function with given fields: crdYamlFilePath
func (_m *CustomResourceDefinitionManager) InstallCustomResourceDefinition(crdYamlFilePath string) error {
	ret := _m.Called(crdYamlFilePath)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(crdYamlFilePath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
