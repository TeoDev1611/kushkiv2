package sri

import "encoding/xml"

// Estructuras para parsear la respuesta de Recepción (ValidarComprobante)

type RespuestaRecepcion struct {
	XMLName xml.Name `xml:"RespuestaSolicitud"`
	Estado  string   `xml:"estado"` // RECIBIDA | DEVUELTA
	Comprobantes struct {
		Comprobante []ComprobanteRecepcion `xml:"comprobante"`
	} `xml:"comprobantes"`
}

type ComprobanteRecepcion struct {
	ClaveAcceso string `xml:"claveAcceso"`
	Mensajes    struct {
		Mensaje []MensajeSRI `xml:"mensaje"`
	} `xml:"mensajes"`
}

type MensajeSRI struct {
	Identificador string `xml:"identificador"`
	Mensaje       string `xml:"mensaje"`
	InformacionAdicional string `xml:"informacionAdicional"`
	Tipo          string `xml:"tipo"` // ERROR, ADVERTENCIA
}

// Estructuras para parsear la respuesta de Autorización

type RespuestaAutorizacion struct {
	XMLName xml.Name `xml:"RespuestaAutorizacionComprobante"`
	Autorizaciones struct {
		Autorizacion []Autorizacion `xml:"autorizacion"`
	} `xml:"autorizaciones"`
}

type Autorizacion struct {
	Estado          string `xml:"estado"` // AUTORIZADO | NO AUTORIZADO
	NumeroAutorizacion string `xml:"numeroAutorizacion"`
	FechaAutorizacion  string `xml:"fechaAutorizacion"` // Formato ISO 8601
	Ambiente        string `xml:"ambiente"`
	Comprobante     string `xml:"comprobante"` // CDATA con el XML original
	Mensajes        struct {
		Mensaje []MensajeSRI `xml:"mensaje"`
	} `xml:"mensajes"`
}
