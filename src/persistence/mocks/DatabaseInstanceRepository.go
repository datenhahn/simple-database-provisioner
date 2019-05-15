// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import persistence "simple-database-provisioner/src/persistence"

// DatabaseInstanceRepository is an autogenerated mock type for the DatabaseInstanceRepository type
type DatabaseInstanceRepository struct {
	mock.Mock
}

// AddDatabaseInstance provides a mock function with given fields: instance
func (_m *DatabaseInstanceRepository) AddDatabaseInstance(instance persistence.DatabaseInstance) error {
	ret := _m.Called(instance)

	var r0 error
	if rf, ok := ret.Get(0).(func(persistence.DatabaseInstance) error); ok {
		r0 = rf(instance)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteDatabaseInstance provides a mock function with given fields: bindingInstance
func (_m *DatabaseInstanceRepository) DeleteDatabaseInstance(bindingInstance persistence.NamespaceUniqueId) error {
	ret := _m.Called(bindingInstance)

	var r0 error
	if rf, ok := ret.Get(0).(func(persistence.NamespaceUniqueId) error); ok {
		r0 = rf(bindingInstance)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAllDatabaseInstances provides a mock function with given fields:
func (_m *DatabaseInstanceRepository) FindAllDatabaseInstances() []persistence.DatabaseInstance {
	ret := _m.Called()

	var r0 []persistence.DatabaseInstance
	if rf, ok := ret.Get(0).(func() []persistence.DatabaseInstance); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]persistence.DatabaseInstance)
		}
	}

	return r0
}

// FindDatabaseInstanceById provides a mock function with given fields: instanceId
func (_m *DatabaseInstanceRepository) FindDatabaseInstanceById(instanceId persistence.NamespaceUniqueId) (persistence.DatabaseInstance, error) {
	ret := _m.Called(instanceId)

	var r0 persistence.DatabaseInstance
	if rf, ok := ret.Get(0).(func(persistence.NamespaceUniqueId) persistence.DatabaseInstance); ok {
		r0 = rf(instanceId)
	} else {
		r0 = ret.Get(0).(persistence.DatabaseInstance)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(persistence.NamespaceUniqueId) error); ok {
		r1 = rf(instanceId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateDatabaseInstanceCredentials provides a mock function with given fields: instanceId, newCredentials
func (_m *DatabaseInstanceRepository) UpdateDatabaseInstanceCredentials(instanceId persistence.NamespaceUniqueId, newCredentials map[string][]byte) error {
	ret := _m.Called(instanceId, newCredentials)

	var r0 error
	if rf, ok := ret.Get(0).(func(persistence.NamespaceUniqueId, map[string][]byte) error); ok {
		r0 = rf(instanceId, newCredentials)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateDatabaseInstanceState provides a mock function with given fields: instanceId, newState
func (_m *DatabaseInstanceRepository) UpdateDatabaseInstanceState(instanceId persistence.NamespaceUniqueId, newState persistence.State) error {
	ret := _m.Called(instanceId, newState)

	var r0 error
	if rf, ok := ret.Get(0).(func(persistence.NamespaceUniqueId, persistence.State) error); ok {
		r0 = rf(instanceId, newState)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}