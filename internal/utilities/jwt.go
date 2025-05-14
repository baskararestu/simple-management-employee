package utilities

import (
	"errors"
	"simple-management-employee/pkg/xlogger"
	"time"

	"github.com/golang-jwt/jwt"
)
type Claims struct {
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
	jwt.StandardClaims
}


func GenerateToken(userID string, roleID string,jwtSecret string) (*string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role_id": roleID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	xlogger.Logger.Info().Msgf("JWT Claims: %v", claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}

	return &signedToken, nil
}

func ExtractClaimsFromToken(tokenStr string, jwtSecret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("could not parse claims")
	}

	return claims, nil
}