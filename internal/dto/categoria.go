package dto

type CreateCategoriaRequest struct {
	Nome string `json:"nome" binding:"required,max=50"`
}
