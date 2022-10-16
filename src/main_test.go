package main

import (
	"database/sql"
	"testing"
)

func TestLoadControllers(t *testing.T) {
	LoadControlles(&sql.DB{})
}
