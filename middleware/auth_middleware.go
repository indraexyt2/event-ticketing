package middleware

import (
	"errors"
	"event-ticketing/config"
	"event-ticketing/entity"
	"event-ticketing/repository"
	"event-ticketing/utils"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID string      `json:"user_id"`
	Email  string      `json:"email"`
	Role   entity.Role `json:"role"`
	jwt.StandardClaims
}

func AuthMiddleware(userRepo repository.UserRepository, config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			utils.UnauthorizedResponse(c, "Unauthorized: Token not provided")
			c.Abort()
			return
		}

		claims, err := validateToken(tokenString, config.JWTSecret)
		if err != nil {
			utils.UnauthorizedResponse(c, "Unauthorized: "+err.Error())
			c.Abort()
			return
		}

		user, err := userRepo.FindByID(claims.UserID)
		if err != nil {
			utils.UnauthorizedResponse(c, "Unauthorized: User not found")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			utils.UnauthorizedResponse(c, "Unauthorized: User role not found")
			c.Abort()
			return
		}

		if role != entity.AdminRole {
			utils.ForbiddenResponse(c, "Forbidden: Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

func GenerateToken(userID string, email string, role entity.Role, secret string, expiresIn time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresIn).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func validateToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
