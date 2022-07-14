package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/vitalygudza/duty-app/internal/model"
	"github.com/vitalygudza/duty-app/internal/repository"
)

const (
	salt       = "jhUYUYadsf6iN5dfb%j3ib"
	signingKey = "cvJU7j3hdsfhbsjf%u3bfa"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (a *AuthService) CreateUser(user model.User) (int, error) {
	user.Password = a.generateHashPassword(user.Password)
	return a.repo.CreateUser(user)
}

func (a *AuthService) CheckCredentials(signInInput model.SignInInput) (model.User, error) {
	signInInput.Password = a.generateHashPassword(signInInput.Password)
	return a.repo.GetUser(signInInput)
}

func (a *AuthService) GenerateToken(expireTime time.Time, userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
	})

	return token.SignedString([]byte(signingKey))
}

func (a *AuthService) ParseToken(cookieValue string) (int, error) {
	claims := &tokenClaims{}
	token, err := jwt.ParseWithClaims(cookieValue, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("cookie token %q is not valid", cookieValue)
	}

	return claims.UserId, nil
}

func (a *AuthService) generateHashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
