package templating

import (
	"embed"
	"fmt"
	"math"
	"path"
	"strings"
	"time"

	"github.com/soerenschneider/alert-forge/internal"
)

const templatesDir = "templates"

var (
	//go:embed templates/*
	templates embed.FS
)

func GetTemplate(name string) (string, error) {
	fileName := path.Join(templatesDir, name)
	data, err := templates.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func IsZeroTime(t time.Time) bool {
	return t.IsZero()
}

func MapSeverityToColor(severity string) string {
	code := internal.DefaultSeverityMapping[strings.ToLower(severity)]
	return MapSeverityValueToColor(code)
}

func MapSeverityValueToColor(severity int64) string {
	// Clamp severity to the range [0, 15]
	if severity < 0 {
		severity = 0
	} else if severity > 15 {
		severity = 15
	}

	// Convert severity to a value in the range [0.0, 1.0]
	normalized := float64(severity) / 15.0

	// Pastel base value to ensure light colors (around 230–255 range)
	base := 230.0

	// Interpolate colors (yellow -> orange -> red) in RGB
	// Red: 255 (constant)
	// Green: 255 -> base (softened to pastel orange and red)
	// Blue: 255 -> base (blue softly fades, keeping pastel tones)
	red := int(base + normalized*(255-base))   // Red increases toward 255
	green := int(base - normalized*(base-150)) // Green decreases toward 150 (pastel orange)
	blue := int(base - normalized*base)        // Blue starts at pastel base and fades softly

	// Ensure RGB values are in the range [0, 255]
	red = int(math.Min(math.Max(float64(red), 0), 255))
	green = int(math.Min(math.Max(float64(green), 0), 255))
	blue = int(math.Min(math.Max(float64(blue), 0), 255))

	// Format RGB values as a pastel hex color
	return fmt.Sprintf("#%02X%02X%02X", red, green, blue)
}
