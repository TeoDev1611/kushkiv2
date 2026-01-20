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
	NombreComercial string // Nuevo campo
	Direccion       string // Dirección Matriz
	P12Path         string
	P12Password     string // Se guardará cifrada
	Ambiente        int    // 1: Pruebas, 2: Producción
	Estab           string // Ej: '001'
	PtoEmi          string // Ej: '001'
	Obligado        bool   // Obligado a llevar contabilidad
	SMTPHost        string
	SMTPUser        string
	SMTPPass        string
	StoragePath     string // Ruta base para guardar archivos (ej: C:/Facturas)
}

// Client representa a los clientes/compradores.
type Client struct {
	ID        string `gorm:"primaryKey"` // RUC o Cédula
	TipoID    string // 04: RUC, 05: Cédula, 06: Pasaporte, 07: Consumidor Final
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
	EstadoSRI       string // PENDIENTE, RECIBIDO, AUTORIZADO, RECHAZADO
	XMLFirmado      []byte `gorm:"type:blob"`
	PDFRIDE         []byte `gorm:"type:blob"` // PDF Generado
	MensajeError    string
	Subtotal15      float64
	Subtotal0       float64
	IVA             float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// FacturaItem almacena el detalle de cada producto en una factura (para reportería rápida).
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

// Product representa el inventario de productos.
type Product struct {
	SKU           string `gorm:"primaryKey"`
	Name          string `gorm:"index"`
	Price         float64
	Stock         int
	TaxCode       int // 2 para IVA
	TaxPercentage int // 0, 15, etc.
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// EmailQueue maneja la cola de envíos de correos electrónicos.
type EmailQueue struct {
	ID         uint      `gorm:"primaryKey"`
	To         string    `gorm:"index"`
	Subject    string
	Body       string    `gorm:"type:text"`
	Attachment []byte    `gorm:"type:blob"` // El PDF de la factura
	AttachName string    // Nombre del archivo adjunto
	Status     string    `gorm:"index;default:'PENDIENTE'"` // PENDIENTE, ENVIADO, ERROR
	RetryCount int       `gorm:"default:0"`
	LastError  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// --- DTOs para transferencia de datos (Sin tipos complejos) ---

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
	SMTPHost        string `json:"SMTPHost"`
	SMTPUser        string `json:"SMTPUser"`
	SMTPPass        string `json:"SMTPPass"`
	StoragePath     string `json:"StoragePath"`
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
	FormaPago        string        `json:"formaPago"` // "01", "16", "19", "20"
	Items            []InvoiceItem `json:"items"`
	ClaveAcceso      string        // Output
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
	// Se reemplaza TieneIVA por CodigoImpuesto explícito para soportar 5%, 15%, etc.
	CodigoIVA     string  `json:"codigoIVA"` // "0", "2", "4" (15%), "5" (5%)
	PorcentajeIVA float64 `json:"porcentajeIVA"` // 0.0, 12.0, 15.0, 5.0
}