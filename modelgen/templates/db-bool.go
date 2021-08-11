package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Bool is a nullable bool.
// It does not consider false values to be null.
// It will decode to null, not false, if null.
type Bool struct {
	Bool  bool
	Valid bool
}

// Scan implements the Scanner interface for reading value from SQL
func (n *Bool) Scan(value interface{}) error {
	if value == nil {
		n.Bool, n.Valid = false, false
		return nil
	}
	s := strings.ToLower(fmt.Sprintf("%v", value))
	switch s {
	case "true", "1", "yes", "y":
		n.Valid = true
		n.Bool = true
	case "false", "0", "no", "n":
		n.Valid = true
		n.Bool = false
	default:
		n.Valid = false
	}
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Bool.
func (b *Bool) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		b.Valid = false
		return nil
	}
	if err := json.Unmarshal(data, &b.Bool); err != nil {
		return fmt.Errorf("null: couldn't unmarshal JSON: %w", err)
	}
	b.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Bool is null.
func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// Equal returns true if both booleans have the same value or are both null.
func (b Bool) Equal(other Bool) bool {
	return b.Valid == other.Valid && (!b.Valid || b.Bool == other.Bool)
}
