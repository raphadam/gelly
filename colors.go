package gelly

import "image/color"

var (
	ColorWhite = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	ColorBlack = color.RGBA{R: 0, G: 0, B: 0, A: 255}

	ColorLightGrey = color.RGBA{R: 160, G: 160, B: 160, A: 255}
	ColorGrey      = color.RGBA{R: 127, G: 127, B: 127, A: 255}
	ColorDarkGrey  = color.RGBA{R: 60, G: 60, B: 60, A: 255}

	ColorRed   = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	ColorGreen = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	ColorBlue  = color.RGBA{R: 0, G: 0, B: 255, A: 255}

	ColorYellow  = color.RGBA{R: 255, G: 255, B: 0, A: 255}
	ColorMagenta = color.RGBA{R: 255, G: 0, B: 255, A: 255}
	ColorCyan    = color.RGBA{R: 0, G: 255, B: 255, A: 255}

	BackgroundColor = color.RGBA{G: 100, B: 120, A: 255}
)
