package system

import (
	"github.com/gin-gonic/gin"
	"prince-x/tools"
	"prince-x/tools/app"
	"prince-x/tools/captcha"
)

func GenerateCaptchaHandler(c *gin.Context) {
	id, b64s, err := captcha.DriverDigitFunc()
	tools.HasError(err, "验证码获取失败", 500)
	app.Custum(c, gin.H{
		"code": 200,
		"data": b64s,
		"id":   id,
		"msg":  "success",
	})
}
