package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Value Marshal
func (a CustomFields) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *CustomFields) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
