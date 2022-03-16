package helpers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

var privateKey = "SINGH"

type testClaimsStruct struct {
	VC        JWTVC  `json:"vc"`
	Nonce     string `json:"nonce"`
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	Id        string `json:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
	jwt.StandardClaims
}

func TestJWTVCEncodingHeaderI_Expect_Correct_Value(t *testing.T) {
	expectToken := "eyJhbGciOiJTRUNQMjU2UjEiLCJraWQiOiJkaWQ6ZXhhbXBsZTphYmZlMTNmNzEyMTIwNDMxYzI3NmUxMmVjYWIja2V5cy0xIiwidHlwIjoiSldUIn0"
	testHeader := struct {
		ALG string `json:"alg"`
		Typ string `json:"typ"`
		Kid string `json:"kid"`
	}{
		ALG: "SECP256R1",
		Typ: "JWT",
		Kid: "did:example:abfe13f712120431c276e12ecab#keys-1",
	}
	token, err := JWTVCEncodingHeaderI(testHeader)
	assert.NoError(t, err)
	assert.Equal(t, token, expectToken, "Error: Miss Match")
}

func TestJWTVCEncodingHeaderS_Expect_Correct_Value(t *testing.T) {
	expectToken := "eyJhbGciOiJTRUNQMjU2UjEiLCJraWQiOiJkaWQ6ZXhhbXBsZTphYmZlMTNmNzEyMTIwNDMxYzI3NmUxMmVjYWIja2V5cy0xIiwidHlwIjoiSldUIn0"
	testHeader := fmt.Sprintf(`{
		"alg": "SECP256R1",
		"typ": "JWT",
		"kid": "did:example:abfe13f712120431c276e12ecab#keys-1"
	}`)
	token, err := JWTVCEncodingHeaderS(testHeader)
	assert.NoError(t, err)
	assert.Equal(t, token, expectToken, "Error: Miss Match")
}

func TestJWTVCEncodingHeaderM_Expect_Correct_Value(t *testing.T) {
	expectToken := "eyJhbGciOiJTRUNQMjU2UjEiLCJraWQiOiJkaWQ6ZXhhbXBsZTphYmZlMTNmNzEyMTIwNDMxYzI3NmUxMmVjYWIja2V5cy0xIiwidHlwIjoiSldUIn0"
	testHeader := map[string]interface{}{
		"alg": "SECP256R1",
		"typ": "JWT",
		"kid": "did:example:abfe13f712120431c276e12ecab#keys-1",
	}
	token, err := JWTVCEncodingHeaderM(testHeader)
	assert.NoError(t, err)
	assert.Equal(t, token, expectToken, "Error: Miss Match")
}

func TestJWTVCEncodingClaimsI_Expect_Correct_Value(t *testing.T) {
	expectToken := "eyJleHAiOjE1NzMwMjk3MjMsImlhdCI6MTU0MTQ5MzcyNCwiaXNzIjoiaHR0cHM6Ly9leGFtcGxlLmNvbS9rZXlzL2Zvby5qd2siLCJqdGkiOiJodHRwOi8vZXhhbXBsZS5lZHUvY3JlZGVudGlhbHMvMzczMiIsIm5iZiI6MTU0MTQ5MzcyNCwibm9uY2UiOiI2NjAhNjM0NUZTZXIiLCJzdWIiOiJkaWQ6ZXhhbXBsZTplYmZlYjFmNzEyZWJjNmYxYzI3NmUxMmVjMjEiLCJ2YyI6eyJAY29udGV4dCI6WyJodHRwczovL3d3dy53My5vcmcvMjAxOC9jcmVkZW50aWFscy92MSIsImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL2V4YW1wbGVzL3YxIl0sImNyZWRlbnRpYWxTdWJqZWN0Ijp7ImRlZ3JlZSI6eyJuYW1lIjoiQmFjaGVsb3Igb2YgU2NpZW5jZSBhbmQgQXJ0cyIsInR5cGUiOiJCYWNoZWxvckRlZ3JlZSJ9fSwidHlwZSI6WyJWZXJpZmlhYmxlQ3JlZGVudGlhbCIsIlVuaXZlcnNpdHlEZWdyZWVDcmVkZW50aWFsIl19fQ"
	testClaims := JWTVCClaim{
		Subject:   "did:example:ebfeb1f712ebc6f1c276e12ec21",
		Id:        "http://example.edu/credentials/3732",
		Issuer:    "https://example.com/keys/foo.jwk",
		NotBefore: 1541493724,
		IssuedAt:  1541493724,
		ExpiresAt: 1573029723,
		Nonce:     "660!6345FSer",
		VC: JWTVC{
			Context: []string{
				"https://www.w3.org/2018/credentials/v1",
				"https://www.w3.org/2018/credentials/examples/v1",
			},
			Type: []string{"VerifiableCredential", "UniversityDegreeCredential"},
			CredentialSubject: map[string]interface{}{
				"degree": map[string]interface{}{
					"type": "BachelorDegree",
					"name": "Bachelor of Science and Arts",
				},
			},
		},
	}

	token, err := JWTVCEncodingClaimsI(testClaims)
	assert.NoError(t, err, fmt.Sprintf(`Error: %s`, err))
	assert.Equal(t, token, expectToken, "Error: Miss Match")
}

func TestJWTVCEncodingClaimsS_Expect_Correct_Value(t *testing.T) {
	expectToken := "eyJleHAiOjE1NzMwMjk3MjMsImlhdCI6MTU0MTQ5MzcyNCwiaXNzIjoiaHR0cHM6Ly9leGFtcGxlLmNvbS9rZXlzL2Zvby5qd2siLCJqdGkiOiJodHRwOi8vZXhhbXBsZS5lZHUvY3JlZGVudGlhbHMvMzczMiIsIm5iZiI6MTU0MTQ5MzcyNCwibm9uY2UiOiI2NjAhNjM0NUZTZXIiLCJzdWIiOiJkaWQ6ZXhhbXBsZTplYmZlYjFmNzEyZWJjNmYxYzI3NmUxMmVjMjEiLCJ2YyI6eyJAY29udGV4dCI6WyJodHRwczovL3d3dy53My5vcmcvMjAxOC9jcmVkZW50aWFscy92MSIsImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL2V4YW1wbGVzL3YxIl0sImNyZWRlbnRpYWxTdWJqZWN0Ijp7ImRlZ3JlZSI6eyJuYW1lIjoiQmFjaGVsb3Igb2YgU2NpZW5jZSBhbmQgQXJ0cyIsInR5cGUiOiJCYWNoZWxvckRlZ3JlZSJ9fSwidHlwZSI6WyJWZXJpZmlhYmxlQ3JlZGVudGlhbCIsIlVuaXZlcnNpdHlEZWdyZWVDcmVkZW50aWFsIl19fQ"
	testClaims := fmt.Sprintf(`{"sub":"did:example:ebfeb1f712ebc6f1c276e12ec21","jti":"http://example.edu/credentials/3732","iss":"https://example.com/keys/foo.jwk","nbf":1541493724,"iat":1541493724,"exp":1573029723,"nonce":"660!6345FSer","vc":{"@context":["https://www.w3.org/2018/credentials/v1","https://www.w3.org/2018/credentials/examples/v1"],"type":["VerifiableCredential","UniversityDegreeCredential"],"credentialSubject":{"degree":{"type":"BachelorDegree","name":"Bachelor of Science and Arts"}}}}`)

	token, err := JWTVCEncodingClaimsS(testClaims)
	assert.NoError(t, err, fmt.Sprintf(`Error: %s`, err))
	assert.Equal(t, token, expectToken, "Error: Miss Match")
}

func TestJWTVCEncodingClaimsM_Expect_Correct_Value(t *testing.T) {
	expectToken := "eyJleHAiOjE1NzMwMjk3MjMsImlhdCI6MTU0MTQ5MzcyNCwiaXNzIjoiaHR0cHM6Ly9leGFtcGxlLmNvbS9rZXlzL2Zvby5qd2siLCJqdGkiOiJodHRwOi8vZXhhbXBsZS5lZHUvY3JlZGVudGlhbHMvMzczMiIsIm5iZiI6MTU0MTQ5MzcyNCwibm9uY2UiOiI2NjAhNjM0NUZTZXIiLCJzdWIiOiJkaWQ6ZXhhbXBsZTplYmZlYjFmNzEyZWJjNmYxYzI3NmUxMmVjMjEiLCJ2YyI6eyJAY29udGV4dCI6WyJodHRwczovL3d3dy53My5vcmcvMjAxOC9jcmVkZW50aWFscy92MSIsImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL2V4YW1wbGVzL3YxIl0sImNyZWRlbnRpYWxTdWJqZWN0Ijp7ImRlZ3JlZSI6eyJuYW1lIjoiQmFjaGVsb3Igb2YgU2NpZW5jZSBhbmQgQXJ0cyIsInR5cGUiOiJCYWNoZWxvckRlZ3JlZSJ9fSwidHlwZSI6WyJWZXJpZmlhYmxlQ3JlZGVudGlhbCIsIlVuaXZlcnNpdHlEZWdyZWVDcmVkZW50aWFsIl19fQ"
	testClaims := map[string]interface{}{
		"sub":   "did:example:ebfeb1f712ebc6f1c276e12ec21",
		"jti":   "http://example.edu/credentials/3732",
		"iss":   "https://example.com/keys/foo.jwk",
		"nbf":   1541493724,
		"iat":   1541493724,
		"exp":   1573029723,
		"nonce": "660!6345FSer",
		"vc": map[string]interface{}{
			"@context": []string{
				"https://www.w3.org/2018/credentials/v1",
				"https://www.w3.org/2018/credentials/examples/v1",
			},
			"type": []string{"VerifiableCredential", "UniversityDegreeCredential"},
			"credentialSubject": map[string]interface{}{
				"degree": map[string]interface{}{
					"type": "BachelorDegree",
					"name": "Bachelor of Science and Arts",
				},
			},
		},
	}

	token, err := JWTVCEncodingClaimsM(testClaims)
	if err != nil {
		t.Fatal(fmt.Sprintf(`Error: %s`, err.Error()))
	}
	if token != expectToken {
		t.Fatal("Error: Miss Match")
	}
}

func TestJWTVCEncodingClaimsC_Expect_Correct_Value(t *testing.T) {
	expectToken := "eyJ2YyI6eyJAY29udGV4dCI6WyJodHRwczovL3d3dy53My5vcmcvMjAxOC9jcmVkZW50aWFscy92MSIsImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL2V4YW1wbGVzL3YxIl0sImNyZWRlbnRpYWxTdWJqZWN0Ijp7ImRlZ3JlZSI6eyJuYW1lIjoiQmFjaGVsb3Igb2YgU2NpZW5jZSBhbmQgQXJ0cyIsInR5cGUiOiJCYWNoZWxvckRlZ3JlZSJ9fSwidHlwZSI6WyJWZXJpZmlhYmxlQ3JlZGVudGlhbCIsIlVuaXZlcnNpdHlEZWdyZWVDcmVkZW50aWFsIl19LCJub25jZSI6IjY2MCE2MzQ1RlNlciIsImV4cCI6MTU3MzAyOTcyMywianRpIjoiaHR0cDovL2V4YW1wbGUuZWR1L2NyZWRlbnRpYWxzLzM3MzIiLCJpYXQiOjE1NDE0OTM3MjQsImlzcyI6Imh0dHBzOi8vZXhhbXBsZS5jb20va2V5cy9mb28uandrIiwibmJmIjoxNTQxNDkzNzI0LCJzdWIiOiJkaWQ6ZXhhbXBsZTplYmZlYjFmNzEyZWJjNmYxYzI3NmUxMmVjMjEifQ"
	testClaims := testClaimsStruct{
		Subject:   "did:example:ebfeb1f712ebc6f1c276e12ec21",
		Id:        "http://example.edu/credentials/3732",
		Issuer:    "https://example.com/keys/foo.jwk",
		NotBefore: 1541493724,
		IssuedAt:  1541493724,
		ExpiresAt: 1573029723,
		Nonce:     "660!6345FSer",
		VC: JWTVC{
			Context: []string{
				"https://www.w3.org/2018/credentials/v1",
				"https://www.w3.org/2018/credentials/examples/v1",
			},
			Type: []string{"VerifiableCredential", "UniversityDegreeCredential"},
			CredentialSubject: map[string]interface{}{
				"degree": map[string]interface{}{
					"type": "BachelorDegree",
					"name": "Bachelor of Science and Arts",
				},
			},
		},
	}

	token, err := JWTVCEncodingClaimsC(testClaims)
	assert.NoError(t, err, fmt.Sprintf(`Error: %s`, err))
	assert.Equal(t, token, expectToken, "Error: Miss Match")
}

func TestJWTVCEncoding_Expect_Correct_Value(t *testing.T) {
	testHeader := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
		"kid": "did:example:abfe13f712120431c276e12ecab#keys-1",
	}
	testClaims := testClaimsStruct{
		Subject:   "did:example:ebfeb1f712ebc6f1c276e12ec21",
		Id:        "http://example.edu/credentials/3732",
		Issuer:    "https://example.com/keys/foo.jwk",
		NotBefore: 1573029723,
		IssuedAt:  1573029723,
		ExpiresAt: 1673029723,
		Nonce:     "660!6345FSer",
		VC: JWTVC{
			Context: []string{
				"https://www.w3.org/2018/credentials/v1",
				"https://www.w3.org/2018/credentials/examples/v1",
			},
			Type: []string{"VerifiableCredential", "UniversityDegreeCredential"},
			CredentialSubject: map[string]interface{}{
				"degree": map[string]interface{}{
					"type": "BachelorDegree",
					"name": "Bachelor of Science and Arts",
				},
			},
		},
	}

	token, err := JWTVCEncoding(testHeader, testClaims, []byte(privateKey))
	assert.NoError(t, err, fmt.Sprintf(`Error: %s`, err))
	assert.NotEmpty(t, token, "Error: Miss Match")
}

func TestJWTVCDecoding_Expect_Collect_Value(t *testing.T) {
	expectHeader := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
		"kid": "did:example:abfe13f712120431c276e12ecab#keys-1",
	}
	expectClaims := testClaimsStruct{
		Subject:   "did:example:ebfeb1f712ebc6f1c276e12ec21",
		Id:        "http://example.edu/credentials/3732",
		Issuer:    "https://example.com/keys/foo.jwk",
		NotBefore: 1573029723,
		IssuedAt:  1573029723,
		ExpiresAt: 1673029723,
		Nonce:     "660!6345FSer",
		VC: JWTVC{
			Context: []string{
				"https://www.w3.org/2018/credentials/v1",
				"https://www.w3.org/2018/credentials/examples/v1",
			},
			Type: []string{"VerifiableCredential", "UniversityDegreeCredential"},
			CredentialSubject: map[string]interface{}{
				"degree": map[string]interface{}{
					"type": "BachelorDegree",
					"name": "Bachelor of Science and Arts",
				},
			},
		},
	}

	token := "eyJhbGciOiJIUzI1NiIsImtpZCI6ImRpZDpleGFtcGxlOmFiZmUxM2Y3MTIxMjA0MzFjMjc2ZTEyZWNhYiNrZXlzLTEiLCJ0eXAiOiJKV1QifQ.eyJ2YyI6eyJAY29udGV4dCI6WyJodHRwczovL3d3dy53My5vcmcvMjAxOC9jcmVkZW50aWFscy92MSIsImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL2V4YW1wbGVzL3YxIl0sImNyZWRlbnRpYWxTdWJqZWN0Ijp7ImRlZ3JlZSI6eyJuYW1lIjoiQmFjaGVsb3Igb2YgU2NpZW5jZSBhbmQgQXJ0cyIsInR5cGUiOiJCYWNoZWxvckRlZ3JlZSJ9fSwidHlwZSI6WyJWZXJpZmlhYmxlQ3JlZGVudGlhbCIsIlVuaXZlcnNpdHlEZWdyZWVDcmVkZW50aWFsIl19LCJub25jZSI6IjY2MCE2MzQ1RlNlciIsImV4cCI6MTY3MzAyOTcyMywianRpIjoiaHR0cDovL2V4YW1wbGUuZWR1L2NyZWRlbnRpYWxzLzM3MzIiLCJpYXQiOjE1NzMwMjk3MjMsImlzcyI6Imh0dHBzOi8vZXhhbXBsZS5jb20va2V5cy9mb28uandrIiwibmJmIjoxNTczMDI5NzIzLCJzdWIiOiJkaWQ6ZXhhbXBsZTplYmZlYjFmNzEyZWJjNmYxYzI3NmUxMmVjMjEifQ.stv6yh2cdBzpSujybQVhw3jhyzA9-pPIsyC6ANtpaTM"
	tokenM, err := JWTVCDecodingM(token, []byte(privateKey))
	assert.NoError(t, err, fmt.Sprintf(`Error: %s`, err))

	if header, ok := tokenM["Header"].(map[string]interface{}); ok {
		assert.Equal(t, header["alg"], expectHeader["alg"], "Error: Header Alg Miss Match")
		assert.Equal(t, header["typ"], expectHeader["typ"], "Error: Header Type Miss Match")
		assert.Equal(t, header["kid"], expectHeader["kid"], "Error: Header KID Miss Match")
	} else {
		t.Fatal("Error: Header Wrong Format")
	}

	if claims, ok := tokenM["Claims"].(map[string]interface{}); ok {
		assert.Equal(t, claims["exp"], float64(expectClaims.ExpiresAt), "Error: Claims ExpiresAt Miss Match")
		assert.Equal(t, claims["iat"], float64(expectClaims.IssuedAt), "Error: Claims IssuedAt Miss Match")
		assert.Equal(t, claims["iss"], expectClaims.Issuer, "Error: Claims Issuer Miss Match")
		assert.Equal(t, claims["jti"], expectClaims.Id, "Error: Claims Id (jti) Miss Match")
		assert.Equal(t, claims["nbf"], float64(expectClaims.NotBefore), "Error: Claims NotBefore Miss Match")
		assert.Equal(t, claims["nonce"], expectClaims.Nonce, "Error: Claims Nonce Miss Match")
		assert.Equal(t, claims["sub"], expectClaims.Subject, "Error: Claims Subject Miss Match")
	} else {
		t.Fatal("Error: Claims Wrong Format")
	}
}
