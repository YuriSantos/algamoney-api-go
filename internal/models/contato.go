package models

type Contato struct {
	Codigo       int64   `gorm:"primaryKey;column:codigo" json:"codigo"`
	Nome         string  `gorm:"column:nome;size:100" json:"nome"`
	Email        *string `gorm:"column:email;size:100" json:"email,omitempty"`
	Telefone     *string `gorm:"column:telefone;size:20" json:"telefone,omitempty"`
	CodigoPessoa int64   `gorm:"column:codigo_pessoa" json:"codigoPessoa"`
}

func (Contato) TableName() string {
	return "contato"
}
