package services

import "github.com/zeelrupapara/trading-api/models"

type UserService struct {
	userModel *models.UserModel
}

func NewUserService(userModel *models.UserModel) *UserService {
	return &UserService{
		userModel: userModel,
	}
}

func (userSvc *UserService) RegisterUser(user models.User) (models.User, error) {
	user, err := userSvc.userModel.InsertUser(user)
	if err != nil {
		return user, err
	}
	return user, err
}

func (userSvc *UserService) GetUser(userId string) (models.User, error) {
	user, err := userSvc.userModel.GetById(userId)
	return user, err
}

func (userSvc *UserService) Authenticate(email, password string) (models.User, error) {
	return userSvc.userModel.GetUserByEmailAndPassword(email, password)
}
