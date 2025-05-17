package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/sidgim/finance_domain/domain"
	"gorm.io/gorm"
	"log"
	"strings"
)

type Repository interface {
	Create(ctx context.Context, user *domain.User) error
	Get(ctx context.Context, id string) (*domain.User, error)
	GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error)
	Delete(ctx context.Context, id string) error
	UpdateContact(ctx context.Context, id string, req UpdateRequest) error
	Count(ctx context.Context, filters Filters) (int, error)
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepository(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		db:  db,
		log: log}
}

func (r *repo) Create(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	r.log.Println("User created:", user)
	return nil
}

func (r *repo) Get(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Println("User not found:", id)
			return nil, nil
		}
		return nil, err
	}
	r.log.Println("User retrieved:", user)
	return &user, nil
}

func (r *repo) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error) {
	var users []domain.User
	db := r.db.Model(&domain.User{})
	db = applyFilters(db, filters)
	db = db.Offset(offset).Limit(limit)
	if err := db.Order("created_at desc").WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}

	r.log.Println("All users retrieved")
	return users, nil
}

func (r *repo) Delete(ctx context.Context, id string) error {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Println("User not found for deletion:", id)
			return nil
		}
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&user).Error; err != nil {
		return err
	}
	r.log.Println("User deleted:", id)
	return nil
}

func (r *repo) UpdateContact(ctx context.Context, id string, req UpdateRequest) error {
	changes := map[string]interface{}{
		"email": req.Email,
		"phone": req.Phone,
	}
	res := r.db.
		Model(&domain.User{}).
		Where("id = ?", id).WithContext(ctx).
		Updates(changes)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func applyFilters(db *gorm.DB, filters Filters) *gorm.DB {
	if filters.FirstName != "" {
		filters.FirstName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.FirstName))
		db = db.Where("lower(first_name) LIKE ?", filters.FirstName)
	}

	if filters.LastName != "" {
		filters.LastName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.LastName))
		db = db.Where("lower(last_name) LIKE ?", filters.LastName)
	}
	return db
}

func (r *repo) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	db := r.db.Model(&domain.User{})
	db = applyFilters(db, filters)
	if err := db.WithContext(ctx).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
