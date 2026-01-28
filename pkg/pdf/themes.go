package pdf

import (
	"os"

	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/barcode"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/linestyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"

	srixml "kushkiv2/pkg/xml"
)

// InvoiceTheme define el contrato para cualquier diseño de factura
type InvoiceTheme interface {
	Build(m core.Maroto, f srixml.FacturaXML, logoPath string)
}

// ============================================================================
// DISEÑO A: MODERN (El original con colores Esmeralda)
// ============================================================================
type ModernTheme struct{}

func (t *ModernTheme) Build(m core.Maroto, f srixml.FacturaXML, logoPath string) {
	// 1. Cabecera
	colLogo := col.New(4)
	if logoPath != "" {
		if _, err := os.Stat(logoPath); err == nil {
			colLogo.Add(image.NewFromFile(logoPath, props.Rect{Center: false, Percent: 100, Left: 0}))
		}
	} else {
		colLogo.Add(text.New("KUSHKI APP", props.Text{Style: fontstyle.Bold, Color: colorEmeraldPrimary, Size: 16}))
	}

	m.AddRow(20,
		colLogo,
		col.New(8).Add(
			text.New("FACTURA ELECTRÓNICA", props.Text{Size: 14, Style: fontstyle.Bold, Align: align.Right, Color: colorEmeraldPrimary, Top: 0}),
			text.New("No. "+f.InfoTributaria.Estab+"-"+f.InfoTributaria.PtoEmi+"-"+f.InfoTributaria.Secuencial, props.Text{Size: 11, Align: align.Right, Top: 8, Style: fontstyle.Bold}),
		),
	)
	m.AddRow(5, col.New(12))

	// Info Emisor y Clave
	m.AddRow(15,
		col.New(6).Add(
			text.New(f.InfoTributaria.RazonSocial, props.Text{Size: 11, Style: fontstyle.Bold, Color: colorDarkGray, Top: 0}),
			text.New("RUC: "+f.InfoTributaria.Ruc, props.Text{Size: 9, Color: colorDarkGray, Top: 6}),
		),
		col.New(6).WithStyle(&props.Cell{BackgroundColor: colorEmeraldLight}).Add(
			text.New("CLAVE DE ACCESO:", props.Text{Size: 7, Style: fontstyle.Bold, Align: align.Center, Color: colorDarkGray, Top: 2}),
			text.New(f.InfoTributaria.ClaveAcceso, props.Text{Size: 7, Align: align.Center, Top: 6, Family: "Courier"}),
		),
	)

	// Direcciones
	m.AddRow(15,
		col.New(6).Add(
			text.New("Matriz: "+f.InfoTributaria.DirMatriz, props.Text{Size: 7, Color: colorDarkGray, Top: 0}),
			text.New("Sucursal: "+f.InfoFactura.DirEstablecimiento, props.Text{Size: 7, Color: colorDarkGray, Top: 4}),
		),
		col.New(6).Add(
			text.New("FECHA: "+f.InfoFactura.FechaEmision, props.Text{Size: 8, Align: align.Right, Top: 0, Style: fontstyle.Bold}),
			text.New("AMBIENTE: "+getAmbienteText(f.InfoTributaria.Ambiente), props.Text{Size: 8, Align: align.Right, Top: 4}),
		),
	)

	// Barcode
	m.AddRow(12,
		col.New(6),
		col.New(6).Add(code.NewBar(f.InfoTributaria.ClaveAcceso, props.Barcode{Type: barcode.Code128, Proportion: props.Proportion{Width: 20, Height: 5}, Center: true})),
	)

	m.AddRow(5, col.New(12))

	// 2. Cliente
	m.AddRow(8, col.New(12).WithStyle(&props.Cell{BackgroundColor: colorEmeraldPrimary}).Add(
		text.New("INFORMACIÓN DEL RECEPTOR", props.Text{Size: 8, Style: fontstyle.Bold, Color: colorWhite, Top: 2, Left: 2}),
	))

	m.AddRow(20,
		col.New(7).WithStyle(&props.Cell{BackgroundColor: colorEmeraldLight}).Add(
			text.New("Razón Social: ", props.Text{Size: 8, Style: fontstyle.Bold, Left: 2, Top: 2}),
			text.New(f.InfoFactura.RazonSocialComprador, props.Text{Size: 9, Left: 2, Top: 6}),
			text.New("Identificación: "+f.InfoFactura.IdentificacionComprador, props.Text{Size: 9, Left: 2, Top: 12}),
		),
		col.New(5).WithStyle(&props.Cell{BackgroundColor: colorEmeraldLight}).Add(
			text.New("Dirección: "+f.InfoFactura.DireccionComprador, props.Text{Size: 8, Align: align.Right, Right: 2, Top: 2}),
		),
	)
	m.AddRow(5, col.New(12).Add(line.New(props.Line{Color: colorEmeraldPrimary, Thickness: 0.5})))

	// 3. Detalles
	m.AddRow(9,
		text.NewCol(2, "CÓDIGO", props.Text{Style: fontstyle.Bold, Size: 8, Color: colorWhite, Top: 1.5, Left: 2}),
		text.NewCol(1, "CANT.", props.Text{Style: fontstyle.Bold, Size: 8, Align: align.Center, Color: colorWhite, Top: 1.5}),
		text.NewCol(5, "DESCRIPCIÓN", props.Text{Style: fontstyle.Bold, Size: 8, Color: colorWhite, Top: 1.5, Left: 2}),
		text.NewCol(2, "P. UNIT", props.Text{Style: fontstyle.Bold, Size: 8, Align: align.Right, Color: colorWhite, Top: 1.5, Right: 2}),
		text.NewCol(2, "TOTAL", props.Text{Style: fontstyle.Bold, Size: 8, Align: align.Right, Color: colorWhite, Top: 1.5, Right: 2}),
	).WithStyle(&props.Cell{BackgroundColor: colorEmeraldPrimary})

	for _, item := range f.Detalles {
		m.AddRow(8,
			text.NewCol(2, item.CodigoPrincipal, props.Text{Size: 8, Top: 2, Left: 2}),
			text.NewCol(1, fmtMoney(item.Cantidad), props.Text{Size: 8, Align: align.Center, Top: 2}),
			text.NewCol(5, item.Descripcion, props.Text{Size: 8, Top: 2, Left: 2}),
			text.NewCol(2, fmtMoney(item.PrecioUnitario), props.Text{Size: 8, Align: align.Right, Top: 2, Right: 2}),
			text.NewCol(2, fmtMoney(item.PrecioTotalSinImpuesto), props.Text{Size: 8, Align: align.Right, Top: 2, Right: 2, Style: fontstyle.Bold}),
		)
		m.AddRow(1, col.New(12).Add(line.New(props.Line{Color: colorLightGray})))
	}

	addTotals(m, f, true)
	addFooter(m)
}

