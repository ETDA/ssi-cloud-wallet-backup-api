package services

import (
	"crypto"
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.finema.co/finema/etda/vc-wallet-api/helpers"
	core "ssi-gitlab.teda.th/ssi/core"
)

type RSAWithHSMServiceTestSuite struct {
	ctx *core.ContextMock
	s   IRSAService
	suite.Suite
}

func TestRSAWithHSMServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RSAWithHSMServiceTestSuite))
}

func (ts *RSAWithHSMServiceTestSuite) SetupTest() {
	ts.ctx = core.NewMockContext()
	ts.s = NewRSAWithHSMService(ts.ctx)
}

func (ts *RSAWithHSMServiceTestSuite) TestEncryptAndDecryptMessageSuccess() {
	originalMessage := "foobar"
	cipher, err := ts.s.EncryptMessage([]byte(originalMessage), nil, nil, crypto.SHA256)
	ts.NoError(err)
	ts.NotEmpty(cipher)
	message, err := ts.s.DecryptCipherText(cipher, nil, nil, crypto.SHA256)
	ts.NoError(err)
	ts.NotEmpty(message)

	expectedOriginalMessage := helpers.ByteArraySeriesToString(message)
	ts.Equal(expectedOriginalMessage, originalMessage)
}
