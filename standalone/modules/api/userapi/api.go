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
		group := e.Group("/api/user")
		{
			// 登录
			group.POST("/login", login)
			// 注册用户
			group.POST("/register", register)
			// 退出登录
			group.Any("/loginOut", apicommon.CheckLogin, loginOut)
			// 新增用户
			group.POST("/insert", apicommon.CheckLogin, insertUser)
			// 删除用户
			group.POST("/delete", apicommon.CheckLogin, deleteUser)
			// 更新用户
			group.POST("/update", apicommon.CheckLogin, updateUser)
			// 展示用户列表
			group.POST("/list", apicommon.CheckLogin, listUser)
			// 更新密码
			group.POST("/changePassword", apicommon.CheckLogin, changePassword)
		}
	})
}

func login(c *gin.Context) {
	var reqVO LoginReqVO
	if util.ShouldBindJSON(&reqVO, c) {
		sessionId, err := usersrv.Login(c.Request.Context(), usersrv.LoginReqDTO{
			Account:  reqVO.Account,
			Password: reqVO.Password,
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, LoginRespVO{
			BaseResp:  ginutil.DefaultSuccessResp,
			SessionId: sessionId,
		})
	}
}

func loginOut(c *gin.Context) {
	sessionId := apicommon.GetSessionId(c)
	if sessionId != "" {
		err := usersrv.LoginOut(c.Request.Context(), usersrv.LoginOutReqDTO{
			SessionId: sessionId,
			Operator:  apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
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
	var reqVO DeleteUserReqVO
	if util.ShouldBindJSON(&reqVO, c) {
		err := usersrv.DeleteUser(c.Request.Context(), usersrv.DeleteUserReqDTO{
			Account:  reqVO.Account,
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

func changePassword(c *gin.Context) {

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
