package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"kushkiv2/pkg/crypto"
)

// CloudService maneja la comunicación con el backend en Deno (Licenciamiento y Reportes)
type CloudService struct {
	BaseURL string
	Client  *http.Client
}

// NewCloudService crea una nueva instancia del servicio
func NewCloudService() *CloudService {
	return &CloudService{
		// URL obtenida de configuración segura
		BaseURL: crypto.GetAPIURL(),
		Client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// Estructuras para Licenciamiento
type LicenseRequest struct {
	LicenseKey string `json:"license_key"`
	MachineID  string `json:"machine_id"`
}

type LicenseResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"` // Algunos backends usan 'message' para errores
}

// GetMachineID genera una huella digital única del hardware.
// Combina Hostname + OS + Arch + User ID (si disponible) y genera un hash SHA256.
func (s *CloudService) GetMachineID() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("error obteniendo hostname: %w", err)
	}

	// Información base del sistema
	platformInfo := fmt.Sprintf("%s-%s-%s", hostname, runtime.GOOS, runtime.GOARCH)

	// Generar Hash
	hash := sha256.Sum256([]byte(platformInfo))
	return hex.EncodeToString(hash[:]), nil
}

// ActivateLicense solicita la activación de la licencia al backend
func (s *CloudService) ActivateLicense(licenseKey string) (*LicenseResponse, error) {
	machineID, err := s.GetMachineID()
	if err != nil {
		return nil, err
	}

	// Limpiar espacios en blanco comunes en copy-paste
	licenseKey = strings.TrimSpace(licenseKey)

	reqBody := LicenseRequest{
		LicenseKey: licenseKey,
		MachineID:  machineID,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error serializando request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/license/activate", s.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creando request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error de conexión: %w", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if jsonErr := json.Unmarshal(bodyBytes, &errResp); jsonErr == nil {
			msg := errResp.Error
			if msg == "" {
				msg = errResp.Message
			}
			return nil, fmt.Errorf("error API (%d): %s", resp.StatusCode, msg)
		}
		return nil, fmt.Errorf("error API (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var successResp LicenseResponse
	if err := json.Unmarshal(bodyBytes, &successResp); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta exitosa: %w", err)
	}

	// Validar la firma del token recibido inmediatamente
	if err := crypto.VerifyLicenseToken(successResp.Token); err != nil {
		return nil, fmt.Errorf("seguridad: el token recibido no es auténtico: %v", err)
	}

	return &successResp, nil
}

// SendPDFReport envía el PDF generado por correo electrónico a través del backend
func (s *CloudService) SendPDFReport(email string, pdfContent []byte, filename string) error {
	// Crear buffer para el multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 1. Campo 'email'
	if err := writer.WriteField("email", email); err != nil {
		return fmt.Errorf("error escribiendo campo email: %w", err)
	}

	// 2. Campo 'archivo' (PDF Binario)
	// CreateFormFile usa el nombre del campo ("archivo") y el nombre del fichero (filename)
	part, err := writer.CreateFormFile("archivo", filename)
	if err != nil {
		return fmt.Errorf("error creando form file: %w", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(pdfContent)); err != nil {
		return fmt.Errorf("error copiando contenido pdf: %w", err)
	}

	// Cerrar el writer para finalizar el boundary
	if err := writer.Close(); err != nil {
		return fmt.Errorf("error cerrando multipart writer: %w", err)
	}

	// Enviar Petición
	url := fmt.Sprintf("%s/api/v1/mail/send-pdf", s.BaseURL)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("error creando request multipart: %w", err)
	}

	// ESTABLECER EL CONTENT-TYPE CON EL BOUNDARY CORRECTO
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando reporte: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error enviando reporte (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
