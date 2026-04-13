package validator

import (
	"testing"
)

func TestValidateStruct(t *testing.T) {
	type testCase struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	t.Run("valid struct returns nil", func(t *testing.T) {
		tc := testCase{Name: "John", Email: "john@example.com", Password: "secret123"}
		errs := ValidateStruct(tc)
		if errs != nil {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		tc := testCase{}
		errs := ValidateStruct(tc)
		if errs == nil {
			t.Fatal("expected errors for empty struct")
		}
		if _, ok := errs["name"]; !ok {
			t.Error("expected error for name field")
		}
		if _, ok := errs["email"]; !ok {
			t.Error("expected error for email field")
		}
		if _, ok := errs["password"]; !ok {
			t.Error("expected error for password field")
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		tc := testCase{Name: "John", Email: "not-an-email", Password: "secret123"}
		errs := ValidateStruct(tc)
		if errs == nil {
			t.Fatal("expected errors for invalid email")
		}
		if _, ok := errs["email"]; !ok {
			t.Error("expected error for email field")
		}
	})

	t.Run("password too short", func(t *testing.T) {
		tc := testCase{Name: "John", Email: "john@example.com", Password: "abc"}
		errs := ValidateStruct(tc)
		if errs == nil {
			t.Fatal("expected errors for short password")
		}
		if _, ok := errs["password"]; !ok {
			t.Error("expected error for password field")
		}
	})
}

func TestPhoneNgUsValidation(t *testing.T) {
	type phoneCase struct {
		Phone string `json:"phone" validate:"required,phone_ng_us"`
	}

	tests := []struct {
		name  string
		phone string
		valid bool
	}{
		{"NG with +234", "+2348012345678", true},
		{"NG with 234", "2348012345678", true},
		{"NG with 0", "08012345678", true},
		{"NG starting with 7", "07012345678", true},
		{"NG starting with 9", "09012345678", true},
		{"US 10 digit", "2125551234", true},
		{"US with +1", "+12125551234", true},
		{"US with 1", "12125551234", true},
		{"too short", "12345", false},
		{"letters", "abcdefghij", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := phoneCase{Phone: tt.phone}
			errs := ValidateStruct(tc)
			if tt.valid && errs != nil {
				t.Fatalf("expected %q to be valid, got errors: %v", tt.phone, errs)
			}
			if !tt.valid && errs == nil {
				t.Fatalf("expected %q to be invalid, but got no errors", tt.phone)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Name", "name"},
		{"Email", "email"},
		{"DateOfBirth", "date_of_birth"},
		{"MaritalStatus", "marital_status"},
		{"ID", "i_d"},
		{"UserType", "user_type"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toSnakeCase(tt.input)
			if result != tt.expected {
				t.Fatalf("toSnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
