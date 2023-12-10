package token

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const minSecretKeySize = 32

// JWT Maker is a JSON Web Token Maker
type JWTMaker struct {
	secretKey string
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	//privateKey := "secret"
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	signedToken, err := jwtToken.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil

}

// VerifyToken checks if the token is valid or not
// if valid - then it will the body of the payload data stored inside the token.
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwToken.Claims.(*Payload)

	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	maker := &JWTMaker{
		secretKey: secretKey,
	}
	return maker, nil
}