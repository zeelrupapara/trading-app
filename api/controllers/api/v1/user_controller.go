package v1

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/doug-martin/goqu/v9"
	"github.com/zeelrupapara/trading-api/constants"
	"github.com/zeelrupapara/trading-api/models"
	"github.com/zeelrupapara/trading-api/services"
	"github.com/zeelrupapara/trading-api/structs"
	"github.com/zeelrupapara/trading-api/utils"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"

	"github.com/gofiber/fiber/v2"
)

// UserController for user controllers
type UserController struct {
	userService *services.UserService
	logger      *zap.Logger
}

// NewUserController returns a user
func NewUserController(goqu *goqu.Database, logger *zap.Logger) (*UserController, error) {
	userModel, err := models.InitUserModel(goqu)
	if err != nil {
		logger.Error("error while initializing user model", zap.Error(err))
		return nil, err
	}

	userSvc := services.NewUserService(&userModel)

	return &UserController{
		userService: userSvc,
		logger:      logger,
	}, nil
}

// UserGet get the user by id
// swagger:route GET /api/v1/users/{userId} Users RequestGetUser
//
// Get a user.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseGetUser
//	      400: GenericResFailNotFound
//		  500: GenericResError
func (ctrl *UserController) GetUser(c *fiber.Ctx) error {
	userID := c.Params(constants.ParamUid)
	user, err := ctrl.userService.GetUser(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctrl.logger.Error("user not found", zap.Any("id", userID))
			return utils.JSONFail(c, http.StatusNotFound, constants.UserNotExist)
		}
		ctrl.logger.Error("error while get user by id", zap.Any("id", userID), zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetUser)
	}
	return utils.JSONSuccess(c, http.StatusOK, user)
}

// CreateUser registers a user
// swagger:route POST /api/v1/users Users RequestCreateUser
//
// Register a user.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  201: ResponseCreateUser
//	      400: GenericResFailBadRequest
//		  500: GenericResError
func (ctrl *UserController) CreateUser(c *fiber.Ctx) error {

	var userReq structs.ReqRegisterUser

	err := json.Unmarshal(c.Body(), &userReq)
	if err != nil {
		ctrl.logger.Error("error while unmarshal", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	err = validate.Struct(userReq)
	if err != nil {
		ctrl.logger.Error("error while validate", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	user, err := ctrl.userService.RegisterUser(models.User{FirstName: userReq.FirstName, LastName: userReq.LastName, Email: userReq.Email, Password: userReq.Password, Roles: userReq.Roles})
	if err != nil {
		ctrl.logger.Error("error while insert user", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrInsertUser)
	}

	ctrl.logger.Debug("user created", zap.Any("user", user))
	return utils.JSONSuccess(c, http.StatusCreated, user)
}
