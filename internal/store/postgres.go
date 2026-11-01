package store

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Dwiyuda/practicum-score-graphql/internal/domain"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Postgres struct{ db *pgxpool.Pool }

func NewPostgres(ctx context.Context, url string) (*Postgres, error) {
	db, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return &Postgres{db: db}, nil
}

func (p *Postgres) Close() { p.db.Close() }

// Migrate applies every embedded migration that has not run yet, in file name order.
func (p *Postgres) Migrate(ctx context.Context) error {
	if _, err := p.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (name TEXT PRIMARY KEY)`); err != nil {
		return err
	}
	names, err := fs.Glob(migrations, "migrations/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)
	for _, name := range names {
		var done bool
		if err := p.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE name = $1)`, name).Scan(&done); err != nil {
			return err
		}
		if done {
			continue
		}
		sql, err := migrations.ReadFile(name)
		if err != nil {
			return err
		}
		tx, err := p.db.Begin(ctx)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, string(sql)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("%s: %w", name, err)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (name) VALUES ($1)`, name); err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (p *Postgres) CreateStudent(ctx context.Context, nim, name, className string) (domain.Student, error) {
	s := domain.Student{Nim: nim, Name: name, ClassName: className}
	err := p.db.QueryRow(ctx,
		`INSERT INTO students (nim, name, class_name) VALUES ($1, $2, $3) RETURNING id`,
		nim, name, className).Scan(&s.ID)
	return s, err
}

func (p *Postgres) ListStudents(ctx context.Context, className *string) ([]domain.Student, error) {
	rows, err := p.db.Query(ctx,
		`SELECT id, nim, name, class_name FROM students
		 WHERE $1::text IS NULL OR class_name = $1 ORDER BY class_name, nim`, className)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (domain.Student, error) {
		var s domain.Student
		err := r.Scan(&s.ID, &s.Nim, &s.Name, &s.ClassName)
		return s, err
	})
}

func (p *Postgres) GetStudent(ctx context.Context, id int64) (domain.Student, error) {
	var s domain.Student
	err := p.db.QueryRow(ctx, `SELECT id, nim, name, class_name FROM students WHERE id = $1`, id).
		Scan(&s.ID, &s.Nim, &s.Name, &s.ClassName)
	return s, notFound(err)
}

func (p *Postgres) CreateSession(ctx context.Context, number int, topic, heldOn string) (domain.Session, error) {
	s := domain.Session{Number: number, Topic: topic}
	err := p.db.QueryRow(ctx,
		`INSERT INTO sessions (number, topic, held_on) VALUES ($1, $2, $3::date) RETURNING id, held_on::text`,
		number, topic, heldOn).Scan(&s.ID, &s.HeldOn)
	return s, err
}

func (p *Postgres) ListSessions(ctx context.Context) ([]domain.Session, error) {
	rows, err := p.db.Query(ctx, `SELECT id, number, topic, held_on::text FROM sessions ORDER BY number`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, scanSession)
}

func (p *Postgres) GetSession(ctx context.Context, id int64) (domain.Session, error) {
	rows, err := p.db.Query(ctx, `SELECT id, number, topic, held_on::text FROM sessions WHERE id = $1`, id)
	if err != nil {
		return domain.Session{}, err
	}
	s, err := pgx.CollectExactlyOneRow(rows, scanSession)
	return s, notFound(err)
}

func (p *Postgres) RecordScore(ctx context.Context, studentID, sessionID int64, value float64, attended bool) (domain.Score, error) {
	s := domain.Score{StudentID: studentID, SessionID: sessionID, Value: value, Attended: attended}
	err := p.db.QueryRow(ctx,
		`INSERT INTO scores (student_id, session_id, value, attended) VALUES ($1, $2, $3, $4)
		 ON CONFLICT (student_id, session_id) DO UPDATE SET value = EXCLUDED.value, attended = EXCLUDED.attended
		 RETURNING id`,
		studentID, sessionID, value, attended).Scan(&s.ID)
	return s, err
}

func (p *Postgres) ScoresByStudent(ctx context.Context, studentID int64) ([]domain.Score, error) {
	rows, err := p.db.Query(ctx,
		`SELECT id, student_id, session_id, value, attended FROM scores WHERE student_id = $1 ORDER BY session_id`, studentID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, scanScore)
}

func (p *Postgres) ScoresByClass(ctx context.Context, className string) ([]domain.Score, error) {
	rows, err := p.db.Query(ctx,
		`SELECT sc.id, sc.student_id, sc.session_id, sc.value, sc.attended
		 FROM scores sc JOIN students st ON st.id = sc.student_id
		 WHERE st.class_name = $1 ORDER BY sc.id`, className)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, scanScore)
}

func scanSession(r pgx.CollectableRow) (domain.Session, error) {
	var s domain.Session
	err := r.Scan(&s.ID, &s.Number, &s.Topic, &s.HeldOn)
	return s, err
}

func scanScore(r pgx.CollectableRow) (domain.Score, error) {
	var s domain.Score
	err := r.Scan(&s.ID, &s.StudentID, &s.SessionID, &s.Value, &s.Attended)
	return s, err
}

func notFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
