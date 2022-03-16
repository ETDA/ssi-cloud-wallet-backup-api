package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"strings"
)

var defaultJWTSigningMethod = jwt.SigningMethodHS256

var signingMethod = map[string]jwt.SigningMethod{
	"secp256r1": jwt.SigningMethodES256,
	"hs256":     jwt.SigningMethodHS256,
}

type JWTVCCredentialSchema struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type JWTVC struct {
	Context           []string               `json:"@context"`
	CredentialSubject map[string]interface{} `json:"credentialSubject"`
	Type              []string               `json:"type"`
	CredentialSchema  *JWTVCCredentialSchema `json:"credentialSchema"`
}

type JWTOptions struct {
	SigningMethod jwt.SigningMethod `json:"signing_method"`
}

type JWTVCClaim struct {
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

func GetDefaultJWTOptions() JWTOptions {
	return JWTOptions{
		SigningMethod: defaultJWTSigningMethod,
	}
}

// External Functions

// JWTVCEncodingHeaderI receive header type interface return string and error
func JWTVCEncodingHeaderI(header interface{}, options ...JWTOptions) (string, error) {
	method := getJWTOptions(options...).SigningMethod
	headerByte, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	headerMap := map[string]interface{}{}
	err = json.Unmarshal(headerByte, &headerMap)
	if err != nil {
		return "", err
	}

	return jwtEncodingHeader(headerMap, method)
}

// JWTVCEncodingHeaderS receive header type string return string and error
func JWTVCEncodingHeaderS(header string, options ...JWTOptions) (string, error) {
	method := getJWTOptions(options...).SigningMethod
	headerMap := map[string]interface{}{}

	err := json.Unmarshal([]byte(header), &headerMap)
	if err != nil {
		return "", err
	}

	return jwtEncodingHeader(headerMap, method)
}

// JWTVCEncodingHeaderM receive header type map string return string and error
func JWTVCEncodingHeaderM(header map[string]interface{}, options ...JWTOptions) (string, error) {
	method := getJWTOptions(options...).SigningMethod
	return jwtEncodingHeader(header, method)
}

// JWTVCEncodingClaimsI receive claims type interface return string and error
func JWTVCEncodingClaimsI(claims interface{}, options ...JWTOptions) (string, error) {
	method := getJWTOptions(options...).SigningMethod

	claimsByte, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	claimsMap := jwt.MapClaims{}
	err = json.Unmarshal(claimsByte, &claimsMap)
	if err != nil {
		return "", err
	}

	return jwtEncodingClaims(claimsMap, method)
}

// JWTVCEncodingClaimsS receive claims type string return string and error
func JWTVCEncodingClaimsS(claims string, options ...JWTOptions) (string, error) {
	method := getJWTOptions(options...).SigningMethod

	claimsMap := jwt.MapClaims{}
	err := json.Unmarshal([]byte(claims), &claimsMap)
	if err != nil {
		return "", err
	}

	return jwtEncodingClaims(claimsMap, method)
}

// JWTVCEncodingClaimsM receive claims type map string return string and error
func JWTVCEncodingClaimsM(claims map[string]interface{}, options ...JWTOptions) (string, error) {
	method := getJWTOptions(options...).SigningMethod

	claimsByte, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	claimsMap := jwt.MapClaims{}
	err = json.Unmarshal(claimsByte, &claimsMap)
	if err != nil {
		return "", err
	}

	return jwtEncodingClaims(claimsMap, method)
}

// JWTVCEncodingClaimsC receive claims type Struct return string and error
func JWTVCEncodingClaimsC(claims jwt.Claims, options ...JWTOptions) (string, error) {
	method := getJWTOptions(options...).SigningMethod
	return jwtEncodingClaims(claims, method)
}

// JWTVCEncoding receive all part of jwt return string error
func JWTVCEncoding(header map[string]interface{}, claims jwt.Claims, secretKey []byte, options ...JWTOptions) (string, error) {
	return jwtVCEncoding(header, claims, secretKey, options...)
}

// JWTVCDecodingT return Token Struct
func JWTVCDecodingT(tokenStr string, secretKey []byte) (*jwt.Token, error) {
	return jwtVCDecoding(tokenStr, secretKey)
}

// JWTVCDecodingM return map string
func JWTVCDecodingM(tokenStr string, secretKey []byte) (map[string]interface{}, error) {
	token, err := jwtVCDecoding(tokenStr, secretKey)
	if err != nil {
		return nil, err
	}

	tokenM := map[string]interface{}{}
	tokenByte, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(tokenByte, &tokenM)
	if err != nil {
		return nil, err
	}

	return tokenM, nil
}

// Internal Functions

func jwtEncodingHeader(header map[string]interface{}, method jwt.SigningMethod) (string, error) {
	token := jwt.Token{
		Header: header,
		Method: method,
	}
	tokenStr, err := token.SigningString()
	if err != nil {
		return "", err
	}

	return removeNullFromToken(tokenStr), nil
}

func jwtEncodingClaims(claims jwt.Claims, method jwt.SigningMethod) (string, error) {
	token := jwt.Token{
		Claims: claims,
		Method: method,
	}
	tokenStr, err := token.SigningString()
	if err != nil {
		return "", err
	}

	return removeNullFromToken(tokenStr), nil
}

func jwtVCEncoding(header map[string]interface{}, claims jwt.Claims, secretKey []byte, options ...JWTOptions) (string, error) {
	method := getMethodFromHeader(header)
	if method == nil {
		log.Fatal("Method is invalid")
		return "", nil
	}

	headerToken, err := JWTVCEncodingHeaderM(header, options...)
	if err != nil {
		return "", nil
	}

	claimsToken, err := JWTVCEncodingClaimsC(claims, options...)
	if err != nil {
		return "", nil
	}

	signature := fmt.Sprintf(`%s.%s`, headerToken, claimsToken)
	token := jwt.Token{
		Header:    header,
		Claims:    claims,
		Signature: signature,
		Method:    method,
	}

	tokenStr, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func jwtVCDecoding(tokenStr string, secretKey []byte) (*jwt.Token, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return token, err
	}

	return token, nil
}

func removeNullFromToken(token string) string {
	return strings.Replace(strings.Replace(token, ".bnVsbA", "", 1), "bnVsbA.", "", 1)
}

func getJWTOptions(options ...JWTOptions) JWTOptions {
	if len(options) == 0 {
		return GetDefaultJWTOptions()
	}

	option := JWTOptions{}
	for _, o := range options {
		if o.SigningMethod != nil {
			option.SigningMethod = o.SigningMethod
		}
	}

	if option.SigningMethod == nil {
		option.SigningMethod = defaultJWTSigningMethod
	}

	return option
}

func getMethodFromHeader(header map[string]interface{}) jwt.SigningMethod {
	if method, ok := header["alg"].(string); ok {
		alg := signingMethod[strings.ToLower(method)]
		return alg
	}

	return nil
}
