package graph

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"github.com/Dwiyuda/practicum-score-graphql/internal/store"
)

type gqlResp struct {
	Data   map[string]any `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// do sends one GraphQL request through the real HTTP handler backed by the in-memory store.
func do(t *testing.T, srv http.Handler, query string) gqlResp {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"query": query})
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	var out gqlResp
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("bad response %q: %v", rec.Body.String(), err)
	}
	return out
}

func newServer() http.Handler {
	h := handler.New(NewExecutableSchema(Config{Resolvers: &Resolver{Store: store.NewMemory()}}))
	h.AddTransport(transport.POST{})
	return h
}

func TestScoreFlow(t *testing.T) {
	srv := newServer()
	mustOK := func(r gqlResp) gqlResp {
		t.Helper()
		if len(r.Errors) > 0 {
			t.Fatalf("unexpected errors: %+v", r.Errors)
		}
		return r
	}

	mustOK(do(t, srv, `mutation { createStudent(input: {nim: "1", name: "Budi", className: "A"}) { id } }`))
	mustOK(do(t, srv, `mutation { createStudent(input: {nim: "2", name: "Siti", className: "A"}) { id } }`))
	mustOK(do(t, srv, `mutation { createSession(input: {number: 1, topic: "SELECT", heldOn: "2026-09-01"}) { id } }`))
	mustOK(do(t, srv, `mutation { createSession(input: {number: 2, topic: "JOIN", heldOn: "2026-09-08"}) { id } }`))
	// memory store ids: students 1, 2; sessions 3, 4
	mustOK(do(t, srv, `mutation { recordScore(input: {studentId: 1, sessionId: 3, value: 80, attended: true}) { id } }`))
	mustOK(do(t, srv, `mutation { recordScore(input: {studentId: 1, sessionId: 4, value: 60, attended: false}) { id } }`))
	mustOK(do(t, srv, `mutation { recordScore(input: {studentId: 2, sessionId: 3, value: 100, attended: true}) { id } }`))

	r := mustOK(do(t, srv, `{ student(id: 1) { name average attendanceRate scores { value session { number } } } }`))
	st := r.Data["student"].(map[string]any)
	if st["average"].(float64) != 70 || st["attendanceRate"].(float64) != 0.5 || len(st["scores"].([]any)) != 2 {
		t.Errorf("student = %+v", st)
	}

	r = mustOK(do(t, srv, `{ classSummary(className: "A") { students average attendanceRate } }`))
	sum := r.Data["classSummary"].(map[string]any)
	if sum["students"].(float64) != 2 || sum["average"].(float64) != 80 || sum["attendanceRate"].(float64) != 0.67 {
		t.Errorf("classSummary = %+v", sum)
	}
}

func TestValidation(t *testing.T) {
	srv := newServer()
	cases := []string{
		`mutation { recordScore(input: {studentId: 1, sessionId: 2, value: 101, attended: true}) { id } }`,
		`mutation { createSession(input: {number: 1, topic: "x", heldOn: "01-09-2026"}) { id } }`,
		`mutation { createStudent(input: {nim: " ", name: "x", className: "A"}) { id } }`,
	}
	for _, q := range cases {
		if r := do(t, srv, q); len(r.Errors) == 0 {
			t.Errorf("expected an error for %s", q)
		}
	}
	if r := do(t, srv, `{ student(id: 42) { id } }`); len(r.Errors) > 0 || r.Data["student"] != nil {
		t.Errorf("missing student should be null without errors, got %+v", r)
	}
}
