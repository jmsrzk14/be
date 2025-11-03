package models

type Status_Aspirations struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	Status int  `gorm:"type:int" json:"status"`
}
