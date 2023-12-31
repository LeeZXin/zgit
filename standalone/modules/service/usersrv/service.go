package usersrv

import (
	"context"
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

var (
	userCache = util.NewGoCache()
)

const (
	LoginSessionExpiry = 2 * time.Hour
)

func GetUserInfoByAccount(ctx context.Context, account string) (usermd.UserInfo, bool, error) {
	uc, b := userCache.Get(account)
	if b {
		u := uc.(usermd.UserInfo)
		// 来自空缓存
		if u.Account == "" {
			return usermd.UserInfo{}, false, nil
		}
		return uc.(usermd.UserInfo), true, nil
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	user, b, err := usermd.GetByAccount(ctx, account)
	if err != nil {
		logger.Logger.Error(err)
		return usermd.UserInfo{}, false, util.InternalError()
	}
	if !b {
		// 设置空缓存
		u := usermd.UserInfo{}
		userCache.Set(account, u, time.Second)
		return usermd.UserInfo{}, false, nil
	}
	// 三分钟缓存
	ret := user.ToUserInfo()
	userCache.Set(account, ret, 3*time.Minute)
	return ret, true, nil
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
		return "", util.InternalError()
	}
	if !b {
		return "", util.NewBizErr(apicode.DataNotExistsCode, i18n.UserNotFound)
	}
	// 校验密码
	if user.Password != util.EncryptUserPassword(reqDTO.Password) {
		return "", util.NewBizErr(apicode.WrongLoginPasswordCode, i18n.UserWrongPassword)
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
		return "", util.InternalError()
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
		return util.InternalError()
	}
	if !b {
		return util.InvalidArgsError()
	}
	if session.UserInfo.Account != reqDTO.Operator.Account {
		return util.InvalidArgsError()
	}
	err = sessionStore.DeleteBySessionId(reqDTO.SessionId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
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
		return util.NewBizErr(apicode.InvalidArgsCode, i18n.UserAlreadyExists)
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
		return util.NewBizErr(apicode.UserAlreadyExistsCode, i18n.UserAlreadyExists)
	}
	userCount, err := usermd.CountUser(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Errorf("RegisterUser err: %v", err)
		return util.NewBizErr(apicode.InvalidArgsCode, i18n.SystemInternalError)
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
		return util.InternalError()
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
		return util.InvalidArgsError()
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

func UpdateUser(ctx context.Context, reqDTO UpdateUserReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	// 系统管理员或本人才能编辑user
	if !reqDTO.Operator.IsAdmin && reqDTO.Account != reqDTO.Operator.Account {
		return util.UnauthorizedError()
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	_, b, err := usermd.GetByAccount(ctx, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	// 账号不存在
	if !b {
		return util.InvalidArgsError()
	}
	if _, err = usermd.UpdateUser(ctx, usermd.UpdateUserReqDTO{
		Account: reqDTO.Account,
		Name:    reqDTO.Name,
		Email:   reqDTO.Email,
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	return nil
}

func UpdateAdmin(ctx context.Context, reqDTO UpdateAdminReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	// 只有系统管理员才能设置系统管理员
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	// 系统管理员不能处理自己
	if reqDTO.Operator.Account == reqDTO.Account {
		return util.InvalidArgsError()
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	_, b, err := usermd.GetByAccount(ctx, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	// 账号不存在
	if !b {
		return util.InvalidArgsError()
	}
	if _, err = usermd.UpdateAdmin(ctx, usermd.UpdateAdminReqDTO{
		Account: reqDTO.Account,
		IsAdmin: reqDTO.IsAdmin,
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	return nil
}

func UpdatePassword(ctx context.Context, reqDTO UpdatePasswordReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	// 如果是系统管理员可以设置任何人密码
	if !reqDTO.Operator.IsAdmin {
		// 否则只能设置自己的密码
		if reqDTO.Account != reqDTO.Operator.Account {
			return util.UnauthorizedError()
		}
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	_, b, err := usermd.GetByAccount(ctx, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	// 账号不存在
	if !b {
		return util.InvalidArgsError()
	}
	if _, err = usermd.UpdatePassword(ctx, usermd.UpdatePasswordReqDTO{
		Account:  reqDTO.Account,
		Password: reqDTO.Password,
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError()
	}
	return nil
}
