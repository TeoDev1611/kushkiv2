package xml

import (
	"encoding/xml"
)

// FacturaXML representa la estructura raíz de una factura electrónica del SRI.
type FacturaXML struct {
	XMLName        xml.Name       `xml:"factura"`
	ID             string         `xml:"id,attr"`
	Version        string         `xml:"version,attr"`
	InfoTributaria InfoTributaria `xml:"infoTributaria"`
	InfoFactura    InfoFactura    `xml:"infoFactura"`
	Detalles       []Detalle      `xml:"detalles>detalle"`
	InfoAdicional  []CampoAdicional `xml:"infoAdicional>campoAdicional,omitempty"`
}

type InfoTributaria struct {
	Ambiente          string `xml:"ambiente"`
	TipoEmision       string `xml:"tipoEmision"`
	RazonSocial       string `xml:"razonSocial"`
	NombreComercial   string `xml:"nombreComercial,omitempty"`
	Ruc               string `xml:"ruc"`
	ClaveAcceso       string `xml:"claveAcceso"`
	CodDoc            string `xml:"codDoc"`
	Estab             string `xml:"estab"`      // Siempre 3 dígitos
	PtoEmi            string `xml:"ptoEmi"`     // Siempre 3 dígitos
	Secuencial        string `xml:"secuencial"` // Siempre 9 dígitos
	DirMatriz         string `xml:"dirMatriz"`
	AgenteRetencion   string `xml:"agenteRetencion,omitempty"`
	ContribuyenteRimpe string `xml:"contribuyenteRimpe,omitempty"`
}

type InfoFactura struct {
	FechaEmision                string          `xml:"fechaEmision"`
	DirEstablecimiento          string          `xml:"dirEstablecimiento,omitempty"`
	ContribuyenteEspecial       string          `xml:"contribuyenteEspecial,omitempty"`
	ObligadoContabilidad        string          `xml:"obligadoContabilidad,omitempty"`
	TipoIdentificacionComprador string          `xml:"tipoIdentificacionComprador"`
	GuiaRemision                string          `xml:"guiaRemision,omitempty"`
	RazonSocialComprador        string          `xml:"razonSocialComprador"`
	IdentificacionComprador     string          `xml:"identificacionComprador"`
	DireccionComprador          string          `xml:"direccionComprador,omitempty"`
	TotalSinImpuestos           float64         `xml:"totalSinImpuestos"`
	TotalDescuento              float64         `xml:"totalDescuento"`
	TotalConImpuestos           []TotalImpuesto `xml:"totalConImpuestos>totalImpuesto"`
	Propina                     float64         `xml:"propina"`
	ImporteTotal                float64         `xml:"importeTotal"`
	Moneda                      string          `xml:"moneda"`
	Pagos                       []Pago          `xml:"pagos>pago"`
}

type TotalImpuesto struct {
	Codigo           string  `xml:"codigo"`
	CodigoPorcentaje string  `xml:"codigoPorcentaje"`
	BaseImponible    float64 `xml:"baseImponible"`
	Valor            float64 `xml:"valor"`
}

type Detalle struct {
	CodigoPrincipal        string     `xml:"codigoPrincipal"`
	CodigoAuxiliar         string     `xml:"codigoAuxiliar,omitempty"`
	Descripcion            string     `xml:"descripcion"`
	UnidadMedida           string     `xml:"unidadMedida,omitempty"`
	Cantidad               float64    `xml:"cantidad"`
	PrecioUnitario         float64    `xml:"precioUnitario"`
	Descuento              float64    `xml:"descuento"`
	PrecioTotalSinImpuesto float64    `xml:"precioTotalSinImpuesto"`
	DetallesAdicionales    *DetallesAdicionales `xml:"detallesAdicionales,omitempty"`
	Impuestos              []Impuesto `xml:"impuestos>impuesto"`
}

type DetallesAdicionales struct {
	DetAdicional []DetAdicional `xml:"detAdicional"`
}

type DetAdicional struct {
	Nombre string `xml:"nombre,attr"`
	Valor  string `xml:"valor,attr"`
}

type Impuesto struct {
	Codigo     string  `xml:"codigo"`
	CodigoPorcentaje string `xml:"codigoPorcentaje"`
	Tarifa     float64 `xml:"tarifa"`
	BaseImponible float64 `xml:"baseImponible"`
	Valor      float64 `xml:"valor"`
}

type Pago struct {
	FormaPago string  `xml:"formaPago"`
	Total     float64 `xml:"total"`
	Plazo     string  `xml:"plazo,omitempty"`
	UnidadTiempo string `xml:"unidadTiempo,omitempty"`
}

type CampoAdicional struct {
	Nombre string `xml:"nombre,attr"`
	Value  string `xml:",chardata"`
}
