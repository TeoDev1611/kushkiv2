package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"math/big"
	mrand "math/rand/v2"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/pkcs12"
)

// Signer maneja la carga del certificado y la firma de documentos.
type Signer struct {
	PrivateKey  *rsa.PrivateKey
	Certificate *x509.Certificate
}

// NewSigner crea un Signer a partir de llaves en memoria.
func NewSigner(privateKey *rsa.PrivateKey, certificate *x509.Certificate) *Signer {
	return &Signer{
		PrivateKey:  privateKey,
		Certificate: certificate,
	}
}

// NewSignerFromFile carga un archivo .p12 y extrae la llave privada y el certificado.
func NewSignerFromFile(p12Path string, password string) (*Signer, error) {
	p12Data, err := os.ReadFile(p12Path)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo p12: %v", err)
	}

	blocks, err := pkcs12.ToPEM(p12Data, password)
	if err != nil {
		return nil, fmt.Errorf("contraseña incorrecta o formato inválido: %v", err)
	}

	var privateKey *rsa.PrivateKey
	var certificate *x509.Certificate

	for _, block := range blocks {
		if block.Type == "PRIVATE KEY" || block.Type == "RSA PRIVATE KEY" {
			if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
				privateKey = key
			} else if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
				if rsaKey, ok := key.(*rsa.PrivateKey); ok {
					privateKey = rsaKey
				}
			}
		}

		if block.Type == "CERTIFICATE" {
			if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
				// Preferimos el certificado que no sea CA, o usamos el primero si no hay opción
				if !cert.IsCA || certificate == nil {
					certificate = cert
				}
			}
		}
	}

	if privateKey == nil {
		return nil, fmt.Errorf("no se encontró llave privada en el archivo P12")
	}
	if certificate == nil {
		return nil, fmt.Errorf("no se encontró certificado en el archivo P12")
	}

	return NewSigner(privateKey, certificate), nil
}

// ValidateCert verifica si la ruta y contraseña son válidas.
func ValidateCert(p12Path, password string) error {
	_, err := NewSignerFromFile(p12Path, password)
	return err
}

