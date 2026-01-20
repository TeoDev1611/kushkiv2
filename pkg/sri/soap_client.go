package sri

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	SRIRecepciónPruebas    = "https://celcer.sri.gob.ec/comprobantes-electronicos-ws/RecepcionComprobantesOffline?wsdl"
	SRIAutorizaciónPruebas = "https://celcer.sri.gob.ec/comprobantes-electronicos-ws/AutorizacionComprobantesOffline?wsdl"
)

type SRIClient struct {
	Client *http.Client
}

func NewSRIClient() *SRIClient {
	return &SRIClient{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Estructuras internas para desempaquetar el Envelope SOAP
type soapResponseRecepcion struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		ValidarResponse struct {
			Respuesta RespuestaRecepcion `xml:"RespuestaSolicitud"`
		} `xml:"validarComprobanteResponse"`
	} `xml:"Body"`
}

type soapResponseAutorizacion struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		AutorizacionResponse struct {
			Respuesta RespuestaAutorizacion `xml:"RespuestaAutorizacionComprobante"`
		} `xml:"autorizacionComprobanteResponse"`
	} `xml:"Body"`
}

// NetworkError indica un fallo de conexión, no un rechazo del SRI.
type NetworkError struct {
	Msg string
}

func (e *NetworkError) Error() string {
	return e.Msg
}

// EnviarComprobante envía el XML firmado al SRI y devuelve la respuesta parseada.
func (s *SRIClient) EnviarComprobante(xmlFirmado []byte) (*RespuestaRecepcion, error) {
	xmlBase64 := base64.StdEncoding.EncodeToString(xmlFirmado)
	
	soapEnvelope := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ecua="http://ec.gob.sri.ws.recepcion">
   <soapenv:Header/>
   <soapenv:Body>
      <ecua:validarComprobante>
         <xml>%s</xml>
      </ecua:validarComprobante>
   </soapenv:Body>
</soapenv:Envelope>`, xmlBase64)

	fmt.Printf("\n--- SRI RECEPCIÓN REQUEST ---\n%s\n-----------------------------\n", soapEnvelope)

	respBody, err := s.doRequest(SRIRecepciónPruebas, soapEnvelope)
	if err != nil {
		return nil, err
	}
	
	fmt.Printf("\n--- SRI RECEPCIÓN RESPONSE ---\n%s\n------------------------------\n", string(respBody))

	var envelope soapResponseRecepcion
	if err := xml.Unmarshal(respBody, &envelope); err != nil {
		return nil, fmt.Errorf("error parsing recepcion response: %v | raw: %s", err, string(respBody))
	}

	return &envelope.Body.ValidarResponse.Respuesta, nil
}

// AutorizarComprobante consulta el estado y devuelve la respuesta parseada.
func (s *SRIClient) AutorizarComprobante(claveAcceso string) (*RespuestaAutorizacion, error) {
	soapEnvelope := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ecua="http://ec.gob.sri.ws.autorizacion">
   <soapenv:Header/>
   <soapenv:Body>
      <ecua:autorizacionComprobante>
         <claveAccesoComprobante>%s</claveAccesoComprobante>
      </ecua:autorizacionComprobante>
   </soapenv:Body>
</soapenv:Envelope>`, claveAcceso)

	fmt.Printf("\n--- SRI AUTORIZACIÓN REQUEST ---\n%s\n--------------------------------\n", soapEnvelope)

	respBody, err := s.doRequest(SRIAutorizaciónPruebas, soapEnvelope)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n--- SRI AUTORIZACIÓN RESPONSE ---\n%s\n---------------------------------\n", string(respBody))

	var envelope soapResponseAutorizacion
	if err := xml.Unmarshal(respBody, &envelope); err != nil {
		return nil, fmt.Errorf("error parsing autorizacion response: %v | raw: %s", err, string(respBody))
	}

	return &envelope.Body.AutorizacionResponse.Respuesta, nil
}

// CheckConnectivity verifica si hay acceso al WSDL.
func (s *SRIClient) CheckConnectivity() bool {
	client := http.Client{Timeout: 5 * time.Second}
	_, err := client.Get(SRIRecepciónPruebas)
	return err == nil
}

func (s *SRIClient) doRequest(url, body string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")

	resp, err := s.Client.Do(req)
	if err != nil {
		// Retornar error de red específico
		return nil, &NetworkError{Msg: fmt.Sprintf("SRI Offline: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return nil, &NetworkError{Msg: fmt.Sprintf("SRI Server Error: %d", resp.StatusCode)}
	}

	return io.ReadAll(resp.Body)
}