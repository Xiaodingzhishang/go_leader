package repository

import (
	"context"
	"github.com/Xiaodingzhishang/go_leader/webook/internal/domain"
	"github.com/Xiaodingzhishang/go_leader/webook/internal/repository/dao"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	ErrUserDuplicateEmail    = dao.ErrUserDuplicateEmail
	ErrUserNotFound          = dao.ErrUserNotFound
	ErrUserDuplicateNickname = dao.ErrUserDuplicateNickname
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}
func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}
func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.dao.FindByID(ctx, id)
	return domain.User{
		Id:           u.Id,
		Email:        u.Email,
		Nickname:     u.Nickname,
		Birthday:     u.Birthday,
		Introduction: u.Introduction,
		Ctime:        time.UnixMilli(u.Ctime),
	}, err

}

func (r *UserRepository) Update(ctx *gin.Context, user domain.User) error {
	return r.dao.Update(ctx, dao.User{
		Id:           user.Id,
		Nickname:     user.Nickname,
		Birthday:     user.Birthday,
		Introduction: user.Introduction,
	})
}
