package dto

import "github.com/algamoney/api/internal/models"

type CreateLancamentoRequest struct {
	Descricao       string                `json:"descricao" binding:"required,max=100"`
	DataVencimento  string                `json:"dataVencimento" binding:"required"`
	DataPagamento   *string               `json:"dataPagamento,omitempty"`
	Valor           float64               `json:"valor" binding:"required,gt=0"`
	Observacao      *string               `json:"observacao,omitempty"`
	Tipo            models.TipoLancamento `json:"tipo" binding:"required,oneof=RECEITA DESPESA"`
	CodigoCategoria int64                 `json:"codigoCategoria" binding:"required"`
	CodigoPessoa    int64                 `json:"codigoPessoa" binding:"required"`
}

type UpdateLancamentoRequest = CreateLancamentoRequest

type LancamentoFilter struct {
	Descricao         string `form:"descricao"`
	DataVencimentoDe  string `form:"dataVencimentoDe"`
	DataVencimentoAte string `form:"dataVencimentoAte"`
}

type ResumoLancamento struct {
	Codigo         int64   `json:"codigo"`
	Descricao      string  `json:"descricao"`
	DataVencimento string  `json:"dataVencimento"`
	DataPagamento  *string `json:"dataPagamento,omitempty"`
	Valor          float64 `json:"valor"`
	Tipo           string  `json:"tipo"`
	Categoria      string  `json:"categoria"`
	Pessoa         string  `json:"pessoa"`
}

type EstatisticaCategoria struct {
	Categoria string  `json:"categoria"`
	Total     float64 `json:"total"`
}

type EstatisticaDia struct {
	Tipo  string  `json:"tipo"`
	Dia   string  `json:"dia"`
	Total float64 `json:"total"`
}
