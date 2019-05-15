// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// EventRepository is an autogenerated mock type for the EventRepository type
type EventRepository struct {
	mock.Mock
}

// MarkProcessed provides a mock function with given fields: eventId
func (_m *EventRepository) MarkProcessed(eventId string) {
	_m.Called(eventId)
}

// WasProcessed provides a mock function with given fields: eventId
func (_m *EventRepository) WasProcessed(eventId string) bool {
	ret := _m.Called(eventId)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(eventId)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
