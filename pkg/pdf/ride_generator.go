package pdf

import (
	"fmt"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	
	srixml "kushkiv2/pkg/xml"
)

// GenerarRIDE crea el PDF de la factura usando el tema seleccionado.
// themeName puede ser: "modern" (default), "minimal", "corporate".
func GenerarRIDE(factura srixml.FacturaXML, logoPath string, themeName string) ([]byte, error) {
	// 1. Configuración Base de Maroto
	cfg := config.NewBuilder().
		WithPageSize(pagesize.A4).
		WithLeftMargin(15).
		WithTopMargin(15).
		WithRightMargin(15).
		WithBottomMargin(15).
		Build()

	m := maroto.New(cfg)

	// 2. Selección de Estrategia (Factory)
	var theme InvoiceTheme

	switch themeName {
	case "minimal":
		theme = &MinimalTheme{}
	case "corporate":
		theme = &CorporateTheme{}
	default:
		// Default a Modern / Emerald si está vacío o es desconocido
		theme = &ModernTheme{}
	}

	// 3. Construcción del PDF
	theme.Build(m, factura, logoPath)

	// 4. Generación de bytes
	document, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("error generando PDF: %w", err)
	}

	return document.GetBytes(), nil
}

// StrPtr Helper para punteros de string (usado en tests o utilidades externas)
func StrPtr(s string) *string {
	return &s
}