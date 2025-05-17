package user

import (
	"context"
	"github.com/sidgim/finance_domain/domain"
	"gorm.io/gorm"
	"log"
)

type (
	Filters struct {
		FirstName string
		LastName  string
	}

	Service interface {
		Create(ctx context.Context, req CreateRequest) (*domain.User, error)
		Get(ctx context.Context, id string) (*domain.User, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error)
		Delete(ctx context.Context, id string) error
		UpdateContact(ctx context.Context, id string, req UpdateRequest) (*domain.User, error)
		Count(ctx context.Context, filters Filters) (int, error)
	}
	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		repo: repo,
		log:  log,
	}
}

func (s *service) Create(ctx context.Context, req CreateRequest) (*domain.User, error) {
	s.log.Println("Creating user:", req)
	user := domain.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
	}
	if err := s.repo.Create(ctx, &user); err != nil {
		s.log.Println("Error creating user:", err)
		return nil, err
	}
	return &user, nil
}

func (s *service) Get(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		s.log.Println("Error getting user:", err)
		return nil, err
	}
	return user, nil
}

func (s *service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error) {
	users, err := s.repo.GetAll(ctx, filters, offset, limit)
	if err != nil {
		s.log.Println("Error getting all users:", err)
		return nil, err
	}
	return users, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.log.Println("Error deleting user:", err)
		return err
	}
	s.log.Println("domain.User deleted:", id)
	return nil
}

func (s *service) UpdateContact(ctx context.Context, id string, req UpdateRequest) (*domain.User, error) {

	// 2.2 Traer el user existente pa’ validar que exista
	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		s.log.Printf("Error fetching user %s: %v", id, err)
		return nil, err
	}
	if existing == nil {
		s.log.Printf("domain.User %s no encontrado", id)
		return nil, gorm.ErrRecordNotFound
	}

	// 2.3 Actualizar los campos en el objeto
	existing.Email = req.Email
	existing.Phone = req.Phone

	// 2.4 LLamar al repo “partial update” (solo email y phone)
	if err := s.repo.UpdateContact(ctx, id, req); err != nil {
		s.log.Printf("Error updating contact for %s: %v", id, err)
		return nil, err
	}

	s.log.Printf("domain.User %s contact updated: email=%s phone=%s", id, req.Email, req.Phone)
	return existing, nil
}
func (s *service) Count(ctx context.Context, filters Filters) (int, error) {
	return s.repo.Count(ctx, filters)
}
