package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/services/pkg/utils/errs"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

// Groups returns the accociated account groups
func Groups() []string {
	return []string{
		User(),
		AdminGroup(),
		SuperAdminGroup(),
	}
}

// Admins returns the administrators group
func Admins() []string {
	return []string{AdminGroup(), SuperAdminGroup()}
}

// User are ordinary app users
func User() string {
	return "USER"
}

// AdminGroup is group for admin users
func AdminGroup() string {
	return "ADMIN"
}

// SuperAdminGroup is group for super admin users
func SuperAdminGroup() string {
	return "SUPER_ADMIN"
}

// API is used for authentication and authorization
type API interface {
	AuthenticateRequest(context.Context) error
	AuthenticateRequestV2(context.Context) (*Payload, error)
	AuthorizeActor(ctx context.Context, actorID string) (*Payload, error)
	AuthorizeActors(ctx context.Context, actorID ...string) (*Payload, error)
	AuthorizeGroups(ctx context.Context, allowedGroups ...string) (*Payload, error)
	AuthorizeStrict(ctx context.Context, actorID string, allowedGroups ...string) (*Payload, error)
	AuthorizeActorOrGroups(ctx context.Context, actorID string, allowedGroups ...string) (*Payload, error)
	GenToken(context.Context, *Payload, time.Time) (string, error)
	GetJwtPayload(context.Context) (*Payload, error)
	AuthFunc(context.Context) (context.Context, error)
}

type authAPI struct {
	signingMethod jwt.SigningMethod
	signingKey    []byte
	issuer        string
	audience      string
}

// NewAPI creates a jwt authentication and authorization API using HS256 algorithm
func NewAPI(signingKey []byte, issuer, audience string) (API, error) {
	// Validation
	switch {
	case signingKey == nil:
		return nil, errs.NilObject("jwt signing key")
	case issuer == "":
		return nil, errs.MissingField("jwt issuer")
	case audience == "":
		return nil, errs.MissingField("jwt audience")
	}

	api := &authAPI{
		signingMethod: jwt.SigningMethodHS256,
		signingKey:    signingKey,
		issuer:        issuer,
		audience:      audience,
	}

	return api, nil
}

func (api *authAPI) GetJwtPayload(ctx context.Context) (*Payload, error) {
	tokenInfo, ok := ctx.Value(ctxKey).(*Claims)
	if !ok {
		return nil, errs.WrapMessage(codes.Unauthenticated, "no claims found in token")
	}

	return tokenInfo.Payload, nil
}

func (api *authAPI) AuthenticateRequestV2(ctx context.Context) (*Payload, error) {
	claims, err := api.ParseFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	return claims.Payload, nil
}

func (api *authAPI) AuthenticateRequest(ctx context.Context) error {
	_, err := api.ParseFromCtx(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (api *authAPI) AuthorizeActor(ctx context.Context, actorID string) (*Payload, error) {
	// Get user claims from token
	claims, ok := ctx.Value(ctxKey).(*Claims)
	if !ok {
		return nil, errs.WrapMessage(codes.Unauthenticated, "no claims found")
	}

	if claims.ID != actorID {
		return nil, errs.TokenCredentialNotMatching("id")
	}

	return claims.Payload, nil
}

func (api *authAPI) AuthorizeActors(ctx context.Context, actorIDs ...string) (*Payload, error) {
	// Get user claims from token
	claims, ok := ctx.Value(ctxKey).(*Claims)
	if !ok {
		return nil, errs.WrapMessage(codes.Unauthenticated, "no claims found")
	}

	for _, actorID := range actorIDs {
		if claims.ID == actorID {
			return claims.Payload, nil
		}
	}

	return nil, errs.TokenCredentialNotMatching("id")
}

func (api *authAPI) AuthorizeGroups(ctx context.Context, allowedGroups ...string) (*Payload, error) {
	// Get user claims from token
	claims, ok := ctx.Value(ctxKey).(*Claims)
	if !ok {
		return nil, errs.WrapMessage(codes.Unauthenticated, "no claims found")
	}

	err := matchGroup(claims.Payload.Group, allowedGroups)
	if err != nil {
		return nil, err
	}

	return claims.Payload, nil
}

func (api *authAPI) AuthorizeStrict(ctx context.Context, actorID string, allowedGroups ...string) (*Payload, error) {
	// Get user claims from token
	claims, ok := ctx.Value(ctxKey).(*Claims)
	if !ok {
		return nil, errs.WrapMessage(codes.Unauthenticated, "no claims found")
	}

	err := matchGroup(claims.Payload.Group, allowedGroups)
	if err != nil {
		return nil, err
	}

	if claims.ID != actorID {
		return nil, err
	}

	return claims.Payload, nil
}

func (api *authAPI) AuthorizeActorOrGroups(
	ctx context.Context, actorID string, allowedGroups ...string,
) (*Payload, error) {
	// Get user claims from token
	claims, ok := ctx.Value(ctxKey).(*Claims)
	if !ok {
		return nil, errs.WrapMessage(codes.Unauthenticated, "no claims found")
	}

	var err error
	if claims.ID != actorID {
		err = errs.TokenCredentialNotMatching("id")
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

func (api *authAPI) GenToken(
	ctx context.Context, payload *Payload, expirationTime time.Time,
) (string, error) {
	return api.genToken(ctx, payload, expirationTime.Unix())
}

func userClaimFromToken(tokenInfo *Claims) *Payload {
	return tokenInfo.Payload
}

type tokenInfo string

// ctxKey holds the context key containing the token information
const ctxKey = tokenInfo("tokenInfo")

func (api *authAPI) AuthFunc(ctx context.Context) (context.Context, error) {
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

// AddMD adds metadata to token
func (api *authAPI) AddMD(
	ctx context.Context, actorID, group string,
) context.Context {
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
			return api.signingKey, nil
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

// AddTokenMD adds token as authorization metadata to context and returns the updated context object
func AddTokenMD(ctx context.Context, token string) context.Context {
	return addTokenMD(ctx, token)
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

func (api *authAPI) genToken(
	ctx context.Context, payload *Payload, expires int64,
) (tokenStr string, err error) {
	// Handling any panic is good trust me!
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("%v", err2)
		}
	}()

	token := jwt.NewWithClaims(api.signingMethod, Claims{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires,
			Issuer:    api.issuer,
			Audience:  api.audience,
		},
	})

	// Generate the token
	return token.SignedString(api.signingKey)
}
