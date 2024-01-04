package usersrv

import (
	"regexp"
	"time"
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
)

var (
	validPasswordPattern = regexp.MustCompile("\\S{6,}")
)

type InsertUserReqDTO struct {
	Account   string
	Name      string
	Email     string
	Password  string
	IsAdmin   bool
	AvatarUrl string
	Operator  usermd.UserInfo
}

func (r *InsertUserReqDTO) IsValid() error {
	if !validateUserAccount(r.Account) {
		return util.InvalidArgsError()
	}
	if !validateUserEmail(r.Email) {
		return util.InvalidArgsError()
	}
	if !validatePassword(r.Password) {
		return util.InvalidArgsError()
	}
	if len(r.Name) > 32 || len(r.Name) == 0 {
		return util.InvalidArgsError()
	}
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type RegisterUserReqDTO struct {
	Account   string
	Name      string
	Email     string
	Password  string
	AvatarUrl string
}

func (r *RegisterUserReqDTO) IsValid() error {
	if !validateUserAccount(r.Account) {
		return util.InvalidArgsError()
	}
	if !validateUserEmail(r.Email) {
		return util.InvalidArgsError()
	}
	if !validPasswordPattern.MatchString(r.Password) {
		return util.InvalidArgsError()
	}
	if !validateUserName(r.Name) {
		return util.InvalidArgsError()
	}
	return nil
}

type LoginReqDTO struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

func (r *LoginReqDTO) IsValid() error {
	if !validateUserAccount(r.Account) {
		return util.InvalidArgsError()
	}
	if !validatePassword(r.Password) {
		return util.InvalidArgsError()
	}
	return nil
}

type LoginOutReqDTO struct {
	SessionId string
	Operator  usermd.UserInfo
}

func (r *LoginOutReqDTO) IsValid() error {
	if r.SessionId == "" {
		return util.InvalidArgsError()
	}
	if validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteUserReqDTO struct {
	Account  string
	Operator usermd.UserInfo
}

func (r *DeleteUserReqDTO) IsValid() error {
	if validateUserAccount(r.Account) {
		return util.InvalidArgsError()
	}
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type ListUserReqDTO struct {
	Account  string
	Offset   int64
	Limit    int
	Operator usermd.UserInfo
}

func (r *ListUserReqDTO) IsValid() error {
	if r.Offset < 0 {
		return util.InvalidArgsError()
	}
	if r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if len(r.Account) > 32 || len(r.Account) == 0 {
		return util.InvalidArgsError()
	}
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type UserDTO struct {
	Account      string
	Name         string
	Email        string
	IsAdmin      bool
	IsProhibited bool
	AvatarUrl    string
	Created      time.Time
	Updated      time.Time
}

type ListUserRespDTO struct {
	UserList []UserDTO
	Cursor   int64
}

type UpdateUserReqDTO struct {
	Account  string
	Name     string
	Email    string
	Operator usermd.UserInfo
}

func (r *UpdateUserReqDTO) IsValid() error {
	if !validateUserAccount(r.Account) {
		return util.InvalidArgsError()
	}
	if !validateUserName(r.Name) {
		return util.InvalidArgsError()
	}
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !validateUserEmail(r.Email) {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateAdminReqDTO struct {
	Account  string
	IsAdmin  bool
	Operator usermd.UserInfo
}

func (r *UpdateAdminReqDTO) IsValid() error {
	if !validateUserAccount(r.Account) {
		return util.InvalidArgsError()
	}
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdatePasswordReqDTO struct {
	Account  string
	Password string
	Operator usermd.UserInfo
}

func (r *UpdatePasswordReqDTO) IsValid() error {
	if !validateUserAccount(r.Account) {
		return util.InvalidArgsError()
	}
	if !validatePassword(r.Password) {
		return util.InvalidArgsError()
	}
	if !validateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

func validateUserAccount(account string) bool {
	return util.ValidUserAccountPattern.MatchString(account)
}

func validateUserName(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func validateOperator(operator usermd.UserInfo) bool {
	return operator.Account != ""
}

func validateUserEmail(email string) bool {
	return util.ValidUserEmailRegPattern.MatchString(email)
}

func validatePassword(password string) bool {
	return validPasswordPattern.MatchString(password)
}
