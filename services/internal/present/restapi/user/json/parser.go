package json

import (
	"encoding/json"
	"io"
)

// Parser handles JSON parsing
type Parser struct {
}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{}
}

// AddOrUpdateUser parses an addOrUpdateUser request JSON data into a user request struct
func (p *Parser) AddOrUpdateUser(b io.Reader) (User, error) {
	var user User
	err := json.NewDecoder(b).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// User represents user request body data
type User struct {
	DisplayName string `json:"displayname"`
}
