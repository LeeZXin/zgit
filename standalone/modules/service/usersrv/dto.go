package usersrv

import (
	"github.com/LeeZXin/zsf-utils/bizerr"
	"regexp"
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
	if r.Account == "" {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n2.GetByKey(i18n2.UserInvalidId))
	}
	if r.Operator.Account == "" {
		return bizerr.NewBizErr(apicode.NotLoginCode.Int(), i18n2.GetByKey(i18n2.SystemNotLogin))
	}
	return nil
}
