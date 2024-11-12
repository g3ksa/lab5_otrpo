package graphApi

import (
	"context"
	"github.com/g3ksa/lab5_otrpo/internal/graph_api/service/model"
)

type Service struct {
	storage Storage
}

type Storage interface {
	GetAllNodes(ctx context.Context) (*[]model.Node, error)
	GetNodeWithRelations(ctx context.Context, id int) (*[]model.Relation, error)
	Insert(ctx context.Context, data model.InsertRequest) error
	Delete(ctx context.Context, id int) error
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) GetAllNodes(ctx context.Context) (*[]model.Node, error) {
	nodes, err := s.storage.GetAllNodes(ctx)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (s *Service) GetNodeWithRelations(ctx context.Context, id int) (*[]model.Relation, error) {
	relations, err := s.storage.GetNodeWithRelations(ctx, id)
	if err != nil {
		return nil, err
	}
	return relations, nil
}

func (s *Service) Insert(ctx context.Context, request model.InsertRequest) error {
	err := s.storage.Insert(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	err := s.storage.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
