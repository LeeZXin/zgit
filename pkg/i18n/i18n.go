package i18n

import (
	"github.com/LeeZXin/zsf-utils/i18n"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"path/filepath"
	"zgit/setting"
)

const (
	ZH_CN = "zh-CN"
	EN_US = "en-US"
)

func init() {
	iniPath := filepath.Join(common.ResourcesDir, "i18n", setting.Lang()+".ini")
	locale, err := i18n.NewImmutableLocaleFromIniFile(iniPath, setting.Lang())
	if err == nil {
		logger.Logger.Infof("init i18n ini path: %s successfully", iniPath)
		i18n.AddLocale(locale)
		i18n.SetDefaultLocale(setting.Lang())
	} else {
		iniPath = filepath.Join(common.ResourcesDir, "i18n", ZH_CN+".ini")
		locale, err = i18n.NewImmutableLocaleFromIniFile(iniPath, ZH_CN)
		if err == nil {
			logger.Logger.Infof("init i18n ini path: %s successfully", iniPath)
			i18n.AddLocale(locale)
			i18n.SetDefaultLocale(ZH_CN)
		}
	}
}

func SupportedLangeList() []string {
	return []string{
		ZH_CN, EN_US,
	}
}

func GetByKey(key Key) string {
	return i18n.GetOrDefault(string(key), key.DefaultRet())
}
