package userapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
)

type LoginReqVO struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type LoginRespVO struct {
	ginutil.BaseResp
	SessionId string `json:"sessionId"`
}

type InsertUserReqVO struct {
	Account   string `json:"account"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	IsAdmin   bool   `json:"isAdmin"`
	AvatarUrl string `json:"avatarUrl"`
}

type RegisterUserReqVO struct {
	Account   string `json:"account"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarUrl string `json:"avatarUrl"`
}

type DeleteUserReqVO struct {
	Account string `json:"account"`
}

type UserVO struct {
	Account      string `json:"account"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	IsAdmin      bool   `json:"isAdmin"`
	IsProhibited bool   `json:"isProhibited"`
	AvatarUrl    string `json:"avatarUrl"`
	Created      string `json:"created"`
	Updated      string `json:"updated"`
}

type ListUserReqVO struct {
	Account string `json:"account"`
	Offset  int64  `json:"offset"`
	Limit   int    `json:"limit"`
}

type ListUserRespVO struct {
	ginutil.BaseResp
	UserList []UserVO `json:"userList"`
	Cursor   int64    `json:"cursor"`
}

type UpdateUserReqVO struct {
	Account string
	Name    string
	Email   string
}

type UpdateAdminReqVO struct {
	Account string
	IsAdmin bool
}

type UpdatePasswordReqVO struct {
	Account  string
	Password string
}
