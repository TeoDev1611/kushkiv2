package db

import (
	"time"

	"gorm.io/gorm"
)

// EmisorConfig almacena la configuración tributaria del emisor.
type EmisorConfig struct {
	gorm.Model
	RUC             string `gorm:"primaryKey;uniqueIndex"`
	RazonSocial     string
	NombreComercial string
	Direccion       string
	// Seguridad / Licenciamiento
	LicenseKey      string // Clave ingresada por el usuario (Ej: PRO-2026-X)
	LicenseToken    string // Token JWT devuelto por el servidor (para autenticar peticiones futuras)
	
	// Firma Electrónica
	P12Path         string
	P12Password     string // Se guardará cifrada
	
	// Configuración SRI
	Ambiente        int    // 1: Pruebas, 2: Producción
	Estab           string // Ej: '001'
	PtoEmi          string // Ej: '001'
	Obligado        bool   // Obligado a llevar contabilidad
	
	// Configuración SMTP (Correo Local)
	SMTPHost        string
	SMTPPort        int
	SMTPUser        string
	SMTPPassword    string

	// Archivos
	StoragePath     string // Ruta base para guardar archivos
}

// Client representa a los clientes/compradores.
type Client struct {
	ID        string `gorm:"primaryKey"`
	TipoID    string
	Nombre    string `gorm:"index"`
	Direccion string
	Email     string
	Telefono  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Factura representa un comprobante electrónico en la base de datos.
type Factura struct {
	ClaveAcceso     string `gorm:"primaryKey;size:49"`
	Secuencial      string `gorm:"size:9"`
	FechaEmision    time.Time
	ClienteID       string
	Total           float64
	EstadoSRI       string
	XMLFirmado      []byte `gorm:"type:blob"`
	PDFRIDE         []byte `gorm:"type:blob"`
	MensajeError    string
	Subtotal15      float64
	Subtotal0       float64
	IVA             float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// FacturaItem almacena el detalle de cada producto.
type FacturaItem struct {
	ID               uint    `gorm:"primaryKey"`
	FacturaClave     string  `gorm:"index"`
	ProductoSKU      string  `gorm:"index"`
	Nombre           string
	Cantidad         float64
	PrecioUnitario   float64
	Subtotal         float64
	PorcentajeIVA    float64
	CreatedAt        time.Time
}

// Product representa el inventario.
type Product struct {
	SKU           string `gorm:"primaryKey"`
	Name          string `gorm:"index"`
	Price         float64
	Stock         int
	TaxCode       int
	TaxPercentage int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// EmailQueue se mantiene temporalmente por si quedan registros antiguos, 
// pero ya no será procesada por el sistema legacy.
type EmailQueue struct {
	ID         uint      `gorm:"primaryKey"`
	To         string    `gorm:"index"`
	Subject    string
	Body       string    `gorm:"type:text"`
	Attachment []byte    `gorm:"type:blob"`
	AttachName string
	Status     string
	RetryCount int
	LastError  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// --- DTOs ---

type EmisorConfigDTO struct {
	RUC             string `json:"RUC"`
	RazonSocial     string `json:"RazonSocial"`
	NombreComercial string `json:"NombreComercial"`
	Direccion       string `json:"Direccion"`
	P12Path         string `json:"P12Path"`
	P12Password     string `json:"P12Password"`
	Ambiente        int    `json:"Ambiente"`
	Estab           string `json:"Estab"`
	PtoEmi          string `json:"PtoEmi"`
	Obligado        bool   `json:"Obligado"`
	StoragePath     string `json:"StoragePath"`
	
	// SMTP
	SMTPHost        string `json:"SMTPHost"`
	SMTPPort        int    `json:"SMTPPort"`
	SMTPUser        string `json:"SMTPUser"`
	SMTPPassword    string `json:"SMTPPassword"`
}

type ClientDTO struct {
	ID        string `json:"ID"`
	TipoID    string `json:"TipoID"`
	Nombre    string `json:"Nombre"`
	Direccion string `json:"Direccion"`
	Email     string `json:"Email"`
	Telefono  string `json:"Telefono"`
}

type ProductDTO struct {
	SKU           string  `json:"SKU"`
	Name          string  `json:"Name"`
	Price         float64 `json:"Price"`
	Stock         int     `json:"Stock"`
	TaxCode       int     `json:"TaxCode"`
	TaxPercentage int     `json:"TaxPercentage"`
}

type FacturaDTO struct {
	Secuencial       string        `json:"secuencial"`
	ClienteID        string        `json:"clienteID"`
	ClienteNombre    string        `json:"clienteNombre"`
	ClienteDireccion string        `json:"clienteDireccion"`
	ClienteEmail     string        `json:"clienteEmail"`
	ClienteTelefono  string        `json:"clienteTelefono"`
	FormaPago        string        `json:"formaPago"`
	Items            []InvoiceItem `json:"items"`
	ClaveAcceso      string
}

type FacturaResumenDTO struct {
	ClaveAcceso   string  `json:"claveAcceso"`
	Secuencial    string  `json:"secuencial"`
	Fecha         string  `json:"fecha"`
	Cliente       string  `json:"cliente"`
	Total         float64 `json:"total"`
	Estado        string  `json:"estado"`
	TienePDF      bool    `json:"tienePDF"`
}

type InvoiceItem struct {
	Codigo        string  `json:"codigo"`
	Nombre        string  `json:"nombre"`
	Cantidad      float64 `json:"cantidad"`
	Precio        float64 `json:"precio"`
	CodigoIVA     string  `json:"codigoIVA"`
	PorcentajeIVA float64 `json:"porcentajeIVA"`
}
