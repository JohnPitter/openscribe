package spreadsheet

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// EvaluateFormula evaluates a simple formula and returns the result
// Supports: SUM, AVERAGE, MIN, MAX, COUNT, ABS, ROUND
func (s *Sheet) EvaluateFormula(formula string) (float64, error) {
	formula = strings.TrimSpace(formula)
	if formula == "" {
		return 0, fmt.Errorf("empty formula")
	}

	// Parse function name and arguments
	parenIdx := strings.Index(formula, "(")
	if parenIdx == -1 {
		// Try to parse as number
		return strconv.ParseFloat(formula, 64)
	}

	funcName := strings.ToUpper(formula[:parenIdx])
	argsStr := formula[parenIdx+1 : len(formula)-1]

	switch funcName {
	case "SUM":
		values, err := s.resolveRange(argsStr)
		if err != nil {
			return 0, err
		}
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		return sum, nil

	case "AVERAGE", "AVG":
		values, err := s.resolveRange(argsStr)
		if err != nil {
			return 0, err
		}
		if len(values) == 0 {
			return 0, fmt.Errorf("AVERAGE: no values")
		}
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		return sum / float64(len(values)), nil

	case "MIN":
		values, err := s.resolveRange(argsStr)
		if err != nil {
			return 0, err
		}
		if len(values) == 0 {
			return 0, fmt.Errorf("MIN: no values")
		}
		min := values[0]
		for _, v := range values[1:] {
			if v < min {
				min = v
			}
		}
		return min, nil

	case "MAX":
		values, err := s.resolveRange(argsStr)
		if err != nil {
			return 0, err
		}
		if len(values) == 0 {
			return 0, fmt.Errorf("MAX: no values")
		}
		max := values[0]
		for _, v := range values[1:] {
			if v > max {
				max = v
			}
		}
		return max, nil

	case "COUNT":
		values, err := s.resolveRange(argsStr)
		if err != nil {
			return 0, err
		}
		return float64(len(values)), nil

	case "ABS":
		values, err := s.resolveRange(argsStr)
		if err != nil {
			return 0, err
		}
		if len(values) != 1 {
			return 0, fmt.Errorf("ABS: expected 1 argument")
		}
		return math.Abs(values[0]), nil

	case "ROUND":
		parts := strings.Split(argsStr, ",")
		if len(parts) != 2 {
			return 0, fmt.Errorf("ROUND: expected 2 arguments")
		}
		values, err := s.resolveRange(strings.TrimSpace(parts[0]))
		if err != nil {
			return 0, err
		}
		if len(values) != 1 {
			return 0, fmt.Errorf("ROUND: first arg must be single value")
		}
		decimals, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return 0, fmt.Errorf("ROUND: invalid decimals: %w", err)
		}
		factor := math.Pow(10, float64(decimals))
		return math.Round(values[0]*factor) / factor, nil

	default:
		return 0, fmt.Errorf("unsupported function: %s", funcName)
	}
}

// resolveRange resolves a cell range like "A1:A5" to a slice of float64 values
func (s *Sheet) resolveRange(rangeStr string) ([]float64, error) {
	rangeStr = strings.TrimSpace(rangeStr)

	// Check if it's a simple number
	if n, err := strconv.ParseFloat(rangeStr, 64); err == nil {
		return []float64{n}, nil
	}

	// Check for range (e.g., "A1:A5")
	parts := strings.Split(rangeStr, ":")
	if len(parts) == 2 {
		startCol, startRow := parseCellRef(parts[0])
		endCol, endRow := parseCellRef(parts[1])

		var values []float64
		for r := startRow; r <= endRow; r++ {
			for c := startCol; c <= endCol; c++ {
				cell := s.Cell(r, c)
				if cell.cellType == CellTypeNumber {
					values = append(values, cell.numVal)
				}
			}
		}
		return values, nil
	}

	// Single cell reference
	col, row := parseCellRef(rangeStr)
	if col > 0 && row > 0 {
		cell := s.Cell(row, col)
		if cell.cellType == CellTypeNumber {
			return []float64{cell.numVal}, nil
		}
		return []float64{0}, nil
	}

	return nil, fmt.Errorf("cannot resolve: %s", rangeStr)
}

// parseCellRef parses "A1" into (col, row) — both 1-based
func parseCellRef(ref string) (col, row int) {
	ref = strings.TrimSpace(ref)
	col = 0
	i := 0
	for i < len(ref) && ref[i] >= 'A' && ref[i] <= 'Z' {
		col = col*26 + int(ref[i]-'A'+1)
		i++
	}
	// Also handle lowercase
	for i < len(ref) && ref[i] >= 'a' && ref[i] <= 'z' {
		col = col*26 + int(ref[i]-'a'+1)
		i++
	}
	if i < len(ref) {
		row, _ = strconv.Atoi(ref[i:])
	}
	return
}
