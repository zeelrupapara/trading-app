package v1

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/zeelrupapara/trading-api/services"
	"github.com/zeelrupapara/trading-api/structs"

	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/zeelrupapara/trading-api/config"
	"github.com/zeelrupapara/trading-api/constants"
	"github.com/zeelrupapara/trading-api/models"
	jwt "github.com/zeelrupapara/trading-api/pkg/jwt"
	"github.com/zeelrupapara/trading-api/utils"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

type AuthController struct {
	userService *services.UserService
	userModel   *models.UserModel
	logger      *zap.Logger
	config      config.AppConfig
}

func NewAuthController(goqu *goqu.Database, logger *zap.Logger, config config.AppConfig) (*AuthController, error) {
	userModel, err := models.InitUserModel(goqu)
	if err != nil {
		return nil, err
	}

	userSvc := services.NewUserService(&userModel)

	return &AuthController{
		userService: userSvc,
		userModel:   &userModel,
		logger:      logger,
		config:      config,
	}, nil
}

// DoAuth authenticate user with email and password
// swagger:route POST /api/v1/login Auth RequestAuthnUser
//
// Authenticate user with email and password.
//
//			Consumes:
//			- application/json
//
//			Schemes: http, https
//
//			Responses:
//			  200: ResponseAuthnUser
//		   400: GenericResFailBadRequest
//	    401: ResForbiddenRequest
//			  500: GenericResError
func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var reqLoginUser structs.ReqLoginUser

	err := json.Unmarshal(c.Body(), &reqLoginUser)
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	err = validate.Struct(reqLoginUser)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	user, err := ctrl.userService.Authenticate(reqLoginUser.Email, reqLoginUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusUnauthorized, constants.InvalidCredentials)
		}
		ctrl.logger.Error("error while get user by email and password", zap.Error(err), zap.Any("email", reqLoginUser.Email))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrLoginUser)
	}

	// token is valid for 1 hour
	token, err := jwt.CreateToken(ctrl.config, user.ID, time.Now().Add(time.Hour*1))
	if err != nil {
		ctrl.logger.Error("error while creating token", zap.Error(err), zap.Any("id", user.ID))
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrLoginUser)
	}

	userCookie := &fiber.Cookie{
		Name:    constants.CookieUser,
		Value:   token,
		Expires: time.Now().Add(1 * time.Hour),
	}
	c.Cookie(userCookie)

	return utils.JSONSuccess(c, http.StatusOK, user)
}

// Logout clear user cookie
// swagger:route POST /api/v1/logout Auth RequestLogout
//
// Clear user cookie.
//
//			Consumes:
//			- application/json
//
//			Schemes: http, https
//
//			Responses:
//			  204: NoContent
//		   400: GenericResFailBadRequest
//	    401: ResForbiddenRequest
//			  500: GenericResError
func (ctrl *AuthController) Logout(c *fiber.Ctx) error {
	userCookie := &fiber.Cookie{
		Name:    constants.CookieUser,
		Value:   "",
		Expires: time.Now(),
	}
	c.Cookie(userCookie)

	return utils.JSONSuccess(c, http.StatusNoContent, nil)
}