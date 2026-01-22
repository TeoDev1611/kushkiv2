package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// Variables ofuscadas
const (
	// URL: Reversed + Base64
	obfuscatedURL = "dmVkLm9uZWQuOTktZ25pbGxpYi1zeXMvLzpzcHR0aA=="
	// Key: Base64
	obfuscatedKey = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUFyV2dxa2t2YnVMcmRzRTVjSG82ckhTakhCSlNENWE3ZzRZOEllbm1maUYrY09KNEF1RFIxejExYVJXUTRRUDU1NGRWU3QwZTkzWE9CdVNxN1ViOVJndW9yaDNkQ1dVTVNzalRvcG83Vi83Q2ZFRnlOVmMvZGU2SUpndXR6Ymdud2dOS1JNaWF4T2NkK1lwMnVHcmhPaHpvZkE3VURqWnlRbGVKTy9IcVdsZUh0YkxxZExDZDYyYk5kRE8yZWcraFVxRGd0cmlwVVRtK3hhbG9lN1BoYlBpK2RQQ2lFVGRJd3VOcGxOWDFrRHRCaEhhM0FVWnlTV1Myc0tuNEhIbUUxVkw1NGF4SWUrSFEwY2hGR0pCbHpXN0tBZFZPeGpYbGtTamhHSEF2VTlMT3RSeHdmcmZTWXFoc0lyZXRHNlpHaW1VbytCOTk1cnU4NXhwVnVHQm5QUndJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0t"
)

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// GetAPIURL revela la URL de la API.
func GetAPIURL() string {
	decoded, _ := base64.StdEncoding.DecodeString(obfuscatedURL)
	return reverse(string(decoded))
}

// GetPublicKey parsea y retorna la llave pública RSA.
func GetPublicKey() (*rsa.PublicKey, error) {
	pemData, err := base64.StdEncoding.DecodeString(obfuscatedKey)
	if err != nil {
		return nil, fmt.Errorf("error decoding key config: %v", err)
	}

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DER encoded public key: %v", err)
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, fmt.Errorf("key type is not RSA")
	}
}
