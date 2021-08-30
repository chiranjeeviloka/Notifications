package middlewares

import (
	"net/http"
	"notification-service/constants"
	"strings"

	"github.com/FreedomCentral/central/secret"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

type Claims struct {
	DocumentId       string `json:"documentId"`
	Username         string `json:"username"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	IsContentCreator bool   `json:"isContentCreator"`
	IsVerified       bool   `json:"isVerified"`
	jwt.StandardClaims
}

func VerifyTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		token, claims, tokenString, err := VerifyToken(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"error":   true,
				"message": err.Error(),
			})
			return
		}
		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"error":   true,
				"message": "unauthorized",
			})
			return
		}
		c.Set("TokenData", claims)
		c.Set("RequestUserEmail", claims.Email)
		c.Set("RequestUserName", claims.Username)
		c.Set("RequestUserPhone", claims.Phone)
		c.Set("RequestUserIsContentCreatorFlag", claims.IsContentCreator)
		c.Set("RequestUserIsVerifiedFlag", claims.IsVerified)
		c.Set("RequestUserID", claims.DocumentId)
		c.Set("AuthToken", tokenString)

		c.Next()
	}
}

func VerifyToken(r *http.Request) (*jwt.Token, *Claims, string, error) {

	sec, err := secret.Open(constants.SERVICE_NAME, secret.UseYAMLPlainText)
	if err != nil {
		logger.Fatalf("Failed to open secrets for %q: %v", constants.SERVICE_NAME, err)
	}

	jwtKey, err := sec.Get("TOKEN_SECRET")
	if err != nil {
		logger.Fatalf("Failed get jwtKey for %q: %v", constants.SERVICE_NAME, err)
	}
	tokenString := ExtractToken(r)
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return token, claims, tokenString, err
}

func ExtractToken(r *http.Request) string {
	var jwtTokenCleaned string
	jwtToken := r.Header.Get("Authorization")
	jwtTokenCleaned = strings.ReplaceAll(jwtToken, "Bearer", "")
	jwtTokenCleaned = strings.Trim(jwtTokenCleaned, " ")
	return jwtTokenCleaned
}
