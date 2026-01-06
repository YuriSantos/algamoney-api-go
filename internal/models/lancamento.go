package models

import "time"

type TipoLancamento string

const (
	TipoReceita TipoLancamento = "RECEITA"
	TipoDespesa TipoLancamento = "DESPESA"
)

type Lancamento struct {
	Codigo          int64          `gorm:"primaryKey;column:codigo" json:"codigo"`
	Descricao       string         `gorm:"column:descricao;size:100" json:"descricao"`
	DataVencimento  time.Time      `gorm:"column:data_vencimento;type:date" json:"dataVencimento"`
	DataPagamento   *time.Time     `gorm:"column:data_pagamento;type:date" json:"dataPagamento,omitempty"`
	Valor           float64        `gorm:"column:valor;type:decimal(10,2)" json:"valor"`
	Observacao      *string        `gorm:"column:observacao;type:text" json:"observacao,omitempty"`
	Tipo            TipoLancamento `gorm:"column:tipo;type:enum('RECEITA','DESPESA')" json:"tipo"`
	Anexo           *string        `gorm:"column:anexo;size:255" json:"anexo,omitempty"`
	CodigoCategoria int64          `gorm:"column:codigo_categoria" json:"codigoCategoria"`
	CodigoPessoa    int64          `gorm:"column:codigo_pessoa" json:"codigoPessoa"`
	Categoria       *Categoria     `gorm:"foreignKey:CodigoCategoria;references:Codigo" json:"categoria,omitempty"`
	Pessoa          *Pessoa        `gorm:"foreignKey:CodigoPessoa;references:Codigo" json:"pessoa,omitempty"`
}

func (Lancamento) TableName() string {
	return "lancamento"
}
