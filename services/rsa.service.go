package services

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"

	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type IRSAService interface {
	EncryptMessage(message []byte, label []byte, publicKey *rsa.PublicKey, hashingAlgorithm crypto.Hash) ([][]byte, core.IError)
	DecryptCipherText(cipherTexts [][]byte, label []byte, privateKey *rsa.PrivateKey, hashingAlgorithm crypto.Hash) ([][]byte, core.IError)
}

func NewRSAService(ctx core.IContext) IRSAService {
	return &rsaService{ctx: ctx}
}

type rsaService struct {
	ctx core.IContext
}

func (s *rsaService) EncryptMessage(message []byte, label []byte, publicKey *rsa.PublicKey, hashingAlgorithm crypto.Hash) ([][]byte, core.IError) {
	keySize := publicKey.Size()
	maxLength := keySize - 2*hashingAlgorithm.Size() - 2
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
		cipherText, err := utils.EncryptMessage(message, publicKey, &utils.MessageEncryptionOptions{
			HashingAlgorithm: hashingAlgorithm,
			Label:            label,
		})
		if err != nil {
			return nil, s.ctx.NewError(err, errmsgs.InternalServerError)
		}
		cipherTexts = append(cipherTexts, cipherText)
	}

	return cipherTexts, nil
}

func (s *rsaService) DecryptCipherText(cipherTexts [][]byte, label []byte, privateKey *rsa.PrivateKey, hashingAlgorithm crypto.Hash) ([][]byte, core.IError) {
	messages := make([][]byte, 0)
	for _, cipherText := range cipherTexts {
		h := hashingAlgorithm.New()
		message, err := rsa.DecryptOAEP(h, rand.Reader, privateKey, cipherText, label)
		if err != nil {
			return nil, s.ctx.NewError(err, errmsgs.InternalServerError)
		}
		messages = append(messages, bytes.Trim(message, "\x00"))
	}

	return messages, nil
}
