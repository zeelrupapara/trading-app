package middlewares

import (
	"errors"
	"net/http"

	"github.com/zeelrupapara/trading-api/constants"
	"github.com/zeelrupapara/trading-api/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	j "github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/zeelrupapara/trading-api/utils"
	"go.uber.org/zap"
)

func (m *Middleware) Authenticated(c *fiber.Ctx) error {
	token := c.Cookies(constants.CookieUser, "")
	if token == "" {
		return utils.JSONFail(c, http.StatusUnauthorized, constants.Unauthenticated)
	}

	claims, err := jwt.ParseToken(m.config, token)
	if err != nil {
		if errors.Is(err, j.ErrInvalidJWT()) || errors.Is(err, j.ErrTokenExpired()) {
			return utils.JSONFail(c, http.StatusUnauthorized, constants.Unauthenticated)
		}

		m.logger.Error("error while checking user identity", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrUnauthenticated)
	}

	c.Locals(constants.ContextUid, claims.Subject())
	return c.Next()
}
