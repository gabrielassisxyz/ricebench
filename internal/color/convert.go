package color

import "math"

// The conversion chain follows CSS Color 4. Linear sRGB converts to D65 XYZ and then to
// OKLab and OKLCH through the matrices the spec pins, and back. The matrices are applied
// as result[i] = sum over k of M[i][k] * v[k].

// JND is the deltaEOK just-noticeable difference used by gamut mapping, and
// searchEpsilon is the chroma search termination epsilon. Both are frozen by the
// versioned color contract rather than inherited from a library default.
const (
	JND           = 0.02
	searchEpsilon = 0.0001
)

// Resolve converts an authored OKLCH color to its resolved sRGB, gamut mapping it when
// it falls outside the sRGB gamut. The result records whether it was mapped and, when it
// was, the pre-mapping color.
func Resolve(c OKLCH) Resolved {
	if inGamut(c) {
		return Resolved{SRGB: toSRGB(c).Hex()}
	}
	mapped := GamutMap(c)
	return Resolved{SRGB: toSRGB(mapped).Hex(), Mapped: true, From: &c}
}

// ToSRGB converts an OKLCH color to rounded 8-bit sRGB, gamut mapping it first, so the
// result is always in gamut.
func ToSRGB(c OKLCH) SRGB {
	return toSRGB(GamutMap(c))
}

// GamutMap maps an OKLCH color into the sRGB gamut using the CSS Color 4 binary search
// with local MINDE. In-gamut colors are returned unchanged. The search holds lightness
// and hue constant while reducing chroma; the local MINDE refinement returns the clipped
// color when its deltaEOK from the current chroma-reduced color is below JND. Lightness
// at or above 1 maps to white, and at or below 0 maps to black.
func GamutMap(c OKLCH) OKLCH {
	if c.Lightness >= 1 {
		return OKLCH{Lightness: 1, Chroma: 0, Hue: 0}
	}
	if c.Lightness <= 0 {
		return OKLCH{Lightness: 0, Chroma: 0, Hue: 0}
	}
	if inGamut(c) {
		return c
	}

	current := c
	clipped := clipToGamma(current)
	if deltaEOK(gammaToOKLab(clipped), oklchToOKLab(current)) < JND {
		return gammaToOKLCH(clipped)
	}

	minChroma := 0.0
	maxChroma := c.Chroma
	minInGamut := true
	for maxChroma-minChroma > searchEpsilon {
		chroma := (minChroma + maxChroma) / 2
		current = OKLCH{Lightness: c.Lightness, Chroma: chroma, Hue: c.Hue}
		if minInGamut && inGamut(current) {
			minChroma = chroma
			continue
		}
		clipped = clipToGamma(current)
		e := deltaEOK(gammaToOKLab(clipped), oklchToOKLab(current))
		if e < JND {
			if JND-e < searchEpsilon {
				return gammaToOKLCH(clipped)
			}
			minInGamut = false
			minChroma = chroma
		} else {
			maxChroma = chroma
		}
	}
	if minInGamut {
		return gammaToOKLCH(clipToGamma(OKLCH{Lightness: c.Lightness, Chroma: minChroma, Hue: c.Hue}))
	}
	return gammaToOKLCH(clipped)
}

// toSRGB converts an in-gamut OKLCH color to rounded 8-bit sRGB.
func toSRGB(c OKLCH) SRGB {
	linear := oklchToLinearSRGB(c)
	return SRGB{
		R: channelToByte(linear[0]),
		G: channelToByte(linear[1]),
		B: channelToByte(linear[2]),
	}
}

// channelToByte encodes one linear sRGB channel to an 8-bit byte. This is the documented
// rounding rule for the OKLCH to sRGB conversion: clamp the channel to [0, 1], gamma
// encode it, scale by 255, and round half away from zero to the nearest integer. The
// serialized precision is fixed at 8 bits per channel, written as "#rrggbb".
func channelToByte(linear float64) uint8 {
	return uint8(math.Round(gammaEncode(clamp01(linear)) * 255))
}

// RelativeLuminance returns the WCAG 2.2 relative luminance of the sRGB color.
func (c SRGB) RelativeLuminance() float64 {
	return relativeLuminance(float64(c.R)/255, float64(c.G)/255, float64(c.B)/255)
}

func relativeLuminance(r, g, b float64) float64 {
	return 0.2126*gammaDecode(r) + 0.7152*gammaDecode(g) + 0.0722*gammaDecode(b)
}

// ContrastRatio returns the WCAG 2.2 contrast ratio of two sRGB colors, always at least
// 1, with the lighter color in the numerator.
func ContrastRatio(a, b SRGB) float64 {
	return contrastRatio(a.RelativeLuminance(), b.RelativeLuminance())
}

func contrastRatio(lighter, darker float64) float64 {
	if lighter < darker {
		lighter, darker = darker, lighter
	}
	return (lighter + 0.05) / (darker + 0.05)
}

func inGamut(c OKLCH) bool {
	linear := oklchToLinearSRGB(c)
	for _, channel := range linear {
		if channel < -gamutTolerance || channel > 1+gamutTolerance {
			return false
		}
	}
	return true
}

