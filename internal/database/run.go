package database

import "github.com/jinzhu/gorm"

type Run struct {
	gorm.Model

	Checks   []Check
	Finished bool
	Message  string
}
