package xml

import (
	"encoding/xml"
	"fmt"
)

// GenerateXML convierte el struct FacturaXML en un arreglo de bytes XML.
func GenerateXML(factura *FacturaXML) ([]byte, error) {
	output, err := xml.MarshalIndent(factura, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error al serializar XML: %v", err)
	}

	header := []byte(xml.Header)
	return append(header, output...), nil
}