func clipToGamma(c OKLCH) [3]float64 {
	linear := oklchToLinearSRGB(c)
	return [3]float64{
		clamp01(gammaEncode(linear[0])),
		clamp01(gammaEncode(linear[1])),
		clamp01(gammaEncode(linear[2])),
	}
}

func gammaToOKLab(gamma [3]float64) [3]float64 {
	return linearSRGBToOKLab([3]float64{
		gammaDecode(gamma[0]),
		gammaDecode(gamma[1]),
		gammaDecode(gamma[2]),
	})
}

func gammaToOKLCH(gamma [3]float64) OKLCH {
	return oklabToOKLCH(gammaToOKLab(gamma))
}

func deltaEOK(reference, sample [3]float64) float64 {
	dl := reference[0] - sample[0]
	da := reference[1] - sample[1]
	db := reference[2] - sample[2]
	return math.Sqrt(dl*dl + da*da + db*db)
}

func oklchToOKLab(c OKLCH) [3]float64 {
	radians := c.Hue * math.Pi / 180
	return [3]float64{
		c.Lightness,
		c.Chroma * math.Cos(radians),
		c.Chroma * math.Sin(radians),
	}
}

func oklabToOKLCH(lab [3]float64) OKLCH {
	chroma := math.Hypot(lab[1], lab[2])
	hue := math.Atan2(lab[2], lab[1]) * 180 / math.Pi
	if hue < 0 {
		hue += 360
	}
	return OKLCH{Lightness: lab[0], Chroma: chroma, Hue: hue}
}

func oklchToLinearSRGB(c OKLCH) [3]float64 {
	return oklabToLinearSRGB(oklchToOKLab(c))
}

func oklabToLinearSRGB(lab [3]float64) [3]float64 {
	return xyzToLinearSRGB(oklabToXYZ(lab))
}

func oklabToXYZ(lab [3]float64) [3]float64 {
	l, a, b := lab[0], lab[1], lab[2]
	lms0 := 1.0*l + 0.3963377773761749*a + 0.2158037573099136*b
	lms1 := 1.0*l - 0.1055613458156586*a - 0.0638541728258133*b
	lms2 := 1.0*l - 0.0894841775298119*a - 1.2914855480194092*b
	lms0 *= lms0 * lms0
	lms1 *= lms1 * lms1
	lms2 *= lms2 * lms2
	return [3]float64{
		1.2268798758459243*lms0 - 0.5578149944602171*lms1 + 0.2813910456659647*lms2,
		-0.0405757452148008*lms0 + 1.1122868032803170*lms1 - 0.0717110580655164*lms2,
		-0.0763729366746601*lms0 - 0.4214933324022432*lms1 + 1.5869240198367816*lms2,
	}
}

func xyzToLinearSRGB(xyz [3]float64) [3]float64 {
	x, y, z := xyz[0], xyz[1], xyz[2]
	return [3]float64{
		(12831.0/3959.0)*x - (329.0/214.0)*y - (1974.0/3959.0)*z,
		(-851781.0/878810.0)*x + (1648619.0/878810.0)*y + (36519.0/878810.0)*z,
		(705.0/12673.0)*x - (2585.0/12673.0)*y + (705.0/667.0)*z,
	}
}

func linearSRGBToOKLab(rgb [3]float64) [3]float64 {
	return xyzToOKLab(linearSRGBToXYZ(rgb))
}

func linearSRGBToXYZ(rgb [3]float64) [3]float64 {
	r, g, b := rgb[0], rgb[1], rgb[2]
	return [3]float64{
		(506752.0/1228815.0)*r + (87881.0/245763.0)*g + (12673.0/70218.0)*b,
		(87098.0/409605.0)*r + (175762.0/245763.0)*g + (12673.0/175545.0)*b,
		(7918.0/409605.0)*r + (87881.0/737289.0)*g + (1001167.0/1053270.0)*b,
	}
}

func xyzToOKLab(xyz [3]float64) [3]float64 {
	x, y, z := xyz[0], xyz[1], xyz[2]
	lms0 := math.Cbrt(0.8190224379967030*x + 0.3619062600528904*y - 0.1288737815209879*z)
	lms1 := math.Cbrt(0.0329836539323885*x + 0.9292868615863434*y + 0.0361446663506424*z)
	lms2 := math.Cbrt(0.0481771893596242*x + 0.2642395317527308*y + 0.6335478284694309*z)
	return [3]float64{
		0.2104542683093140*lms0 + 0.7936177747023054*lms1 - 0.0040720430116193*lms2,
		1.9779985324311684*lms0 - 2.4285922420485799*lms1 + 0.4505937096174110*lms2,
		0.0259040424655478*lms0 + 0.7827717124575296*lms1 - 0.8086757549230774*lms2,
	}
}

func gammaEncode(linear float64) float64 {
	if linear <= 0.0031308 {
		return 12.92 * linear
	}
	return 1.055*math.Pow(linear, 1/2.4) - 0.055
}

func gammaDecode(encoded float64) float64 {
	if encoded <= 0.04045 {
		return encoded / 12.92
	}
	return math.Pow((encoded+0.055)/1.055, 2.4)
}

func clamp01(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}
