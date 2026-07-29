package example

import (
	"context"
	"errors"

	"go-api-boilerplate/internal/shared"
)

type Service struct {
	repo *Repository
}

var ErrExampleNotFound = errors.New("example not found")

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Status      string `json:"status" validate:"omitempty,oneof=active inactive"`
}

type UpdateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Status      string `json:"status" validate:"omitempty,oneof=active inactive"`
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (Response, error) {
	status := req.Status
	if status == "" {
		status = "active"
	}

	item, err := s.repo.Create(ctx, Example{
		Name:        req.Name,
		Description: req.Description,
		Status:      status,
	})
	if err != nil {
		return Response{}, err
	}

	return toResponse(item), nil
}

func (s *Service) GetAll(ctx context.Context, search string, limit, offset int32) ([]Response, int64, error) {
	items, total, err := s.repo.FindAll(ctx, shared.NormalizeSearch(search), limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]Response, len(items))
	for i, item := range items {
		result[i] = toResponse(item)
	}
	return result, total, nil
}

func (s *Service) GetOne(ctx context.Context, id string) (Response, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if shared.IsNotFound(err) {
			return Response{}, ErrExampleNotFound
		}
		return Response{}, err
	}
	return toResponse(item), nil
}

func (s *Service) Update(ctx context.Context, id string, req UpdateRequest) (Response, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if shared.IsNotFound(err) {
			return Response{}, ErrExampleNotFound
		}
		return Response{}, err
	}

	item.Name = req.Name
	item.Description = req.Description
	if req.Status != "" {
		item.Status = req.Status
	}

	item, err = s.repo.Update(ctx, item)
	if err != nil {
		return Response{}, err
	}
	return toResponse(item), nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		if shared.IsNotFound(err) {
			return ErrExampleNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}
