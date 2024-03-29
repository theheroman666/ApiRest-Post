package user

import (
	"goweb/internal/domain"
	"log"
)

type (
	Filters struct {
		FirstnameF string
		LastnameF  string
		domain.User
	}
	Service interface {
		Create(firstName, lastName, email, phone string) (*domain.User, error)
		Get(id string) (*domain.User, error)
		GetAll(filters Filters, offset, limit int) ([]domain.User, error)
		Delete(id string) error
		Update(id string, firstname *string, lastname *string, email *string, phone *string) error
		Count(filters Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (serv service) Create(firstName, lastName, email, phone string) (*domain.User, error) {
	serv.log.Println("Create user service")
	user := domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}
	if err := serv.repo.Create(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (serv service) GetAll(filters Filters, offset, limit int) ([]domain.User, error) {

	users, err := serv.repo.GetAll(filters, offset, limit)

	if err != nil {
		return nil, err
	}
	return users, nil
}

func (serv service) Get(id string) (*domain.User, error) {
	user, err := serv.repo.Get(id)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (serv service) Delete(id string) error {
	return serv.repo.Delete(id)
}

func (serv service) Update(id string, firstname *string, lastname *string, email *string, phone *string) error {
	return serv.repo.Update(id, firstname, lastname, email, phone)
}

func (serv service) Count(filters Filters) (int, error) {
	return serv.repo.Count(filters)
}
