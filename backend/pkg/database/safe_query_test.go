package database

import (
	"testing"
)

func TestEscapeLike(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal string",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "String with percent",
			input:    "50%",
			expected: "50\\%",
		},
		{
			name:     "String with underscore",
			input:    "test_value",
			expected: "test\\_value",
		},
		{
			name:     "String with backslash",
			input:    "path\\to\\file",
			expected: "path\\\\to\\\\file",
		},
		{
			name:     "String with all special chars",
			input:    "test_%\\value",
			expected: "test\\_\\%\\\\value",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeLike(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeLike(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsValidTableName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid table name",
			input:    "users",
			expected: true,
		},
		{
			name:     "Valid table with underscore",
			input:    "user_accounts",
			expected: true,
		},
		{
			name:     "Valid table with numbers",
			input:    "table123",
			expected: true,
		},
		{
			name:     "Table with semicolon",
			input:    "users; DROP TABLE users",
			expected: false,
		},
		{
			name:     "Table with space",
			input:    "user accounts",
			expected: false,
		},
		{
			name:     "Table with dash",
			input:    "user-accounts",
			expected: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "Table starting with number",
			input:    "123table",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidTableName(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidTableName(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsValidColumnName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid column name",
			input:    "id",
			expected: true,
		},
		{
			name:     "Valid column with underscore",
			input:    "created_at",
			expected: true,
		},
		{
			name:     "Column with special chars",
			input:    "created-at",
			expected: false,
		},
		{
			name:     "Column with SQL injection",
			input:    "id); DROP TABLE users--",
			expected: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidColumnName(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidColumnName(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSafeOrderBy(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected string
	}{
		{
			name:     "Valid field",
			field:    "id",
			expected: "id DESC",
		},
		{
			name:     "Valid field with underscore",
			field:    "created_at",
			expected: "created_at DESC",
		},
		{
			name:     "Invalid field with special chars",
			field:    "id; DROP TABLE users",
			expected: "created_at DESC", // Should return default
		},
		{
			name:     "Empty field",
			field:    "",
			expected: "created_at DESC", // Should return default
		},
		{
			name:     "Field with dash",
			field:    "created-at",
			expected: "created_at DESC", // Should return default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeOrderBy(tt.field)
			if result != tt.expected {
				t.Errorf("SafeOrderBy(%q) = %q, want %q", tt.field, result, tt.expected)
			}
		})
	}
}