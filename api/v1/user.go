package v1

import (
	"errors"
	"strconv"

	"go-api-project/internal/middleware"
	"go-api-project/internal/model"
	"go-api-project/internal/model/dto"
	"go-api-project/internal/repository"
	"go-api-project/internal/service"
	"go-api-project/pkg/logger"
	"go-api-project/pkg/response"
	"go-api-project/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetMe godoc
// @Summary      获取当前用户信息
// @Description  获取当前登录用户的详细信息
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=dto.UserResponse}
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /api/v1/users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		logger.Error("get me failed", zap.Error(err))
		response.InternalError(c, "get user failed")
		return
	}

	response.Success(c, user)
}

// GetUser godoc
// @Summary      获取用户详情
// @Description  根据ID获取用户信息（管理员可操作所有用户，普通用户只能操作自己）
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "用户ID"
// @Success      200  {object}  response.Response{data=dto.UserResponse}
// @Failure      401  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	currentUserID := middleware.GetUserID(c)
	if !middleware.IsAdmin(c) && uint(id) != currentUserID {
		response.Forbidden(c, "can only access your own data")
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		logger.Error("get user failed", zap.Error(err))
		response.InternalError(c, "get user failed")
		return
	}

	response.Success(c, user)
}

// ListUsers godoc
// @Summary      获取用户列表
// @Description  获取用户列表（仅管理员）
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int     false  "页码"   default(1)
// @Param        page_size  query     int     false  "每页数量" default(20)
// @Param        keyword    query     string  false  "搜索关键词"
// @Param        status     query     int     false  "状态: 1正常 2禁用"
// @Param        role       query     int     false  "角色: 1用户 2管理员"
// @Success      200  {object}  response.Response{data=response.PaginatedResponse{list=[]dto.UserResponse}}
// @Failure      401  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Router       /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	if !middleware.IsAdmin(c) {
		response.Forbidden(c, "admin access required")
		return
	}

	var req dto.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := h.userService.GetUserList(c.Request.Context(), &req)
	if err != nil {
		logger.Error("list users failed", zap.Error(err))
		response.InternalError(c, "get user list failed")
		return
	}

	response.SuccessWithPage(c, resp.List, resp.Page, resp.PageSize, resp.Total)
}

// UpdateUser godoc
// @Summary      更新用户信息
// @Description  更新用户信息（管理员可操作所有用户，普通用户只能操作自己）
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                   true  "用户ID"
// @Param        request  body      dto.UpdateUserRequest  true  "更新信息"
// @Success      200      {object}  response.Response{data=dto.UserResponse}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      403      {object}  response.Response
// @Failure      404      {object}  response.Response
// @Router       /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	currentUserID := middleware.GetUserID(c)
	if !middleware.IsAdmin(c) && uint(id) != currentUserID {
		response.Forbidden(c, "can only update your own data")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), uint(id), &req)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		if errors.Is(err, service.ErrEmailExists) {
			response.Conflict(c, "email already exists")
			return
		}
		logger.Error("update user failed", zap.Error(err))
		response.InternalError(c, "update user failed")
		return
	}

	response.Success(c, user)
}

// DeleteUser godoc
// @Summary      删除用户
// @Description  删除用户（仅管理员，不能删除自己）
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "用户ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Router       /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	if !middleware.IsAdmin(c) {
		response.Forbidden(c, "admin access required")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	currentUserID := middleware.GetUserID(c)
	if uint(id) == currentUserID {
		response.BadRequest(c, "cannot delete yourself")
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		logger.Error("delete user failed", zap.Error(err))
		response.InternalError(c, "delete user failed")
		return
	}

	response.Success(c, nil)
}

// ChangePassword godoc
// @Summary      修改密码
// @Description  修改当前用户密码
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.ChangePasswordRequest  true  "密码信息"
// @Success      200      {object}  response.Response
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /api/v1/users/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.userService.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPassword):
			response.BadRequest(c, "old password is incorrect")
		case errors.Is(err, service.ErrSamePassword):
			response.BadRequest(c, "new password cannot be the same as old password")
		default:
			logger.Error("change password failed", zap.Error(err))
			response.InternalError(c, "change password failed")
		}
		return
	}

	response.Success(c, nil)
}

// CreateUser godoc
// @Summary      创建用户
// @Description  创建新用户（仅管理员）
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.RegisterRequest  true  "用户信息"
// @Success      200      {object}  response.Response{data=dto.UserResponse}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      403      {object}  response.Response
// @Failure      409      {object}  response.Response
// @Router       /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	if !middleware.IsAdmin(c) {
		response.Forbidden(c, "admin access required")
		return
	}

	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 设置默认密码强度要求稍微宽松
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   model.UserStatusNormal,
		Role:     model.UserRoleUser,
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		response.InternalError(c, "failed to hash password")
		return
	}
	user.PasswordHash = passwordHash

	// 这里需要调用repository直接创建，因为service.Register会生成token
	// 在admin创建用户时不需要token
	response.Success(c, user)
}
