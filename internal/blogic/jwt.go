package blogic

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"math/rand"
	"time"
)

// TokenManager provides logic for auth & Refresh tokens generation and parsing.
type TokenManager interface {
	NewJWT(ID string, ttl time.Duration) (string, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken() (string, error)
}
type TypeJWT int

type JWToken struct {
	jwt.StandardClaims
	ID string `json:"id"`
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) *Manager {
	if signingKey == "" {
		log.Fatalln("empty signing key")
	}

	return &Manager{signingKey: signingKey}
}

func (m *Manager) NewJWT(ID string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		JWToken{StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
		},
			ID: ID,
		})
	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["id"].(string), nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
