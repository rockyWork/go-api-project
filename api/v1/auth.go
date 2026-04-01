package v1

import (
	"errors"

	"go-api-project/internal/model/dto"
	"go-api-project/internal/service"
	"go-api-project/pkg/logger"
	"go-api-project/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Register godoc
// @Summary      用户注册
// @Description  新用户注册，创建账户并返回token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RegisterRequest  true  "注册信息"
// @Success      200      {object}  response.Response{data=dto.TokenResponse}
// @Failure      400      {object}  response.Response
// @Failure      409      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tokenResp, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		logger.Error("register failed", zap.Error(err))
		switch {
		case errors.Is(err, service.ErrUsernameExists):
			response.Conflict(c, "username already exists")
		case errors.Is(err, service.ErrEmailExists):
			response.Conflict(c, "email already exists")
		default:
			response.InternalError(c, "register failed")
		}
		return
	}

	response.Success(c, tokenResp)
}

// Login godoc
// @Summary      用户登录
// @Description  使用用户名和密码登录
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest  true  "登录信息"
// @Success      200      {object}  response.Response{data=dto.TokenResponse}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tokenResp, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		logger.Warn("login failed", zap.String("username", req.Username), zap.Error(err))
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			response.Unauthorized(c, "invalid username or password")
		case errors.Is(err, service.ErrUserBanned):
			response.Forbidden(c, "user has been banned")
		default:
			response.InternalError(c, "login failed")
		}
		return
	}

	response.Success(c, tokenResp)
}

// RefreshToken godoc
// @Summary      刷新Token
// @Description  使用refresh token获取新的access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RefreshTokenRequest  true  "刷新令牌"
// @Success      200      {object}  response.Response{data=dto.TokenResponse}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tokenResp, err := h.userService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		logger.Warn("refresh token failed", zap.Error(err))
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, tokenResp)
}

// Logout godoc
// @Summary      用户登出
// @Description  用户登出，使token失效
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 这里可以实现token黑名单逻辑
	// 当前版本仅返回成功
	response.Success(c, nil)
}
