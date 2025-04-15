package service

import (
	"TaskBooker/internal/storage"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenInteractor interface {
	Create(userID int) (string, string, error)
	Refresh(token string) (string, string, error)
	DeleteByUserID(userID int) error
}

type TokenService struct {
	storage storage.TokenStorageInteractor
}

func NewTokenService(storage storage.TokenStorageInteractor) *TokenService {
	return &TokenService{storage: storage}
}

func (t *TokenService) Create(userID int) (string, string, error) {
	claims := jwt.MapClaims{
		"user_id":    userID,
		"exp":        time.Now().Add(time.Minute * 15).Unix(),
		"token_type": "refresh",
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", "", err
	}

	claims = jwt.MapClaims{
		"user_id":    userID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
		"token_type": "refresh",
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (t *TokenService) Refresh(token string) (string, string, error) {
	userID, err := t.storage.GetUserIDByRefreshToken(token)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	err = t.storage.DeleteByUserID(userID)
	if err != nil {
		return "", "", err
	}

	return t.Create(userID)
}

func (t *TokenService) DeleteByUserID(userID int) error {
	err := t.storage.DeleteByUserID(userID)
	if err != nil {
		return err
	}

	return nil
}
