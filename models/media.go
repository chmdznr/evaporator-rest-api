package models

import "time"

// Media represents the structure of the media table
type Media struct {
	ID                   uint      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ModelType            string    `json:"model_type"`
	ModelID              uint      `json:"model_id"`
	UUID                 string    `json:"uuid"`
	CollectionName       string    `json:"collection_name"`
	Name                 string    `json:"name"`
	FileName             string    `json:"file_name"`
	MimeType             string    `json:"mime_type"`
	Disk                 string    `json:"disk"`
	Size                 uint      `json:"size"`
	Manipulations        string    `gorm:"type:json" json:"manipulations"`
	CustomProperties     string    `gorm:"type:json" json:"custom_properties"`
	GeneratedConversions string    `gorm:"type:json" json:"generated_conversions"`
	ResponsiveImages     string    `gorm:"type:json" json:"responsive_images"`
	OrderColumn          uint      `json:"order_column"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// TableName returns the name of the media table
func (*Media) TableName() string {
	return "media"
}
