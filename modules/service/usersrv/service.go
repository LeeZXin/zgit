package usersrv

import (
	"context"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"time"
	"zgit/modules/model/usermd"
	"zgit/pkg/apicode"
	"zgit/pkg/apisession"
	"zgit/pkg/i18n"
	"zgit/util"
)

const (
	LoginSessionExpiry = 2 * time.Hour
)

func GetUserInfoByUserId(ctx context.Context, userId string) (usermd.UserInfo, bool, error) {
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	user, b, err := usermd.GetByUserId(ctx, userId)
	if err != nil {
		logger.Logger.Errorf("GetUserInfoByUserId err: %v", err)
		return usermd.UserInfo{}, false, bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SystemInvalidArgs))
	}
	if !b {
		return usermd.UserInfo{}, false, nil
	}
	return user.ToUserInfo(), true, nil
}

func Login(ctx context.Context, reqDTO LoginReqDTO) (string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	user, b, err := usermd.GetByAccountAndCorpId(ctx, reqDTO.Account, reqDTO.CorpId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	if !b {
		return "", bizerr.NewBizErr(apicode.DataNotExistsCode.Int(), i18n.GetByKey(i18n.UserNotFound))
	}
	// 校验密码
	if user.Password != util.EncryptUserPassword(reqDTO.Password) {
		return "", bizerr.NewBizErr(apicode.WrongLoginPasswordCode.Int(), i18n.GetByKey(i18n.UserWrongPassword))
	}
	sessionStore := apisession.GetStore()
	// 删除原有的session
	sessionStore.DeleteByUserId(user.UserId)
	// 生成sessionId
	sessionId := apisession.GenSessionId()
	err = sessionStore.PutSession(apisession.Session{
		SessionId: sessionId,
		UserInfo:  user.ToUserInfo(),
		ExpireAt:  time.Now().Add(LoginSessionExpiry).UnixMilli(),
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	return sessionId, nil
}

func LoginOut(ctx context.Context, reqDTO LoginOutReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	sessionStore := apisession.GetStore()
	// 删除原有的session
	session, b, err := sessionStore.GetBySessionId(reqDTO.SessionId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	if !b {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SystemInvalidArgs))
	}
	if session.UserInfo.UserId != reqDTO.Operator.UserId {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SystemInvalidArgs))
	}
	err = sessionStore.DeleteBySessionId(reqDTO.SessionId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	return nil
}

func InsertUser(ctx context.Context, reqDTO InsertUserReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if !reqDTO.Operator.IsAdmin {
		return bizerr.NewBizErr(apicode.NotAdminCode.Int(), i18n.GetByKey(i18n.SystemNotAdmin))
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	_, err := usermd.InsertUser(ctx, usermd.InsertUserReqDTO{
		Name:      reqDTO.Name,
		Email:     reqDTO.Email,
		Password:  reqDTO.Password,
		IsAdmin:   reqDTO.IsAdmin,
		CorpId:    reqDTO.Operator.CorpId,
		AvatarUrl: reqDTO.AvatarUrl,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Errorf("InsertUser err: %v", err)
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	return nil
}

func DeleteUser(ctx context.Context, reqDTO DeleteUserReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if !reqDTO.Operator.IsAdmin {
		return bizerr.NewBizErr(apicode.NotAdminCode.Int(), i18n.GetByKey(i18n.SystemNotAdmin))
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	user, b, err := usermd.GetByUserId(ctx, reqDTO.UserId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	if !b {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.UserNotFound))
	}
	// 校验corpId是否和登录人相等
	if user.CorpId != reqDTO.Operator.CorpId {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.UserInvalidCorpId))
	}
	// 数据库删除用户
	_, err = usermd.DeleteUser(ctx, user)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	// 删除用户登录状态
	apisession.GetStore().DeleteByUserId(reqDTO.UserId)
	return nil
}
