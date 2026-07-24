package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
	"github.com/99katedegree/bunkasairpg2/api/internal/domain/repository"
)

type AuthUsecase struct {
	userRepo  repository.UserRepository
	adminRepo repository.AdminRepository
	jwtSecret []byte
}

func NewAuthUsecase(userRepo repository.UserRepository, adminRepo repository.AdminRepository, jwtSecret string) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo, adminRepo: adminRepo, jwtSecret: []byte(jwtSecret)}
}

// Login はユーザー UUID を確認し JWT を返す
func (u *AuthUsecase) Login(ctx context.Context, userID uuid.UUID) (string, error) {
	_, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", entity.ErrNotFound
	}
	return u.generateToken("user", userID.String())
}

// AdminLogin は email/password で管理者認証し JWT を返す
func (u *AuthUsecase) AdminLogin(ctx context.Context, email, password string) (string, error) {
	admin, err := u.adminRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", entity.ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return "", entity.ErrUnauthorized
	}
	return u.generateToken("admin", fmt.Sprintf("%d", admin.ID))
}

func (u *AuthUsecase) generateToken(role, sub string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  sub,
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour * 30).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(u.jwtSecret)
}
