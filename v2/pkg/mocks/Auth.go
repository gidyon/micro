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
	AuthAPI.On("AuthorizeFunc", mock.Anything).Return(context.Background(), nil)
	AuthAPI.On("AuthenticateRequest", mock.Anything).Return(nil)
	AuthAPI.On("AuthenticateRequestV2", mock.Anything).Return(&auth.Payload{Group: auth.DefaultAdminGroup()}, nil)
	AuthAPI.On("AuthorizeGroup", mock.Anything, mock.Anything).Return(&auth.Payload{Group: auth.DefaultAdminGroup()}, nil)
	AuthAPI.On("AuthorizeActor", mock.Anything, mock.Anything).Return(&auth.Payload{Group: auth.DefaultAdminGroup()}, nil)
	AuthAPI.On("AuthorizeActors", mock.Anything, mock.Anything).Return(&auth.Payload{Group: auth.DefaultAdminGroup()}, nil)
	AuthAPI.On("AuthorizeActorAndGroup", mock.Anything, mock.Anything, mock.Anything).Return(&auth.Payload{Group: auth.DefaultAdminGroup()}, nil)
	AuthAPI.On("AuthorizeActorOrGroup", mock.Anything, mock.Anything, mock.Anything).Return(&auth.Payload{Group: auth.DefaultAdminGroup()}, nil)
	AuthAPI.On("AuthorizeAdmin", mock.Anything).Return(&auth.Payload{Group: auth.DefaultAdminGroup()}, nil)
	AuthAPI.On("AuthorizeAdminStrict", mock.Anything, mock.Anything).Return(&auth.Payload{Group: auth.DefaultAdminGroup()}, nil)
	AuthAPI.On("AdminGroups").Return([]string{auth.DefaultAdminGroup(), auth.DefaultAdminGroup()})
	AuthAPI.On("IsAdmin", mock.Anything).Return(true)
	AuthAPI.On("GenToken", mock.Anything, mock.Anything, mock.Anything).Return("token", nil)
	AuthAPI.On("GetJwtPayload", mock.Anything).Return(&auth.Payload{Group: auth.DefaultAdminGroup()}, nil)
}
