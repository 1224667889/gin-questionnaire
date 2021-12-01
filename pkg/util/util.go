package util

import "main/pkg/setting"

// Setup 初始化工具类
func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}