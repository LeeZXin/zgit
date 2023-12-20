package lfs

import "github.com/golang-jwt/jwt/v5"

// Claims is a JWT Token Claims
type Claims struct {
	RepoId string
	Op     string
	UserId string
	jwt.RegisteredClaims
}

type TokenRespVO struct {
	Header map[string]string `json:"header"`
	Href   string            `json:"href"`
}
