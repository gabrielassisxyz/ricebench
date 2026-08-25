package color

import (
	"encoding/json"
	"flag"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"testing"
)

var updateGolden = flag.Bool("update", false, "update the resolved golden file")

type referenceVector struct {
	Name  string `json:"name"`
	OKLCH struct {
		Lightness float64 `json:"lightness"`
		Chroma    float64 `json:"chroma"`
		Hue       float64 `json:"hue"`
	} `json:"oklch"`
	SRGB   string `json:"srgb"`
	Mapped bool   `json:"mapped"`
	From   *OKLCH `json:"from"`
}

type contrastVector struct {
	A     string  `json:"a"`
	B     string  `json:"b"`
	Ratio float64 `json:"ratio"`
}

func readReferenceVectors(t *testing.T) []referenceVector {
	t.Helper()
	return decodeJSON[[]referenceVector](t, "reference-vectors.json")
}

func readContrastVectors(t *testing.T) []contrastVector {
	t.Helper()
	return decodeJSON[[]contrastVector](t, "contrast-vectors.json")
}

func decodeJSON[T any](t *testing.T, name string) T {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	var value T
	if err := json.Unmarshal(contents, &value); err != nil {
		t.Fatalf("unmarshal %s: %v", name, err)
	}
	return value
}

func mustOKLCH(t *testing.T, lightness, chroma, hue float64) OKLCH {
	t.Helper()
	color, err := NewOKLCH(lightness, chroma, hue)
	if err != nil {
		t.Fatalf("NewOKLCH(%v, %v, %v): %v", lightness, chroma, hue, err)
	}
	return color
}

// TestReferenceVectorsResolveToFixedHex converts every committed reference vector and
// asserts the exact resolved hexadecimal and mapping metadata. The table is test data,
// committed rather than generated at test time, so a conversion change breaks a number
// that was fixed outside this test.
func TestReferenceVectorsResolveToFixedHex(t *testing.T) {
	for _, vector := range readReferenceVectors(t) {
		t.Run(vector.Name, func(t *testing.T) {
			c := mustOKLCH(t, vector.OKLCH.Lightness, vector.OKLCH.Chroma, vector.OKLCH.Hue)
			resolved := Resolve(c)

			if resolved.SRGB != vector.SRGB {
				t.Fatalf("resolved %q, want %q", resolved.SRGB, vector.SRGB)
			}
			if resolved.Mapped != vector.Mapped {
				t.Fatalf("mapped = %v, want %v", resolved.Mapped, vector.Mapped)
			}
			if (resolved.From == nil) != (vector.From == nil) {
				t.Fatalf("from = %v, want %v", resolved.From, vector.From)
			}
			if vector.From != nil && *resolved.From != *vector.From {
				t.Fatalf("from = %+v, want %+v", *resolved.From, *vector.From)
			}

			if got := ToSRGB(c).Hex(); got != vector.SRGB {
				t.Fatalf("ToSRGB hex = %q, want %q", got, vector.SRGB)
			}
		})
	}
}

// TestResolvedOutputMatchesGolden snapshots the serialized resolved output of the
// reference table, so a change in rounding shows as a reviewable diff rather than a
// silently different number.
func TestResolvedOutputMatchesGolden(t *testing.T) {
	resolved := make([]Resolved, 0, len(readReferenceVectors(t)))
	for _, vector := range readReferenceVectors(t) {
		c := mustOKLCH(t, vector.OKLCH.Lightness, vector.OKLCH.Chroma, vector.OKLCH.Hue)
		resolved = append(resolved, Resolve(c))
	}

	serialized, err := json.MarshalIndent(resolved, "", "  ")
	if err != nil {
		t.Fatalf("marshal resolved: %v", err)
	}
	serialized = append(serialized, '\n')

	path := filepath.Join("testdata", "resolved.golden.json")
	if *updateGolden {
		if err := os.WriteFile(path, serialized, 0o644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
		return
	}
	golden, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	if string(serialized) != string(golden) {
		t.Fatalf("resolved output differs from golden:\n%s", serialized)
	}
}

func TestNewOKLCHRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name      string
		lightness float64
		chroma    float64
		hue       float64
		wantErr   string
	}{
		{name: "NaN lightness", lightness: math.NaN(), chroma: 0.1, hue: 0, wantErr: "oklch lightness: NaN is not finite"},
		{name: "infinite lightness", lightness: math.Inf(1), chroma: 0.1, hue: 0, wantErr: "oklch lightness: +Inf is not finite"},
		{name: "lightness below zero", lightness: -0.01, chroma: 0.1, hue: 0, wantErr: "oklch lightness -0.01 out of range [0, 1]"},
		{name: "lightness above one", lightness: 1.01, chroma: 0.1, hue: 0, wantErr: "oklch lightness 1.01 out of range [0, 1]"},
		{name: "NaN chroma", lightness: 0.5, chroma: math.NaN(), hue: 0, wantErr: "oklch chroma: NaN is not finite"},
		{name: "negative chroma", lightness: 0.5, chroma: -0.1, hue: 0, wantErr: "oklch chroma -0.1 out of range [0, inf)"},
		{name: "infinite hue", lightness: 0.5, chroma: 0.1, hue: math.Inf(1), wantErr: "oklch hue: +Inf is not finite"},
		{name: "NaN hue", lightness: 0.5, chroma: 0.1, hue: math.NaN(), wantErr: "oklch hue: NaN is not finite"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewOKLCH(test.lightness, test.chroma, test.hue)
			if err == nil {
				t.Fatal("accepted invalid OKLCH input")
			}
			if err.Error() != test.wantErr {
				t.Fatalf("error %q, want %q", err, test.wantErr)
			}
		})
	}
}

