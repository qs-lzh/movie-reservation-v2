package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qs-lzh/movie-reservation/internal/app"
	"github.com/qs-lzh/movie-reservation/internal/dto"
	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/service"
)

type AuthHandler struct {
	App *app.App
}

func NewAuthHandler(app *app.App) *AuthHandler {
	return &AuthHandler{
		App: app,
	}
}

type RegisterRequest struct {
	UserName          string         `json:"username" binding:"required"`
	Password          string         `json:"password" binding:"required"`
	Role              model.UserRole `json:"user_role" binding:"required"`
	CacheKey          string         `json:"key" binding:"required"`
	AdminRolePassword string         `json:"admin_role_password"`
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		dto.BadRequest(ctx, "Invalid request body")
		return
	}

	valid, err := h.App.Cache.GetBool(req.CacheKey)
	if err != nil {
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to get from cache")
		return
	}
	if !valid {
		ctx.Error(err)
		dto.Unauthorized(ctx, "Captcha not passed")
		return
	}

	// If the user want to register an admin account, need admin-role-password
	if req.Role == model.RoleAdmin {
		if req.AdminRolePassword == "" {
			ctx.Error(err)
			dto.BadRequest(ctx, "Invalid request body: need admin role password to register an admin")
			return
		}
		if req.AdminRolePassword != h.App.Config.AdminRolePassword {
			ctx.Error(err)
			dto.Unauthorized(ctx, "You do not have the permission to register an admin account: Wrong admin role password")
			return
		}
	}

	if err := h.App.UserService.CreateUser(req.UserName, req.Password, req.Role); err != nil {
		if errors.Is(err, service.ErrAlreadyExists) {
			ctx.Error(err)
			dto.Conflict(ctx, "USER_CONFLICTS", fmt.Sprintf("User named %s already exists", req.UserName))
			return
		}
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to create user")
		return
	}

	dto.Success(ctx, 201, fmt.Sprintf("Created user named %s successfully", req.UserName))
}

type LoginRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	CacheKey string `json:"key" binding:"required"`
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		dto.BadRequest(ctx, "Invalid request body")
		return
	}

	valid, err := h.App.Cache.GetBool(req.CacheKey)
	if err != nil {
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to get from cache")
		return
	}
	if !valid {
		ctx.Error(err)
		dto.Unauthorized(ctx, "Captcha not passed")
		return
	}
	tokenStr, err := h.App.AuthService.Login(req.UserName, req.Password, req.CacheKey)
	if err != nil {
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to login")
		return
	}
	// change the parameter secure to true when deploy
	ctx.SetCookie("jwt", tokenStr, 3600, "/", "", false, true)

	// Get user role to return in response
	userRole, err := h.App.UserService.GetUserRoleByName(req.UserName)
	if err != nil {
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to get user role")
		return
	}

	dto.Success(ctx, http.StatusOK, gin.H{
		"status":   "Login successfully",
		"username": req.UserName,
		"role":     userRole,
	})
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	ctx.SetCookie("jwt", "", -1, "/", "", false, true)
	dto.Success(ctx, http.StatusOK, "Logged out successfully")
}
