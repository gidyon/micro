package mocks

import (
	"context"

	"github.com/gidyon/micro/v2/pkg/middleware/grpc/auth"
	"github.com/gidyon/micro/v2/pkg/mocks/mocks"
	"github.com/stretchr/testify/mock"
)

// AuthAPIMock is auth API
type AuthAPIMock interface {
	auth.API
}

// AuthAPI is a fake authentication API
var AuthAPI = &mocks.AuthAPIMock{}

func init() {
	AuthAPI.On("AuthenticateRequest", mock.Anything).Return(nil)
	AuthAPI.On("AuthenticateRequestV2", mock.Anything).Return(&auth.Payload{Group: auth.DefaultAdminGroup()}, nil)
	AuthAPI.On("AuthorizeFunc", anything()...).Return(context.Background(), nil)
	AuthAPI.On("AuthorizeGroup", anything()...).Return(&auth.Payload{Group: auth.DefaultAdminGroup()}, nil)
	AuthAPI.On("AuthorizeActor", anything()...).Return(&auth.Payload{Group: auth.DefaultAdminGroup()}, nil)
	AuthAPI.On("AuthorizeActors", anything()...).Return(&auth.Payload{Group: auth.DefaultAdminGroup()}, nil)
	AuthAPI.On("AuthorizeActorAndGroup", anything()...).Return(&auth.Payload{Group: auth.DefaultAdminGroup()}, nil)
	AuthAPI.On("AuthorizeActorOrGroup", anything()...).Return(&auth.Payload{Group: auth.DefaultAdminGroup()}, nil)
	AuthAPI.On("AuthorizeAdmin", anything()...).Return(&auth.Payload{Group: auth.DefaultAdminGroup()}, nil)
	AuthAPI.On("AuthorizeAdminStrict", anything()...).Return(&auth.Payload{Group: auth.DefaultAdminGroup()}, nil)
	AuthAPI.On("AdminGroups").Return([]string{auth.DefaultAdminGroup(), auth.DefaultAdminGroup()})
	AuthAPI.On("IsAdmin", mock.Anything).Return(true)
	AuthAPI.On("GenToken", anything()...).Return("token", nil)
	AuthAPI.On("GetJwtPayload", anything()...).Return(&auth.Payload{Group: auth.DefaultAdminGroup()}, nil)
}

func anything() []interface{} {
	v := make([]interface{}, 0, 10)
	for i := 0; i < 10; i++ {
		v = append(v, mock.Anything)
	}
	return v
}
