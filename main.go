package main

import (
	"github.com/CyrusRoshan/Golf/game"
	"github.com/CyrusRoshan/Golf/screen"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	width, height := screen.ScreenBounds()
	cfg := pixelgl.WindowConfig{
		Title:     "Golf",
		Bounds:    pixel.R(0, 0, width/2, height/2),
		VSync:     true,
		Resizable: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	golf := game.NewGame(win)
	golf.Run()
}
