package repository

import (
	"math"
	"testing"
)

func TestHaversineDistance(t *testing.T) {
	t.Run("Lagos to Abuja is approximately 534 km (great-circle)", func(t *testing.T) {
		// Lagos: 6.5244, 3.3792
		// Abuja: 9.0579, 7.4951
		// Great-circle distance is ~534 km (road distance is ~764 km but haversine gives straight-line).
		dist := haversineDistance(6.5244, 3.3792, 9.0579, 7.4951, earthRadiusKm)
		if dist < 520 || dist > 550 {
			t.Fatalf("expected ~534km, got %.2f km", dist)
		}
	})

	t.Run("same point returns zero", func(t *testing.T) {
		dist := haversineDistance(40.7128, -74.0060, 40.7128, -74.0060, earthRadiusKm)
		if dist > 0.001 {
			t.Fatalf("expected ~0, got %.6f", dist)
		}
	})

	t.Run("km vs miles conversion", func(t *testing.T) {
		distKm := haversineDistance(6.5244, 3.3792, 9.0579, 7.4951, earthRadiusKm)
		distMi := haversineDistance(6.5244, 3.3792, 9.0579, 7.4951, earthRadiusMi)
		// 1 km ≈ 0.621371 miles
		ratio := distMi / distKm
		expectedRatio := earthRadiusMi / earthRadiusKm
		if math.Abs(ratio-expectedRatio) > 0.001 {
			t.Fatalf("km/mi ratio: expected %.4f, got %.4f", expectedRatio, ratio)
		}
	})

	t.Run("London to New York is approximately 5570 km", func(t *testing.T) {
		// London: 51.5074, -0.1278
		// New York: 40.7128, -74.0060
		dist := haversineDistance(51.5074, -0.1278, 40.7128, -74.0060, earthRadiusKm)
		if dist < 5500 || dist > 5600 {
			t.Fatalf("expected ~5570km, got %.2f km", dist)
		}
	})

	t.Run("antipodal points are approximately half circumference", func(t *testing.T) {
		// (0,0) to (0,180) should be approximately pi * R
		dist := haversineDistance(0, 0, 0, 180, earthRadiusKm)
		expected := math.Pi * earthRadiusKm
		if math.Abs(dist-expected) > 1 {
			t.Fatalf("expected %.2f, got %.2f", expected, dist)
		}
	})
}

func TestDegToRad(t *testing.T) {
	tests := []struct {
		deg float64
		rad float64
	}{
		{0, 0},
		{90, math.Pi / 2},
		{180, math.Pi},
		{360, 2 * math.Pi},
		{-90, -math.Pi / 2},
	}

	for _, tt := range tests {
		result := degToRad(tt.deg)
		if math.Abs(result-tt.rad) > 1e-10 {
			t.Errorf("degToRad(%f) = %f, want %f", tt.deg, result, tt.rad)
		}
	}
}

func TestIsAllowedSortColumn(t *testing.T) {
	allowed := []string{"id", "name", "address", "created_at", "updated_at"}
	for _, col := range allowed {
		if !isAllowedSortColumn(col) {
			t.Errorf("%q should be allowed", col)
		}
	}

	forbidden := []string{"password", "email", "phone", "deleted_at", "distance", "'; DROP TABLE users;--"}
	for _, col := range forbidden {
		if isAllowedSortColumn(col) {
			t.Errorf("%q should not be allowed", col)
		}
	}
}
