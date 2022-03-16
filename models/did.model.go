package models

type DIDDocument struct {
	Context            string                          `json:"@context"`
	ID                 string                          `json:"id"`
	Controller         string                          `json:"controller"`
	VerificationMethod []DIDDocumentVerificationMethod `json:"VerificationMethod"`
}

type DIDDocumentVerificationMethod struct {
	ID           string `json:"id"`
	Controller   string `json:"controller"`
	PublicKeyPem string `json:"publicKeyPem"`
	Type         string `json:"type"`
}
