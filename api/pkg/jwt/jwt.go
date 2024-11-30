package jwt

import (
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/zeelrupapara/trading-api/config"
)

var (
	issuer = "zeeltrader.com"
)

func ParseToken(config config.AppConfig, token string) (jwt.Token, error) {
	key, err := jwk.FromRaw([]byte(config.Secret))
	if err != nil {
		return nil, err
	}

	claims, err := jwt.Parse([]byte(token), jwt.WithKey(jwa.HS256, key), jwt.WithIssuer(issuer))
	return claims, err
}

func CreateToken(config config.AppConfig, sub string, exp time.Time) (string, error) {
	stringToken := ""
	token, err := jwt.NewBuilder().Subject(sub).Expiration(exp).Issuer(issuer).Build()
	if err != nil {
		return stringToken, err
	}
	key, err := jwk.FromRaw([]byte(config.Secret))
	if err != nil {
		return stringToken, err
	}
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, key))
	if err != nil {
		return stringToken, err
	}
	stringToken = string(signed)
	return stringToken, nil
}
