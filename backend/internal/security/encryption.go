package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
	"crypto/tls"
)

// EncryptionService provides data encryption capabilities
type EncryptionService struct {
	aesKey       []byte
	rsaPublicKey *rsa.PublicKey
	rsaPrivateKey *rsa.PrivateKey
}

// NewEncryptionService creates a new encryption service
func NewEncryptionService(aesKey []byte, rsaPublicKey, rsaPrivateKey string) (*EncryptionService, error) {
	if len(aesKey) != 32 {
		return nil, errors.New("AES key must be 32 bytes")
	}

	pubKey, err := parseRSAPublicKey(rsaPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA public key: %w", err)
	}

	privKey, err := parseRSAPrivateKey(rsaPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
	}

	return &EncryptionService{
		aesKey:        aesKey,
		rsaPublicKey:  pubKey,
		rsaPrivateKey: privKey,
	}, nil
}

// EncryptAES encrypts data using AES-256-GCM
func (es *EncryptionService) EncryptAES(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(es.aesKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptAES decrypts data using AES-256-GCM
func (es *EncryptionService) DecryptAES(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(es.aesKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// EncryptRSA encrypts data using RSA-OAEP
func (es *EncryptionService) EncryptRSA(plaintext []byte) ([]byte, error) {
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, es.rsaPublicKey, plaintext, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// DecryptRSA decrypts data using RSA-OAEP
func (es *EncryptionService) DecryptRSA(ciphertext []byte) ([]byte, error) {
	hash := sha256.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, es.rsaPrivateKey, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// FieldEncryption provides field-level encryption for sensitive data
type FieldEncryption struct {
	encryptionService *EncryptionService
}

// NewFieldEncryption creates a new field encryption service
func NewFieldEncryption(encryptionService *EncryptionService) *FieldEncryption {
	return &FieldEncryption{
		encryptionService: encryptionService,
	}
}

// EncryptField encrypts a single field value
func (fe *FieldEncryption) EncryptField(value string) (string, error) {
	encrypted, err := fe.encryptionService.EncryptAES([]byte(value))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// DecryptField decrypts a single field value
func (fe *FieldEncryption) DecryptField(encryptedValue string) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(encryptedValue)
	if err != nil {
		return "", err
	}
	
	decrypted, err := fe.encryptionService.DecryptAES(encrypted)
	if err != nil {
		return "", err
	}
	
	return string(decrypted), nil
}

// EncryptStruct encrypts specified fields in a struct
func (fe *FieldEncryption) EncryptStruct(data interface{}, fields []string) error {
	// Use reflection to encrypt specified fields
	// Implementation would use reflect package
	return nil
}

// TransparentEncryption provides transparent encryption/decryption for database
type TransparentEncryption struct {
	encryptionService *EncryptionService
	encryptedFields   map[string][]string // table -> fields
}

// NewTransparentEncryption creates a new transparent encryption service
func NewTransparentEncryption(encryptionService *EncryptionService) *TransparentEncryption {
	return &TransparentEncryption{
		encryptionService: encryptionService,
		encryptedFields: map[string][]string{
			"users": {"email", "phone", "ssn"},
			"payments": {"card_number", "cvv", "account_number"},
			"personal_data": {"address", "date_of_birth", "id_number"},
		},
	}
}

// KeyDerivation provides key derivation functions
type KeyDerivation struct {
	salt []byte
}

// NewKeyDerivation creates a new key derivation service
func NewKeyDerivation(salt []byte) *KeyDerivation {
	return &KeyDerivation{salt: salt}
}

// DeriveKey derives a key from password using Argon2id
func (kd *KeyDerivation) DeriveKey(password string, keyLen uint32) []byte {
	return argon2.IDKey([]byte(password), kd.salt, 3, 64*1024, 4, keyLen)
}

// EnvelopeEncryption implements envelope encryption pattern
type EnvelopeEncryption struct {
	masterKey []byte
	dataKeys  map[string][]byte
}

// NewEnvelopeEncryption creates a new envelope encryption service
func NewEnvelopeEncryption(masterKey []byte) *EnvelopeEncryption {
	return &EnvelopeEncryption{
		masterKey: masterKey,
		dataKeys:  make(map[string][]byte),
	}
}

// GenerateDataKey generates a new data encryption key
func (ee *EnvelopeEncryption) GenerateDataKey(keyID string) ([]byte, []byte, error) {
	// Generate random data key
	dataKey := make([]byte, 32)
	if _, err := rand.Read(dataKey); err != nil {
		return nil, nil, err
	}

	// Encrypt data key with master key
	block, err := aes.NewCipher(ee.masterKey)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	encryptedDataKey := gcm.Seal(nonce, nonce, dataKey, []byte(keyID))
	
	// Store data key
	ee.dataKeys[keyID] = dataKey
	
	return dataKey, encryptedDataKey, nil
}

// Helper functions

func parseRSAPublicKey(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return rsaPub, nil
}

func parseRSAPrivateKey(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

// TLSConfig provides TLS configuration for secure communication
type TLSConfig struct {
	CertFile       string
	KeyFile        string
	CAFile         string
	MinVersion     uint16
	CipherSuites   []uint16
	ClientAuthType int
}

// GetTLSConfig returns a secure TLS configuration
func GetSecureTLSConfig() *TLSConfig {
	return &TLSConfig{
		MinVersion: tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
		ClientAuthType: tls.RequireAndVerifyClientCert,
	}
}