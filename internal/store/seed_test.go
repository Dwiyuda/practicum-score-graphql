package store

import (
	"context"
	"testing"
)

func TestSeedIsIdempotent(t *testing.T) {
	ctx := context.Background()
	db := NewMemory()
	for range 2 {
		if err := Seed(ctx, db); err != nil {
			t.Fatal(err)
		}
	}
	students, err := db.ListStudents(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(students) != 6 {
		t.Fatalf("got %d students, want 6", len(students))
	}
	// README examples use studentId 1 and sessionId 2, so IDs must count per table as in Postgres.
	if _, err := db.RecordScore(ctx, 1, 2, 87.5, true); err != nil {
		t.Fatal(err)
	}
	scores, err := db.ScoresByClass(ctx, "A")
	if err != nil {
		t.Fatal(err)
	}
	if len(scores) != 15 {
		t.Fatalf("got %d scores for class A, want 15 (3 students x 5 sessions)", len(scores))
	}
}