// SignXML firma un XML usando XAdES-BES.
func (s *Signer) SignXML(xmlData []byte) ([]byte, error) {
	// 1. Preparar Datos Aleatorios y Tiempos
	signatureID := fmt.Sprintf("Signature-%d", mrand.IntN(1000000))
	signedPropsID := fmt.Sprintf("SignedProperties-%s", signatureID)
	objectID := fmt.Sprintf("Object-%s", signatureID)
	referenceID := fmt.Sprintf("Reference-%s", signatureID)
	signedInfoID := fmt.Sprintf("SignedInfo-%s", signatureID)
	keyInfoID := fmt.Sprintf("KeyInfo-%s", signatureID)

	// 2. Canonicalizar y Hashear el Documento (Comprobante)
	// Asumimos que xmlData es el XML limpio.
	// El transform 'enveloped-signature' implica que hasheamos el doc original (sin firma).
	// SRI requiere que el comprobante sea canonicalizado.
	
	// FIX: El Reference apunta a "#comprobante" (la etiqueta raíz), no al documento entero.
	// Por ende, el hash no debe incluir la declaración XML (<?xml ...?>).
	xmlStrForHash := string(xmlData)
	startIndex := strings.Index(xmlStrForHash, "<factura")
	if startIndex != -1 {
		xmlStrForHash = xmlStrForHash[startIndex:]
	}
	
	// HASH del Documento
	hashDocumento := sha1.Sum([]byte(xmlStrForHash))
	digestDocumento := base64.StdEncoding.EncodeToString(hashDocumento[:])

	// 3. Construir SignedProperties (XAdES)
	// Este bloque contiene la fecha de firma y el certificado firmante.
	// Es CRÍTICO que la estructura de bytes sea exacta para el hash.
	
	certDigest := sha1.Sum(s.Certificate.Raw)
	certDigestB64 := base64.StdEncoding.EncodeToString(certDigest[:])
	issuerName := s.Certificate.Issuer.String() // Esto puede requerir formateo específico RFC4514 si SRI es estricto
	serialNumber := s.Certificate.SerialNumber.String()
	
	signingTime := time.Now().Format("2006-01-02T15:04:05")

	// Construcción Manual de SignedProperties para control total de C14N
	// Nota: Los namespaces deben coincidir con los usados en SignedInfo
	signedPropertiesXML := fmt.Sprintf(`<etsi:SignedProperties Id="%s" xmlns:etsi="http://uri.etsi.org/01903/v1.3.2#"><etsi:SignedSignatureProperties><etsi:SigningTime>%s</etsi:SigningTime><etsi:SigningCertificate><etsi:Cert><etsi:CertDigest><ds:DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1" xmlns:ds="http://www.w3.org/2000/09/xmldsig#"></ds:DigestMethod><ds:DigestValue xmlns:ds="http://www.w3.org/2000/09/xmldsig#">%s</ds:DigestValue></etsi:CertDigest><etsi:IssuerSerial><ds:X509IssuerName xmlns:ds="http://www.w3.org/2000/09/xmldsig#">%s</ds:X509IssuerName><ds:X509SerialNumber xmlns:ds="http://www.w3.org/2000/09/xmldsig#">%s</ds:X509SerialNumber></etsi:IssuerSerial></etsi:Cert></etsi:SigningCertificate></etsi:SignedSignatureProperties><etsi:SignedDataObjectProperties><etsi:DataObjectFormat ObjectReference="#%s"><etsi:Description>contenido comprobante</etsi:Description><etsi:MimeType>text/xml</etsi:MimeType></etsi:DataObjectFormat></etsi:SignedDataObjectProperties></etsi:SignedProperties>`,
		signedPropsID,
		signingTime,
		certDigestB64,
		issuerName,
		serialNumber,
		referenceID,
	)

	// Hash de SignedProperties
	hashSignedProps := sha1.Sum([]byte(signedPropertiesXML))
	digestSignedProps := base64.StdEncoding.EncodeToString(hashSignedProps[:])

	// 4. Construir SignedInfo
	// Contiene las referencias al documento y a las propiedades firmadas.
	kushkiModulus := base64.StdEncoding.EncodeToString(s.PrivateKey.N.Bytes())
	kushkiExponent := base64.StdEncoding.EncodeToString(big.NewInt(int64(s.PrivateKey.E)).Bytes())

	signedInfoXML := fmt.Sprintf(`<ds:SignedInfo Id="%s" xmlns:ds="http://www.w3.org/2000/09/xmldsig#" xmlns:etsi="http://uri.etsi.org/01903/v1.3.2#"><ds:CanonicalizationMethod Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"></ds:CanonicalizationMethod><ds:SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"></ds:SignatureMethod><ds:Reference Id="%s" URI="#comprobante"><ds:Transforms><ds:Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"></ds:Transform></ds:Transforms><ds:DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"></ds:DigestMethod><ds:DigestValue>%s</ds:DigestValue></ds:Reference><ds:Reference URI="#%s"><ds:Transforms><ds:Transform Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"></ds:Transform></ds:Transforms><ds:DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"></ds:DigestMethod><ds:DigestValue>%s</ds:DigestValue></ds:Reference></ds:SignedInfo>`,
		signedInfoID,
		referenceID,
		digestDocumento,
		signedPropsID,
		digestSignedProps,
	)

	// 5. Calcular Firma (RSA-SHA1) sobre SignedInfo
	// IMPORTANTE: Se firma el hash del SignedInfo CANONICALIZADO.
	// Como lo construimos string-by-string, asumimos que "signedInfoXML" YA ES la forma canónica deseada.
	hashSignedInfo := sha1.Sum([]byte(signedInfoXML))
	signatureBytes, err := rsa.SignPKCS1v15(rand.Reader, s.PrivateKey, crypto.SHA1, hashSignedInfo[:])
	if err != nil {
		return nil, fmt.Errorf("error generando firma RSA: %v", err)
	}
	signatureValue := base64.StdEncoding.EncodeToString(signatureBytes)

	// 6. Construir Bloque Completo <ds:Signature>
	certificateB64 := base64.StdEncoding.EncodeToString(s.Certificate.Raw)
	// Formatear certificado con saltos de línea cada 76 caracteres es buena práctica pero no obligatoria en XMLDSig puro, 
	// aunque algunos validadores lo prefieren. Lo dejaremos en una línea por simplicidad.

	fullSignature := fmt.Sprintf(`<ds:Signature Id="%s" xmlns:ds="http://www.w3.org/2000/09/xmldsig#" xmlns:etsi="http://uri.etsi.org/01903/v1.3.2#">%s<ds:SignatureValue Id="SignatureValue-%s">%s</ds:SignatureValue><ds:KeyInfo Id="%s"><ds:X509Data><ds:X509Certificate>%s</ds:X509Certificate></ds:X509Data><ds:KeyValue><ds:RSAKeyValue><ds:Modulus>%s</ds:Modulus><ds:Exponent>%s</ds:Exponent></ds:RSAKeyValue></ds:KeyValue></ds:KeyInfo><ds:Object Id="%s"><etsi:QualifyingProperties Target="#%s"><etsi:SignedProperties Id="%s"><etsi:SignedSignatureProperties><etsi:SigningTime>%s</etsi:SigningTime><etsi:SigningCertificate><etsi:Cert><etsi:CertDigest><ds:DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"></ds:DigestMethod><ds:DigestValue>%s</ds:DigestValue></etsi:CertDigest><etsi:IssuerSerial><ds:X509IssuerName>%s</ds:X509IssuerName><ds:X509SerialNumber>%s</ds:X509SerialNumber></etsi:IssuerSerial></etsi:Cert></etsi:SigningCertificate></etsi:SignedSignatureProperties><etsi:SignedDataObjectProperties><etsi:DataObjectFormat ObjectReference="#%s"><etsi:Description>contenido comprobante</etsi:Description><etsi:MimeType>text/xml</etsi:MimeType></etsi:DataObjectFormat></etsi:SignedDataObjectProperties></etsi:SignedProperties></etsi:QualifyingProperties></ds:Object></ds:Signature>`,
		signatureID,
		signedInfoXML, // Ya incluye los xmlns
		signatureID,
		signatureValue,
		keyInfoID,
		certificateB64,
		kushkiModulus,
		kushkiExponent,
		objectID,
		signatureID,
		signedPropsID,
		signingTime,
		certDigestB64,
		issuerName,
		serialNumber,
		referenceID,
	)

	// 7. Insertar Firma en el XML Original
	// Buscamos la etiqueta de cierre </factura> y anteponemos la firma.
	// Asumimos que xmlData termina en </factura> o similar.
	xmlStr := string(xmlData)
	endTag := "</factura>"
	if strings.Contains(xmlStr, endTag) {
		// Insertar antes del cierre
		finalXML := strings.Replace(xmlStr, endTag, fullSignature+endTag, 1)
		return []byte(finalXML), nil
	}

	return nil, fmt.Errorf("no se encontró la etiqueta de cierre </factura> para insertar la firma")
}
