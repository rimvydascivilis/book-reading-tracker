// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// OAuth2Service is an autogenerated mock type for the OAuth2Service type
type OAuth2Service struct {
	mock.Mock
}

type OAuth2Service_Expecter struct {
	mock *mock.Mock
}

func (_m *OAuth2Service) EXPECT() *OAuth2Service_Expecter {
	return &OAuth2Service_Expecter{mock: &_m.Mock}
}

// ValidateToken provides a mock function with given fields: token
func (_m *OAuth2Service) ValidateToken(token string) (string, error) {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for ValidateToken")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(token)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OAuth2Service_ValidateToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateToken'
type OAuth2Service_ValidateToken_Call struct {
	*mock.Call
}

// ValidateToken is a helper method to define mock.On call
//   - token string
func (_e *OAuth2Service_Expecter) ValidateToken(token interface{}) *OAuth2Service_ValidateToken_Call {
	return &OAuth2Service_ValidateToken_Call{Call: _e.mock.On("ValidateToken", token)}
}

func (_c *OAuth2Service_ValidateToken_Call) Run(run func(token string)) *OAuth2Service_ValidateToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *OAuth2Service_ValidateToken_Call) Return(_a0 string, _a1 error) *OAuth2Service_ValidateToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OAuth2Service_ValidateToken_Call) RunAndReturn(run func(string) (string, error)) *OAuth2Service_ValidateToken_Call {
	_c.Call.Return(run)
	return _c
}

// NewOAuth2Service creates a new instance of OAuth2Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOAuth2Service(t interface {
	mock.TestingT
	Cleanup(func())
}) *OAuth2Service {
	mock := &OAuth2Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
