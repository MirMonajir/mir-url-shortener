package application

import (
	"github.com/MirMonajir/mir-url-shortener/internal_logic/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HTTPHandler struct {
	shortener domain.Shortener
	validator *domain.URLValidator
}

func NewHTTPHandler(s domain.Shortener) *HTTPHandler {
	return &HTTPHandler{
		shortener: s,
		validator: domain.NewURLValidator(),
	}
}

type shortenReq struct {
	URL string `json:"url" binding:"required"`
}

type shortenResp struct {
	ShortURL string `json:"short_url"`
}

type errorResp struct {
	ErrorType string `json:"error_type"`
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
	Code      int    `json:"code"`
}

func (h *HTTPHandler) ShortenURL(c *gin.Context) {
	var req shortenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := domain.NewAppErrorWithDetails(
			domain.ErrInvalidInput,
			"Request body is invalid",
			err.Error(),
			http.StatusBadRequest,
		)
		c.JSON(appErr.Code, errorResp{
			ErrorType: string(appErr.Type),
			Message:   appErr.Message,
			Details:   appErr.Details,
			Code:      appErr.Code,
		})
		return
	}

	// Validate URL format, length, scheme, and SSRF
	if validationErr := h.validator.Validate(req.URL); validationErr != nil {
		c.JSON(validationErr.Code, errorResp{
			ErrorType: string(validationErr.Type),
			Message:   validationErr.Message,
			Details:   validationErr.Details,
			Code:      validationErr.Code,
		})
		return
	}

	short, err := h.shortener.Shorten(req.URL)
	if err != nil {
		// Check if it's an app error
		if appErr, ok := err.(*domain.AppError); ok {
			c.JSON(appErr.Code, errorResp{
				ErrorType: string(appErr.Type),
				Message:   appErr.Message,
				Details:   appErr.Details,
				Code:      appErr.Code,
			})
		} else {
			// Generic error
			c.JSON(http.StatusInternalServerError, errorResp{
				ErrorType: string(domain.ErrInternal),
				Message:   "Failed to shorten URL",
				Details:   err.Error(),
				Code:      http.StatusInternalServerError,
			})
		}
		return
	}

	c.JSON(http.StatusOK, shortenResp{ShortURL: short})
}

func (h *HTTPHandler) Redirect(c *gin.Context) {
	code := c.Param("shortenedurl")

	if code == "" {
		appErr := domain.NewAppError(
			domain.ErrInvalidInput,
			"Shortened URL code is required",
			http.StatusBadRequest,
		)
		c.JSON(appErr.Code, errorResp{
			ErrorType: string(appErr.Type),
			Message:   appErr.Message,
			Code:      appErr.Code,
		})
		return
	}

	orig, err := h.shortener.Resolve(code)
	if err != nil {
		appErr := domain.NewAppError(
			domain.ErrURLNotFound,
			"Shortened URL not found",
			http.StatusNotFound,
		)
		c.JSON(appErr.Code, errorResp{
			ErrorType: string(appErr.Type),
			Message:   appErr.Message,
			Code:      appErr.Code,
		})
		return
	}

	c.Redirect(http.StatusFound, orig)
}

func (h *HTTPHandler) Metrics(c *gin.Context) {
	top := h.shortener.TopDomains(3)
	c.JSON(http.StatusOK, gin.H{"top_domains": top})
}
