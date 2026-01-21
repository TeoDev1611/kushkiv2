package main

import (
	"context"
	"embed"
	"fmt"
	"kushkiv2/internal/db"
	"kushkiv2/pkg/crypto"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 1. Inicializar Seguridad (Cargar o Generar Llave Maestra)
	cwd, _ := os.Getwd()
	keyPath := filepath.Join(cwd, "master.key")
	if err := crypto.InitSecurity(keyPath); err != nil {
		panic(fmt.Sprintf("FATAL: No se pudo inicializar la seguridad criptográfica: %v", err))
	}

	// Initialize Database and Run Migrations
	db.Migrate(db.GetDB())

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "kushkiv2",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown: func(ctx context.Context) {
			// Cerrar conexión DB para liberar archivo
			if err := db.CloseDB(); err != nil {
				println("Error cerrando DB:", err.Error())
			}
			
			// Ejecutar Backup Automático al cerrar
			println("Creando respaldo automático...")
			if err := app.CreateBackup(); err != nil {
				println("Error en respaldo:", err.Error())
			} else {
				println("Respaldo completado.")
			}
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
