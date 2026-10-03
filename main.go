package main

import (
	"embed"
	"log"

	"desktopcal/internal/db"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	database, err := db.New()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	app := NewApp(database)

	if err := wails.Run(&options.App{
		Title:             "桌面日历",
		Width:             1100,
		Height:            750,
		MinWidth:          900,
		MinHeight:         600,
		AlwaysOnTop:       true,
		HideWindowOnClose: true,
		AssetServer:       &assetserver.Options{Assets: assets},
		BackgroundColour:  &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:         app.startup,
		Mac: &mac.Options{
			DisableEscapeExitsFullscreen: true,
		},
		Bind: []interface{}{
			app,
		},
	}); err != nil {
		log.Fatal(err)
	}
}
