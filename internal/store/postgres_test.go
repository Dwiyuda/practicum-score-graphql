package store

import (
	"context"
	"os"
	"testing"
)

// Runs only when DATABASE_URL points at a disposable database (CI provides one).
func TestPostgresRoundTrip(t *testing.T) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}
	ctx := context.Background()
	p, err := NewPostgres(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if err := p.Migrate(ctx); err != nil {
		t.Fatal(err)
	}
	if err := p.Migrate(ctx); err != nil { // second run must be a no-op
		t.Fatal(err)
	}
	if _, err := p.db.Exec(ctx, `TRUNCATE scores, sessions, students RESTART IDENTITY CASCADE`); err != nil {
		t.Fatal(err)
	}

	st, err := p.CreateStudent(ctx, "223510001", "Budi Santoso", "A")
	if err != nil {
		t.Fatal(err)
	}
	se, err := p.CreateSession(ctx, 1, "SELECT dasar", "2026-09-01")
	if err != nil {
		t.Fatal(err)
	}
	if se.HeldOn != "2026-09-01" {
		t.Errorf("HeldOn = %q", se.HeldOn)
	}
	first, err := p.RecordScore(ctx, st.ID, se.ID, 70, true)
	if err != nil {
		t.Fatal(err)
	}
	again, err := p.RecordScore(ctx, st.ID, se.ID, 85.5, false)
	if err != nil {
		t.Fatal(err)
	}
	if again.ID != first.ID {
		t.Errorf("upsert created a new row: %d vs %d", again.ID, first.ID)
	}
	scores, err := p.ScoresByClass(ctx, "A")
	if err != nil {
		t.Fatal(err)
	}
	if len(scores) != 1 || scores[0].Value != 85.5 || scores[0].Attended {
		t.Errorf("scores = %+v", scores)
	}
	if _, err := p.GetStudent(ctx, 999); err != ErrNotFound {
		t.Errorf("GetStudent(999) err = %v, want ErrNotFound", err)
	}
	if _, err := p.RecordScore(ctx, st.ID, se.ID, 120, true); err == nil {
		t.Error("database accepted a score above 100")
	}
}
