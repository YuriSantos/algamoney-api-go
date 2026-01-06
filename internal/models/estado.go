package models

type Estado struct {
	Codigo int64  `gorm:"primaryKey;column:codigo" json:"codigo"`
	Nome   string `gorm:"column:nome;size:50" json:"nome"`
}

func (Estado) TableName() string {
	return "estado"
}
