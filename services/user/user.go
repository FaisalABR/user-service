package services

import (
	"context"
	"strings"
	"time"
	errWrap "user-service/common/error"
	"user-service/config"
	"user-service/constants"
	errConstant "user-service/constants/error"
	"user-service/domain/dto"
	"user-service/domain/models"
	"user-service/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repositories.IRepositoryRegistry
}

type IUserService interface {
	Login(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(context.Context, *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Update(context.Context, *dto.UpdateRequest, string) (*dto.UserResponse, error)
	GetUserLogin(context.Context) (*dto.UserResponse, error)
	GetUserByUUID(context.Context, string) (*dto.UserResponse, error)
}

type Claims struct {
	User *dto.UserResponse
	jwt.RegisteredClaims
}

func NewUserService(repository repositories.IRepositoryRegistry) IUserService {
	return &UserService{
		repository: repository,
	}
}

func (u *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.repository.GetUser().FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	// encrypt password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrPasswordIncorrect)
	}

	// create expiration time
	data := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        strings.ToLower(user.Role.Code),
	}
	expirationTime := time.Now().Add(time.Duration(config.Config.JwtExpirationTime) * time.Minute).Unix()

	// create claims
	claims := &Claims{
		User: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
		},
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// get token string
	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))
	if err != nil {
		return nil, err
	}

	// return response
	response := &dto.LoginResponse{
		User:  *data,
		Token: tokenString,
	}

	return response, nil

}

func (u *UserService) isUsernameExist(ctx context.Context, username string) bool {
	user, err := u.repository.GetUser().FindByUsername(ctx, username)
	if err != nil {
		return false
	}

	if user != nil {
		return true
	}

	return false
}

func (u *UserService) isEmailExist(ctx context.Context, email string) bool {
	user, err := u.repository.GetUser().FindByEmail(ctx, email)
	if err != nil {
		return false
	}

	if user != nil {
		return true
	}

	return false
}

func (u *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if u.isUsernameExist(ctx, req.Username) {
		return nil, errWrap.WrapError(errConstant.ErrUsernameExist)
	}

	if u.isEmailExist(ctx, req.Email) {
		return nil, errWrap.WrapError(errConstant.ErrEmailExist)
	}

	if req.Password != req.ConfirmPassword {
		return nil, errWrap.WrapError(errConstant.ErrPasswordDoesNotMatch)
	}

	data := &dto.RegisterRequest{
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		Password:    string(hashPassword),
		PhoneNumber: req.PhoneNumber,
		RoleID:      constants.User,
	}

	user, err := u.repository.GetUser().Register(ctx, data)
	if err != nil {
		return nil, err
	}

	response := &dto.RegisterResponse{
		User: dto.UserResponse{
			UUID:        user.UUID,
			Name:        user.Name,
			Username:    user.Username,
			Email:       user.Email,
			Role:        user.Role.Code,
			PhoneNumber: user.PhoneNumber,
		},
	}

	return response, nil
}

func (u *UserService) Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*dto.UserResponse, error) {
	var (
		password         string
		hashedPassword   []byte
		user, userResult *models.User
		err              error
		data             dto.UserResponse
	)

	user, err = u.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	// check if user already input others username but already taken
	isUsernameExist := u.isUsernameExist(ctx, req.Username)
	if isUsernameExist && user.Username != req.Username {
		return nil, errConstant.ErrUsernameExist
	}

	// check if user already input others emails but already taken
	isEmailExist := u.isEmailExist(ctx, req.Email)
	if isEmailExist && user.Email != req.Email {
		return nil, errConstant.ErrEmailExist
	}

	if req.Password != nil {
		if *req.Password != *req.ConfirmPassword {
			return nil, errConstant.ErrPasswordDoesNotMatch
		}

		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		password = string(hashedPassword)
	}

	userResult, err = u.repository.GetUser().Update(ctx,
		&dto.UpdateRequest{
			Name:        req.Name,
			Username:    req.Username,
			Password:    &password,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
		}, uuid)

	if err != nil {
		return nil, err
	}

	data = dto.UserResponse{
		UUID:        userResult.UUID,
		Name:        userResult.Username,
		Username:    userResult.Username,
		Email:       userResult.Email,
		Role:        userResult.Role.Code,
		PhoneNumber: userResult.PhoneNumber,
	}

	return &data, nil
}

func (u *UserService) GetUserLogin(ctx context.Context) (*dto.UserResponse, error) {
	var (
		userLogin = ctx.Value(constants.UserLogin).(*dto.UserResponse)
		data      dto.UserResponse
	)

	data = dto.UserResponse{
		UUID:        userLogin.UUID,
		Name:        userLogin.Name,
		Username:    userLogin.Username,
		Email:       userLogin.Email,
		PhoneNumber: userLogin.PhoneNumber,
		Role:        userLogin.Role,
	}

	return &data, nil
}

func (u *UserService) GetUserByUUID(ctx context.Context, uuid string) (*dto.UserResponse, error) {
	user, err := u.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	data := dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}

	return &data, nil
}
