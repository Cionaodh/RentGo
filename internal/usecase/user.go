package usecase

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo"
	"EasyRentGo/internal/repo/repoerrors"
	repotype "EasyRentGo/internal/repo/repotypes"
	"EasyRentGo/pkg/logger"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo repo.User
	l    logger.Interface
}

func NewUserUsecase(r repo.User, l logger.Interface) *UserUsecase {
	return &UserUsecase{
		repo: r,
		l:    l,
	}
}

func (uc *UserUsecase) Register(ctx context.Context, input RegisterUserInput) (entity.User, error) {
	// Хэшируем пароль перед сохранением
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.l.Error("UserUsecase - Register - bcrypt: %v", err)
		return entity.User{}, fmt.Errorf("failed to hash password")
	}

	// Генерируем юзернейм (например, user_1694200000)
	username := fmt.Sprintf("user_%d", time.Now().Unix())

	user, err := uc.repo.Create(ctx, repotype.CreateUserInput{
		Email:        input.Email,
		Username:     username,
		PasswordHash: string(hashedBytes),
	})
	if err != nil {
		if errors.Is(err, repoerrors.ErrAlreadyExists) {
			return entity.User{}, ErrUserAlreadyExists
		}
		uc.l.Error("UserUsecase - Register - repo.Create: %v", err)
		return entity.User{}, fmt.Errorf("failed to create user")
	}

	return user, nil
}

func (uc *UserUsecase) Login(ctx context.Context, input LoginUserInput) (entity.User, error) {
	// Ищем юзера по email
	user, err := uc.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return entity.User{}, ErrInvalidCredentials // Скрываем, что именно не так (email или пароль)
		}
		uc.l.Error("UserUsecase - Login - repo.GetByEmail: %v", err)
		return entity.User{}, fmt.Errorf("failed to get user")
	}

	// Сравниваем хэш из базы с присланным паролем
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return entity.User{}, ErrInvalidCredentials
	}

	return user, nil
}

func (uc *UserUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (entity.UserProfile, error) {

	profile, err := uc.repo.GetProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return entity.UserProfile{}, ErrUserNotFound
		}
		uc.l.Error("UserUsecase - GetProfile: %v", err)
		return entity.UserProfile{}, fmt.Errorf("failed to get profile")
	}

	return profile, nil
}
