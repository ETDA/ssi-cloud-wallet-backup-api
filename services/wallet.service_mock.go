package services

import (
	"github.com/stretchr/testify/mock"
	"gitlab.finema.co/finema/etda/vc-wallet-api/models"
	"gitlab.finema.co/finema/etda/vc-wallet-api/views"
	core "ssi-gitlab.teda.th/ssi/core"
)

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) CreateWallet(payload *WalletCreatePayload) (*models.Wallet, core.IError) {
	args := m.Called(payload)
	return args.Get(0).(*models.Wallet), core.MockIError(args, 1)
}

func (m *MockWalletService) FindWallet(id string) (*models.Wallet, core.IError) {
	args := m.Called(id)
	return args.Get(0).(*models.Wallet), core.MockIError(args, 1)
}

func (m *MockWalletService) AddVC(did string, payload *WalletAddVCPayload) (*models.VC, core.IError) {
	args := m.Called(did, payload)
	return args.Get(0).(*models.VC), core.MockIError(args, 1)
}

func (m *MockWalletService) FindVC(did string, cid string, options *WalletFindVCOptions) (*models.VC, core.IError) {
	args := m.Called(did, cid, options)
	return args.Get(0).(*models.VC), core.MockIError(args, 1)
}

func (m *MockWalletService) VCPagination(did string, pageOptions *core.PageOptions, options *VCPaginationOptions) ([]models.VC, *core.PageResponse, core.IError) {
	args := m.Called(did, pageOptions, options)
	return args.Get(0).([]models.VC), args.Get(1).(*core.PageResponse), core.MockIError(args, 2)
}

func (m *MockWalletService) Summary(did string) (*views.WalletSummary, core.IError) {
	args := m.Called(did)
	return args.Get(0).(*views.WalletSummary), core.MockIError(args, 1)
}

func NewMockWalletService() *MockWalletService {
	return &MockWalletService{}
}
