package auth

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type AuthHandler struct {
	svc *Service
}

func NewAuthHandler(svc *Service) *AuthHandler { return &AuthHandler{svc: svc} }

type loginReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type tokenResp struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type registerReq struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req loginReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	access, refresh, exp, err := h.svc.Login(c.Request().Context(), req.Phone, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid credentials"})
	}

	return c.JSON(http.StatusOK, tokenResp{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    exp,
	})
}

// POST /register
func (h *AuthHandler) Register(c echo.Context) error {
	var req registerReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	id, err := h.svc.Register(c.Request().Context(), req.Phone, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"user_id": id})
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	var req refreshReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	access, refresh, exp, err := h.svc.Refresh(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid or expired refresh token"})
	}

	return c.JSON(http.StatusOK, tokenResp{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    exp,
	})
}

// POST /logout
func (h *AuthHandler) Logout(c echo.Context) error {
	var req logoutReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.svc.Logout(c.Request().Context(), req.RefreshToken); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "logout failed"})
	}

	return c.NoContent(http.StatusNoContent)
}

type logoutReq struct {
	RefreshToken string `json:"refresh_token"`
}
