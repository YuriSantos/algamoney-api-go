package dto

type PaginatedResponse struct {
	Content       interface{} `json:"content"`
	TotalElements int64       `json:"totalElements"`
	TotalPages    int         `json:"totalPages"`
	Size          int         `json:"size"`
	Number        int         `json:"number"`
}

type PaginationParams struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

func (p *PaginationParams) GetOffset() int {
	return p.Page * p.GetSize()
}

func (p *PaginationParams) GetSize() int {
	if p.Size <= 0 {
		return 20
	}
	return p.Size
}
