package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"

	models "VEDA95/open_board/api/internal/db/model"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&models.User{},
		&models.Session{},
		&models.Role{},
		&models.Permission{},
		&models.PasswordResetToken{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
