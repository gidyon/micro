// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	auth "github.com/gidyon/micro/v2/pkg/middleware/grpc/auth"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// AuthAPIMock is an autogenerated mock type for the AuthAPIMock type
type AuthAPIMock struct {
	mock.Mock
}

// AdminGroups provides a mock function with given fields:
func (_m *AuthAPIMock) AdminGroups() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// AuthenticateRequest provides a mock function with given fields: ctx
func (_m *AuthAPIMock) AuthenticateRequest(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthenticateRequestV2 provides a mock function with given fields: ctx
func (_m *AuthAPIMock) AuthenticateRequestV2(ctx context.Context) (*auth.Payload, error) {
	ret := _m.Called(ctx)

	var r0 *auth.Payload
	if rf, ok := ret.Get(0).(func(context.Context) *auth.Payload); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Payload)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthorizeActor provides a mock function with given fields: ctx, actorID
func (_m *AuthAPIMock) AuthorizeActor(ctx context.Context, actorID string) (*auth.Payload, error) {
	ret := _m.Called(ctx, actorID)

	var r0 *auth.Payload
	if rf, ok := ret.Get(0).(func(context.Context, string) *auth.Payload); ok {
		r0 = rf(ctx, actorID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Payload)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, actorID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthorizeActorAndGroup provides a mock function with given fields: ctx, actorID, allowedGroups
func (_m *AuthAPIMock) AuthorizeActorAndGroup(ctx context.Context, actorID string, allowedGroups ...string) (*auth.Payload, error) {
	_va := make([]interface{}, len(allowedGroups))
	for _i := range allowedGroups {
		_va[_i] = allowedGroups[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, actorID)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *auth.Payload
	if rf, ok := ret.Get(0).(func(context.Context, string, ...string) *auth.Payload); ok {
		r0 = rf(ctx, actorID, allowedGroups...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Payload)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, ...string) error); ok {
		r1 = rf(ctx, actorID, allowedGroups...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthorizeActorOrGroup provides a mock function with given fields: ctx, actorID, allowedGroups
func (_m *AuthAPIMock) AuthorizeActorOrGroup(ctx context.Context, actorID string, allowedGroups ...string) (*auth.Payload, error) {
	_va := make([]interface{}, len(allowedGroups))
	for _i := range allowedGroups {
		_va[_i] = allowedGroups[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, actorID)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *auth.Payload
	if rf, ok := ret.Get(0).(func(context.Context, string, ...string) *auth.Payload); ok {
		r0 = rf(ctx, actorID, allowedGroups...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Payload)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, ...string) error); ok {
		r1 = rf(ctx, actorID, allowedGroups...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthorizeActors provides a mock function with given fields: ctx, actorID
func (_m *AuthAPIMock) AuthorizeActors(ctx context.Context, actorID ...string) (*auth.Payload, error) {
	_va := make([]interface{}, len(actorID))
	for _i := range actorID {
		_va[_i] = actorID[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *auth.Payload
	if rf, ok := ret.Get(0).(func(context.Context, ...string) *auth.Payload); ok {
		r0 = rf(ctx, actorID...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Payload)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...string) error); ok {
		r1 = rf(ctx, actorID...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthorizeAdmin provides a mock function with given fields: ctx
func (_m *AuthAPIMock) AuthorizeAdmin(ctx context.Context) (*auth.Payload, error) {
	ret := _m.Called(ctx)

	var r0 *auth.Payload
	if rf, ok := ret.Get(0).(func(context.Context) *auth.Payload); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Payload)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthorizeAdminStrict provides a mock function with given fields: ctx, adminID
func (_m *AuthAPIMock) AuthorizeAdminStrict(ctx context.Context, adminID string) (*auth.Payload, error) {
	ret := _m.Called(ctx, adminID)

	var r0 *auth.Payload
	if rf, ok := ret.Get(0).(func(context.Context, string) *auth.Payload); ok {
		r0 = rf(ctx, adminID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Payload)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, adminID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthorizeFunc provides a mock function with given fields: ctx
func (_m *AuthAPIMock) AuthorizeFunc(ctx context.Context) (context.Context, error) {
	ret := _m.Called(ctx)

	var r0 context.Context
	if rf, ok := ret.Get(0).(func(context.Context) context.Context); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthorizeGroup provides a mock function with given fields: ctx, allowedGroups
func (_m *AuthAPIMock) AuthorizeGroup(ctx context.Context, allowedGroups ...string) (*auth.Payload, error) {
	_va := make([]interface{}, len(allowedGroups))
	for _i := range allowedGroups {
		_va[_i] = allowedGroups[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *auth.Payload
	if rf, ok := ret.Get(0).(func(context.Context, ...string) *auth.Payload); ok {
		r0 = rf(ctx, allowedGroups...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Payload)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...string) error); ok {
		r1 = rf(ctx, allowedGroups...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenToken provides a mock function with given fields: ctx, payload, expires
func (_m *AuthAPIMock) GenToken(ctx context.Context, payload *auth.Payload, expires time.Time) (string, error) {
	ret := _m.Called(ctx, payload, expires)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, *auth.Payload, time.Time) string); ok {
		r0 = rf(ctx, payload, expires)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *auth.Payload, time.Time) error); ok {
		r1 = rf(ctx, payload, expires)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetJwtPayload provides a mock function with given fields: ctx
func (_m *AuthAPIMock) GetJwtPayload(ctx context.Context) (*auth.Payload, error) {
	ret := _m.Called(ctx)

	var r0 *auth.Payload
	if rf, ok := ret.Get(0).(func(context.Context) *auth.Payload); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Payload)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPayloadFromJwt provides a mock function with given fields: jwt
func (_m *AuthAPIMock) GetPayloadFromJwt(jwt string) (*auth.Payload, error) {
	ret := _m.Called(jwt)

	var r0 *auth.Payload
	if rf, ok := ret.Get(0).(func(string) *auth.Payload); ok {
		r0 = rf(jwt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Payload)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(jwt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsAdmin provides a mock function with given fields: group
func (_m *AuthAPIMock) IsAdmin(group string) bool {
	ret := _m.Called(group)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(group)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
