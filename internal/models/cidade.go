package models

type Cidade struct {
	Codigo       int64   `gorm:"primaryKey;column:codigo" json:"codigo"`
	Nome         string  `gorm:"column:nome;size:50" json:"nome"`
	CodigoEstado int64   `gorm:"column:codigo_estado" json:"codigoEstado"`
	Estado       *Estado `gorm:"foreignKey:CodigoEstado;references:Codigo" json:"estado,omitempty"`
}

func (Cidade) TableName() string {
	return "cidade"
}
