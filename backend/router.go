package backend

import (
	"duifene_auto_sign/backend/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 启用CORS
	r.Use(CORSMiddleware())

	// 静态文件
	r.Static("/static", "./frontend/static")
	r.StaticFile("/", "./frontend/index.html")

	// 创建处理器
	h := handler.NewHandler()

	// 主路由组
	dfysign := r.Group("/dfysign")
	{
		// 认证路由
		auth := dfysign.Group("/auth")
		{
			auth.POST("/login/wechat", h.WechatLogin)
			auth.POST("/login/password", h.PasswordLogin)
			auth.GET("/check", h.CheckLogin)
		}

		// 课程路由
		course := dfysign.Group("/course")
		{
			course.GET("/list", h.GetCourseList)
		}

		// 签到路由
		sign := dfysign.Group("/sign")
		{
			sign.POST("/submit", h.Sign)
			sign.POST("/location", h.SignWithLocation)
			sign.GET("/status", h.CheckSignStatus)
		}
	}

	return r
}

// CORSMiddleware CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
