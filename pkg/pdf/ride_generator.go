package pdf

import (
	"fmt"
	"os"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/barcode"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/props"

	srixml "kushkiv2/pkg/xml"
)

// COLORES ACTUALIZADOS A TEMA ESMERALDA (#34d399)
var (
	// #34d399 -> RGB(52, 211, 153)
	brandPrimary = &props.Color{Red: 52, Green: 211, Blue: 153}
	// Un verde muy suave para fondos (#ecfdf5)
	brandLight = &props.Color{Red: 236, Green: 253, Blue: 245}

	textGray  = &props.Color{Red: 100, Green: 100, Blue: 100}
	textDark  = &props.Color{Red: 50, Green: 50, Blue: 50}
	textWhite = &props.Color{Red: 255, Green: 255, Blue: 255}
)

func fmtMoney(val float64) string {
	return fmt.Sprintf("%.2f", val)
}

func StrPtr(s string) *string {
	return &s
}

func GenerarRIDE(factura srixml.FacturaXML, logoPath string) ([]byte, error) {
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

	ambienteStr := "PRUEBAS"
	if factura.InfoTributaria.Ambiente == "2" {
		ambienteStr = "PRODUCCIÓN"
	}

	// --- FILA 1: Logo y Título ---
	colLogo := col.New(4)
	if logoPath != "" {
		if _, err := os.Stat(logoPath); err == nil {
			colLogo.Add(image.NewFromFile(logoPath, props.Rect{Center: false, Percent: 100, Left: 0}))
		}
	} else {
		colLogo.Add(text.New("KUSHKI APP", props.Text{Style: fontstyle.Bold, Color: brandPrimary, Size: 16}))
	}

	m.AddRow(20,
		colLogo,
		col.New(8).Add(
			text.New("FACTURA ELECTRÓNICA", props.Text{Size: 14, Style: fontstyle.Bold, Align: align.Right, Color: brandPrimary, Top: 0}),
			text.New("No. "+fmt.Sprintf("%s-%s-%s", factura.InfoTributaria.Estab, factura.InfoTributaria.PtoEmi, factura.InfoTributaria.Secuencial), props.Text{Size: 11, Align: align.Right, Top: 8, Style: fontstyle.Bold}),
		),
	)

	m.AddRow(5, col.New(12))

	// --- FILA 2: Razón Social y Clave Acceso ---
	m.AddRow(15,
		col.New(6).Add(
			text.New(factura.InfoTributaria.RazonSocial, props.Text{Size: 11, Style: fontstyle.Bold, Color: textDark, Top: 0}),
			text.New("RUC: "+factura.InfoTributaria.Ruc, props.Text{Size: 9, Color: textGray, Top: 6}),
		),
		col.New(6).WithStyle(&props.Cell{BackgroundColor: brandLight}).Add(
			text.New("CLAVE DE ACCESO:", props.Text{Size: 7, Style: fontstyle.Bold, Align: align.Center, Color: textGray, Top: 2}),
			text.New(factura.InfoTributaria.ClaveAcceso, props.Text{Size: 7, Align: align.Center, Top: 6, Family: "Courier"}),
		),
	)

	// --- FILA 3: Direcciones y Fechas ---
	m.AddRow(15,
		col.New(6).Add(
			text.New("Matriz: "+factura.InfoTributaria.DirMatriz, props.Text{Size: 7, Color: textGray, Top: 0}),
			text.New("Sucursal: "+factura.InfoFactura.DirEstablecimiento, props.Text{Size: 7, Color: textGray, Top: 4}),
		),
		col.New(6).Add(
			text.New("FECHA DE EMISIÓN: "+factura.InfoFactura.FechaEmision, props.Text{Size: 8, Align: align.Right, Top: 0, Style: fontstyle.Bold}),
			text.New("AMBIENTE: "+ambienteStr, props.Text{Size: 8, Align: align.Right, Top: 4}),
		),
	)

	// --- FILA 4: CÓDIGO DE BARRAS ---
	m.AddRow(12,
		col.New(6), // Espacio vacío a la izquierda
		col.New(6).Add(
			code.NewBar(factura.InfoTributaria.ClaveAcceso, props.Barcode{
				Type:       barcode.Code128,
				Proportion: props.Proportion{Width: 20, Height: 5},
				Center:     true,
			}),
		),
	)

	m.AddRow(5, col.New(12))

	// =========================================================================
	// 2. DATOS CLIENTE (Diseño de Tarjeta)
	// =========================================================================

	m.AddRow(8,
		col.New(12).WithStyle(&props.Cell{BackgroundColor: brandPrimary}).Add(
			text.New("INFORMACIÓN DEL RECEPTOR", props.Text{Size: 8, Style: fontstyle.Bold, Color: textWhite, Top: 2, Left: 2}),
		),
	)

	phone := "-"
	for _, info := range factura.InfoAdicional {
		if info.Nombre == "Telefono" {
			phone = info.Value
			break
		}
	}

	m.AddRow(20,
		col.New(7).WithStyle(&props.Cell{BackgroundColor: brandLight}).Add(
			text.New("Razón Social: ", props.Text{Size: 8, Style: fontstyle.Bold, Left: 2, Top: 2}),
			text.New(factura.InfoFactura.RazonSocialComprador, props.Text{Size: 9, Left: 2, Top: 6}),
			text.New("Identificación: "+factura.InfoFactura.IdentificacionComprador, props.Text{Size: 9, Left: 2, Top: 12}),
		),
		col.New(5).WithStyle(&props.Cell{BackgroundColor: brandLight}).Add(
			text.New("Dirección: "+factura.InfoFactura.DireccionComprador, props.Text{Size: 8, Align: align.Right, Right: 2, Top: 2}),
			text.New("Teléfono: "+phone, props.Text{Size: 8, Align: align.Right, Right: 2, Top: 8}),
		),
	)

	m.AddRow(5, col.New(12).Add(line.New(props.Line{Color: brandPrimary, Thickness: 0.5})))

	// =========================================================================
	// 3. TABLA DETALLES (Diseño Moderno)
	// =========================================================================

	m.AddRow(9,
		text.NewCol(2, "CÓDIGO", props.Text{Style: fontstyle.Bold, Size: 8, Align: align.Left, Color: textWhite, Top: 1.5, Left: 2}),
		text.NewCol(1, "CANT.", props.Text{Style: fontstyle.Bold, Size: 8, Align: align.Center, Color: textWhite, Top: 1.5}),
		text.NewCol(5, "DESCRIPCIÓN", props.Text{Style: fontstyle.Bold, Size: 8, Align: align.Left, Color: textWhite, Top: 1.5, Left: 2}),
		text.NewCol(2, "P. UNIT", props.Text{Style: fontstyle.Bold, Size: 8, Align: align.Right, Color: textWhite, Top: 1.5, Right: 2}),
		text.NewCol(2, "TOTAL", props.Text{Style: fontstyle.Bold, Size: 8, Align: align.Right, Color: textWhite, Top: 1.5, Right: 2}),
	).WithStyle(&props.Cell{BackgroundColor: brandPrimary})

	for _, item := range factura.Detalles {
		m.AddRow(8,
			text.NewCol(2, item.CodigoPrincipal, props.Text{Size: 8, Align: align.Left, Top: 2, Color: textGray, Left: 2}),
			text.NewCol(1, fmtMoney(item.Cantidad), props.Text{Size: 8, Align: align.Center, Top: 2}),
			text.NewCol(5, item.Descripcion, props.Text{Size: 8, Align: align.Left, Top: 2, Left: 2}),
			text.NewCol(2, fmtMoney(item.PrecioUnitario), props.Text{Size: 8, Align: align.Right, Top: 2, Right: 2}),
			text.NewCol(2, fmtMoney(item.PrecioTotalSinImpuesto), props.Text{Size: 8, Align: align.Right, Top: 2, Style: fontstyle.Bold, Right: 2}),
		)
		m.AddRow(1, col.New(12).Add(line.New(props.Line{Color: &props.Color{Red: 240, Green: 240, Blue: 240}})))
	}

	m.AddRow(10, col.New(12))

	// =========================================================================
	// 4. FOOTER Y TOTALES (Mejorado)
	// =========================================================================

	var subtotal15, subtotal0, iva15 float64
	for _, tax := range factura.InfoFactura.TotalConImpuestos {
		if tax.Codigo == "2" {
			if tax.CodigoPorcentaje == "2" || tax.CodigoPorcentaje == "3" || tax.CodigoPorcentaje == "4" {
				subtotal15 += tax.BaseImponible
				iva15 += tax.Valor
			}
			if tax.CodigoPorcentaje == "0" {
				subtotal0 += tax.BaseImponible
			}
		}
	}

	type totalRow struct {
		label string
		val   string
	}
	totals := []totalRow{
		{"Subtotal 15%", fmtMoney(subtotal15)},
		{"Subtotal 0%", fmtMoney(subtotal0)},
		{"Descuento", fmtMoney(factura.InfoFactura.TotalDescuento)},
		{"IVA 15%", fmtMoney(iva15)},
	}

	rowHeight := 5.0
	totalRowsNeeded := len(totals)
	if len(factura.InfoAdicional) > totalRowsNeeded {
		totalRowsNeeded = len(factura.InfoAdicional)
	}

	m.AddRow(6,
		col.New(7).Add(text.New("INFORMACIÓN ADICIONAL", props.Text{Style: fontstyle.Bold, Size: 8})),
		col.New(5),
	)

	for i := 0; i < totalRowsNeeded+2; i++ {
		colIzq := col.New(7)
		colDerLbl := col.New(3)
		colDerVal := col.New(2)

		if i < len(factura.InfoAdicional) {
			info := factura.InfoAdicional[i]
			label := info.Nombre
			if label == "Observacion" {
				label = "Observación"
			}
			colIzq.Add(text.New(fmt.Sprintf("%s: %s", label, info.Value), props.Text{Size: 7, Color: textGray, Top: 0}))
		} else if i == len(factura.InfoAdicional) {
			colIzq.Add(text.New("Forma Pago: Sistema Financiero", props.Text{Size: 7, Style: fontstyle.Bold, Top: 2}))
		}

		if i < len(totals) {
			t := totals[i]
			colDerLbl.Add(text.New(t.label, props.Text{Size: 8, Align: align.Right, Color: textGray, Right: 2}))
			colDerVal.Add(text.New(t.val, props.Text{Size: 8, Align: align.Right, Color: textDark}))
		}

		m.AddRow(rowHeight, colIzq, colDerLbl, colDerVal)
	}

	m.AddRow(5, col.New(12))

	// --- Total Final Verde ---
	m.AddRow(12,
		col.New(7),
		col.New(5).WithStyle(&props.Cell{BackgroundColor: brandPrimary}).Add(
			text.New("TOTAL A PAGAR", props.Text{Size: 10, Style: fontstyle.Bold, Color: textWhite, Align: align.Left, Left: 4, Top: 3.5}),
			text.New("$ "+fmtMoney(factura.InfoFactura.ImporteTotal), props.Text{Size: 12, Style: fontstyle.Bold, Color: textWhite, Align: align.Right, Right: 4, Top: 3}),
		),
	)

	// =========================================================================
	// 5. PIE DE PÁGINA
	// =========================================================================
	m.AddRow(15, col.New(12))
	m.AddRow(10,
		text.NewCol(12, "Documento generado por Kushki App - Tecnología hecha en Ecuador", props.Text{
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
