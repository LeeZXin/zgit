package usermd

import (
	"time"
)

const (
	UserTableName = "user"
)

type User struct {
	Id           int64     `json:"id" xorm:"pk autoincr"`
	Account      string    `json:"account"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	IsAdmin      bool      `json:"isAdmin"`
	IsProhibited bool      `json:"isProhibited"`
	AvatarUrl    string    `json:"avatarUrl"`
	Created      time.Time `json:"created" xorm:"created"`
	Updated      time.Time `json:"updated" xorm:"updated"`
}

func (*User) TableName() string {
	return UserTableName
}

func (u *User) ToUserInfo() UserInfo {
	return UserInfo{
		Account:      u.Account,
		Name:         u.Name,
		Email:        u.Email,
		IsAdmin:      u.IsAdmin,
		IsProhibited: u.IsProhibited,
		AvatarUrl:    u.AvatarUrl,
	}
}
