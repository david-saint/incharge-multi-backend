package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullString is a wrapper around sql.NullString that marshals/unmarshals
// correctly for JSON APIs. It serializes as a plain string or null instead
// of the default {"String":"...","Valid":true} format.
type NullString struct {
	sql.NullString
}

// MarshalJSON serializes NullString as a plain string or null.
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON deserializes a plain string or null into NullString.
func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.String = *s
		ns.Valid = true
	} else {
		ns.Valid = false
	}
	return nil
}

// Scan implements the sql.Scanner interface.
func (ns *NullString) Scan(value interface{}) error {
	return ns.NullString.Scan(value)
}

// Value implements the driver.Valuer interface.
func (ns NullString) Value() (driver.Value, error) {
	return ns.NullString.Value()
}

// NewNullString creates a NullString from a string pointer.
func NewNullString(s string, valid bool) NullString {
	return NullString{sql.NullString{String: s, Valid: valid}}
}
