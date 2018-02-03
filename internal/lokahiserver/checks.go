package lokahiserver

import "github.com/jinzhu/gorm"

// Checks implements service Checks.
type Checks struct {
	db *gorm.DB
}
