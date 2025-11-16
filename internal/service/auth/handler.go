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

type LoginReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type tokenResp struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type RegisterReq struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Login
//
// @Description Auth
// @Summary	login
// @Tags auth
// @Accept json
// @Produce	json
// @Param		request	body		LoginReq	true	"body param"
// @Success	200				{object}		LoginReq
// @Router /login	[post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginReq
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

// Register
//
// @Description Auth
// @Summary	register
// @Tags auth
// @Accept json
// @Produce	json
// @Param		request	body		RegisterReq	true	"body param"
// @Success	200			{object}		RegisterReq
// @Router /register	[post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterReq
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
