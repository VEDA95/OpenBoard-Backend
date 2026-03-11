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
		&models.GeneralSettings{},
		&models.AuthSettings{},
		&models.EmailSettings{},
		&models.User{},
		&models.ExternalAuthProvider{},
		&models.MultiAuthMethod{},
		&models.Session{},
		&models.Role{},
		&models.Permission{},
		&models.PasswordResetToken{},
		&models.EmailVerificationToken{},
		&models.Workspace{},
		&models.Board{},
		&models.List{},
		&models.Card{},
		&models.CheckListItem{},
		&models.CardActivity{},
		&models.Comment{},
		&models.Label{},
		&models.FileUpload{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
