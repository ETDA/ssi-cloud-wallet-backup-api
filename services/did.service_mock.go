package services

import (
	"github.com/stretchr/testify/mock"
	"gitlab.finema.co/finema/etda/vc-wallet-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
)

type MockDIDService struct {
	mock.Mock
}

func (m *MockDIDService) Find(address string) (*models.DIDDocument, core.IError) {
	args := m.Called(address)
	return args.Get(0).(*models.DIDDocument), core.MockIError(args, 1)
}

func NewMockDIDService() *MockDIDService {
	return &MockDIDService{}
}
