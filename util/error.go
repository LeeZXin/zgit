package util

import (
	"fmt"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"zgit/pkg/apicode"
	"zgit/pkg/i18n"
)

func InvalidArgsError() error {
	return NewBizErr(apicode.InvalidArgsCode, i18n.SystemInvalidArgs)
}

func InternalError() error {
	return NewBizErr(apicode.InternalErrorCode, i18n.SystemInternalError)
}

func UnauthorizedError() error {
	return NewBizErr(apicode.UnauthorizedCode, i18n.SystemUnauthorized)
}

func AlreadyExistsError() error {
	return NewBizErr(apicode.DataAlreadyExistsCode, i18n.SystemAlreadyExists)
}

func NewBizErr(code apicode.Code, key i18n.Key, msg ...string) error {
	if len(msg) == 0 {
		return bizerr.NewBizErr(code.Int(), i18n.GetByKey(key))
	}
	return bizerr.NewBizErr(code.Int(), fmt.Sprintf(i18n.GetByKey(key), msg))
}
