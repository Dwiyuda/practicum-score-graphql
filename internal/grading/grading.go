// Package grading holds the score rules, kept apart from storage so they are easy to test.
package grading

import (
	"fmt"
	"math"

	"github.com/Dwiyuda/practicum-score-graphql/internal/domain"
)

// Validate rejects scores outside 0..100 and NaN.
func Validate(value float64) error {
	if math.IsNaN(value) || value < 0 || value > 100 {
		return fmt.Errorf("score must be between 0 and 100, got %v", value)
	}
	return nil
}

// Average is the mean of all recorded values, 0 for an empty slice.
func Average(scores []domain.Score) float64 {
	if len(scores) == 0 {
		return 0
	}
	var sum float64
	for _, s := range scores {
		sum += s.Value
	}
	return round2(sum / float64(len(scores)))
}

// AttendanceRate is attended sessions divided by recorded sessions, 0 for an empty slice.
func AttendanceRate(scores []domain.Score) float64 {
	if len(scores) == 0 {
		return 0
	}
	attended := 0
	for _, s := range scores {
		if s.Attended {
			attended++
		}
	}
	return round2(float64(attended) / float64(len(scores)))
}

// Summarize builds the class-level view from its students and their scores.
func Summarize(className string, students int, scores []domain.Score) domain.ClassSummary {
	return domain.ClassSummary{
		ClassName:      className,
		Students:       students,
		Average:        Average(scores),
		AttendanceRate: AttendanceRate(scores),
	}
}

func round2(x float64) float64 { return math.Round(x*100) / 100 }
