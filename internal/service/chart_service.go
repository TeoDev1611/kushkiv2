package service

import (
	"bytes"
	"kushkiv2/internal/db"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

type ChartService struct{}

func NewChartService() *ChartService {
	return &ChartService{}
}

// GenerateRevenueChart genera una gráfica de barras de ventas mensuales
func (s *ChartService) GenerateRevenueChart() (string, error) {
	var facturas []db.Factura
	// Cargamos facturas que tengan valor, sin importar el estado
	db.GetDB().Where("total > 0").Order("fecha_emision ASC").Find(&facturas)

	if len(facturas) == 0 {
		return "", nil
	}

	trendMap := make(map[string]float64)
	var keys []string
	
	for _, f := range facturas {
		mes := f.FechaEmision.Format("2006-01")
		if _, ok := trendMap[mes]; !ok {
			keys = append(keys, mes)
		}
		trendMap[mes] += f.Total
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Ingresos Mensuales",
			Left: "center",
			TitleStyle: &opts.TextStyle{Color: "#34d399", FontSize: 18, FontWeight: "bold"},
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeWesteros, 
			Height: "380px", 
			BackgroundColor: "transparent",
		}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(false)}), // ELIMINAR LEYENDA (Punto blanco)
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true), Trigger: "axis"}),
		charts.WithXAxisOpts(opts.XAxis{
			AxisLabel: &opts.AxisLabel{Color: "#94a3b8", Rotate: 35, FontSize: 10},
			SplitLine: &opts.SplitLine{Show: opts.Bool(false)},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			AxisLabel: &opts.AxisLabel{Color: "#94a3b8", FontSize: 10},
			SplitLine: &opts.SplitLine{Show: opts.Bool(true), LineStyle: &opts.LineStyle{Color: "rgba(255,255,255,0.05)", Type: "dashed"}},
		}),
		charts.WithGridOpts(opts.Grid{
			Top: "20%",
			Bottom: "20%",
			Left: "10%",
			Right: "5%",
			ContainLabel: opts.Bool(true),
		}),
	)

	y := make([]opts.LineData, 0)
	for _, k := range keys {
		y = append(y, opts.LineData{Value: trendMap[k], Symbol: "circle", SymbolSize: 6})
	}

	line.SetXAxis(keys).AddSeries("Ventas", y).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				Smooth: opts.Bool(true),
			}),
			charts.WithAreaStyleOpts(opts.AreaStyle{
				Opacity: opts.Float(0.1),
				Color:   "#34d399",
			}),
			charts.WithLabelOpts(opts.Label{Show: opts.Bool(true), Position: "top", Color: "#eee"}),
		)

	var buf bytes.Buffer
	line.Render(&buf)
	return buf.String(), nil
}

// GenerateClientsPie genera una gráfica de pastel de los mejores clientes
func (s *ChartService) GenerateClientsPie() (string, error) {
	type DataPoint struct {
		Nombre string
		Total  float64
	}
	var results []DataPoint

	// Consulta simplificada para asegurar resultados
	db.GetDB().Table("facturas").
		Select("clients.nombre as nombre, SUM(facturas.total) as total").
		Joins("JOIN clients ON clients.id = facturas.cliente_id").
		Where("facturas.total > 0").
		Group("facturas.cliente_id").
		Order("total DESC").
		Limit(5).
		Scan(&results)

	if len(results) == 0 {
		return "", nil
	}

	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Top 5 Clientes", 
			Left: "center",
			TitleStyle: &opts.TextStyle{Color: "#eee", FontSize: 16},
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeMacarons, 
			Height: "300px",
			BackgroundColor: "transparent",
		}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(false)}), // Ocultar leyenda para dar espacio
	)

	items := make([]opts.PieData, 0)
	for _, r := range results {
		items = append(items, opts.PieData{Name: r.Nombre, Value: r.Total})
	}

	pie.AddSeries("Ventas", items).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:      opts.Bool(true),
				Formatter: "{b}: {d}%", // Nombre y porcentaje
				Color:     "#fff",
			}),
			charts.WithPieChartOpts(opts.PieChart{
				Radius: []string{"30%", "60%"},
				Center: []string{"50%", "55%"},
			}),
		)

	var buf bytes.Buffer
	pie.Render(&buf)
	return buf.String(), nil
}
