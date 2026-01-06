package dto

type CreateContatoRequest struct {
	Codigo   *int64  `json:"codigo,omitempty"`
	Nome     string  `json:"nome" binding:"required,max=100"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email,max=100"`
	Telefone *string `json:"telefone,omitempty" binding:"omitempty,max=20"`
}

type CreatePessoaRequest struct {
	Nome         string                 `json:"nome" binding:"required,max=100"`
	Ativo        *bool                  `json:"ativo,omitempty"`
	Logradouro   *string                `json:"logradouro,omitempty" binding:"omitempty,max=100"`
	Numero       *string                `json:"numero,omitempty" binding:"omitempty,max=20"`
	Complemento  *string                `json:"complemento,omitempty" binding:"omitempty,max=100"`
	Bairro       *string                `json:"bairro,omitempty" binding:"omitempty,max=100"`
	Cep          *string                `json:"cep,omitempty" binding:"omitempty,max=10"`
	CodigoCidade *int64                 `json:"codigoCidade,omitempty"`
	Contatos     []CreateContatoRequest `json:"contatos,omitempty"`
}

type UpdatePessoaRequest = CreatePessoaRequest

type UpdateAtivoRequest struct {
	Ativo bool `json:"ativo" binding:"required"`
}
