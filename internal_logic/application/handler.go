package application

import (
    "net/http"
"log"
    "github.com/gin-gonic/gin"
    "github.com/MirMonajir/mir-url-shortener/internal_logic/domain"
)

type HTTPHandler struct {
    shortener domain.Shortener
}

func NewHTTPHandler(s domain.Shortener) *HTTPHandler {
    return &HTTPHandler{shortener: s}
}

type shortenReq struct {
    URL string `json:"url" binding:"required"`
}
type shortenResp struct {
    ShortURL string `json:"short_url"`
}

func (h *HTTPHandler) ShortenURL(c *gin.Context) {
    var req shortenReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    short, err := h.shortener.Shorten(req.URL)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, shortenResp{ShortURL: short})
}

func (h *HTTPHandler) Redirect(c *gin.Context) {
    code := c.Param("shortenedurl")
	log.Printf("Redirect requested for code: %s", code)
    orig, err := h.shortener.Resolve(code)
    if err != nil {
        c.JSON(http.StatusNotFound, err.Error())
        return
    }
	log.Printf("Redirecting code %s to %s", code, orig)
    c.Redirect(http.StatusFound, orig)
}

func (h *HTTPHandler) Metrics(c *gin.Context) {
    top := h.shortener.TopDomains(3)
    c.JSON(http.StatusOK, gin.H{"top_domains": top})
}
