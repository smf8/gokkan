package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4/middleware"
	"github.com/smf8/gokkan/internal/app/gokkan/config"
)

// Issuer is jwt issuer for token generation.
const Issuer = "gokkan.io"

const defaultExpiration = 24 * time.Hour

var (
	// ErrInvalidIssuer occurs with invalid ISS field.
	ErrInvalidIssuer = errors.New("invalid issuer")
	// ErrTokenExpired occurs when token is expired.
	ErrTokenExpired = errors.New("token expired")
	// ErrInvalidIssuedAt occurs when token issue time is in future.
	ErrInvalidIssuedAt = errors.New("token issue date is not valid")
	// ErrInvalidClaimsType occurs when claims types are other than GokkanClaims.
	ErrInvalidClaimsType = errors.New("invalid claims type")
)

// GokkanClaims store a combination of registered and public
// jwt claims. check https://www.iana.org/assignments/jwt/jwt.xhtml
// for description about claims.
type GokkanClaims struct {
	// Issuer
	Iss string `json:"iss,omitempty"`
	// subject claim (username)
	Sub string `json:"sub,omitempty"`
	// Expiration
	Exp int64 `json:"exp,omitempty"`
	// Issued at
	Iat int64 `json:"iat,omitempty"`
	// JWT ID
	// we're not implementing JTI
	// we should do so in a production environment
	// since this is the field that let's us control
	// tokens from server side (i.e. block a specific token)
	JTI string `json:"jti,omitempty"`
	// is user privileged (private claim)
	Privieged bool `json:"privieged,omitempty"`
}

// Valid checks if a GokkanClaims is valid jwt token or not.
func (g GokkanClaims) Valid() error {
	if g.Iss != Issuer {
		return ErrInvalidIssuer
	}

	now := time.Now().Unix()

	if g.Exp < now {
		return ErrTokenExpired
	}

	if now < g.Iat {
		return ErrInvalidIssuedAt
	}

	return nil
}

// MiddlewareConfig returns echo's default jwt middleware config
// with custom signing key and claims struct.
func MiddlewareConfig(cfg config.Server) middleware.JWTConfig {
	config := middleware.JWTConfig{
		SigningKey: cfg.Secret,
		Claims:     &GokkanClaims{},
	}

	return config
}

// ExtractClaims extracts GokkanClaims from given jwt token.
func ExtractClaims(token jwt.Token) (*GokkanClaims, error) {
	claims, ok := token.Claims.(*GokkanClaims)

	if !ok {
		return nil, fmt.Errorf("type conversion failed: %t, %w", token.Claims, ErrInvalidClaimsType)
	}

	return claims, nil
}

// Generate is used to create a jwt string for given user.
func Generate(secret, username string, isAdmin bool) (string, error) {
	claims := &GokkanClaims{
		Iss:       Issuer,
		Sub:       username,
		Exp:       time.Now().Add(defaultExpiration).Unix(),
		Iat:       time.Now().Unix(),
		JTI:       fmt.Sprintf("%d", time.Now().UnixNano()),
		Privieged: isAdmin,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}
