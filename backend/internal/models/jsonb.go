package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONB is a custom type for PostgreSQL jsonb fields
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	
	return json.Unmarshal(bytes, j)
}

// ConvertMapToJSONB converts map[string]string to JSONB
func ConvertMapToJSONB(m map[string]string) JSONB {
	result := make(JSONB)
	for k, v := range m {
		result[k] = v
	}
	return result
}