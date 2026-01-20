package main

import (
	"context"
	"embed"
	"kushkiv2/internal/db"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
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
