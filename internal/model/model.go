// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime(3);comment:Creation time" json:"created_at"` // Timestamp when the record was created
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime(3);comment:Update time" json:"updated_at"`   // Timestamp when the record was last updated
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;comment:Deletion time" json:"-"`                     // Timestamp when the record was deleted (nullable, used for soft deletes)
}