// ============================================================================
// DISEÑO B: MINIMAL (Blanco y Negro, Sin Logo, Ahorro Tinta)
// ============================================================================
type MinimalTheme struct{}

func (t *MinimalTheme) Build(m core.Maroto, f srixml.FacturaXML, logoPath string) {
	// Sin Logo, puro texto
	m.AddRow(10,
		col.New(12).Add(text.New("FACTURA ELECTRÓNICA", props.Text{Size: 16, Style: fontstyle.Bold, Align: align.Center})),
	)
	m.AddRow(8,
		col.New(12).Add(text.New(f.InfoTributaria.RazonSocial, props.Text{Size: 12, Align: align.Center})),
	)
	m.AddRow(5, col.New(12).Add(line.New(props.Line{Style: linestyle.Dashed, Thickness: 0.5})))

	// Datos en lista simple
	m.AddRow(25,
		col.New(6).Add(
			text.New("RUC: "+f.InfoTributaria.Ruc, props.Text{Size: 9, Top: 2}),
			text.New("No: "+f.InfoTributaria.Estab+"-"+f.InfoTributaria.PtoEmi+"-"+f.InfoTributaria.Secuencial, props.Text{Size: 9, Top: 6}),
			text.New("Fecha: "+f.InfoFactura.FechaEmision, props.Text{Size: 9, Top: 10}),
			text.New("Ambiente: "+getAmbienteText(f.InfoTributaria.Ambiente), props.Text{Size: 9, Top: 14}),
		),
		col.New(6).Add(
			text.New("CLAVE DE ACCESO:", props.Text{Size: 8, Style: fontstyle.Bold, Top: 2}),
			text.New(f.InfoTributaria.ClaveAcceso, props.Text{Size: 8, Family: "Courier", Top: 6}),
			code.NewBar(f.InfoTributaria.ClaveAcceso, props.Barcode{Type: barcode.Code128, Proportion: props.Proportion{Width: 20, Height: 8}, Top: 10}),
		),
	)

	m.AddRow(5, col.New(12).Add(line.New(props.Line{Style: linestyle.Dashed, Thickness: 0.5})))

	// Cliente (Texto plano)
	m.AddRow(15,
		col.New(12).Add(
			text.New("CLIENTE: "+f.InfoFactura.RazonSocialComprador, props.Text{Size: 9, Style: fontstyle.Bold, Top: 2}),
			text.New("RUC/CI: "+f.InfoFactura.IdentificacionComprador+"   |   Dir: "+f.InfoFactura.DireccionComprador, props.Text{Size: 9, Top: 6}),
		),
	)

	// Headers Detalles (Bordes simples superior e inferior)
	m.AddRow(1, col.New(12).Add(line.New(props.Line{Thickness: 1})))
	m.AddRow(6,
		text.NewCol(2, "COD", props.Text{Size: 8, Style: fontstyle.Bold}),
		text.NewCol(6, "DESCRIPCIÓN", props.Text{Size: 8, Style: fontstyle.Bold}),
		text.NewCol(2, "CANT", props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(2, "TOTAL", props.Text{Size: 8, Style: fontstyle.Bold, Align: align.Right}),
	)
	m.AddRow(1, col.New(12).Add(line.New(props.Line{Thickness: 1})))

	// Items
	for _, item := range f.Detalles {
		m.AddRow(6,
			text.NewCol(2, item.CodigoPrincipal, props.Text{Size: 8}),
			text.NewCol(6, item.Descripcion, props.Text{Size: 8}),
			text.NewCol(2, fmtMoney(item.Cantidad), props.Text{Size: 8, Align: align.Center}),
			text.NewCol(2, fmtMoney(item.PrecioTotalSinImpuesto), props.Text{Size: 8, Align: align.Right}),
		)
	}
	m.AddRow(1, col.New(12).Add(line.New(props.Line{Thickness: 0.5})))

	addTotals(m, f, false) // False = sin fondo de color
	addFooter(m)
}

// ============================================================================
// DISEÑO C: CORPORATE (Blanco y Negro, Con Logo, Formal)
// ============================================================================
type CorporateTheme struct{}

func (t *CorporateTheme) Build(m core.Maroto, f srixml.FacturaXML, logoPath string) {
	// Header estilo membrete
	colLogo := col.New(3)
	if logoPath != "" {
		if _, err := os.Stat(logoPath); err == nil {
			colLogo.Add(image.NewFromFile(logoPath, props.Rect{Center: false, Percent: 100}))
		}
	}
	
	m.AddRow(25,
		colLogo,
		col.New(9).Add(
			text.New(f.InfoTributaria.RazonSocial, props.Text{Size: 14, Style: fontstyle.Bold, Align: align.Right}),
			text.New("RUC: "+f.InfoTributaria.Ruc, props.Text{Size: 10, Align: align.Right, Top: 6}),
			text.New("Dir: "+f.InfoTributaria.DirMatriz, props.Text{Size: 8, Align: align.Right, Top: 10}),
			text.New(f.InfoTributaria.Estab+"-"+f.InfoTributaria.PtoEmi+"-"+f.InfoTributaria.Secuencial, props.Text{Size: 10, Align: align.Right, Top: 15, Style: fontstyle.Bold}),
		),
	)
	m.AddRow(2, col.New(12).Add(line.New(props.Line{Thickness: 2, Color: colorBlack})))

	// Datos Documento Boxed
	m.AddRow(25,
		col.New(6).Add(
			text.New("CLIENTE:", props.Text{Size: 8, Style: fontstyle.Bold, Top: 2}),
			text.New(f.InfoFactura.RazonSocialComprador, props.Text{Size: 10, Top: 6}),
			text.New("RUC/CI: "+f.InfoFactura.IdentificacionComprador, props.Text{Size: 9, Top: 11}),
			text.New("Dir: "+f.InfoFactura.DireccionComprador, props.Text{Size: 9, Top: 15}),
		),
		col.New(6).Add(
			text.New("AUTORIZACIÓN / CLAVE DE ACCESO:", props.Text{Size: 8, Style: fontstyle.Bold, Top: 2}),
			text.New(f.InfoTributaria.ClaveAcceso, props.Text{Size: 8, Family: "Courier", Top: 6}),
			text.New("FECHA EMISIÓN: "+f.InfoFactura.FechaEmision, props.Text{Size: 9, Top: 11, Style: fontstyle.Bold}),
			text.New("AMBIENTE: "+getAmbienteText(f.InfoTributaria.Ambiente), props.Text{Size: 9, Top: 15}),
		),
	)

	m.AddRow(5, col.New(12))

	// Tabla con Cabecera Negra
	m.AddRow(8,
		text.NewCol(2, "CÓDIGO", props.Text{Style: fontstyle.Bold, Size: 8, Color: colorWhite, Top: 1.5, Align: align.Center}),
		text.NewCol(1, "CANT", props.Text{Style: fontstyle.Bold, Size: 8, Color: colorWhite, Top: 1.5, Align: align.Center}),
		text.NewCol(5, "DESCRIPCIÓN", props.Text{Style: fontstyle.Bold, Size: 8, Color: colorWhite, Top: 1.5}),
		text.NewCol(2, "P. UNIT", props.Text{Style: fontstyle.Bold, Size: 8, Color: colorWhite, Top: 1.5, Align: align.Right, Right: 2}),
		text.NewCol(2, "TOTAL", props.Text{Style: fontstyle.Bold, Size: 8, Color: colorWhite, Top: 1.5, Align: align.Right, Right: 2}),
	).WithStyle(&props.Cell{BackgroundColor: colorBlack})

	for i, item := range f.Detalles {
		bg := colorWhite
		if i%2 != 0 {
			bg = colorBackground // Zebra striping
		}
		m.AddRow(7,
			text.NewCol(2, item.CodigoPrincipal, props.Text{Size: 8, Align: align.Center, Top: 1.5}),
			text.NewCol(1, fmtMoney(item.Cantidad), props.Text{Size: 8, Align: align.Center, Top: 1.5}),
			text.NewCol(5, item.Descripcion, props.Text{Size: 8, Top: 1.5}),
			text.NewCol(2, fmtMoney(item.PrecioUnitario), props.Text{Size: 8, Align: align.Right, Top: 1.5, Right: 2}),
			text.NewCol(2, fmtMoney(item.PrecioTotalSinImpuesto), props.Text{Size: 8, Align: align.Right, Top: 1.5, Right: 2}),
		).WithStyle(&props.Cell{BackgroundColor: bg})
	}

	m.AddRow(2, col.New(12).Add(line.New(props.Line{Thickness: 1, Color: colorBlack})))

	addTotals(m, f, false)
	addFooter(m)
}

// ============================================================================
// HELPERS PRIVADOS DE RENDERIZADO
// ============================================================================

func addTotals(m core.Maroto, f srixml.FacturaXML, colorful bool) {
	var subtotal15, subtotal0, iva15 float64
	for _, tax := range f.InfoFactura.TotalConImpuestos {
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

	totals := []struct {
		label string
		val   string
	}{
		{"Subtotal 15%", fmtMoney(subtotal15)},
		{"Subtotal 0%", fmtMoney(subtotal0)},
		{"Descuento", fmtMoney(f.InfoFactura.TotalDescuento)},
		{"IVA 15%", fmtMoney(iva15)},
		{"TOTAL", fmtMoney(f.InfoFactura.ImporteTotal)},
	}

	m.AddRow(5, col.New(12))

	for i, t := range totals {
		// Estilo para el TOTAL final
		isLast := i == len(totals)-1
		
		lblStyle := props.Text{Size: 8, Align: align.Right, Right: 2}
		valStyle := props.Text{Size: 8, Align: align.Right}
		
		var cellStyle *props.Cell
		if isLast && colorful {
			cellStyle = &props.Cell{BackgroundColor: colorEmeraldPrimary}
			lblStyle.Color = colorWhite
			valStyle.Color = colorWhite
			lblStyle.Style = fontstyle.Bold
			valStyle.Style = fontstyle.Bold
		} else if isLast {
			lblStyle.Style = fontstyle.Bold
			valStyle.Style = fontstyle.Bold
		}

		m.AddRow(5,
			col.New(8),
			col.New(2).Add(text.New(t.label, lblStyle)).WithStyle(cellStyle),
			col.New(2).Add(text.New(t.val, valStyle)).WithStyle(cellStyle),
		)
	}
}

func addFooter(m core.Maroto) {
	m.AddRow(10, col.New(12))
	m.AddRow(8, text.NewCol(12, footerText, props.Text{
		Size:  6,
		Align: align.Center,
		Color: colorLightGray,
		Style: fontstyle.Italic,
	}))
}
