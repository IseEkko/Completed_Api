package middlewares

import (
	"Completed_Api/user-web/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clamis, _ := ctx.Get("claims")
		currentUser := clamis.(*models.CustomClaims)

		if currentUser.AuthorityId == 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
