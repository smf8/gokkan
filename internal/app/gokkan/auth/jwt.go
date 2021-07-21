package auth

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/smf8/gokkan/internal/app/gokkan/config"
	"github.com/smf8/gokkan/internal/app/gokkan/model"
)

const (
	// Issuer is jwt issuer for token generation.
	Issuer = "gokkan.io"

	// DefaultExpiration specifies JWT token expiration period.
	DefaultExpiration = 24 * time.Hour

	// ContextKey defines the key by which token is stored in echo context.
	// it's defined by default by echo. we can change it in JWTConfig.
	ContextKey = "user"
)

var (
	// ErrInvalidIssuer occurs with invalid ISS field.
	ErrInvalidIssuer = errors.New("invalid issuer")
	// ErrTokenExpired occurs when token is expired.
	ErrTokenExpired = errors.New("token expired")
	// ErrInvalidIssuedAt occurs when token issue time is in future.
	ErrInvalidIssuedAt = errors.New("token issue date is not valid")
	// ErrInvalidClaimsType occurs when claims types are other than GokkanClaims.
	ErrInvalidClaimsType = errors.New("invalid claims type")
	// ErrBlockedToken indicates that token with given JTI is blocked.
	ErrBlockedToken = errors.New("jwt token is blocked")
	// ErrTokenNotFound indicates there was no token found.
	ErrTokenNotFound = errors.New("token not found")
)

//nolint:gochecknoglobals
// probably not a good call to use a global variable here.
// but since our repo is thread-safe and by using this design
// we have JTI validation in GokkanClaims.Valid, I chose this
// design.
var (
	blacklistRepo model.TokenBlacklistRepo
	once          sync.Once
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

	if blacklistRepo.Check(g.JTI) {
		return ErrBlockedToken
	}

	return nil
}

// SetTokenBlacklistRepo sets global token blacklist repo in auth package.
// it should be called only once.
func SetTokenBlacklistRepo(repo model.TokenBlacklistRepo) {
	once.Do(func() {
		blacklistRepo = repo
	})
}

// MiddlewareConfig returns echo's default jwt middleware config
// with custom signing key and claims struct.
func MiddlewareConfig(cfg config.Server) middleware.JWTConfig {
	errHandler := func(err error) error {
		logrus.Errorf("jwt middleware failed: %s", err.Error())

		return err
	}

	config := middleware.JWTConfig{
		SigningKey:   []byte(cfg.Secret),
		Claims:       &GokkanClaims{},
		ErrorHandler: errHandler,
	}

	return config
}

// ExtractClaims extracts GokkanClaims from given context.
func ExtractClaims(c echo.Context) (*GokkanClaims, error) {
	token, ok := c.Get(ContextKey).(*jwt.Token)
	if !ok {
		return nil, fmt.Errorf("failed to extract cliams: %w", ErrTokenNotFound)
	}

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
		Exp:       time.Now().Add(DefaultExpiration).Unix(),
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
