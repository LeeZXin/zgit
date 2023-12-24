package lfs

import "github.com/golang-jwt/jwt/v5"

// Claims is a JWT Token Claims
type Claims struct {
	RepoPath string
	Op       string
	Account  string
	jwt.RegisteredClaims
}

type TokenRespVO struct {
	Header map[string]string `json:"header"`
	Href   string            `json:"href"`
}
