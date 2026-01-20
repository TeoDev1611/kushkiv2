package util

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ZipSource agrega un archivo o directorio al escritor zip.
func ZipSource(source, prefix string, writer *zip.Writer) error {
	info, err := os.Stat(source)
	if err != nil {
		return nil // Si no existe, lo ignoramos para no romper el backup
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			// Mantener estructura relativa
			relPath, _ := filepath.Rel(source, path)
			header.Name = filepath.Join(prefix, baseDir, relPath)
		} else {
			header.Name = filepath.Join(prefix, filepath.Base(path))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			// Si no podemos abrir el archivo (lock), lo saltamos
			return nil 
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
}

// CreateBackupZip genera un archivo zip con múltiples fuentes.
// sources: mapa de ruta_origen -> prefijo_en_zip
func CreateBackupZip(destPath string, sources map[string]string) error {
	// Crear archivo zip de destino
	outFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)
	defer w.Close()

	for src, prefix := range sources {
		// Normalizar paths
		src = strings.TrimSpace(src)
		if src == "" {
			continue
		}
		if err := ZipSource(src, prefix, w); err != nil {
			// Loggear pero continuar con otros archivos
			// fmt.Printf("Error comprimiendo %s: %v\n", src, err)
		}
	}

	return nil
}
