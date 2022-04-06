package mocks

import (
	"context"
	"fmt"

	"github.com/Pallinder/go-randomdata"
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
	AuthAPI.On("AuthenticateRequestV2", mock.Anything).Return(mockPayload(), nil)
	AuthAPI.On("AuthorizeFunc", anything()...).Return(context.Background(), nil)
	AuthAPI.On("AuthorizeGroup", anything()...).Return(mockPayload(), nil)
	AuthAPI.On("AuthorizeActor", anything()...).Return(mockPayload(), nil)
	AuthAPI.On("AuthorizeActors", anything()...).Return(mockPayload(), nil)
	AuthAPI.On("AuthorizeActorAndGroup", anything()...).Return(mockPayload(), nil)
	AuthAPI.On("AuthorizeActorOrGroup", anything()...).Return(mockPayload(), nil)
	AuthAPI.On("AuthorizeAdmin", anything()...).Return(mockPayload(), nil)
	AuthAPI.On("AuthorizeAdminStrict", anything()...).Return(mockPayload(), nil)
	AuthAPI.On("AdminGroups").Return([]string{auth.DefaultAdminGroup(), auth.DefaultAdminGroup()})
	AuthAPI.On("AddAdminGroups").Return()
	AuthAPI.On("IsAdmin", mock.Anything).Return(true)
	AuthAPI.On("GenToken", anything()...).Return("token", nil)
	AuthAPI.On("GetJwtPayload", anything()...).Return(mockPayload(), nil)
	AuthAPI.On("GetPayloadFromJwt", anything()...).Return(mockPayload(), nil)
}

func anything() []interface{} {
	v := make([]interface{}, 0, 10)
	for i := 0; i < 10; i++ {
		v = append(v, mock.Anything)
	}
	return v
}

func mockPayload() *auth.Payload {
	return &auth.Payload{
		Group:        auth.DefaultAdminGroup(),
		Names:        randomdata.SillyName(),
		ProjectID:    "test",
		EmailAddress: randomdata.Email(),
		ID:           fmt.Sprint(randomdata.Number(1, 10)),
	}
}
