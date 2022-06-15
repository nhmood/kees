package helpers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"

	"kees/server/config"
)

var jwtConfiguration config.JWTConfig

func ConfigureJWT(config config.JWTConfig) {
	jwtConfiguration = config
}

func GenerateJWT(data map[string]interface{}) (string, int64, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	// Standard claims
	// TODO - create proper claims struct
	exp := time.Now().AddDate(0, 0, 1).Unix()
	claims["exp"] = exp
	claims["iss"] = jwtConfiguration.Issuer

	// TODO: replace with function pointer to specific claim updater function
	claims["kees"] = data

	Debug(claims)

	// Generate JWT with signing key from config
	jwt, err := token.SignedString([]byte(jwtConfiguration.SigningKey))
	if err != nil {
		// TODO: don't love empty string vs nil, maybe change to str ptr
		return "", 0, err
	}

	expiresIn := exp - time.Now().Unix()

	return jwt, expiresIn, nil
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtConfiguration.SigningKey), nil
	})

	// TODO - handle individual token failures (expired, etc)
	if err != nil {
		Debug(err)
		return nil, err
	}

	return claims, nil
}
