// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import time "time"

// Cache is an autogenerated mock type for the Cache type
type Cache struct {
	mock.Mock
}

// Del provides a mock function with given fields: key
func (_m *Cache) Del(key string) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: key
func (_m *Cache) Get(key string) (string, error) {
	ret := _m.Called(key)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Set provides a mock function with given fields: key, value, duration
func (_m *Cache) Set(key string, value string, duration time.Duration) (string, error) {
	ret := _m.Called(key, value, duration)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, time.Duration) string); ok {
		r0 = rf(key, value, duration)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, time.Duration) error); ok {
		r1 = rf(key, value, duration)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
