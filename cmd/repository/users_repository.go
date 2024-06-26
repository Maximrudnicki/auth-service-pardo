package repository

import "auth_service/cmd/model"

type UsersRepository interface {
	Save(users model.Users) error
	Delete(usersId int)
	FindById(usersId int) (model.Users, error)
	FindAll() []model.Users
	FindByEmail(email string) (model.Users, error)
}