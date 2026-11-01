package main

import (
	"context"
	"fmt"

	"github.com/Dwiyuda/practicum-score-graphql/internal/store"
)

// seedDemo inserts synthetic students, sessions, and scores. Names are fictional.
// It is skipped when students already exist, so restarting with -seed is safe.
func seedDemo(ctx context.Context, db *store.Postgres) error {
	existing, err := db.ListStudents(ctx, nil)
	if err != nil || len(existing) > 0 {
		return err
	}
	type st struct{ nim, name, class string }
	students := []st{
		{"2201", "Budi Santoso", "A"}, {"2202", "Siti Rahma", "A"}, {"2203", "Andi Pratama", "A"},
		{"2204", "Dewi Lestari", "B"}, {"2205", "Rizky Hidayat", "B"}, {"2206", "Nur Aini", "B"},
	}
	topics := []string{"SELECT dasar", "Filter dan sorting", "JOIN", "Subquery", "View"}

	var sIDs, sesIDs []int64
	for _, s := range students {
		x, err := db.CreateStudent(ctx, s.nim, s.name, s.class)
		if err != nil {
			return err
		}
		sIDs = append(sIDs, x.ID)
	}
	for i, t := range topics {
		x, err := db.CreateSession(ctx, i+1, t, fmt.Sprintf("2026-09-%02d", 1+7*i))
		if err != nil {
			return err
		}
		sesIDs = append(sesIDs, x.ID)
	}
	for i, sid := range sIDs {
		for j, ses := range sesIDs {
			attended := (i+j)%4 != 0
			value := 0.0
			if attended {
				value = float64(60 + (i*7+j*5)%40)
			}
			if _, err := db.RecordScore(ctx, sid, ses, value, attended); err != nil {
				return err
			}
		}
	}
	return nil
}
