package db

import (
	"VEDA95/open_board/api/internal/auth"
	models "VEDA95/open_board/api/internal/db/model"
	"os"

	"gorm.io/gorm"
)

type Seeder struct {
	db *gorm.DB
}

func NewSeeder(dbInstance *gorm.DB) *Seeder {
	return &Seeder{
		db: dbInstance,
	}
}

func (seeder *Seeder) SeedPermissions() error {
	defaultPermissions := []string{
		"auth:superuser",
		"users:create",
		"users:read",
		"users:update",
		"users:delete",
		"users:manage",
		"roles:create",
		"roles:read",
		"roles:update",
		"roles:delete",
		"roles:manage",
		"workspaces:create",
		"workspaces:read",
		"workspaces:update",
		"workspaces:delete",
		"workspaces:manage_all",
		"boards:create",
		"boards:read",
		"boards:update",
		"boards:delete",
		"boards:manage_all",
		"system:settings",
		"system:maintenance",
		"system:logs",
	}
	currentPermissions := make([]*models.Permission, 0)
	if err := seeder.db.Where("path IN ?", defaultPermissions).Find(&currentPermissions).Error; err != nil {
		return err
	}

	if len(currentPermissions) == len(defaultPermissions) {
		return nil
	}

	currentPermissionMap := make(map[string]bool)
	permissionsToCreate := make([]*models.Permission, 0)

	for _, currentPermission := range currentPermissions {
		currentPermissionMap[currentPermission.Path] = true
	}

	for _, defaultPermission := range defaultPermissions {
		if !currentPermissionMap[defaultPermission] {
			permissionsToCreate = append(permissionsToCreate, &models.Permission{Path: defaultPermission})
		}
	}

	if err := seeder.db.Create(&permissionsToCreate).Error; err != nil {
		return err
	}

	return nil
}

func (seeder *Seeder) SeedRoles() error {
	defaultRoles := []struct {
		Name        string
		Permissions []string
	}{
		{
			Name: "superuser",
			Permissions: []string{
				"auth:superuser",
			},
		},
		{
			Name: "admin",
			Permissions: []string{
				"users:manage",
				"roles:manage",
				"workspaces:manage_all",
				"boards:manage_all",
				"system:settings",
			},
		},
		{
			Name: "moderator",
			Permissions: []string{
				"users:read",
				"users:update",
				"workspaces:read",
				"workspaces:update",
				"boards:read",
				"boards:update",
				"boards:delete",
			},
		},
		{
			Name: "user",
			Permissions: []string{
				"workspaces:create",
				"workspaces:read",
				"workspaces:update",
				"boards:create",
				"boards:read",
				"boards:update",
			},
		},
		{
			Name: "viewer",
			Permissions: []string{
				"workspaces:read",
				"boards:read",
			},
		},
	}
	currentRoles := make([]*models.Role, 0)
	roleNames := make([]string, len(defaultRoles))

	for index := range roleNames {
		roleNames[index] = defaultRoles[index].Name
	}

	if err := seeder.db.Omit("Permissions").Where("name in ?", roleNames).Find(&currentRoles).Error; err != nil {
		return err
	}

	if len(currentRoles) == len(defaultRoles) {
		return nil
	}

	roleMap := make(map[string]bool)
	rolesToCreate := make([]*models.Role, 0)

	for _, currentRole := range currentRoles {
		roleMap[currentRole.Name] = true
	}

	for _, defaultRole := range defaultRoles {
		if roleMap[defaultRole.Name] {
			continue
		}

		role := &models.Role{Name: defaultRole.Name}
		permissions := make([]models.Permission, 0)
		if err := seeder.db.Where("path IN ?", defaultRole.Permissions).Find(&permissions).Error; err != nil {
			return err
		}

		role.Permissions = permissions
		rolesToCreate = append(rolesToCreate, role)
	}

	if len(rolesToCreate) > 0 {
		if err := seeder.db.Create(&rolesToCreate).Error; err != nil {
			return err
		}
	}

	return nil
}

func (seeder *Seeder) SeedInitialUser() error {
	var userCount int64
	if err := seeder.db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return err
	}

	if userCount > 0 {
		return nil
	}

	initialUsername := os.Getenv("INITIAL_USER_USERNAME")
	initialEmail := os.Getenv("INITIAL_USER_EMAIL")
	initialPassword := os.Getenv("INITIAL_USER_PASSWORD")

	if len(initialUsername) == 0 {
		initialUsername = "admin"
	}

	if len(initialEmail) == 0 {
		initialEmail = "admin@openboard.local"
	}

	if len(initialPassword) == 0 {
		initialPassword = auth.GenerateRandomPassword()
	}

	role := new(models.Role)
	if err := seeder.db.Where("name = ?", "superuser").First(role).Error; err != nil {
		return err
	}

	user := &models.User{
		Username: initialUsername,
		Email:    initialEmail,
	}
	if err := user.HashPassword(initialPassword); err != nil {
		return err
	}

	if err := seeder.db.Create(user).Error; err != nil {
		return err
	}

	if err := seeder.db.Model(user).Association("Roles").Append(role); err != nil {
		return err
	}

	if err := auth.SaveUserCredentials(user.Username, user.Email, initialPassword, "initial_credentials"); err != nil {
		return err
	}

	return nil
}
