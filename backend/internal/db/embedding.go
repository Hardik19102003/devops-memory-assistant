package db

import (
	"github.com/pgvector/pgvector-go"
)

// ConvertFloat64ToFloat32 converts []float64 to []float32 for pgvector compatibility
func ConvertFloat64ToFloat32(input []float64) []float32 {
	if input == nil {
		return nil
	}
	output := make([]float32, len(input))
	for i, v := range input {
		output[i] = float32(v)
	}
	return output
}

// ConvertToPgVector converts []float64 to *pgvector.Vector
func ConvertToPgVector(input []float64) *pgvector.Vector {
	if input == nil {
		return nil
	}
	vec := pgvector.NewVector(ConvertFloat64ToFloat32(input))
	return &vec
}