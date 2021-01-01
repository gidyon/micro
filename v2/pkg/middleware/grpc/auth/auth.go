package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/micro/v2/utils/errs"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

// DefaultAdminGroups returns the default administrators group
func DefaultAdminGroups() []string {
	return []string{DefaultAdminGroup(), DefaultSuperAdminGroup()}
}

// DefaultUserGroup is the default user group
func DefaultUserGroup() string {
	return "USER"
}

// DefaultAdminGroup is the default admin group
func DefaultAdminGroup() string {
	return "ADMIN"
}

// DefaultSuperAdminGroup is the default super admin group
func DefaultSuperAdminGroup() string {
	return "SUPER_ADMIN"
}

// API is the interface used for authentication and authorization
type API interface {
	AuthenticateRequest(ctx context.Context) error
	AuthenticateRequestV2(ctx context.Context) (*Payload, error)
	AuthorizeGroup(ctx context.Context, allowedGroups ...string) (*Payload, error)
	AuthorizeActor(ctx context.Context, actorID string) (*Payload, error)
	AuthorizeActors(ctx context.Context, actorID ...string) (*Payload, error)
	AuthorizeActorAndGroup(ctx context.Context, actorID string, allowedGroups ...string) (*Payload, error)
	AuthorizeActorOrGroup(ctx context.Context, actorID string, allowedGroups ...string) (*Payload, error)
	AuthorizeAdmin(ctx context.Context) (*Payload, error)
	AuthorizeAdminStrict(ctx context.Context, adminID string) (*Payload, error)
	AdminGroups() []string
	IsAdmin(group string) bool
	GenToken(ctx context.Context, payload *Payload, expires time.Time) (string, error)
	GetJwtPayload(ctx context.Context) (*Payload, error)
	AuthorizeFunc(ctx context.Context) (context.Context, error)
}

type authAPI struct {
	*Options
}

// Options contains parameters for instantiating new API
type Options struct {
	SigningMethod jwt.SigningMethod
	SigningKey    []byte
	Issuer        string
	Audience      string
	AdminsGroup   []string
}

// NewAPI creates a jwt authentication and authorization API using HS256 algorithm
func NewAPI(opt *Options) (API, error) {

	// Validation
	switch {
	case opt.SigningKey == nil:
		return nil, errs.NilObject("jwt signing key")
	case opt.Issuer == "":
		return nil, errs.MissingField("jwt issuer")
	case opt.Audience == "":
		return nil, errs.MissingField("jwt audience")
	}

	if len(opt.AdminsGroup) == 0 {
		opt.AdminsGroup = DefaultAdminGroups()
	}

	opt.SigningMethod = jwt.SigningMethodHS256

	optVal := *opt

	api := &authAPI{Options: &optVal}

	return api, nil
}

