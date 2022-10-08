package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"net/http"
)

var store = base64Captcha.DefaultMemStore

//TODO [GetCaptcha] 验证码生成服务
func GetCaptcha(ctx *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	id, n64s, err := cp.Generate()
	if err != nil {
		zap.S().Info("生成验证码错误:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成验证码错误",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"catchaId": id,
		"picPath":  n64s,
	})
}

func Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})
}
