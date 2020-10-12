package modulorgo

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//ComparedPassword should be exported
func ComparedPassword(hashed string, str string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(str))
	if err != nil {
		return false
	}
	return true
}

//GenerateToken should be exported
func GenerateToken(payload map[string]string) (string, error) {
	mySigningKey := []byte(os.Getenv("JWT_KEY"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["public"] = payload
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//VerifyToken should be exported
func VerifyToken(next http.Handler) http.Handler {
	type key int
	const (
		keyPrincipalID key = iota
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			ResponseOnError(401, errors.New("JWT Invalid to Access this Endpoint"), w)
		} else {
			jwtToken := authHeader[1]
			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					ResponseOnError(401, errors.New("SigningMethodHMAC not OK"), w)
					return "", nil
				}
				return []byte(os.Getenv("JWT_KEY")), nil
			})
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				ctx := context.WithValue(r.Context(), keyPrincipalID, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				ResponseOnError(401, err, w)
			}
		}
	})
}

//IndetifyToken should be exported
func IndetifyToken(r *http.Request, w http.ResponseWriter) jwt.MapClaims {
	if len(strings.Split(r.Header.Get("Authorization"), "Bearer ")) != 2 {
		ResponseOnError(401, errors.New("JWT Invalid to Access this Endpoint"), w)
		return nil
	}
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	token, _, err := new(jwt.Parser).ParseUnverified(authHeader, jwt.MapClaims{})
	if err != nil {
		ResponseOnError(401, err, w)
		return nil
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims
	}

	ResponseOnError(401, err, w)
	return nil
}
