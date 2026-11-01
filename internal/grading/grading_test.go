package grading

import (
	"math"
	"testing"

	"github.com/Dwiyuda/practicum-score-graphql/internal/domain"
)

func TestValidate(t *testing.T) {
	for _, ok := range []float64{0, 55.5, 100} {
		if err := Validate(ok); err != nil {
			t.Errorf("Validate(%v) = %v, want nil", ok, err)
		}
	}
	for _, bad := range []float64{-1, 100.01, math.NaN()} {
		if Validate(bad) == nil {
			t.Errorf("Validate(%v) = nil, want error", bad)
		}
	}
}

func TestAverageAndAttendance(t *testing.T) {
	scores := []domain.Score{
		{Value: 80, Attended: true},
		{Value: 90, Attended: true},
		{Value: 0, Attended: false},
	}
	if got := Average(scores); got != 56.67 {
		t.Errorf("Average = %v, want 56.67", got)
	}
	if got := AttendanceRate(scores); got != 0.67 {
		t.Errorf("AttendanceRate = %v, want 0.67", got)
	}
}

func TestEmpty(t *testing.T) {
	if Average(nil) != 0 || AttendanceRate(nil) != 0 {
		t.Error("empty slices must give 0")
	}
}

func TestSummarize(t *testing.T) {
	s := Summarize("A", 2, []domain.Score{{Value: 70, Attended: true}, {Value: 90, Attended: false}})
	if s.ClassName != "A" || s.Students != 2 || s.Average != 80 || s.AttendanceRate != 0.5 {
		t.Errorf("Summarize = %+v", s)
	}
}
