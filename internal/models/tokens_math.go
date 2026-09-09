package models

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// HexToRGB parses a hex or rgb/rgba color string to normalized float64 values [0, 1]
func HexToRGB(colorStr string) (r, g, b float64, err error) {
	clean := strings.TrimSpace(colorStr)
	if strings.HasPrefix(clean, "rgb") {
		start := strings.Index(clean, "(")
		end := strings.LastIndex(clean, ")")
		if start != -1 && end != -1 && end > start {
			parts := strings.Split(clean[start+1:end], ",")
			if len(parts) >= 3 {
				pr, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
				pg, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
				pb, err3 := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
				if err1 == nil && err2 == nil && err3 == nil {
					return pr / 255.0, pg / 255.0, pb / 255.0, nil
				}
			}
		}
	}

	clean = strings.TrimPrefix(clean, "#")
	if len(clean) == 3 {
		clean = fmt.Sprintf("%c%c%c%c%c%c", clean[0], clean[0], clean[1], clean[1], clean[2], clean[2])
	}
	if len(clean) != 6 {
		return 0, 0, 0, fmt.Errorf("invalid hex color format: %s", colorStr)
	}

	val, err := strconv.ParseUint(clean, 16, 32)
	if err != nil {
		return 0, 0, 0, err
	}

	r = float64((val>>16)&0xFF) / 255.0
	g = float64((val>>8)&0xFF) / 255.0
	b = float64(val&0xFF) / 255.0
	return r, g, b, nil
}

// RGBToHex formats normalized RGB [0, 1] to #RRGGBB
func RGBToHex(r, g, b float64) string {
	ir := int(math.Round(math.Max(0, math.Min(1, r)) * 255.0))
	ig := int(math.Round(math.Max(0, math.Min(1, g)) * 255.0))
	ib := int(math.Round(math.Max(0, math.Min(1, b)) * 255.0))
	return fmt.Sprintf("#%02x%02x%02x", ir, ig, ib)
}

// RelativeLuminance computes standard WCAG 2.1 relative luminance
func RelativeLuminance(r, g, b float64) float64 {
	convert := func(c float64) float64 {
		if c <= 0.04045 {
			return c / 12.92
		}
		return math.Pow((c+0.055)/1.055, 2.4)
	}
	return 0.2126*convert(r) + 0.7152*convert(g) + 0.0722*convert(b)
}

// ContrastRatio calculates WCAG 2.1 contrast ratio between two colors
func ContrastRatio(hex1, hex2 string) (float64, error) {
	r1, g1, b1, err := HexToRGB(hex1)
	if err != nil {
		return 1.0, err
	}
	r2, g2, b2, err := HexToRGB(hex2)
	if err != nil {
		return 1.0, err
	}

	l1 := RelativeLuminance(r1, g1, b1)
	l2 := RelativeLuminance(r2, g2, b2)

	if l1 < l2 {
		l1, l2 = l2, l1
	}

	ratio := (l1 + 0.05) / (l2 + 0.05)
	return math.Round(ratio*100.0) / 100.0, nil
}

// RGBToHSL converts normalized RGB to HSL [0, 1]
func RGBToHSL(r, g, b float64) (h, s, l float64) {
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	l = (max + min) / 2.0

	if max == min {
		h = 0
		s = 0
		return
	}

	d := max - min
	if l > 0.5 {
		s = d / (2.0 - max - min)
	} else {
		s = d / (max + min)
	}

	switch max {
	case r:
		h = (g - b) / d
		if g < b {
			h += 6.0
		}
	case g:
		h = (b-r)/d + 2.0
	case b:
		h = (r-g)/d + 4.0
	}
	h /= 6.0
	return
}

func hue2rgb(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6.0*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6.0
	}
	return p
}

// HSLToRGB converts normalized HSL [0, 1] to RGB [0, 1]
func HSLToRGB(h, s, l float64) (r, g, b float64) {
	if s == 0 {
		r = l
		g = l
		b = l
		return
	}

	var q float64
	if l < 0.5 {
		q = l * (1.0 + s)
	} else {
		q = l + s - l*s
	}
	p := 2.0*l - q

	r = hue2rgb(p, q, h+1.0/3.0)
	g = hue2rgb(p, q, h)
	b = hue2rgb(p, q, h-1.0/3.0)
	return
}

