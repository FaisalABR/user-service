package seeders

import "gorm.io/gorm"

type Registery struct {
	db *gorm.DB
}

type ISeederRegistry interface {
	Run()
}

func (r *Registery) Run() {
	RunRoleSeeder(r.db)
	RunUserSeeder(r.db)
}

func NewSeederRegistry(db *gorm.DB) ISeederRegistry {
	return &Registery{
		db: db,
	}
}
