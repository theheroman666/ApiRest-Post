package user

import (
	"fmt"
	"goweb/internal/domain"
	"log"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *domain.User) error
	GetAll(filters Filters, offset, limit int) ([]domain.User, error)
	Get(id string) (*domain.User, error)
	Delete(id string) error
	Update(id string, firstname *string, lastname *string, email *string, phone *string) error
	Count(filters Filters) (int, error)
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(user *domain.User) error {
	if err := repo.db.Create(user).Error; err != nil {
		repo.log.Println(err)
		return err
	}
	repo.log.Println("user created:", user.ID)
	return nil
}

func (repo *repo) GetAll(filters Filters, offset, limit int) ([]domain.User, error) {
	var user []domain.User
	// if err := repo.db.Model(&user).Order("create_at desc").Find(&user).Error; err != nil {
	// 	return nil, err
	// }

	tx := repo.db.Model(&user)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("create_at desc").Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *repo) Get(id string) (*domain.User, error) {
	user := domain.User{Email: id}
	// consulta con uso de primary key o id
	if err := repo.db.First(&user).Error; err != nil {
		return nil, err
	}

	// consulta por campos
	// if err := repo.db.Model(&User{}).Where("Email = ?", id).First(&user).Error; err != nil {
	// 	return nil, err
	// }
	return &user, nil
}

func (repo *repo) Delete(id string) error {
	user := domain.User{ID: id}
	if err := repo.db.Delete(&user).Error; err != nil {
		return err
	}
	return nil

}

func (repo *repo) Update(id string, firstname *string, lastname *string, email *string, phone *string) error {

	values := make(map[string]interface{})
	if firstname != nil {
		values["first_name"] = *firstname
	}
	if lastname != nil {
		values["last_name"] = *lastname
	}
	if email != nil {
		values["email"] = *email
	}
	if phone != nil {
		values["phone"] = *phone
	}
	if err := repo.db.Model(&domain.User{}).Where("id = ?", id).Updates(values).Error; err != nil {
		return err
	}
	return nil
}

func (repo *repo) Count(filters Filters) (int, error) {
	var count int64
	tx := repo.db.Model(domain.User{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if strings.TrimSpace(filters.FirstnameF) != "" {
		filters.FirstnameF = fmt.Sprintf("%%%s%%", strings.ToLower(filters.FirstnameF))
		tx = tx.Where("lower(first_name) like ?", filters.FirstnameF)
	}
	if strings.TrimSpace(filters.LastnameF) != "" {
		filters.LastnameF = fmt.Sprintf("%%%s%%", strings.ToLower(filters.LastnameF))
		tx = tx.Where("lower(last_name) like ?", filters.LastnameF)
	}
	return tx
}
