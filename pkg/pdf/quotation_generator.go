package pdf

import (
	"os"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/props"

	"kushkiv2/internal/db"
)

// GenerarCotizacionPDF crea el PDF de una cotización
func GenerarCotizacionPDF(cotizacion db.Quotation, items []db.QuotationItem, configEmisor db.EmisorConfig) ([]byte, error) {
	cfg := config.NewBuilder().
		WithPageSize(pagesize.A4).
		WithLeftMargin(15).
		WithTopMargin(15).
		WithRightMargin(15).
		WithBottomMargin(15).
		Build()

	m := maroto.New(cfg)

	// =========================================================================
	// 1. CABECERA
	// =========================================================================

	// --- FILA 1: Logo y Título ---
	colLogo := col.New(4)
	if configEmisor.LogoPath != "" {
		if _, err := os.Stat(configEmisor.LogoPath); err == nil {
			colLogo.Add(image.NewFromFile(configEmisor.LogoPath, props.Rect{Center: false, Percent: 100, Left: 0}))
		}
	} else {
		colLogo.Add(text.New("KUSHKI APP", props.Text{Style: fontstyle.Bold, Color: brandPrimary, Size: 16}))
	}

	m.AddRow(20,
		colLogo,
		col.New(8).Add(
			text.New("COTIZACIÓN", props.Text{Size: 14, Style: fontstyle.Bold, Align: align.Right, Color: brandPrimary, Top: 0}),
			text.New("No. "+cotizacion.Secuencial, props.Text{Size: 11, Align: align.Right, Top: 8, Style: fontstyle.Bold}),
		),
	)

	m.AddRow(5, col.New(12))

	// --- FILA 2: Datos Emisor ---
	m.AddRow(15,
		col.New(6).Add(
			text.New(configEmisor.RazonSocial, props.Text{Size: 11, Style: fontstyle.Bold, Color: textDark, Top: 0}),
			text.New("RUC: "+configEmisor.RUC, props.Text{Size: 9, Color: textGray, Top: 6}),
			text.New(configEmisor.Direccion, props.Text{Size: 8, Color: textGray, Top: 10}),
		),
		col.New(6).Add(
			text.New("FECHA DE EMISIÓN: "+cotizacion.FechaEmision.Format("02/01/2006"), props.Text{Size: 9, Align: align.Right, Top: 0, Style: fontstyle.Bold}),
			text.New("VÁLIDO HASTA: "+cotizacion.FechaEmision.AddDate(0, 0, 15).Format("02/01/2006"), props.Text{Size: 8, Align: align.Right, Top: 6, Color: textGray}), // Validez 15 días hardcoded por ahora
		),
	)

	m.AddRow(10, col.New(12))

	// =========================================================================
	// 2. DATOS CLIENTE
	// =========================================================================

	m.AddRow(8,
		col.New(12).WithStyle(&props.Cell{BackgroundColor: brandPrimary}).Add(
			text.New("CLIENTE", props.Text{Size: 8, Style: fontstyle.Bold, Color: textWhite, Top: 2, Left: 2}),
		),
	)

	m.AddRow(18,
		col.New(7).WithStyle(&props.Cell{BackgroundColor: brandLight}).Add(
			text.New("Nombre: "+cotizacion.ClienteNombre, props.Text{Size: 9, Style: fontstyle.Bold, Left: 2, Top: 2}),
			text.New("Identificación: "+cotizacion.ClienteID, props.Text{Size: 8, Left: 2, Top: 8}),
		),
		col.New(5).WithStyle(&props.Cell{BackgroundColor: brandLight}).Add(
			text.New("Dirección: "+cotizacion.ClienteDireccion, props.Text{Size: 8, Align: align.Right, Right: 2, Top: 2}),
			text.New("Teléfono: "+cotizacion.ClienteTelefono, props.Text{Size: 8, Align: align.Right, Right: 2, Top: 8}),
		),
	)

	m.AddRow(5, col.New(12).Add(line.New(props.Line{Color: brandPrimary, Thickness: 0.5})))

	// =========================================================================
	// 3. TABLA DETALLES
	// =========================================================================

	m.AddRow(9,
		text.NewCol(2, "CÓDIGO", props.Text{Style: fontstyle.Bold, Size: 8, Align: align.Left, Color: textWhite, Top: 1.5, Left: 2}),
		text.NewCol(1, "CANT.", props.Text{Style: fontstyle.Bold, Size: 8, Align: align.Center, Color: textWhite, Top: 1.5}),
		text.NewCol(5, "DESCRIPCIÓN", props.Text{Style: fontstyle.Bold, Size: 8, Align: align.Left, Color: textWhite, Top: 1.5, Left: 2}),
		text.NewCol(2, "P. UNIT", props.Text{Style: fontstyle.Bold, Size: 8, Align: align.Right, Color: textWhite, Top: 1.5, Right: 2}),
		text.NewCol(2, "TOTAL", props.Text{Style: fontstyle.Bold, Size: 8, Align: align.Right, Color: textWhite, Top: 1.5, Right: 2}),
	).WithStyle(&props.Cell{BackgroundColor: brandPrimary})

	for _, item := range items {
		m.AddRow(8,
			text.NewCol(2, item.ProductoSKU, props.Text{Size: 8, Align: align.Left, Top: 2, Color: textGray, Left: 2}),
			text.NewCol(1, fmtMoney(item.Cantidad), props.Text{Size: 8, Align: align.Center, Top: 2}),
			text.NewCol(5, item.Nombre, props.Text{Size: 8, Align: align.Left, Top: 2, Left: 2}),
			text.NewCol(2, fmtMoney(item.PrecioUnitario), props.Text{Size: 8, Align: align.Right, Top: 2, Right: 2}),
			text.NewCol(2, fmtMoney(item.Subtotal), props.Text{Size: 8, Align: align.Right, Top: 2, Style: fontstyle.Bold, Right: 2}),
		)
		m.AddRow(1, col.New(12).Add(line.New(props.Line{Color: &props.Color{Red: 240, Green: 240, Blue: 240}})))
	}

	m.AddRow(10, col.New(12))

	// =========================================================================
	// 4. TOTALES
	// =========================================================================

	// Cálculos simples (ya vienen en el struct pero por si acaso)
	subtotal := cotizacion.Subtotal15 + cotizacion.Subtotal0
	iva := cotizacion.IVA
	total := cotizacion.Total

	type totalRow struct {
		label string
		val   string
	}
	totals := []totalRow{
		{"Subtotal", fmtMoney(subtotal)},
		{"IVA", fmtMoney(iva)},
	}

	// Observación
	colIzq := col.New(7)
	if cotizacion.Observacion != "" {
		colIzq.Add(text.New("Observaciones:", props.Text{Style: fontstyle.Bold, Size: 8}))
		colIzq.Add(text.New(cotizacion.Observacion, props.Text{Size: 8, Top: 5, Color: textGray}))
	}

	for i, t := range totals {
		colDerLbl := col.New(3)
		colDerVal := col.New(2)

		colDerLbl.Add(text.New(t.label, props.Text{Size: 8, Align: align.Right, Color: textGray, Right: 2}))
		colDerVal.Add(text.New(t.val, props.Text{Size: 8, Align: align.Right, Color: textDark}))

		// Truco para solo renderizar la columna izquierda una vez
		if i == 0 {
			m.AddRow(5, colIzq, colDerLbl, colDerVal)
		} else {
			m.AddRow(5, col.New(7), colDerLbl, colDerVal)
		}
	}

	m.AddRow(5, col.New(12))

	// --- Total Final Verde ---
	m.AddRow(12,
		col.New(7),
		col.New(5).WithStyle(&props.Cell{BackgroundColor: brandPrimary}).Add(
			text.New("TOTAL", props.Text{Size: 10, Style: fontstyle.Bold, Color: textWhite, Align: align.Left, Left: 4, Top: 3.5}),
			text.New("$ "+fmtMoney(total), props.Text{Size: 12, Style: fontstyle.Bold, Color: textWhite, Align: align.Right, Right: 4, Top: 3}),
		),
	)

	// =========================================================================
	// 5. PIE DE PÁGINA
	// =========================================================================
	m.AddRow(15, col.New(12))
	m.AddRow(10,
		text.NewCol(12, "Documento de Cotización - Sin validez tributaria", props.Text{
			Size:  7,
			Align: align.Center,
			Color: &props.Color{Red: 150, Green: 150, Blue: 150},
			Style: fontstyle.Italic,
		}),
	)

	document, err := m.Generate()
	if err != nil {
		return nil, err
	}

	return document.GetBytes(), nil
}
