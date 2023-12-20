package util

import (
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"zgit/pkg/apicode"
	"zgit/pkg/i18n"
)

func HandleApiErr(err error, c *gin.Context) {
	if err != nil {
		berr, ok := err.(*bizerr.Err)
		if !ok {
			c.JSON(http.StatusInternalServerError, ginutil.BaseResp{
				Code:    apicode.InternalErrorCode.Int(),
				Message: i18n.GetByKey(i18n.SystemInternalError),
			})
		} else {
			c.JSON(http.StatusOK, ginutil.BaseResp{
				Code:    berr.Code,
				Message: berr.Message,
			})
		}
	}
}

func ShouldBindJSON(obj any, c *gin.Context) bool {
	err := c.ShouldBindJSON(obj)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginutil.BaseResp{
			Code:    apicode.BadRequestCode.Int(),
			Message: i18n.GetByKey(i18n.SystemInvalidArgs),
		})
		return false
	}
	return true
}
