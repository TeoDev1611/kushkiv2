package crypto

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
)

// VerifyLicenseToken verifica la firma de un token JWT usando la llave pública incrustada.
// Solo soporta RS256.
func VerifyLicenseToken(tokenString string) error {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return errors.New("token format invalid: expected 3 parts")
	}

	header := parts[0]
	payload := parts[1]
	signature := parts[2]

	// Reconstruir el mensaje firmado (header + "." + payload)
	signedContent := header + "." + payload

	// Decodificar la firma (Base64 URL)
	// JWT usa Base64URL sin padding
	sigBytes, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		// Intentar con padding si falla, por robustez
		sigBytes, err = base64.URLEncoding.DecodeString(signature)
		if err != nil {
			return errors.New("invalid signature encoding")
		}
	}

	// Obtener Llave Pública
	pubKey, err := GetPublicKey()
	if err != nil {
		return err
	}

	// Hashear el contenido
	hashed := sha256.Sum256([]byte(signedContent))

	// Verificar Firma
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], sigBytes)
	if err != nil {
		return errors.New("invalid token signature")
	}

	return nil
}
