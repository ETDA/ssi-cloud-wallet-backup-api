package services

import (
	"crypto"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gitlab.finema.co/finema/etda/vc-wallet-api/consts"
	"gitlab.finema.co/finema/etda/vc-wallet-api/emsgs"
	"gitlab.finema.co/finema/etda/vc-wallet-api/helpers"
	"gitlab.finema.co/finema/etda/vc-wallet-api/models"
	"gitlab.finema.co/finema/etda/vc-wallet-api/requests"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"gorm.io/gorm"
)

type WalletCreatePayload struct {
	DIDAddress string `json:"did_address"`
}

type WalletAddVCPayloadVC struct {
	Context           []string `json:"@context"`
	Type              []string `json:"type"`
	CredentialSubject core.Map `json:"credentialSubject"`
	Proof             core.Map `json:"proof"`
}

type WalletAddVCClaim struct {
	VC *WalletAddVCPayloadVC `json:"vc"`
	jwt.StandardClaims
}

type WalletAddVCPayload struct {
	JWT      string `json:"jwt"`
	Operator string `json:"operator"`
}

type WalletFindVCOptions struct {
}

type VCPaginationOptions struct {
	SchemaType   string
	IssuanceDate string
	Holder       string
	Issuer       string
}

type VPPaginationOptions struct {
	SchemaType   string
	IssuanceDate string
	Audience     string
	Issuer       string
}

type IWalletService interface {
	CheckWallet(did string) core.IError
	CreateWallet(payload *WalletCreatePayload) (*models.Wallet, core.IError)
	FindWallet(id string) (*models.Wallet, core.IError)
	AddVC(did string, payload *WalletAddVCPayload) core.IError
	CheckVC(cid string) core.IError
	VCPagination(did string, pageOptions *core.PageOptions, options *VCPaginationOptions) ([]models.VC, *core.PageResponse, core.IError)
}

type walletService struct {
	ctx        core.IContext
	rsaService IRSAService
}

func NewWalletService(ctx core.IContext) IWalletService {
	return &walletService{
		ctx:        ctx,
		rsaService: NewRSAWithHSMService(ctx),
	}
}

func (s walletService) CreateWallet(payload *WalletCreatePayload) (*models.Wallet, core.IError) {
	w, ierr := s.FindWallet(payload.DIDAddress)
	if ierr != nil && !(emsgs.WalletNotFound.GetCode() == ierr.GetCode()) {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	if w != nil {
		return nil, s.ctx.NewError(emsgs.DuplicatedWallet, emsgs.DuplicatedWallet)
	}

	err := s.ctx.DB().Create(&models.Wallet{
		ID:         payload.DIDAddress,
		DIDAddress: payload.DIDAddress,
		CreatedAt:  utils.GetCurrentDateTime(),
	}).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return s.FindWallet(payload.DIDAddress)
}

func (s walletService) CheckWallet(did string) core.IError {
	_, ierr := s.FindWallet(did)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}
	return nil
}

func (s walletService) CheckVC(cid string) core.IError {
	item := &models.VC{}
	err := s.ctx.DB().Where("cid = ?", cid).First(item).Error
	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}
	return nil
}

func (s walletService) FindWallet(id string) (*models.Wallet, core.IError) {
	item := &models.Wallet{}
	err := s.ctx.DB().Where("id = ?", id).First(item).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, emsgs.WalletNotFound)
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return item, nil
}

func (s walletService) AddVC(walletAddress string, payload *WalletAddVCPayload) core.IError {
	_, ierr := s.FindWallet(walletAddress)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	vc, ierr := s.decodeVCJWT(payload.JWT)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	if payload.Operator != vc.Claims.Iss && payload.Operator != vc.Claims.Sub {
		return s.ctx.NewError(emsgs.IssuerMismatched, emsgs.IssuerMismatched, walletAddress, vc.Claims.Iss, vc.Claims.Sub)
	}

	issuanceDate := time.Unix(vc.Claims.Iat, 0)

	privateKeyBytes, err := os.ReadFile(s.ctx.ENV().String(consts.ENVRSAPrivateKeyFile))
	if err != nil {
		return s.ctx.NewError(err, errmsgs.InternalServerError)
	}

	privateKey, err := utils.LoadRSAPrivateKey(string(privateKeyBytes))
	if err != nil {
		return s.ctx.NewError(err, errmsgs.InternalServerError)
	}

	encryptedSchemaType, ierr := s.rsaService.EncryptMessage(utils.StringToBytes(utils.GetString(vc.Claims.VC.CredentialSchema.Type)), nil, &privateKey.PublicKey, crypto.SHA256)
	if ierr != nil {
		return s.ctx.NewError(ierr, errmsgs.InternalServerError)
	}

	encryptedJWT, ierr := s.rsaService.EncryptMessage(utils.StringToBytes(payload.JWT), nil, &privateKey.PublicKey, crypto.SHA256)
	if err != nil {
		return s.ctx.NewError(ierr, errmsgs.InternalServerError)
	}

	err = s.ctx.DB().Create(&models.VC{
		ID:           vc.Claims.Jti,
		CID:          vc.Claims.Jti,
		SchemaType:   helpers.ByteArraySeriesToBase64StringJoined(encryptedSchemaType, "."),
		IssuanceDate: &issuanceDate,
		Issuer:       vc.Claims.Iss,
		Holder:       vc.Claims.Sub,
		JWT:          helpers.ByteArraySeriesToBase64StringJoined(encryptedJWT, "."),
	}).Error

	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}

	return nil
}

