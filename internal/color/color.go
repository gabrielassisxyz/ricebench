// Package color is the single authoritative implementation of color arithmetic for
// RiceBench. Every contrast ratio, gamut mapping decision and resolved sRGB value that a
// validity gate or a candidate palette depends on is computed here and nowhere else.
// The frontend consumes resolved values and performs no color math of its own.
package color

import (
	"fmt"
	"math"
)

// gamutTolerance is the tolerance for the in-gamut membership test. The OKLCH to linear
// sRGB round trip through the conversion matrices is exact only to float precision, so a
// color whose channel lands a few units in the last place outside [0, 1] is still in
// gamut. The value sits far below the 8-bit channel step of 1/255, so it can never admit
// a color that is visibly out of gamut.
const gamutTolerance = 1e-6

// OKLCH is a bounded color in the OKLCH color space. Lightness is in [0, 1], chroma is
// non-negative, and hue is in degrees, normalized to [0, 360).
type OKLCH struct {
	Lightness float64 `json:"lightness"`
	Chroma    float64 `json:"chroma"`
	Hue       float64 `json:"hue"`
}

// NewOKLCH constructs an OKLCH color, rejecting non-finite and out-of-range input
// instead of silently clamping it. Hue is periodic, so a finite hue wraps into [0, 360).
func NewOKLCH(lightness, chroma, hue float64) (OKLCH, error) {
	if math.IsNaN(lightness) || math.IsInf(lightness, 0) {
		return OKLCH{}, fmt.Errorf("oklch lightness: %v is not finite", lightness)
	}
	if lightness < 0 || lightness > 1 {
		return OKLCH{}, fmt.Errorf("oklch lightness %v out of range [0, 1]", lightness)
	}
	if math.IsNaN(chroma) || math.IsInf(chroma, 0) {
		return OKLCH{}, fmt.Errorf("oklch chroma: %v is not finite", chroma)
	}
	if chroma < 0 {
		return OKLCH{}, fmt.Errorf("oklch chroma %v out of range [0, inf)", chroma)
	}
	if math.IsNaN(hue) || math.IsInf(hue, 0) {
		return OKLCH{}, fmt.Errorf("oklch hue: %v is not finite", hue)
	}
	return OKLCH{Lightness: lightness, Chroma: chroma, Hue: normalizeHue(hue)}, nil
}

// normalizeHue wraps a finite hue into [0, 360).
func normalizeHue(hue float64) float64 {
	hue = math.Mod(hue, 360)
	if hue < 0 {
		hue += 360
	}
	return hue
}

// SRGB is a bounded 8-bit sRGB color. The byte range is the bound, so an SRGB value can
// never hold an out-of-range channel.
type SRGB struct {
	R uint8
	G uint8
	B uint8
}

// Hex returns the canonical six-digit hexadecimal form, "#rrggbb".
func (c SRGB) Hex() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// ParseHex parses a six-digit hexadecimal color. It rejects malformed input: a missing
// leading '#', a wrong length, or a non-hexadecimal digit.
func ParseHex(value string) (SRGB, error) {
	if len(value) != 7 || value[0] != '#' {
		return SRGB{}, fmt.Errorf("srgb %q: want six hexadecimal digits with a leading #", value)
	}
	var channels [3]uint8
	for index := 0; index < 3; index++ {
		high := hexDigit(value[1+2*index])
		low := hexDigit(value[2+2*index])
		if high < 0 || low < 0 {
			return SRGB{}, fmt.Errorf("srgb %q: channel %d is not hexadecimal", value, index)
		}
		channels[index] = uint8(high<<4 | low)
	}
	return SRGB{R: channels[0], G: channels[1], B: channels[2]}, nil
}

func hexDigit(digit byte) int {
	switch {
	case digit >= '0' && digit <= '9':
		return int(digit - '0')
	case digit >= 'a' && digit <= 'f':
		return int(digit-'a') + 10
	case digit >= 'A' && digit <= 'F':
		return int(digit-'A') + 10
	default:
		return -1
	}
}

// Resolved is the authoritative resolved sRGB of an authored OKLCH color. Mapped records
// whether gamut mapping changed the color, and From records the pre-mapping color when it
// did.
type Resolved struct {
	SRGB   string `json:"srgb"`
	Mapped bool   `json:"mapped"`
	From   *OKLCH `json:"from,omitempty"`
}
