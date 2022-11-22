package core

import (
	"context"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAllUsers(c context.Context) ([]User, error)
	FindUserByID(c context.Context, userID uint) (*User, error)
	FindUserByUserName(c context.Context, userName string) (*User, error)
	DeleteUserByID(c context.Context, userID uint) error
	SaveUser(c context.Context, user *User) error
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryMysql(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) FindAllUsers(c context.Context) ([]User, error) {
	var result []User
	err := r.db.WithContext(c).Model(&User{}).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *userRepositoryImpl) FindUserByID(c context.Context, userID uint) (*User, error) {
	result := &User{}
	err := r.db.WithContext(c).Model(&User{}).Where("id = ?", userID).First(result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *userRepositoryImpl) FindUserByUserName(c context.Context, userName string) (*User, error) {
	result := &User{}
	err := r.db.WithContext(c).Model(&User{}).Where("username = ?", userName).First(result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *userRepositoryImpl) DeleteUserByID(c context.Context, userID uint) error {
	return r.db.WithContext(c).Delete(&User{ID: userID}).Error
}

func (r *userRepositoryImpl) SaveUser(c context.Context, user *User) error {
	return r.db.WithContext(c).Save(user).Error
}
