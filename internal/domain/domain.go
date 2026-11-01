// Package domain holds the types shared by the store and the GraphQL layer.
package domain

type Student struct {
	ID        int64
	Nim       string
	Name      string
	ClassName string
}

type Session struct {
	ID     int64
	Number int
	Topic  string
	HeldOn string // YYYY-MM-DD
}

type Score struct {
	ID        int64
	StudentID int64
	SessionID int64
	Value     float64
	Attended  bool
}

type ClassSummary struct {
	ClassName      string
	Students       int
	Average        float64
	AttendanceRate float64
}
