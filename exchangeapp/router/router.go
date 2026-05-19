package router //路由器

import (
	"exchangeapp/controllers"
	"exchangeapp/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, //允许的请求头,包括Authorization用于携带JWT token
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, //设置最大缓存时间为12小时
	}))

	auth := r.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login) //登录

		auth.POST("/register", controllers.Register) //注册
	}

	api := r.Group("/api")
	api.GET("/exchangerates", controllers.GetExchangeRate)
	api.Use(middlewares.AuthMiddleware()) //使用AuthMiddleware中间件保护/exchangerates路由,只有通过验证的请求才能访问这个路由
	{
		api.POST("/exchangerates", controllers.CreateExchangeRate) //创建汇率
		api.POST("/articles", controllers.CreateAritcle)           //创建文章
		api.GET("/articles", controllers.GetArticles)              //获取所有文章
		api.GET("/articles/:id", controllers.GetArticleById)       //获取指定文章

		api.POST("/articles/:id/like", controllers.LikeArticles)   //点赞文章
		api.GET("/articles/:id/like", controllers.GetArticleLikes) //获取文章点赞数
	}
	return r
}
