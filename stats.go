package main

import "math"

func CalculateStats(attempts int, prefixLength int, elapsedSeconds float64) (float64, float64) {
	expected := math.Pow(32, float64(prefixLength))
	hashrate := 0.0

	if elapsedSeconds > 0 {
		hashrate = float64(attempts) / elapsedSeconds
	}

	effort := (float64(attempts) / expected) * 100.0

	return hashrate, effort
}
