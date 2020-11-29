package auth

// JWTCookie is name of cookie holding jwt
func JWTCookie() string {
	return "x12HIDs"
}

// RefreshCookie is name of cookie holding jwt refresh token
func RefreshCookie() string {
	return "x12HIDv"
}

// AccountIDCookie is name of cookie holding account ID of signed in user
func AccountIDCookie() string {
	return "x12HIDt"
}

// SessionIDCookie is name of cookie holding session ID of signed in user
func SessionIDCookie() string {
	return "x12HIDu"
}

// Scheme returns authentication scheme
func Scheme() string {
	return "Bearer"
}

// Header returns authentication header
func Header() string {
	return "authorization"
}
