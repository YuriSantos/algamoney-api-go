package models

type Pessoa struct {
	Codigo       int64     `gorm:"primaryKey;column:codigo" json:"codigo"`
	Nome         string    `gorm:"column:nome;size:100" json:"nome"`
	Ativo        bool      `gorm:"column:ativo;default:true" json:"ativo"`
	Logradouro   *string   `gorm:"column:logradouro;size:100" json:"logradouro,omitempty"`
	Numero       *string   `gorm:"column:numero;size:20" json:"numero,omitempty"`
	Complemento  *string   `gorm:"column:complemento;size:100" json:"complemento,omitempty"`
	Bairro       *string   `gorm:"column:bairro;size:100" json:"bairro,omitempty"`
	Cep          *string   `gorm:"column:cep;size:10" json:"cep,omitempty"`
	CodigoCidade *int64    `gorm:"column:codigo_cidade" json:"codigoCidade,omitempty"`
	Cidade       *Cidade   `gorm:"foreignKey:CodigoCidade;references:Codigo" json:"cidade,omitempty"`
	Contatos     []Contato `gorm:"foreignKey:CodigoPessoa;references:Codigo" json:"contatos,omitempty"`
}

func (Pessoa) TableName() string {
	return "pessoa"
}
