package dao

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail    = errors.New("邮箱冲突")
	ErrUserDuplicateNickname = errors.New("昵称冲突")
	ErrDataTooLong           = errors.New("数据太长")
	ErrUserNotFound          = gorm.ErrRecordNotFound
)

// ErrDataNotFound 通用的数据没找到
var ErrDataNotFound = gorm.ErrRecordNotFound

const uniqueConflictsErrNo uint16 = 1062
const dataTooLongErrNo uint16 = 1406

type UserDAO struct {
	db *gorm.DB
}
type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 全部用户唯一
	Email    string `gorm:"unique;size:50;comment:邮箱"`
	Password string `gorm:"comment:密码"`
	// 往这面加
	Nickname     string `gorm:"unique;size:30;comment:昵称"`
	Birthday     string `gorm:"size:10;comment:生日"`
	Introduction string `gorm:"size:150;comment:个人介绍"`

	// 创建时间，毫秒数
	Ctime int64 `gorm:"comment:创建时间"`
	// 更新时间，毫秒数
	Utime int64 `gorm:"comment:更新时间"`
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}
func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	// 存毫秒数
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {

		if mysqlErr.Number == uniqueConflictsErrNo {
			// 邮箱冲突
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) Update(ctx *gin.Context, user User) error {
	user.Utime = time.Now().UnixMilli()
	err := dao.db.WithContext(ctx).Updates(&user).Error
	var m *mysql.MySQLError
	if errors.As(err, &m) {
		switch m.Number {
		case uniqueConflictsErrNo:
			return ErrUserDuplicateNickname
		case dataTooLongErrNo:
			return ErrDataTooLong
		}
	}
	return err
}

func (dao *UserDAO) FindByID(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Find(&u, "id = ?", id).Error
	if err == nil && u.Id == 0 {
		return u, ErrDataNotFound
	}
	return u, err
}
