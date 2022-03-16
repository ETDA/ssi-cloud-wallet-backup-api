package services

import (
	"bytes"
	"crypto"
	"crypto/rsa"

	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
)

func NewRSAWithHSMService(ctx core.IContext) IRSAService {
	return &rsaWithHSMService{
		ctx:        ctx,
		hsmService: NewHSMService(ctx),
	}
}

type rsaWithHSMService struct {
	ctx        core.IContext
	hsmService IHSMService
}

func (s *rsaWithHSMService) EncryptMessage(message []byte, label []byte, publicKey *rsa.PublicKey, hashingAlgorithm crypto.Hash) ([][]byte, core.IError) {
	maxLength := 190 // statically set for RSA 2048 with SHA-256
	messages := make([][]byte, 0)

	i := 0
	for {
		newMassage := make([]byte, maxLength)
		if ((i + 1) * maxLength) > len(message) {
			copy(newMassage, message[i*maxLength:])
			messages = append(messages, newMassage)
			break
		}
		copy(newMassage, message[i*maxLength:((i+1)*maxLength)])
		messages = append(messages, newMassage)
		i++
	}

	cipherTexts := make([][]byte, 0)
	for _, message := range messages {
		cipherText, err := s.hsmService.Encrypt(message)
		if err != nil {
			return nil, s.ctx.NewError(err, errmsgs.InternalServerError)
		}
		cipherTexts = append(cipherTexts, cipherText)
	}

	return cipherTexts, nil
}

func (s *rsaWithHSMService) DecryptCipherText(cipherTexts [][]byte, label []byte, privateKey *rsa.PrivateKey, hashingAlgorithm crypto.Hash) ([][]byte, core.IError) {
	messages := make([][]byte, 0)
	for _, cipherText := range cipherTexts {
		message, err := s.hsmService.Decrypt(cipherText)
		if err != nil {
			return nil, s.ctx.NewError(err, errmsgs.InternalServerError)
		}
		messages = append(messages, bytes.Trim(message, "\x00"))
	}

	return messages, nil
}
