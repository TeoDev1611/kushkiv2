package pdf

import (
	"fmt"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

// Colores Globales
var (
	// Modern / Emerald Theme
	colorEmeraldPrimary = &props.Color{Red: 52, Green: 211, Blue: 153}
	colorEmeraldLight   = &props.Color{Red: 236, Green: 253, Blue: 245}

	// Corporate / Grayscale Theme
	colorBlack      = &props.Color{Red: 0, Green: 0, Blue: 0}
	colorDarkGray   = &props.Color{Red: 50, Green: 50, Blue: 50}
	colorGray       = &props.Color{Red: 100, Green: 100, Blue: 100}
	colorLightGray  = &props.Color{Red: 200, Green: 200, Blue: 200}
	colorWhite      = &props.Color{Red: 255, Green: 255, Blue: 255}
	colorBackground = &props.Color{Red: 245, Green: 245, Blue: 245}
)

// Utilidades compartidas
func fmtMoney(val float64) string {
	return fmt.Sprintf("%.2f", val)
}

func getAmbienteText(codigo string) string {
	if codigo == "2" {
		return "PRODUCCIÓN"
	}
	return "PRUEBAS"
}

// Footer común para cumplir con el requerimiento de marca
const footerText = "Documento generado por Kushki App - Tecnología hecha en Ecuador"
