package usersrv

import (
	"github.com/LeeZXin/zsf-utils/bizerr"
	"regexp"
	"zgit/modules/model/usermd"
	"zgit/pkg/apicode"
	"zgit/pkg/i18n"
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
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.UserInvalidAccount))
	}
	if !util.ValidEmailRegPattern.MatchString(r.Email) {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.UserInvalidEmail))
	}
	if !validPasswordPattern.MatchString(r.Password) {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.UserInvalidPassword))
	}
	if r.Operator.UserId == "" {
		return bizerr.NewBizErr(apicode.NotLoginCode.Int(), i18n.GetByKey(i18n.SystemNotLogin))
	}
	return nil
}

type LoginReqDTO struct {
	CorpId   string `json:"corpId"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

func (r *LoginReqDTO) IsValid() error {
	if !util.ValidUserAccountPattern.MatchString(r.Account) {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.UserInvalidAccount))
	}
	if !util.ValidCorpIdPattern.MatchString(r.CorpId) {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.UserInvalidCorpId))
	}
	if !validPasswordPattern.MatchString(r.Password) {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.UserInvalidPassword))
	}
	return nil
}

type LoginOutReqDTO struct {
	SessionId string
	Operator  usermd.UserInfo
}

func (r *LoginOutReqDTO) IsValid() error {
	if r.SessionId == "" {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.UserInvalidSessionId))
	}
	if r.Operator.UserId == "" {
		return bizerr.NewBizErr(apicode.NotLoginCode.Int(), i18n.GetByKey(i18n.SystemNotLogin))
	}
	return nil
}

type DeleteUserReqDTO struct {
	UserId   string
	Operator usermd.UserInfo
}

func (r *DeleteUserReqDTO) IsValid() error {
	if r.UserId == "" {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.UserInvalidId))
	}
	if r.Operator.UserId == "" {
		return bizerr.NewBizErr(apicode.NotLoginCode.Int(), i18n.GetByKey(i18n.SystemNotLogin))
	}
	return nil
}
