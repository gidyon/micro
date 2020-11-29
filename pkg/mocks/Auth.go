package mocks

import (
	"context"

	"github.com/gidyon/micro/pkg/grpc/auth"
	"github.com/gidyon/micro/pkg/mocks/mocks"
	"github.com/stretchr/testify/mock"
)

// AuthAPIMock is auth API
type AuthAPIMock interface {
	auth.API
}

// AuthAPI is a fake authentication API
var AuthAPI = &mocks.AuthAPIMock{}

func init() {
	AuthAPI.On("AuthFunc", mock.Anything).
		Return(context.Background(), nil)
	AuthAPI.On("AuthenticateRequestV2", mock.Anything).
		Return(&auth.Payload{Group: auth.AdminGroup()}, nil)
	AuthAPI.On("AuthenticateRequest", mock.Anything).
		Return(nil)
	AuthAPI.On("AuthorizeActor", mock.Anything, mock.Anything).
		Return(&auth.Payload{Group: auth.AdminGroup()}, nil)
	AuthAPI.On("AuthorizeActors", mock.Anything, mock.Anything).
		Return(&auth.Payload{Group: auth.AdminGroup()}, nil)
	AuthAPI.On("AuthorizeGroups",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&auth.Payload{Group: auth.AdminGroup()}, nil)
	AuthAPI.On("AuthorizeStrict",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&auth.Payload{Group: auth.AdminGroup()}, nil)
	AuthAPI.On("AuthorizeActorOrGroups",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&auth.Payload{Group: auth.AdminGroup()}, nil)
	AuthAPI.On("GenToken", mock.Anything, mock.Anything, mock.Anything).
		Return("token", nil)
}
