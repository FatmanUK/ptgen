package main

import (
	"testing"
	pt "ptgen/internal/protracker"
)

func TestReadPattern(t *testing.T) {
	// func readPattern(file io.Reader, logs chan string, isHexRows bool) (pt.Pattern, error)
}

func TestReadMetadata(t *testing.T) {
	// func readMetadata(proj *pt.ModProject, logs chan string, file *os.File) (*pt.ModProject, error)
}

func TestDefaultOrderlist(t *testing.T) {
	proj := pt.ModProjectFactory()
	pttn := pt.PatternFactory()
	proj.Patterns = []pt.Pattern{ pttn, pttn, pttn, pttn, pttn }
	proj.OrderList = defaultOrderlist(&proj)
	if proj.IsOrderListValid() {
		t.Logf("Order list is ok,")
	} else {
		t.Errorf("Order list is wrong,")
	}
}

func TestCountPatterns(t *testing.T) {
	// func countPatterns(path string) (uint8, error)
}

func TestIsSequentialPatterns(t *testing.T) {
	// func isSequentialPatterns(path string, count uint8) error
}

func TestIsHexRowNotationDetected(t *testing.T) {
	// func isHexRowNotationDetected(path string) (bool, error)
}

func TestValidatePatterns(t *testing.T) {
	// func validatePatterns(path string) (uint8, error)
}
