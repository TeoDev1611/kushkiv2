package util

import (
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

// ProcessAndSaveLogo toma una ruta de imagen, la redimensiona si es necesario y la guarda como PNG optimizado.
func ProcessAndSaveLogo(inputPath string, outputDir string) (string, error) {
	// 1. Abrir archivo original
	file, err := os.Open(inputPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 2. Decodificar (Soporta PNG y JPG)
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// 3. Calcular nuevas dimensiones (Max ancho 400px para logos de facturas)
	// Manteniendo aspecto
	maxWidth := 400
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	var finalImg image.Image

	if width > maxWidth {
		ratio := float64(maxWidth) / float64(width)
		newHeight := int(float64(height) * ratio)
		
		// Crear lienzo destino
		dst := image.NewRGBA(image.Rect(0, 0, maxWidth, newHeight))
		
		// Redimensionar usando interpolación de alta calidad (si es posible, si no, básica)
		// Usamos "draw.CatmullRom" si estuviera disponible, pero usaremos ApproxBiLinear del paquete x/image/draw
		// Si no se tiene x/image, usaremos un método manual simple, pero asumiremos standard draw.
		// Para asegurar compatibilidad sin `go get`, usaremos una aproximación simple si draw.Scaler no está.
		// NOTA: Para evitar problemas de dependencias externas en este entorno, usaré un resize simple.
		
		draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)
		finalImg = dst
	} else {
		finalImg = img
	}

	// 4. Crear directorio si no existe
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	// 5. Guardar como PNG (Transparencia soportada)
	outputPath := filepath.Join(outputDir, "company_logo.png")
	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	// Si el origen era JPG, convertimos a PNG. Si era PNG, optimizamos.
	// Usamos PNG para preservar calidad en texto/logos.
	err = png.Encode(outFile, finalImg)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

// IsImage verifica extensiones
func IsImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg"
}
