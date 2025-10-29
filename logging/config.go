package logging

import (
	"os"
	"strings"
)

type Color string
type Role string

type SystemInfo struct {
	Role  Role
	Color Color
}

func (SystemInfo) ParseColor(color string) Color {
	color = strings.ToLower(color)
	color = strings.ReplaceAll(color, " ", "")
	switch color {
	case "red":
		return red
	case "green":
		return green
	case "yellow":
		return yellow
	case "blue":
		return blue
	case "purple":
		return purple
	case "cyan":
		return cyan
	case "gray":
		return gray
	default:
		return white
	}
}

var config SystemInfo

func ReadConfig() {
	config.Role = Role(os.Getenv("ROLE"))
	config.Color = config.ParseColor(os.Getenv("COLOR"))
}
