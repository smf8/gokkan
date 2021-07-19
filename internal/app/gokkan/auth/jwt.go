package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4/middleware"
	"github.com/smf8/gokkan/internal/app/gokkan/config"
)

// Issuer is jwt issuer for token generation
const Issuer = "gokkan.io"

var defaultExpiration = 24 * time.Hour

// GokkanClaims store a combination of registered and public
// jwt claims. check https://www.iana.org/assignments/jwt/jwt.xhtml
// for description about claims
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

// Valid checks if a GokkanClaims is valid jwt token or not
func (g GokkanClaims) Valid() error {
	if g.Iss != Issuer {
		return fmt.Errorf("invalid issuer")
	}

	now := time.Now().Unix()

	if g.Exp < now {
		return fmt.Errorf("token expired")
	}

	if now < g.Iat {
		return fmt.Errorf("token is not ready to use")
	}

	return nil
}

// MiddlewareConfig returns echo's default jwt middleware config
// with custom signing key and claims struct
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
		return nil, fmt.Errorf("invalid token claims type: %t", token.Claims)
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

	return token.SignedString([]byte(secret))
}