func TestNewOKLCHAcceptsBoundsAndWrapsHue(t *testing.T) {
	tests := []struct {
		name      string
		lightness float64
		chroma    float64
		hue       float64
		want      OKLCH
	}{
		{name: "lightness at zero", lightness: 0, chroma: 0.1, hue: 50, want: OKLCH{Lightness: 0, Chroma: 0.1, Hue: 50}},
		{name: "lightness at one", lightness: 1, chroma: 0.1, hue: 50, want: OKLCH{Lightness: 1, Chroma: 0.1, Hue: 50}},
		{name: "chroma at zero", lightness: 0.5, chroma: 0, hue: 50, want: OKLCH{Lightness: 0.5, Chroma: 0, Hue: 50}},
		{name: "hue wraps from 360", lightness: 0.5, chroma: 0.1, hue: 360, want: OKLCH{Lightness: 0.5, Chroma: 0.1, Hue: 0}},
		{name: "hue wraps negative", lightness: 0.5, chroma: 0.1, hue: -30, want: OKLCH{Lightness: 0.5, Chroma: 0.1, Hue: 330}},
		{name: "hue wraps over full turns", lightness: 0.5, chroma: 0.1, hue: 725, want: OKLCH{Lightness: 0.5, Chroma: 0.1, Hue: 5}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewOKLCH(test.lightness, test.chroma, test.hue)
			if err != nil {
				t.Fatalf("NewOKLCH: %v", err)
			}
			if got != test.want {
				t.Fatalf("got %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestParseHexRejectsMalformedInput(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr string
	}{
		{name: "missing hash", value: "123456", wantErr: `srgb "123456": want six hexadecimal digits with a leading #`},
		{name: "too short", value: "#12345", wantErr: `srgb "#12345": want six hexadecimal digits with a leading #`},
		{name: "too long", value: "#1234567", wantErr: `srgb "#1234567": want six hexadecimal digits with a leading #`},
		{name: "non-hex digit", value: "#12g456", wantErr: `srgb "#12g456": channel 1 is not hexadecimal`},
		{name: "empty", value: "", wantErr: `srgb "": want six hexadecimal digits with a leading #`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ParseHex(test.value)
			if err == nil {
				t.Fatal("accepted malformed hex")
			}
			if err.Error() != test.wantErr {
				t.Fatalf("error %q, want %q", err, test.wantErr)
			}
		})
	}
}

func TestSRGBHexRoundTrip(t *testing.T) {
	for _, hex := range []string{"#000000", "#ffffff", "#ff0000", "#00ff00", "#0000ff", "#636363", "#abcdef"} {
		color, err := ParseHex(hex)
		if err != nil {
			t.Fatalf("ParseHex(%q): %v", hex, err)
		}
		if got := color.Hex(); got != hex {
			t.Fatalf("Hex() = %q, want %q", got, hex)
		}
	}
}

func TestOKLCHSerializationRoundTrip(t *testing.T) {
	for _, c := range []OKLCH{
		mustOKLCH(t, 0.6279553606145516, 0.2576833074911568, 29.233885192342633),
		mustOKLCH(t, 0.5, 0, 0),
		mustOKLCH(t, 1, 0.2, 50),
	} {
		serialized, err := json.Marshal(c)
		if err != nil {
			t.Fatalf("marshal %+v: %v", c, err)
		}
		var roundTripped OKLCH
		if err := json.Unmarshal(serialized, &roundTripped); err != nil {
			t.Fatalf("unmarshal %s: %v", serialized, err)
		}
		// Documented precision: lightness, chroma and hue are float64 and round-trip
		// through JSON exactly.
		if roundTripped != c {
			t.Fatalf("round trip %+v -> %s -> %+v", c, serialized, roundTripped)
		}
	}
}

func TestInGamutPassesThroughUnchanged(t *testing.T) {
	for _, c := range []OKLCH{
		mustOKLCH(t, 1, 0, 0),   // white
		mustOKLCH(t, 0, 0, 0),   // black
		mustOKLCH(t, 0.5, 0, 0), // mid gray
		mustOKLCH(t, 0.7, 0.1, 210),
	} {
		if mapped := GamutMap(c); mapped != c {
			t.Fatalf("in-gamut color %+v was mapped to %+v", c, mapped)
		}
	}
}

func TestGamutMapBoundaryLightness(t *testing.T) {
	white := GamutMap(mustOKLCH(t, 1, 0.2, 50))
	if white != (OKLCH{Lightness: 1, Chroma: 0, Hue: 0}) {
		t.Fatalf("lightness 1 mapped to %+v, want white", white)
	}
	black := GamutMap(mustOKLCH(t, 0, 0.2, 50))
	if black != (OKLCH{Lightness: 0, Chroma: 0, Hue: 0}) {
		t.Fatalf("lightness 0 mapped to %+v, want black", black)
	}
}

// TestGamutMapProperties checks the gamut mapping invariants over random colors. The
// local MINDE refinement may return a clipped color whose lightness deviates by at most
// one deltaEOK JND and whose chroma shifts by a float-sized clip artifact, so the bounds
// below are the algorithm's documented guarantees rather than exact equalities.
func TestGamutMapProperties(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for iteration := 0; iteration < 20000; iteration++ {
		c, err := NewOKLCH(rng.Float64(), rng.Float64()*0.6, rng.Float64()*360)
		if err != nil {
			t.Fatalf("NewOKLCH: %v", err)
		}
		if c.Chroma == 0 {
			continue
		}

		mapped := GamutMap(c)
		if inGamut(c) {
			if mapped != c {
				t.Fatalf("in-gamut color %+v was not preserved: %+v", c, mapped)
			}
			continue
		}

		if !inGamut(mapped) {
			t.Fatalf("mapped %+v of out-of-gamut %+v is not in gamut", mapped, c)
		}
		if mapped.Chroma > c.Chroma+1e-3 {
			t.Fatalf("mapped chroma %v exceeds input chroma %v", mapped.Chroma, c.Chroma)
		}
		if drift := math.Abs(mapped.Lightness - c.Lightness); drift > JND+1e-6 {
			t.Fatalf("mapped lightness %v drifted %v from %v, beyond JND %v", mapped.Lightness, drift, c.Lightness, JND)
		}
	}
}

func TestContrastRatioMatchesWCAG(t *testing.T) {
	for _, vector := range readContrastVectors(t) {
		t.Run(vector.A+"-vs-"+vector.B, func(t *testing.T) {
			a, err := ParseHex(vector.A)
			if err != nil {
				t.Fatalf("ParseHex(%q): %v", vector.A, err)
			}
			b, err := ParseHex(vector.B)
			if err != nil {
				t.Fatalf("ParseHex(%q): %v", vector.B, err)
			}
			if got := ContrastRatio(a, b); math.Abs(got-vector.Ratio) > 1e-9 {
				t.Fatalf("contrast = %.9f, want %.9f", got, vector.Ratio)
			}
		})
	}
}

// TestContrastRatioBoundaries verifies the WCAG threshold ratios at exactly 4.5:1 and
// 3:1. A relative luminance of 0.175 against black is exactly 4.5:1, and 0.10 against
// black is exactly 3:1.
func TestContrastRatioBoundaries(t *testing.T) {
	tests := []struct {
		name    string
		lighter float64
		darker  float64
		want    float64
	}{
		{name: "4.5 boundary", lighter: 0.175, darker: 0, want: 4.5},
		{name: "4.5 boundary reversed", lighter: 0, darker: 0.175, want: 4.5},
		{name: "3.0 boundary", lighter: 0.10, darker: 0, want: 3.0},
		{name: "white on black", lighter: 1, darker: 0, want: 21.0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := contrastRatio(test.lighter, test.darker); math.Abs(got-test.want) > 1e-9 {
				t.Fatalf("contrastRatio(%v, %v) = %v, want %v", test.lighter, test.darker, got, test.want)
			}
		})
	}
}

func TestDeltaEOKIsEuclideanInOKLab(t *testing.T) {
	reference := [3]float64{0.5, 0.1, 0.2}
	sample := [3]float64{0.5, 0.3, 0.2}
	want := 0.2
	if got := deltaEOK(reference, sample); math.Abs(got-want) > 1e-12 {
		t.Fatalf("deltaEOK = %v, want %v", got, want)
	}
}
