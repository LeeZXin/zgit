package usersrv

import (
	"github.com/LeeZXin/zsf-utils/bizerr"
	"regexp"
	"time"
	"zgit/pkg/apicode"
	i18n2 "zgit/pkg/i18n"
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
	if !util.ValidUserAccountPattern.MatchString(r.Account) {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n2.GetByKey(i18n2.UserInvalidAccount))
	}
	if !util.ValidEmailRegPattern.MatchString(r.Email) {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n2.GetByKey(i18n2.UserInvalidEmail))
	}
	if !validPasswordPattern.MatchString(r.Password) {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n2.GetByKey(i18n2.UserInvalidPassword))
	}
	if r.Operator.Account == "" {
		return bizerr.NewBizErr(apicode.NotLoginCode.Int(), i18n2.GetByKey(i18n2.SystemNotLogin))
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
	if !util.ValidUserAccountPattern.MatchString(r.Account) {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n2.GetByKey(i18n2.UserInvalidAccount))
	}
	if !util.ValidEmailRegPattern.MatchString(r.Email) {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n2.GetByKey(i18n2.UserInvalidEmail))
	}
	if !validPasswordPattern.MatchString(r.Password) {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n2.GetByKey(i18n2.UserInvalidPassword))
	}
	return nil
}

type LoginReqDTO struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

func (r *LoginReqDTO) IsValid() error {
	if !util.ValidUserAccountPattern.MatchString(r.Account) {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n2.GetByKey(i18n2.UserInvalidAccount))
	}
	if !validPasswordPattern.MatchString(r.Password) {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n2.GetByKey(i18n2.UserInvalidPassword))
	}
	return nil
}

type LoginOutReqDTO struct {
	SessionId string
	Operator  usermd.UserInfo
}

func (r *LoginOutReqDTO) IsValid() error {
	if r.SessionId == "" {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n2.GetByKey(i18n2.UserInvalidSessionId))
	}
	if r.Operator.Account == "" {
		return bizerr.NewBizErr(apicode.NotLoginCode.Int(), i18n2.GetByKey(i18n2.SystemNotLogin))
	}
	return nil
}

type DeleteUserReqDTO struct {
	Account  string
	Operator usermd.UserInfo
}

func (r *DeleteUserReqDTO) IsValid() error {
	if len(r.Account) > 32 || len(r.Account) == 0 {
		return util.InternalError()
	}
	if r.Operator.Account == "" {
		return util.InternalError()
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
	if r.Limit < 0 {
		return util.InvalidArgsError()
	}
	if len(r.Account) > 32 || len(r.Account) == 0 {
		return util.InternalError()
	}
	if r.Operator.Account == "" {
		return util.InternalError()
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
