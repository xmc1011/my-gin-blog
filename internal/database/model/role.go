package model

type Role struct {
	Model
	Name      string `gorm:"unique" json:"name"`
	Label     string `gorm:"unique" json:"label"`
	IsDisable bool   `json:"is_disable"`
}
