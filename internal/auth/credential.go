package auth

import (
	"fmt"
	"os"
)

func SaveUserCredentials(username string, email string, password string, fileName string) error {
	content := fmt.Sprintf(
		"OpenBoard Initial User Credentials\n"+
			"====================================\n"+
			"Username: %s\n"+
			"Email: %s\n"+
			"Password: %s\n"+
			"====================================\n"+
			"Please save these credentials and delete this file after noting them down.\n"+
			"You should change the password after first login.\n",
		username, email, password,
	)

	return os.WriteFile(fmt.Sprintf("%s.txt", fileName), []byte(content), 0600)
}
