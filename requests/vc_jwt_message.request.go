package requests

import (
	core "ssi-gitlab.teda.th/ssi/core"
)

type VCJWTMessage struct {
	core.BaseValidator
	Header    *VCJWTMessageHeader `json:"Header"`
	Claims    *VCJWTMessageClaims `json:"Claims"`
	Signature *string             `json:"Signature"`
}

type VCJWTMessageHeader struct {
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	Typ string `json:"typ"`
}

type VCJWTMessageClaims struct {
	Exp   int64                 `json:"exp"`
	Iat   int64                 `json:"iat"`
	Iss   string                `json:"iss"`
	Jti   string                `json:"jti"`
	Nbf   int64                 `json:"nbf"`
	Nonce string                `json:"nonce"`
	Sub   string                `json:"sub"`
	Aud   string                `json:"aud"`
	VC    *VCJWTMessageClaimsVC `json:"vc"`
}

type VCJWTMessageClaimsVC struct {
	Context           []string                      `json:"@context"`
	Type              []string                      `json:"type"`
	CredentialSubject core.Map                      `json:"credentialSubject"`
	CredentialSchema  *VCJWTMessageCredentialSchema `json:"credentialSchema"`
}

type VCJWTMessageCredentialSchema struct {
	ID   *string `json:"id"`
	Type *string `json:"type"`
}

func (r *VCJWTMessage) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsRequired(r.Header, "Header"))
	r.Must(r.IsRequired(r.Claims, "Claims"))

	r.Must(r.IsStrRequired(&r.Header.Alg, "Header.alg"))
	r.Must(r.IsStrRequired(&r.Header.Typ, "Header.typ"))
	r.Must(r.IsStrRequired(&r.Header.Kid, "Header.kid"))
	r.Must(r.IsStrRequired(&r.Claims.Jti, "Claims.jti"))

	if r.Claims.VC != nil {
		r.Must(r.IsStrRequired(&r.Claims.Sub, "Claims.sub"))
		r.Must(r.IsStrRequired(&r.Claims.Iss, "Claims.iss"))
		r.Must(r.IsRequiredArray(r.Claims.VC.Type, "Claims.vc.type"))
		r.Must(r.IsArrayMin(r.Claims.VC.Type, 1, "Claims.vc.type"))
		r.Must(r.IsRequired(r.Claims.VC.CredentialSubject, "Claims.vc.credentialSubject"))
		if r.Must(r.IsRequired(r.Claims.VC.CredentialSchema, "Claims.vc.credentialSchema")) {
			r.Must(r.IsStrRequired(r.Claims.VC.CredentialSchema.ID, "Claims.vc.credentialSchema.id"))
		}
	}

	if r.Claims.VC == nil {
		r.Must(false, &core.IValidMessage{
			Name:    "Claims.vc",
			Code:    "REQUIRED",
			Message: "The Claims.vc field is required",
		})
	}

	return r.Error()
}
