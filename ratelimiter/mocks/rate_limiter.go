// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// RateLimiter is an autogenerated mock type for the RateLimiter type
type RateLimiter struct {
	mock.Mock
}

// AllowRequest provides a mock function with given fields: userID
func (_m *RateLimiter) AllowRequest(userID string) bool {
	ret := _m.Called(userID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}