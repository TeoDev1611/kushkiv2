package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"strings"
	"testing"
	"time"
)

func createTestCredentials(t *testing.T) (*rsa.PrivateKey, *x509.Certificate) {
	// Generate Key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Generate Cert
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Org"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour),
		KeyUsage:  x509.KeyUsageDigitalSignature,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		t.Fatalf("Failed to parse certificate: %v", err)
	}

	return priv, cert
}

func TestSignXML(t *testing.T) {
	priv, cert := createTestCredentials(t)
	signer := NewSigner(priv, cert)

	// XML de prueba con la etiqueta de cierre esperada
	xmlInput := `<?xml version="1.0" encoding="UTF-8"?>
<factura id="comprobante" version="1.0.0">
	<infoTributaria>
		<claveAcceso>1234567890</claveAcceso>
	</infoTributaria>
</factura>`
	
	signedXML, err := signer.SignXML([]byte(xmlInput))
	if err != nil {
		t.Fatalf("SignXML failed: %v", err)
	}

	signedString := string(signedXML)

	// Validations
	checks := []string{
		`<ds:Signature`,
		`<ds:SignedInfo`,
		`<ds:KeyInfo`,
		`<etsi:SignedProperties`,
		`SignatureValue`,
		`URI="#comprobante"`, // Verificar referencia al ID
		`</factura>`, // Verificar que cierra bien
	}

	for _, check := range checks {
		if !strings.Contains(signedString, check) {
			t.Errorf("Signed XML missing %s", check)
		}
	}
}