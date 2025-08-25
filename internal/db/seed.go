package db

import "gorm.io/gorm"

type Seeder struct {
	db *gorm.DB
}

func NewSeeder(dbInstance *gorm.DB) *Seeder {
	return &Seeder{
		db: dbInstance,
	}
}

func (seeder *Seeder) SeedAll() error {
	return nil
}

func (seeder *Seeder) SeedPermissions() error {
	return nil
}

func (seeder *Seeder) SeedRoles() error {
	return nil
}

func (seeder *Seeder) SeedSuperUser() error {
	return nil
}
