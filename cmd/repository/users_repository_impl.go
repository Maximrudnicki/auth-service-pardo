package repository

import (
	"errors"
	"auth_service/cmd/model"

	"gorm.io/gorm"
)

type UsersRepositoryImpl struct {
	Db *gorm.DB
}

func NewUsersRepositoryImpl(Db *gorm.DB) UsersRepository {
	return &UsersRepositoryImpl{Db: Db}
}

// Save implements UsersRepository
func (u *UsersRepositoryImpl) Save(user model.Users) error {
	result := u.Db.Create(&user)
	if result.Error != nil {
		return errors.New("please use different email")
	}
	return nil
}

// Delete implements UsersRepository
func (u *UsersRepositoryImpl) Delete(usersId int) {
	var users model.Users
	result := u.Db.Where("id = ?", usersId).Delete(&users)
	if result.Error != nil {
		panic(result.Error)
	}
}

// FindAll implements UsersRepository
func (u *UsersRepositoryImpl) FindAll() []model.Users {
	var users []model.Users
	results := u.Db.Find(&users)
	if results.Error != nil {
		panic(results.Error)
	}
	return users
}

// FindById implements UsersRepository
func (u *UsersRepositoryImpl) FindById(userId int) (model.Users, error) {
	var user model.Users
	result := u.Db.Find(&user, userId)
	if result != nil {
		return user, nil
	} else {
		return user, errors.New("users is not found")
	}
}

// FindByUsername implements UsersRepository
func (u *UsersRepositoryImpl) FindByEmail(email string) (model.Users, error) {
	var users model.Users
	result := u.Db.First(&users, "email = ?", email)

	if result.Error != nil {
		return users, errors.New("invalid email or Password")
	}
	return users, nil
}