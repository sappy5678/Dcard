package transport

import (
	"net/http"
	"time"

	"github.com/sappy5678/dcard/pkg/domain"

	"github.com/labstack/echo"
)

type HTTP struct {
	Service domain.ShortURLService
}

func NewHTTP(svc domain.ShortURLService, r *echo.Group) {
	h := HTTP{Service: svc}

	// Get short URL
	// GET /{shortCode}
	r.GET("/:shortCode", h.get)

	ur := r.Group("/api/v1")

	// Create short url
	// POST /api/v1/urls/
	ur.POST("/urls", h.create)
}

type createReq struct {
	OriginalURL string `json:"url"`
	ExpireTime  string `json:"expireAt"`
}

type createResp struct {
	ShortCode string `json:"id"`
	ShortURL  string `json:"shortUrl"`
}

func (h HTTP) create(c echo.Context) error {
	req := createReq{}
	if err := c.Bind(&req); err != nil {
		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: err.Error()})
		return err
	}

	expireTime, err := time.Parse(time.RFC3339, req.ExpireTime)
	if err != nil || expireTime.IsZero() {
		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: domain.ErrShortURLInvalid.Error()})
		return err
	}

	short, err := h.Service.Create(c.Request().Context(), req.OriginalURL, uint64(expireTime.Unix()))
	if err != nil {
		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: domain.ErrShortURLInvalid.Error()})
		return err
	}

	resp := createResp{
		ShortCode: short.ShortCode,
		ShortURL:  short.ShortURL,
	}
	return c.JSON(http.StatusOK, resp)
}

func (h HTTP) get(c echo.Context) error {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		err := c.HTML(http.StatusOK, "mock home page") // should redirect to home page
		return err
	}

	short, err := h.Service.Get(c.Request().Context(), shortCode)
	if err != nil {
		err := c.JSON(http.StatusNotFound, domain.ErrorRespond{Error: domain.ErrShortURLNotFound.Error()})
		return err
	}

	return c.Redirect(http.StatusTemporaryRedirect, short.OriginalURL)
	// c.Response().Header().Set(echo.HeaderLocation, short.OriginalURL)
	// return c.JSON(http.StatusTemporaryRedirect, short)
}
