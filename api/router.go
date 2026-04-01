package api

import (
	"go-api-project/api/v1"
	"go-api-project/config"
	"go-api-project/internal/middleware"
	"go-api-project/internal/service"
	"go-api-project/pkg/jwt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(
	cfg *config.Config,
	jwtManager *jwt.Manager,
	userService *service.UserService,
) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 全局中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(&cfg.CORS))

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	apiV1 := r.Group("/api/v1")

	// 认证模块（不需要JWT）
	authHandler := v1.NewAuthHandler(userService)
	authGroup := apiV1.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)
	}

	// 需要JWT认证的路由
	authorized := apiV1.Group("")
	authorized.Use(middleware.JWTAuth(jwtManager))
	{
		// 认证相关
		authGroup := authorized.Group("/auth")
		{
			authGroup.POST("/logout", authHandler.Logout)
		}

		// 用户模块
		userHandler := v1.NewUserHandler(userService)
		userGroup := authorized.Group("/users")
		{
			userGroup.GET("/me", userHandler.GetMe)
			userGroup.GET("", middleware.AdminAuth(), userHandler.ListUsers)
			userGroup.POST("", middleware.AdminAuth(), userHandler.CreateUser)
			userGroup.GET("/:id", userHandler.GetUser)
			userGroup.PUT("/:id", userHandler.UpdateUser)
			userGroup.DELETE("/:id", middleware.AdminAuth(), userHandler.DeleteUser)
			userGroup.POST("/change-password", userHandler.ChangePassword)
		}
	}

	return r
}
