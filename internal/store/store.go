// Package store persists students, sessions, and scores.
package store

import (
	"context"
	"errors"

	"github.com/Dwiyuda/practicum-score-graphql/internal/domain"
)

var ErrNotFound = errors.New("not found")

// Store is implemented by Postgres (production) and Memory (tests and the browser demo).
type Store interface {
	CreateStudent(ctx context.Context, nim, name, className string) (domain.Student, error)
	ListStudents(ctx context.Context, className *string) ([]domain.Student, error)
	GetStudent(ctx context.Context, id int64) (domain.Student, error)

	CreateSession(ctx context.Context, number int, topic, heldOn string) (domain.Session, error)
	ListSessions(ctx context.Context) ([]domain.Session, error)
	GetSession(ctx context.Context, id int64) (domain.Session, error)

	// RecordScore inserts the score or replaces the existing one for the same student and session.
	RecordScore(ctx context.Context, studentID, sessionID int64, value float64, attended bool) (domain.Score, error)
	ScoresByStudent(ctx context.Context, studentID int64) ([]domain.Score, error)
	ScoresByClass(ctx context.Context, className string) ([]domain.Score, error)
}
