package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

// encryptionKey almacena la llave en memoria RAM una vez cargada.
var encryptionKey []byte

// InitSecurity inicializa el sistema criptográfico.
// Busca un archivo 'master.key'. Si no existe, genera uno nuevo aleatorio de 32 bytes (AES-256).
func InitSecurity(keyPath string) error {
	// 1. Intentar leer la llave existente
	key, err := os.ReadFile(keyPath)
	if err == nil {
		if len(key) != 32 {
			return fmt.Errorf("la llave maestra en %s está corrupta (longitud incorrecta)", keyPath)
		}
		encryptionKey = key
		return nil
	}

	// 2. Si no existe, generar una nueva llave segura
	if os.IsNotExist(err) {
		newKey := make([]byte, 32) // 32 bytes = 256 bits
		if _, err := io.ReadFull(rand.Reader, newKey); err != nil {
			return fmt.Errorf("error generando entropía para la llave: %v", err)
		}

		// 3. Guardar la nueva llave en el disco (Permisos 0600: solo lectura para el usuario)
		if err := os.WriteFile(keyPath, newKey, 0600); err != nil {
			return fmt.Errorf("error guardando master.key: %v", err)
		}

		encryptionKey = newKey
		return nil
	}

	return err
}

func Encrypt(text string) (string, error) {
	if len(encryptionKey) == 0 {
		return "", errors.New("sistema de seguridad no inicializado")
	}

	c, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(text), nil)), nil
}

func Decrypt(ciphertext string) (string, error) {
	if len(encryptionKey) == 0 {
		return "", errors.New("sistema de seguridad no inicializado")
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}