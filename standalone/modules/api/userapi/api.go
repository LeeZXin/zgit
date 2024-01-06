package userapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/timeutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"zgit/standalone/modules/api/apicommon"
	"zgit/standalone/modules/service/usersrv"
	"zgit/util"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/login")
		{
			// 登录
			group.POST("/login", login)
			// 注册用户
			group.POST("/register", register)
			// 退出登录
			group.Any("/loginOut", apicommon.CheckLogin, loginOut)
		}
		group = e.Group("/api/user", apicommon.CheckLogin)
		{
			// 新增用户
			group.POST("/insert", insertUser)
			// 删除用户
			group.POST("/delete", deleteUser)
			// 更新用户
			group.POST("/update", updateUser)
			// 展示用户列表
			group.POST("/list", listUser)
			// 更新密码
			group.POST("/updatePassword", updatePassword)
			// 系统管理员设置
			group.POST("/setAdmin", updateAdmin)
		}
	})
}

func login(c *gin.Context) {
	var req LoginReqVO
	if util.ShouldBindJSON(&req, c) {
		sessionId, err := usersrv.Login(c.Request.Context(), usersrv.LoginReqDTO{
			Account:  req.Account,
			Password: req.Password,
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.SetCookie(apicommon.LoginCookie, sessionId, int(usersrv.LoginSessionExpiry.Seconds()), "/", "", false, true)
		c.JSON(http.StatusOK, LoginRespVO{
			BaseResp:  ginutil.DefaultSuccessResp,
			SessionId: sessionId,
		})
	}
}

func loginOut(c *gin.Context) {
	err := usersrv.LoginOut(c.Request.Context(), usersrv.LoginOutReqDTO{
		SessionId: apicommon.GetSessionId(c),
		Operator:  apicommon.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
}

func insertUser(c *gin.Context) {
	var reqVO InsertUserReqVO
	if util.ShouldBindJSON(&reqVO, c) {
		err := usersrv.InsertUser(c.Request.Context(), usersrv.InsertUserReqDTO{
			Account:   reqVO.Account,
			Name:      reqVO.Name,
			Email:     reqVO.Email,
			Password:  reqVO.Password,
			IsAdmin:   reqVO.IsAdmin,
			AvatarUrl: reqVO.AvatarUrl,
			Operator:  apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
		} else {
			c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
		}
	}
}

func deleteUser(c *gin.Context) {
	var req DeleteUserReqVO
	if util.ShouldBindJSON(&req, c) {
		err := usersrv.DeleteUser(c.Request.Context(), usersrv.DeleteUserReqDTO{
			Account:  req.Account,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
		} else {
			c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
		}
	}
}

func updateUser(c *gin.Context) {
	var req UpdateUserReqVO
	if util.ShouldBindJSON(&req, c) {
		err := usersrv.UpdateUser(c.Request.Context(), usersrv.UpdateUserReqDTO{
			Account:  req.Account,
			Name:     req.Name,
			Email:    req.Email,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
		} else {
			c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
		}
	}
}

func listUser(c *gin.Context) {
	var req ListUserReqVO
	if util.ShouldBindJSON(&req, c) {
		respDTO, err := usersrv.ListUser(c.Request.Context(), usersrv.ListUserReqDTO{
			Account:  req.Account,
			Offset:   req.Offset,
			Limit:    req.Limit,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		ret := ListUserRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Cursor:   respDTO.Cursor,
		}
		ret.UserList, _ = listutil.Map(respDTO.UserList, func(t usersrv.UserDTO) (UserVO, error) {
			return UserVO{
				Account:      t.Account,
				Name:         t.Name,
				Email:        t.Email,
				IsAdmin:      t.IsAdmin,
				IsProhibited: t.IsProhibited,
				AvatarUrl:    t.AvatarUrl,
				Created:      t.Created.Format(timeutil.DefaultTimeFormat),
				Updated:      t.Updated.Format(timeutil.DefaultTimeFormat),
			}, nil
		})
		c.JSON(http.StatusOK, ret)
	}
}

func updatePassword(c *gin.Context) {
	var req UpdatePasswordReqVO
	if util.ShouldBindJSON(&req, c) {
		err := usersrv.UpdatePassword(c.Request.Context(), usersrv.UpdatePasswordReqDTO{
			Account:  req.Account,
			Password: req.Password,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
		} else {
			c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
		}
	}
}

func register(c *gin.Context) {
	var req RegisterUserReqVO
	if util.ShouldBindJSON(&req, c) {
		err := usersrv.RegisterUser(c.Request.Context(), usersrv.RegisterUserReqDTO{
			Account:   req.Account,
			Name:      req.Name,
			Email:     req.Email,
			Password:  req.Password,
			AvatarUrl: req.AvatarUrl,
		})
		if err != nil {
			util.HandleApiErr(err, c)
		} else {
			c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
		}
	}
}

func updateAdmin(c *gin.Context) {
	var req UpdateAdminReqVO
	if util.ShouldBindJSON(&req, c) {
		err := usersrv.UpdateAdmin(c.Request.Context(), usersrv.UpdateAdminReqDTO{
			Account:  req.Account,
			IsAdmin:  req.IsAdmin,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
		} else {
			c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
		}
	}
}
