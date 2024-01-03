package usersrv

import (
	"context"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"time"
	"zgit/pkg/apicode"
	"zgit/pkg/apisession"
	"zgit/pkg/i18n"
	"zgit/standalone/modules/model/usermd"
	"zgit/standalone/modules/service/cfgsrv"
	"zgit/util"
)

const (
	LoginSessionExpiry = 2 * time.Hour
)

func GetUserInfoByAccount(ctx context.Context, account string) (usermd.UserInfo, bool, error) {
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	user, b, err := usermd.GetByAccount(ctx, account)
	if err != nil {
		logger.Logger.Error(err)
		return usermd.UserInfo{}, false, util.InternalError()
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
	user, b, err := usermd.GetByAccount(ctx, reqDTO.Account)
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
	sessionStore.DeleteByAccount(user.Account)
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
	if session.UserInfo.Account != reqDTO.Operator.Account {
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
	// 只有系统管理员才能操作
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	_, b, err := usermd.GetByAccount(ctx, reqDTO.Account)
	if b {
		return bizerr.NewBizErr(apicode.UserAlreadyExistsCode.Int(), i18n.GetByKey(i18n.UserAlreadyExists))
	}
	_, err = usermd.InsertUser(ctx, usermd.InsertUserReqDTO{
		Name:      reqDTO.Name,
		Email:     reqDTO.Email,
		Password:  reqDTO.Password,
		IsAdmin:   reqDTO.IsAdmin,
		AvatarUrl: reqDTO.AvatarUrl,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Errorf("InsertUser err: %v", err)
		return util.InternalError()
	}
	return nil
}

func RegisterUser(ctx context.Context, reqDTO RegisterUserReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	// 检查是否已禁用该规则
	cfg, err := cfgsrv.GetSysCfgWithCache(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return err
	}
	if cfg.DisableSelfRegisterUser {
		return util.UnauthorizedError()
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	_, b, err := usermd.GetByAccount(ctx, reqDTO.Account)
	if b {
		return bizerr.NewBizErr(apicode.UserAlreadyExistsCode.Int(), i18n.GetByKey(i18n.UserAlreadyExists))
	}
	userCount, err := usermd.CountUser(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Errorf("RegisterUser err: %v", err)
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	// 如果用户表为空 就是管理员
	_, err = usermd.InsertUser(ctx, usermd.InsertUserReqDTO{
		Name:      reqDTO.Name,
		Email:     reqDTO.Email,
		Password:  reqDTO.Password,
		IsAdmin:   userCount == 0,
		AvatarUrl: reqDTO.AvatarUrl,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Errorf("RegisterUser err: %v", err)
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
	}
	return nil
}

func DeleteUser(ctx context.Context, reqDTO DeleteUserReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	// 只有系统管理员才能操作
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	user, b, err := usermd.GetByAccount(ctx, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	if !b {
		return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.UserNotFound))
	}
	// 数据库删除用户
	_, err = usermd.DeleteUser(ctx, user)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	// 删除用户登录状态
	apisession.GetStore().DeleteByAccount(reqDTO.Account)
	return nil
}

// ListUser 展示用户列表
func ListUser(ctx context.Context, reqDTO ListUserReqDTO) (ListUserRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return ListUserRespDTO{}, err
	}
	// 只有系统管理员才能操作
	if !reqDTO.Operator.IsAdmin {
		return ListUserRespDTO{}, util.UnauthorizedError()
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	userList, err := usermd.ListUser(ctx, usermd.ListUserReqDTO{
		Account: reqDTO.Account,
		Offset:  reqDTO.Offset,
		Limit:   reqDTO.Limit,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return ListUserRespDTO{}, util.InternalError()
	}
	ret := ListUserRespDTO{}
	ret.UserList, _ = listutil.Map(userList, func(t usermd.User) (UserDTO, error) {
		return UserDTO{
			Account:      t.Account,
			Name:         t.Name,
			Email:        t.Email,
			IsAdmin:      t.IsAdmin,
			IsProhibited: t.IsProhibited,
			AvatarUrl:    t.AvatarUrl,
			Created:      t.Created,
			Updated:      t.Updated,
		}, nil
	})
	if len(userList) > 0 {
		ret.Cursor = userList[len(userList)-1].Id
	}
	return ret, nil
}
