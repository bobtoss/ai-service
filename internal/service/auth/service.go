package auth

import (
	"ai-service/internal/repository/postgres"
	user2 "ai-service/internal/service/user"
	"context"
	"errors"
	"github.com/google/uuid"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	users      *postgres.UserRepository
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewService(u *postgres.UserRepository, jwtSecret []byte) *Service {
	return &Service{
		users:      u,
		jwtSecret:  jwtSecret,
		accessTTL:  15 * time.Minute,
		refreshTTL: 7 * 24 * time.Hour,
	}
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func (s *Service) HashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

func (s *Service) ComparePassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

func (s *Service) createAccessToken(userID string) (string, time.Time, error) {
	now := time.Now().UTC()
	exp := time.Now().Add(s.accessTTL)
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(s.accessTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtSecret)
	return signed, exp, err
}

func (s *Service) createRefreshTokenString(userID string) (string, error) {
	exp := time.Now().Add(s.refreshTTL)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"registered_claims": jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	signed, err := token.SignedString(s.jwtSecret)
	return signed, err
}

func (s *Service) issueTokens(userID string) (string, string, time.Time, error) {
	access, exp, err := s.createAccessToken(userID)
	if err != nil {
		return "", "", time.Time{}, err
	}

	refresh, err := s.createRefreshTokenString(userID)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return access, refresh, exp, nil
}

func (s *Service) Login(ctx context.Context, phone, password string) (accessToken, refreshToken string, refreshExp time.Time, err error) {
	u, err := s.users.GetByPhone(ctx, phone)
	if err != nil {
		return "", "", time.Time{}, ErrInvalidCredentials
	}
	if err := s.ComparePassword(u.Password, password); err != nil {
		return "", "", time.Time{}, ErrInvalidCredentials
	}

	return s.issueTokens(u.ID)
}

func (s *Service) Register(ctx context.Context, phone, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &user2.User{
		ID:       uuid.New().String(),
		Phone:    phone,
		Password: string(hash),
	}

	id, err := s.users.Create(ctx, user)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Service) Refresh(ctx context.Context, oldRefresh string) (accessToken, newRefresh string, newRefreshExp time.Time, err error) {
	userID, err := s.ParseAccessToken(oldRefresh)
	if err != nil {
		return "", "", time.Time{}, err
	}
	return s.issueTokens(userID)
}

func (s *Service) Logout(ctx context.Context, refresh string) error {
	return nil
}

// Helper to validate access token and return user id
func (s *Service) ParseAccessToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		// sometimes numeric -> try convert
		return "", errors.New("invalid subject claim")
	}
	return sub, nil
}
