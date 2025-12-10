package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/qs-lzh/movie-reservation/internal/app"
	"github.com/qs-lzh/movie-reservation/internal/dto"
	"github.com/qs-lzh/movie-reservation/internal/service"
)

type CaptchaHandler struct {
	App *app.App
}

func NewCaptchaHandler(app *app.App) *CaptchaHandler {
	return &CaptchaHandler{
		App: app,
	}
}

func (h *CaptchaHandler) GenerateCaptcha(ctx *gin.Context) {
	mBase64, tBase64, cacheKey, err := h.App.CaptchaService.Generate()
	if err != nil {
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to generate captcha")
		return
	}
	dto.Success(ctx, 200, gin.H{
		"image": mBase64,
		"thumb": tBase64,
		"key":   cacheKey,
	})
}

func (h *CaptchaHandler) VerifyCaptcha(ctx *gin.Context) {
	type CaptchaVerifyRequest struct {
		Dots []service.Dot `json:"dots" binding:"required"`
		Key  string        `json:"key" binding:"required"`
	}

	var req CaptchaVerifyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		dto.BadRequest(ctx, "Invalid request body: "+err.Error())
		return
	}

	valid, err := h.App.CaptchaService.VerifyWithKey(req.Dots, req.Key)
	if err != nil {
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to verify captcha: "+err.Error())
		return
	}

	if err := h.App.Cache.SetBool(req.Key, valid); err != nil {
		ctx.Error(err)
		dto.InternalServerError(ctx, "Failed to set bool in cache")
		return
	}

	dto.Success(ctx, 200, gin.H{
		"success": valid,
	})
}
