package service

import (
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/program"
	"context"
	"errors"
	"github.com/kelseyhightower/envconfig"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	NewAccessToken(userID uint) (string, error)
	ParseAccessToken(tokenString string) (uint, error) // returns userID
}

type jwtService struct {
	secret    []byte
	issuer    string
	accessTTL time.Duration
}

func NewJWTService() JWTService {
	var wrapper struct {
		JWT jwtService `envconfig:"jwt"`
	}

	prefix := program.AppPrefix
	if prefix == "" {
		prefix = "API"
	}

	err := envconfig.Process(prefix, &wrapper)
	if err != nil {
		panic(logger.Error(context.Background(), err, "failed to process env vars"))
	}

	return &jwtService{
		secret:    wrapper.JWT.secret,
		issuer:    wrapper.JWT.issuer,
		accessTTL: wrapper.JWT.accessTTL,
	}
}

func (s *jwtService) NewAccessToken(userID uint) (string, error) {
	now := time.Now().UTC()

	claims := jwt.MapClaims{
		"sub": userID,
		"iss": s.issuer,
		"iat": now.Unix(),
		"exp": now.Add(s.accessTTL).Unix(),
		"typ": "access",
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.secret)
}

func (s *jwtService) ParseAccessToken(tokenString string) (uint, error) {
	parsed, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil || !parsed.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	if iss, _ := claims["iss"].(string); iss != s.issuer {
		return 0, errors.New("invalid issuer")
	}

	if typ, _ := claims["typ"].(string); typ != "access" {
		return 0, errors.New("wrong token type")
	}

	sub, ok := claims["sub"]
	if !ok {
		return 0, errors.New("missing sub")
	}

	switch v := sub.(type) {
	case float64:
		return uint(v), nil
	case uint:
		return v, nil
	default:
		return 0, errors.New("invalid sub type")
	}
}
