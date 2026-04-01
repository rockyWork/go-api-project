package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   int    `json:"role"`
	Type   string `json:"type"` // access or refresh
	jwt.RegisteredClaims
}

type Manager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessExpire  time.Duration
	refreshExpire time.Duration
}

func NewManager(accessSecret, refreshSecret string, accessExpire, refreshExpire time.Duration) *Manager {
	return &Manager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessExpire:  accessExpire,
		refreshExpire: refreshExpire,
	}
}

func (m *Manager) GenerateTokenPair(userID uint, role int) (accessToken, refreshToken string, err error) {
	accessToken, err = m.generateToken(userID, role, "access", m.accessSecret, m.accessExpire)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = m.generateToken(userID, role, "refresh", m.refreshSecret, m.refreshExpire)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (m *Manager) generateToken(userID uint, role int, tokenType string, secret []byte, expire time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expire)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (m *Manager) ParseAccessToken(tokenString string) (*Claims, error) {
	return m.parseToken(tokenString, m.accessSecret)
}

func (m *Manager) ParseRefreshToken(tokenString string) (*Claims, error) {
	return m.parseToken(tokenString, m.refreshSecret)
}

func (m *Manager) parseToken(tokenString string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (m *Manager) RefreshTokens(refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	claims, err := m.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	if claims.Type != "refresh" {
		return "", "", ErrInvalidToken
	}

	return m.GenerateTokenPair(claims.UserID, claims.Role)
}