func (s walletService) VCPagination(did string, pageOptions *core.PageOptions, options *VCPaginationOptions) ([]models.VC, *core.PageResponse, core.IError) {
	if pageOptions == nil {
		iHttpContext, ok := s.ctx.(core.IHTTPContext)
		if !ok {
			return nil, nil, s.ctx.NewError(errmsgs.InternalServerError, errmsgs.InternalServerError)
		}
		pageOptions = iHttpContext.GetPageOptions()
	}
	if options == nil {
		options = &VCPaginationOptions{}
	}
	if options.Holder == "" {
		options.Holder = did
	}
	if options.Issuer == "" {
		options.Issuer = did
	}

	_, ierr := s.FindWallet(did)
	if ierr != nil && !(emsgs.WalletNotFound.GetCode() == ierr.GetCode()) {
		return nil, nil, s.ctx.NewError(emsgs.WalletNotFound, emsgs.WalletNotFound)
	}
	if ierr != nil {
		return nil, nil, s.ctx.NewError(ierr, ierr)
	}

	items := make([]models.VC, 0)
	keys := append([]string{"id", "cid"}, consts.VCSearchQueryKeywords...)

	db := s.ctx.DB()
	db = core.SetSearchSimple(db, pageOptions.Q, keys)

	keyWrapper, ierr := s.getKeyConditionFromVCPaginationOptions(did, options)
	if ierr != nil {
		return nil, nil, s.ctx.NewError(ierr, ierr)
	}

	core.SetSearch(db, keyWrapper)
	s.setDateFromString(db, options.IssuanceDate)

	pageRes, err := core.Paginate(db, &items, pageOptions)
	if err != nil {
		return nil, nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	privateKeyBytes, err := os.ReadFile(s.ctx.ENV().String(consts.ENVRSAPrivateKeyFile))
	if err != nil {
		return nil, nil, s.ctx.NewError(err, errmsgs.InternalServerError)
	}

	privateKey, err := utils.LoadRSAPrivateKey(string(privateKeyBytes))
	if err != nil {
		return nil, nil, s.ctx.NewError(err, errmsgs.InternalServerError)
	}

	for i, item := range items {
		encryptedSchemaType, err := helpers.Base64StringJoinedToByteArraySeries(item.SchemaType, ".")
		if err != nil {
			return nil, nil, s.ctx.NewError(err, errmsgs.InternalServerError)
		}

		decryptedSchemaType, ierr := s.rsaService.DecryptCipherText(encryptedSchemaType, nil, privateKey, crypto.SHA256)
		if ierr != nil {
			return nil, nil, s.ctx.NewError(ierr, ierr)
		}

		encryptedJWT, err := helpers.Base64StringJoinedToByteArraySeries(item.JWT, ".")
		if err != nil {
			return nil, nil, s.ctx.NewError(err, errmsgs.InternalServerError)
		}

		decryptedJWT, err := s.rsaService.DecryptCipherText(encryptedJWT, nil, privateKey, crypto.SHA256)
		if ierr != nil {
			return nil, nil, s.ctx.NewError(ierr, ierr)
		}

		items[i].SchemaType = helpers.ByteArraySeriesToString(decryptedSchemaType)
		items[i].JWT = helpers.ByteArraySeriesToString(decryptedJWT)
	}

	return items, pageRes, nil
}

func (s walletService) decodeVCJWT(jwtToken string) (*requests.VCJWTMessage, core.IError) {
	token, _ := helpers.JWTVCDecodingT(jwtToken, []byte(""))
	if token == nil || token.Header == nil || token.Claims == nil || token.Signature == "" {
		return nil, s.ctx.NewError(emsgs.JWTInValid, emsgs.JWTInValid)
	}

	msgPayload := &requests.VCJWTMessage{}
	_ = utils.MapToStruct(token, &msgPayload)
	if ierr := msgPayload.Valid(s.ctx); ierr != nil {
		utils.LogStruct(ierr.JSON())
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return msgPayload, nil
}

func (s walletService) getKeyConditionFromVCPaginationOptions(did string, options *VCPaginationOptions) (*core.KeywordConditionWrapper, core.IError) {
	if options == nil {
		return &core.KeywordConditionWrapper{}, nil
	}

	keywords := make([]core.KeywordOptions, 0)
	if options.SchemaType != "" {
		keywords = append(keywords, core.KeywordOptions{
			Type:  core.MustMatch,
			Key:   "type",
			Value: options.SchemaType,
		})
	}

	if options.Holder != "" {
		if did != options.Holder {
			return nil, emsgs.VCHolderInvalid
		}

		keywords = append(keywords, core.KeywordOptions{
			Type:  core.MustMatch,
			Key:   "holder",
			Value: options.Holder,
		})
	}

	if options.Issuer != "" {
		if did != options.Issuer {
			return nil, emsgs.VCIssuerInvalid
		}

		keywords = append(keywords, core.KeywordOptions{
			Type:  core.MustMatch,
			Key:   "issuer",
			Value: options.Issuer,
		})
	}

	condition := &core.KeywordConditionWrapper{
		Condition:      core.Or,
		KeywordOptions: keywords,
	}
	if len(condition.KeywordOptions) == 0 {
		return &core.KeywordConditionWrapper{}, nil
	}

	return condition, nil
}

func (s walletService) setDateFromString(db *gorm.DB, date string) {
	if date != "" {
		if valid, _ := regexp.MatchString(`^\w{4}\-\w{2}\-\w{2}$`, date); valid {
			db.Where(fmt.Sprintf("DATE(issuance_date) = DATE('%s')", date))
		}
	}
}
