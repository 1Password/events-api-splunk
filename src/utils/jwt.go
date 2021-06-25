package utils

import (
	"fmt"

	"gopkg.in/square/go-jose.v2/jwt"
)

type JWTClaims struct {
	Audience []string `json:"aud"`
}

func ParseJWTClaims(token string) (*JWTClaims, error) {
	t, err := jwt.ParseSigned(token)

	if err != nil {
		return nil, err
	}

	claims := &JWTClaims{}

	// We don't have the ECDSA Pub Key for verification at this point
	// as we need the URL to retrieve the metadata
	// This is fine as the server will properly verify the token
	err = t.UnsafeClaimsWithoutVerification(claims)

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (t *JWTClaims) GetEventsURL() string {
	if t.Audience[0] == "com.1password.streamingservice" {
		return ""
	}

	return fmt.Sprintf("https://%s", t.Audience[0])
}
