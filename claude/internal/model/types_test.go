package model

import (
	"encoding/json"
	"testing"
)

func TestNullStringMarshalJSON(t *testing.T) {
	t.Run("valid string marshals as plain string", func(t *testing.T) {
		ns := NewNullString("hello", true)
		data, err := json.Marshal(ns)
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		if string(data) != `"hello"` {
			t.Fatalf("expected %q, got %q", `"hello"`, string(data))
		}
	})

	t.Run("invalid NullString marshals as null", func(t *testing.T) {
		ns := NewNullString("", false)
		data, err := json.Marshal(ns)
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		if string(data) != "null" {
			t.Fatalf("expected null, got %q", string(data))
		}
	})

	t.Run("empty valid string marshals as empty string", func(t *testing.T) {
		ns := NewNullString("", true)
		data, err := json.Marshal(ns)
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		if string(data) != `""` {
			t.Fatalf("expected %q, got %q", `""`, string(data))
		}
	})
}

func TestNullStringUnmarshalJSON(t *testing.T) {
	t.Run("plain string", func(t *testing.T) {
		var ns NullString
		if err := json.Unmarshal([]byte(`"world"`), &ns); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !ns.Valid || ns.String != "world" {
			t.Fatalf("expected Valid=true, String=world; got Valid=%v, String=%q", ns.Valid, ns.String)
		}
	})

	t.Run("null value", func(t *testing.T) {
		var ns NullString
		if err := json.Unmarshal([]byte("null"), &ns); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if ns.Valid {
			t.Fatal("expected Valid=false for null")
		}
	})
}

func TestNullStringInStruct(t *testing.T) {
	type Example struct {
		Name  string     `json:"name"`
		Notes NullString `json:"notes"`
	}

	t.Run("struct with valid NullString", func(t *testing.T) {
		e := Example{Name: "test", Notes: NewNullString("some notes", true)}
		data, err := json.Marshal(e)
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		expected := `{"name":"test","notes":"some notes"}`
		if string(data) != expected {
			t.Fatalf("expected %s, got %s", expected, string(data))
		}
	})

	t.Run("struct with null NullString", func(t *testing.T) {
		e := Example{Name: "test", Notes: NewNullString("", false)}
		data, err := json.Marshal(e)
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		expected := `{"name":"test","notes":null}`
		if string(data) != expected {
			t.Fatalf("expected %s, got %s", expected, string(data))
		}
	})

	t.Run("unmarshal into struct", func(t *testing.T) {
		input := `{"name":"test","notes":"hello"}`
		var e Example
		if err := json.Unmarshal([]byte(input), &e); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if e.Name != "test" || !e.Notes.Valid || e.Notes.String != "hello" {
			t.Fatalf("unexpected: %+v", e)
		}
	})

	t.Run("unmarshal null notes", func(t *testing.T) {
		input := `{"name":"test","notes":null}`
		var e Example
		if err := json.Unmarshal([]byte(input), &e); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if e.Notes.Valid {
			t.Fatal("expected Notes.Valid=false")
		}
	})
}
