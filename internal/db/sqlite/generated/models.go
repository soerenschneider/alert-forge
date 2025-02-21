// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package generated

import (
	"database/sql"
)

type Alert struct {
	ID          int64         `json:"id"`
	Fingerprint string        `json:"fingerprint"`
	Started     int64         `json:"started"`
	Ended       sql.NullInt64 `json:"ended"`
	Severity    int64         `json:"severity"`
	Instance    string        `json:"instance"`
	Name        string        `json:"name"`
	Data        interface{}   `json:"data"`
}

type Report struct {
	ID            int64          `json:"id"`
	Date          sql.NullString `json:"date"`
	SeverityCount sql.NullString `json:"severity_count"`
}

type SchemaVersion struct {
	Version int64 `json:"version"`
}
