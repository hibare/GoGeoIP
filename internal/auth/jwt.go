package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hibare/GoGeoIP/internal/config"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidClaimsType       = errors.New("invalid claims type")
	ErrMissingDataClaim        = errors.New("missing data claim")
	ErrFailedJWTSigning        = errors.New("failed to sign JWT token")
)

func GenerateJWT[T any](
	claims T,
	issuer string,
	subject string,
	audience string,
	ttl time.Duration,
	method jwt.SigningMethod,
) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(method, jwt.MapClaims{
		"iss": issuer,
		"sub": subject,
		"aud": audience,
		"iat": now.Unix(),
		"exp": now.Add(ttl).Unix(),
		"nbf": now.Unix(),
		"dat": claims, // custom payload
	})

	secret := []byte(config.Current.Core.SecretKey)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedJWTSigning, err)
	}

	return tokenString, nil
}

func VerifyJWT[T any](
	tokenString string,
	issuer string,
	audience string,
	method jwt.SigningMethod,
) (*T, error) {
	token, err := jwt.Parse(
		tokenString,
		func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != method.Alg() {
				return nil, ErrUnexpectedSigningMethod
			}
			return []byte(config.Current.Core.SecretKey), nil
		},
		jwt.WithIssuer(issuer),
		jwt.WithAudience(audience),
		jwt.WithValidMethods([]string{method.Alg()}),
	)

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaimsType
	}

	raw, ok := mapClaims["dat"]
	if !ok {
		return nil, ErrMissingDataClaim
	}

	// Convert map → struct
	bytes, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	var out T
	if err := json.Unmarshal(bytes, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
