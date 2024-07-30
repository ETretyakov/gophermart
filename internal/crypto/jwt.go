package crypto

import (
	"gophermart/internal/exceptions"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

var JWT *JWTSigner

type JWTSigner struct {
	secretKey string
	expire    time.Duration
}

type UserClaim struct {
	jwt.StandardClaims
	Sub    string
	UserID string
}

func InitJWTSigner(secretKey string, expire time.Duration) {
	JWT = &JWTSigner{secretKey: secretKey, expire: expire}
}

func (j JWTSigner) GetToken(userID string) (string, error) {
	payload := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(j.expire).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", errors.Wrapf(err, "failed to sign token")
	}

	return signedToken, nil
}

func (j JWTSigner) ExtractToken(rawToken string) (*jwt.RegisteredClaims, *jwt.Token, error) {
	splitToken := strings.Split(rawToken, "Bearer ")
	if len(splitToken) != 2 {
		return nil, nil, errors.New("failed to parse bearer token")
	}
	rawToken = splitToken[1]

	token, err := jwt.ParseWithClaims(
		rawToken,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.secretKey), nil
		},
	)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to parse token")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !(ok && token.Valid) {
		return nil, nil, errors.Wrapf(err, "failed to parse claims")
	}

	return claims, token, nil
}

func (j JWTSigner) AuthToken(h http.Header) (*jwt.RegisteredClaims, error) {
	rawToken := h.Get("Authorization")

	userClaim, token, err := j.ExtractToken(rawToken)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to extract token")
	}

	if !token.Valid {
		return nil, exceptions.ErrNotAuthorised
	}

	return userClaim, nil
}
