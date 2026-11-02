package store

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/Dwiyuda/practicum-score-graphql/internal/domain"
)

// Memory is an in-process Store used by the resolver tests and the browser demo.
type Memory struct {
	mu sync.Mutex
	// One sequence per table, like the SERIAL columns in Postgres.
	nextStudent, nextSession, nextScore int64
	students                            map[int64]domain.Student
	sessions                            map[int64]domain.Session
	scores                              map[[2]int64]domain.Score
}

func NewMemory() *Memory {
	return &Memory{
		students: map[int64]domain.Student{},
		sessions: map[int64]domain.Session{},
		scores:   map[[2]int64]domain.Score{},
	}
}

func (m *Memory) CreateStudent(_ context.Context, nim, name, className string) (domain.Student, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, s := range m.students {
		if s.Nim == nim {
			return domain.Student{}, fmt.Errorf("nim %s already exists", nim)
		}
	}
	m.nextStudent++
	s := domain.Student{ID: m.nextStudent, Nim: nim, Name: name, ClassName: className}
	m.students[s.ID] = s
	return s, nil
}

func (m *Memory) ListStudents(_ context.Context, className *string) ([]domain.Student, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []domain.Student
	for _, s := range m.students {
		if className == nil || s.ClassName == *className {
			out = append(out, s)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (m *Memory) GetStudent(_ context.Context, id int64) (domain.Student, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.students[id]
	if !ok {
		return s, ErrNotFound
	}
	return s, nil
}

func (m *Memory) CreateSession(_ context.Context, number int, topic, heldOn string) (domain.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextSession++
	s := domain.Session{ID: m.nextSession, Number: number, Topic: topic, HeldOn: heldOn}
	m.sessions[s.ID] = s
	return s, nil
}

func (m *Memory) ListSessions(_ context.Context) ([]domain.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []domain.Session
	for _, s := range m.sessions {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Number < out[j].Number })
	return out, nil
}

func (m *Memory) GetSession(_ context.Context, id int64) (domain.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.sessions[id]
	if !ok {
		return s, ErrNotFound
	}
	return s, nil
}

func (m *Memory) RecordScore(_ context.Context, studentID, sessionID int64, value float64, attended bool) (domain.Score, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.students[studentID]; !ok {
		return domain.Score{}, ErrNotFound
	}
	if _, ok := m.sessions[sessionID]; !ok {
		return domain.Score{}, ErrNotFound
	}
	key := [2]int64{studentID, sessionID}
	s, ok := m.scores[key]
	if !ok {
		m.nextScore++
		s = domain.Score{ID: m.nextScore, StudentID: studentID, SessionID: sessionID}
	}
	s.Value, s.Attended = value, attended
	m.scores[key] = s
	return s, nil
}

func (m *Memory) ScoresByStudent(_ context.Context, studentID int64) ([]domain.Score, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []domain.Score
	for k, s := range m.scores {
		if k[0] == studentID {
			out = append(out, s)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].SessionID < out[j].SessionID })
	return out, nil
}

func (m *Memory) ScoresByClass(_ context.Context, className string) ([]domain.Score, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []domain.Score
	for _, s := range m.scores {
		if m.students[s.StudentID].ClassName == className {
			out = append(out, s)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