// RecommendedAction represents an automated contrast repair suggestion
type RecommendedAction struct {
	TargetToken  string  `json:"target_token"`
	CurrentHex   string  `json:"current_hex"`
	SuggestedHex string  `json:"suggested_hex"`
	RatioBefore  float64 `json:"ratio_before"`
	RatioAfter   float64 `json:"ratio_after"`
	Description  string  `json:"description"`
}

// ContrastReport provides accessibility metrics for a specific token pair
type ContrastReport struct {
	FgToken           string             `json:"fg_token"`
	BgToken           string             `json:"bg_token"`
	FgHex             string             `json:"fg_hex"`
	BgHex             string             `json:"bg_hex"`
	Ratio             float64            `json:"ratio"`
	PassesAA          bool               `json:"passes_aa"`
	PassesAAA         bool               `json:"passes_aaa"`
	RecommendedAction *RecommendedAction `json:"recommended_action,omitempty"`
}

// GenerateRecommendedAction finds the minimal lightness delta to satisfy WCAG AA (>= 4.5:1)
func GenerateRecommendedAction(fgHex, bgHex, targetToken string) *RecommendedAction {
	curRatio, err := ContrastRatio(fgHex, bgHex)
	if err != nil || curRatio >= 4.5 {
		return nil
	}

	rBg, gBg, bBg, _ := HexToRGB(bgHex)
	lumBg := RelativeLuminance(rBg, gBg, bBg)

	rFg, gFg, bFg, _ := HexToRGB(fgHex)
	h, s, l := RGBToHSL(rFg, gFg, bFg)

	// Binary search for optimal lightness
	low, high := 0.0, 1.0
	if lumBg < 0.5 {
		low = l // Background is dark, search lighter
	} else {
		high = l // Background is light, search darker
	}

	bestHex := fgHex
	bestRatio := curRatio

	for i := 0; i < 16; i++ {
		mid := (low + high) / 2.0
		rTest, gTest, bTest := HSLToRGB(h, s, mid)
		candHex := RGBToHex(rTest, gTest, bTest)
		ratio, err := ContrastRatio(candHex, bgHex)
		if err == nil {
			if ratio >= 4.55 {
				bestHex = candHex
				bestRatio = ratio
				if lumBg < 0.5 {
					high = mid // Try closer to original
				} else {
					low = mid
				}
			} else {
				if lumBg < 0.5 {
					low = mid // Need more lightness
				} else {
					high = mid // Need more darkness
				}
			}
		}
	}

	direction := "Lighten"
	if lumBg >= 0.5 {
		direction = "Darken"
	}

	return &RecommendedAction{
		TargetToken:  targetToken,
		CurrentHex:   fgHex,
		SuggestedHex: bestHex,
		RatioBefore:  curRatio,
		RatioAfter:   bestRatio,
		Description:  fmt.Sprintf("%s %s to %s (brings contrast from %.1f:1 to %.1f:1)", direction, targetToken, bestHex, curRatio, bestRatio),
	}
}

// EvaluateThemeContrast audits all essential token pairs in a palette
func EvaluateThemeContrast(p ThemePalette) []ContrastReport {
	pairs := []struct {
		fgToken string
		bgToken string
		fgHex   string
		bgHex   string
	}{
		{"textPrimary", "bgCanvas", p.TextPrimary, p.BgCanvas},
		{"textPrimary", "bgSurface", p.TextPrimary, p.BgSurface},
		{"textSecondary", "bgSurface", p.TextSecondary, p.BgSurface},
		{"accentPrimary", "bgCanvas", p.AccentPrimary, p.BgCanvas},
	}

	var reports []ContrastReport
	for _, pair := range pairs {
		ratio, err := ContrastRatio(pair.fgHex, pair.bgHex)
		if err != nil {
			ratio = 1.0
		}

		passesAA := ratio >= 4.5
		passesAAA := ratio >= 7.0

		var rec *RecommendedAction
		if !passesAA {
			rec = GenerateRecommendedAction(pair.fgHex, pair.bgHex, pair.fgToken)
		}

		reports = append(reports, ContrastReport{
			FgToken:           pair.fgToken,
			BgToken:           pair.bgToken,
			FgHex:             pair.fgHex,
			BgHex:             pair.bgHex,
			Ratio:             ratio,
			PassesAA:          passesAA,
			PassesAAA:         passesAAA,
			RecommendedAction: rec,
		})
	}

	return reports
}
