package middlewares

import (
	"exchangeapp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware函数是一个Gin中间件,用于验证请求中的JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization") //从请求头中获取Authorization字段的值,这个值应该是Bearer token
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "wrong credentials",
			})
			ctx.Abort() //终止请求的处理,不再继续执行后续的中间件
			return
		}
		username, err := utils.ParseJWT(token) //解析JWT token,获取用户名
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err,
			})
			ctx.Abort()
			return
		}
		ctx.Set("username", username) //将解析出的用户名存储在上下文中,以便后续的处理函数使用
		ctx.Next() //继续执行后续的中间件或处理函数

	}
}
