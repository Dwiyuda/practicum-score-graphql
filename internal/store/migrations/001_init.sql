CREATE TABLE students (
    id         BIGSERIAL PRIMARY KEY,
    nim        TEXT NOT NULL UNIQUE,
    name       TEXT NOT NULL,
    class_name TEXT NOT NULL
);

CREATE TABLE sessions (
    id      BIGSERIAL PRIMARY KEY,
    number  INT  NOT NULL UNIQUE CHECK (number > 0),
    topic   TEXT NOT NULL,
    held_on DATE NOT NULL
);

CREATE TABLE scores (
    id         BIGSERIAL PRIMARY KEY,
    student_id BIGINT  NOT NULL REFERENCES students (id) ON DELETE CASCADE,
    session_id BIGINT  NOT NULL REFERENCES sessions (id) ON DELETE CASCADE,
    value      NUMERIC(5, 2) NOT NULL CHECK (value BETWEEN 0 AND 100),
    attended   BOOLEAN NOT NULL,
    UNIQUE (student_id, session_id)
);

CREATE INDEX scores_student_idx ON scores (student_id);
CREATE INDEX students_class_idx ON students (class_name);
