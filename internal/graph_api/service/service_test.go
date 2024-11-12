package graphApi

import (
	"context"
	"testing"

	"github.com/g3ksa/lab5_otrpo/internal/graph_api/service/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetAllNodes(ctx context.Context) (*[]model.Node, error) {
	args := m.Called(ctx)
	return args.Get(0).(*[]model.Node), args.Error(1)
}

func (m *MockStorage) GetNodeWithRelations(ctx context.Context, id int) (*[]model.Relation, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*[]model.Relation), args.Error(1)
}

func (m *MockStorage) Insert(ctx context.Context, request model.InsertRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

func (m *MockStorage) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGetAllNodes(t *testing.T) {
	mockStorage := new(MockStorage)
	svc := NewService(mockStorage)

	expectedNodes := &[]model.Node{{ID: 1}, {ID: 2}}
	mockStorage.On("GetAllNodes", mock.Anything).Return(expectedNodes, nil)

	nodes, err := svc.GetAllNodes(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, nodes)
	assert.Equal(t, expectedNodes, nodes)
	mockStorage.AssertExpectations(t)
}

func TestGetNodeWithRelations(t *testing.T) {
	mockStorage := new(MockStorage)
	svc := NewService(mockStorage)

	expectedRelations := &[]model.Relation{{Node: model.Node{ID: 1}, RelType: "FOLLOWS", EndNode: model.Node{ID: 2}}}
	mockStorage.On("GetNodeWithRelations", mock.Anything, 1).Return(expectedRelations, nil)

	relations, err := svc.GetNodeWithRelations(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedRelations, relations)
	mockStorage.AssertExpectations(t)
}

func TestInsert(t *testing.T) {
	mockStorage := new(MockStorage)
	svc := NewService(mockStorage)

	request := model.InsertRequest{Node: model.Node{ID: 1}}
	mockStorage.On("Insert", mock.Anything, request).Return(nil)

	err := svc.Insert(context.Background(), request)

	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	mockStorage := new(MockStorage)
	svc := NewService(mockStorage)

	mockStorage.On("Delete", mock.Anything, 1).Return(nil)

	err := svc.Delete(context.Background(), 1)

	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}
