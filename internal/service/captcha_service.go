package service

import (
	"fmt"
	"image"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/google/uuid"
	"github.com/wenlng/go-captcha-assets/resources/fonts/fzshengsksjw"
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/click"

	"github.com/qs-lzh/movie-reservation/internal/cache"
)

type CaptchaService interface {
	Generate() (mBase64, tBase64, cacheKey string, err error)
	Verify(clickData []Dot, dotAnswerData map[int]*click.Dot) bool
	VerifyWithKey(clickData []Dot, cacheKey string) (bool, error)
}

type captchaService struct {
	cache *cache.RedisCache
}

func NewCaptchaService(cache *cache.RedisCache) *captchaService {
	return &captchaService{
		cache: cache,
	}
}

var _ CaptchaService = (*captchaService)(nil)

var textCapt click.Captcha

func (s *captchaService) Generate() (mBase64, tBase64, cacheKey string, err error) {
	builder := click.NewBuilder(
		click.WithRangeLen(option.RangeVal{Min: 4, Max: 6}),
		click.WithRangeVerifyLen(option.RangeVal{Min: 2, Max: 4}),
	)

	fontN, err := fzshengsksjw.GetFont()
	if err != nil {
		return "", "", "", fmt.Errorf("Failed to load font: %w", err)
	}
	bgImage, err := imagesv2.GetImages()
	if err != nil {
		return "", "", "", fmt.Errorf("Failed to load background images: %w", err)
	}

	builder.SetResources(
		click.WithChars([]string{"1A", "5E", "3d", "0p", "78", "DL", "CB", "9M"}),
		click.WithFonts([]*truetype.Font{
			fontN,
		}),
		click.WithBackgrounds([]image.Image{
			bgImage[0],
		}),
	)

	textCapt = builder.Make()

	captData, err := textCapt.Generate()
	if err != nil {
		return "", "", "", fmt.Errorf("Failed to generate captcha: %w", err)
	}

	dotAnswerData := captData.GetData()
	captchaID := uuid.New().String()

	err = s.cache.Set(captchaID, dotAnswerData, 5*time.Minute)
	if err != nil {
		return "", "", "", fmt.Errorf("Failed to save captcha answer to redis: %w", err)
	}

	mBase64, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		return "", "", "", fmt.Errorf("Failed to get master image of captcha")
	}
	tBase64, err = captData.GetThumbImage().ToBase64()
	if err != nil {
		return "", "", "", fmt.Errorf("Failed to get thumb image of captcha")
	}
	return mBase64, tBase64, captchaID, nil
}

// for parse data from frontend
type Dot struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (s *captchaService) Verify(clickData []Dot, dotAnswerData map[int]*click.Dot) bool {

	if len(clickData) != len(dotAnswerData) {
		return false
	}
	fmt.Println("1")

	// the key of dotAnswerData begin with 0
	chkRet := false
	for idx, dot := range clickData {
		answerDot := dotAnswerData[idx]
		// noticing that the answerDot.Y is always larger than actual, I subtract 30 from it
		chkRet = click.Validate(dot.X, dot.Y, answerDot.X, max(0, answerDot.Y-45), answerDot.Width+10, answerDot.Height+10, 5)
		if !chkRet {
			return false
		}
	}
	return true
}

func (s *captchaService) VerifyWithKey(clickData []Dot, cacheKey string) (bool, error) {
	dotAnswerData := make(map[int]*click.Dot)
	if err := s.cache.Get(cacheKey, &dotAnswerData); err != nil {
		return false, fmt.Errorf("failed to get captcha answer data from cache: %v", err)
	}

	valid := s.Verify(clickData, dotAnswerData)

	return valid, nil
}
