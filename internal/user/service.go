package user

import (
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
		Create(req CreateRequest) (*domain.User, error)
		Get(id string) (*domain.User, error)
		GetAll(filters Filters, offset, limit int) ([]domain.User, error)
		Delete(id string) error
		UpdateContact(id string, req UpdateRequest) (*domain.User, error)
		Count(filters Filters) (int, error)
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

func (s *service) Create(req CreateRequest) (*domain.User, error) {
	s.log.Println("Creating user:", req)
	user := domain.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
	}
	if err := s.repo.Create(&user); err != nil {
		s.log.Println("Error creating user:", err)
		return nil, err
	}
	return &user, nil
}

func (s *service) Get(id string) (*domain.User, error) {
	user, err := s.repo.Get(id)
	if err != nil {
		s.log.Println("Error getting user:", err)
		return nil, err
	}
	return user, nil
}

func (s *service) GetAll(filters Filters, offset, limit int) ([]domain.User, error) {
	users, err := s.repo.GetAll(filters, offset, limit)
	if err != nil {
		s.log.Println("Error getting all users:", err)
		return nil, err
	}
	return users, nil
}

func (s *service) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		s.log.Println("Error deleting user:", err)
		return err
	}
	s.log.Println("domain.User deleted:", id)
	return nil
}

func (s *service) UpdateContact(id string, req UpdateRequest) (*domain.User, error) {

	// 2.2 Traer el user existente pa’ validar que exista
	existing, err := s.repo.Get(id)
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
	if err := s.repo.UpdateContact(id, req); err != nil {
		s.log.Printf("Error updating contact for %s: %v", id, err)
		return nil, err
	}

	s.log.Printf("domain.User %s contact updated: email=%s phone=%s", id, req.Email, req.Phone)
	return existing, nil
}
func (s *service) Count(filters Filters) (int, error) {
	return s.repo.Count(filters)
}
