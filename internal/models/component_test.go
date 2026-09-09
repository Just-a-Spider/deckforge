package models

import (
	"strings"
	"testing"
)

func TestComponentRender(t *testing.T) {
	comps := GetBuiltinComponents()
	if len(comps) != 7 {
		t.Fatalf("expected 7 builtin components, got %d", len(comps))
	}

	// Test df-metric-card
	metricCard := comps[0]
	if metricCard.Selector != "df-metric-card" {
		t.Errorf("expected df-metric-card, got %s", metricCard.Selector)
	}

	html, err := metricCard.Render(map[string]interface{}{
		"value": "+45%",
		"label": "Throughput Boost",
	}, nil)

	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(html, "+45%") {
		t.Errorf("rendered html missing value +45%%: %s", html)
	}
	if !strings.Contains(html, "Throughput Boost") {
		t.Errorf("rendered html missing label Throughput Boost: %s", html)
	}
	if !strings.Contains(html, "glass-panel") {
		t.Errorf("rendered html missing class glass-panel: %s", html)
	}
}
