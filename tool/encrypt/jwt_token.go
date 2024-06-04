package encrypt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	Data string
	jwt.RegisteredClaims
}

func BuildClaims(data string, ttl int64) Claims {
	return Claims{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ttl*24) * time.Hour)), //Expiration time
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                        //Issuing time
			NotBefore: jwt.NewNumericDate(time.Now()),                                        //Begin Effective time
		}}
}

func CreateToken(data string, ttl int64, accessSecret string) (string, error) {
	claims := BuildClaims(data, ttl)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(accessSecret))
	return tokenString, err
}

func secret(accessSecret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	}
}

func ParseToken(tokenStr, accessSecret string) (claims *Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, secret(accessSecret))

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("that's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("token not active yet")
			} else {
				return nil, errors.New("couldn't handle this token")
			}
		}
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("couldn't handle this token")
}