func (api *authAPI) AuthenticateRequest(ctx context.Context) error {
	_, err := api.ParseFromCtx(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (api *authAPI) AuthenticateRequestV2(ctx context.Context) (*Payload, error) {
	claims, err := api.ParseFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	return claims.Payload, nil
}

func (api *authAPI) AuthorizeGroup(ctx context.Context, allowedGroups ...string) (*Payload, error) {
	claims, ok := ctx.Value(ctxKey).(*Claims)
	if !ok {
		return nil, errs.WrapMessage(codes.Unauthenticated, "no claims found in token")
	}

	err := matchGroup(claims.Payload.Group, allowedGroups)
	if err != nil {
		return nil, err
	}

	return claims.Payload, nil
}

func (api *authAPI) AuthorizeActor(ctx context.Context, actorID string) (*Payload, error) {
	claims, ok := ctx.Value(ctxKey).(*Claims)
	if !ok {
		return nil, errs.WrapMessage(codes.Unauthenticated, "no claims found in token")
	}

	if claims.ID != actorID {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied for actor with id %s", claims.ID)
	}

	return claims.Payload, nil
}

func (api *authAPI) AuthorizeActors(ctx context.Context, actorIDs ...string) (*Payload, error) {
	claims, ok := ctx.Value(ctxKey).(*Claims)
	if !ok {
		return nil, errs.WrapMessage(codes.Unauthenticated, "no claims found in token")
	}

	for _, actorID := range actorIDs {
		if claims.ID == actorID {
			return claims.Payload, nil
		}
	}

	return nil, status.Errorf(codes.PermissionDenied, "permission denied for actors ids [%s]", strings.Join(actorIDs, ", "))
}

func (api *authAPI) AuthorizeActorAndGroup(ctx context.Context, actorID string, allowedGroups ...string) (*Payload, error) {
	claims, ok := ctx.Value(ctxKey).(*Claims)
	if !ok {
		return nil, errs.WrapMessage(codes.Unauthenticated, "no claims found in token")
	}

	err := matchGroup(claims.Payload.Group, allowedGroups)
	if err != nil {
		return nil, err
	}

	if claims.ID != actorID {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied for actor with id %s", claims.ID)
	}

	return claims.Payload, nil
}

func (api *authAPI) AuthorizeActorOrGroup(ctx context.Context, actorID string, allowedGroups ...string) (*Payload, error) {
	claims, ok := ctx.Value(ctxKey).(*Claims)
	if !ok {
		return nil, errs.WrapMessage(codes.Unauthenticated, "no claims found in token")
	}

	var err error
	if claims.ID != actorID {
		err = status.Errorf(codes.PermissionDenied, "permission denied for actor with id %s", claims.ID)
	}

	err2 := matchGroup(claims.Payload.Group, allowedGroups)

	switch {
	case err2 == nil && err == nil:
	case err2 != nil && err == nil:
	case err2 == nil && err != nil:
	case err != nil:
		return nil, err
	default:
		return nil, err2
	}

	return claims.Payload, nil
}

func (api *authAPI) AuthorizeAdmin(ctx context.Context) (*Payload, error) {
	claims, ok := ctx.Value(ctxKey).(*Claims)
	if !ok {
		return nil, errs.WrapMessage(codes.Unauthenticated, "no claims found in token")
	}

	err := matchGroup(claims.Payload.Group, api.AdminsGroup)
	if err != nil {
		return nil, err
	}

	return claims.Payload, nil
}

func (api *authAPI) AuthorizeAdminStrict(ctx context.Context, adminID string) (*Payload, error) {
	claims, ok := ctx.Value(ctxKey).(*Claims)
	if !ok {
		return nil, errs.WrapMessage(codes.Unauthenticated, "no claims found in context")
	}

	err := matchGroup(claims.Payload.Group, api.AdminsGroup)
	if err != nil {
		return nil, err
	}

	if claims.ID != adminID {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied for admin with id %s", claims.ID)
	}

	return claims.Payload, nil
}

func (api *authAPI) AdminGroups() []string {
	v := make([]string, 0, len(api.AdminsGroup))
	for _, g := range api.AdminsGroup {
		v = append(v, g)
	}
	return v
}

func (api *authAPI) IsAdmin(group string) bool {
	err := matchGroup(group, api.AdminsGroup)
	if err != nil {
		return false
	}
	return true
}

func (api *authAPI) GenToken(ctx context.Context, payload *Payload, expirationTime time.Time) (string, error) {
	return api.genToken(ctx, payload, expirationTime.Unix())
}

func (api *authAPI) GetJwtPayload(ctx context.Context) (*Payload, error) {
	tokenInfo, ok := ctx.Value(ctxKey).(*Claims)
	if !ok {
		return nil, errs.WrapMessage(codes.Unauthenticated, "no claims found in token")
	}

	return tokenInfo.Payload, nil
}

func (api *authAPI) AuthorizeFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	tokenInfo, err := api.ParseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	grpc_ctxtags.Extract(ctx).Set("auth.sub", userClaimFromToken(tokenInfo))

	return context.WithValue(ctx, ctxKey, tokenInfo), nil
}

// ParseFromCtx jwt token from context
func (api *authAPI) ParseFromCtx(ctx context.Context) (*Claims, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, status.Errorf(
			codes.PermissionDenied, "failed to get Bearer token from authorization header: %v", err,
		)
	}

	return api.ParseToken(token)
}

// ParseToken parses a jwt token and return claims or error if token is invalid
func (api *authAPI) ParseToken(tokenString string) (claims *Claims, err error) {
	// Handling any panic is good trust me!
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("%v", err2)
		}
	}()

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return api.SigningKey, nil
		},
	)
	if err != nil {
		return nil, status.Errorf(
			codes.Unauthenticated, "failed to parse token with claims: %v", err,
		)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, status.Error(codes.Unauthenticated, "JWT is not valid")
	}
	return claims, nil
}

func userClaimFromToken(tokenInfo *Claims) *Payload {
	return tokenInfo.Payload
}

type tokenInfo string

// ctxKey holds the context key containing the token information
const ctxKey = tokenInfo("tokenInfo")

// AddMD adds metadata to token
func (api *authAPI) AddMD(ctx context.Context, actorID, group string) context.Context {
	payload := &Payload{
		ID:           actorID,
		Names:        randomdata.SillyName(),
		EmailAddress: randomdata.Email(),
		Group:        group,
	}
	token, err := api.genToken(ctx, payload, 0)
	if err != nil {
		panic(err)
	}

	return addTokenMD(ctx, token)
}

// AddTokenMD adds token as authorization metadata to context and returns the updated context object
func AddTokenMD(ctx context.Context, token string) context.Context {
	return addTokenMD(ctx, token)
}

func (api *authAPI) genToken(ctx context.Context, payload *Payload, expires int64) (tokenStr string, err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("%v", err2)
		}
	}()

	token := jwt.NewWithClaims(api.SigningMethod, Claims{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires,
			Issuer:    api.Issuer,
			Audience:  api.Audience,
		},
	})

	return token.SignedString(api.SigningKey)
}

func addTokenMD(ctx context.Context, token string) context.Context {
	return metadata.NewIncomingContext(
		ctx, metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", token)),
	)
}

func matchGroup(claimGroup string, allowedGroups []string) error {
	for _, group := range allowedGroups {
		if claimGroup == group {
			return nil
		}
	}
	return status.Errorf(codes.PermissionDenied, "permission denied for group %s", claimGroup)
}
