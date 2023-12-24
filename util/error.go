package util

import (
	"github.com/LeeZXin/zsf-utils/bizerr"
	"zgit/pkg/apicode"
	"zgit/pkg/i18n"
)

func InvalidArgsError() error {
	return bizerr.NewBizErr(apicode.InvalidArgsCode.Int(), i18n.GetByKey(i18n.SystemInvalidArgs))
}

func InternalError() error {
	return bizerr.NewBizErr(apicode.InternalErrorCode.Int(), i18n.GetByKey(i18n.SystemInternalError))
}

func UnauthorizedError() error {
	return bizerr.NewBizErr(apicode.UnauthorizedCode.Int(), i18n.GetByKey(i18n.SystemUnauthorized))
}
