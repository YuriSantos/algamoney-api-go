package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/algamoney/api/internal/dto"
	"github.com/algamoney/api/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// Auth
func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.svc.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, token)
}

// Categoria
func (h *Handler) GetCategorias(c *gin.Context) {
	categorias, err := h.svc.FindAllCategorias()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categorias)
}

func (h *Handler) GetCategoriaByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("codigo"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "código inválido"})
		return
	}

	categoria, err := h.svc.FindCategoriaByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "categoria não encontrada"})
		return
	}
	c.JSON(http.StatusOK, categoria)
}

func (h *Handler) CreateCategoria(c *gin.Context) {
	var req dto.CreateCategoriaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	categoria, err := h.svc.CreateCategoria(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, categoria)
}

// Estado
func (h *Handler) GetEstados(c *gin.Context) {
	estados, err := h.svc.FindAllEstados()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, estados)
}

// Cidade
func (h *Handler) GetCidades(c *gin.Context) {
	var codigoEstado *int64
	if estadoParam := c.Query("estado"); estadoParam != "" {
		id, err := strconv.ParseInt(estadoParam, 10, 64)
		if err == nil {
			codigoEstado = &id
		}
	}

	cidades, err := h.svc.FindCidadesByEstado(codigoEstado)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cidades)
}

// Pessoa
func (h *Handler) GetPessoas(c *gin.Context) {
	nome := c.Query("nome")
	pagination := getPaginationParams(c)

	result, err := h.svc.FindAllPessoas(nome, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetPessoaByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("codigo"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "código inválido"})
		return
	}

	pessoa, err := h.svc.FindPessoaByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pessoa não encontrada"})
		return
	}
	c.JSON(http.StatusOK, pessoa)
}

func (h *Handler) CreatePessoa(c *gin.Context) {
	var req dto.CreatePessoaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pessoa, err := h.svc.CreatePessoa(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pessoa)
}

func (h *Handler) UpdatePessoa(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("codigo"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "código inválido"})
		return
	}

	var req dto.UpdatePessoaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pessoa, err := h.svc.UpdatePessoa(id, req)
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "pessoa não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pessoa)
}

func (h *Handler) DeletePessoa(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("codigo"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "código inválido"})
		return
	}

	if err := h.svc.DeletePessoa(id); err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "pessoa não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) UpdatePessoaAtivo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("codigo"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "código inválido"})
		return
	}

	var body struct {
		Ativo bool `json:"ativo"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdatePessoaAtivo(id, body.Ativo); err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "pessoa não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// Lancamento
func (h *Handler) GetLancamentos(c *gin.Context) {
	filter := getLancamentoFilter(c)
	pagination := getPaginationParams(c)

	result, err := h.svc.FindAllLancamentos(filter, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetLancamentosResumo(c *gin.Context) {
	filter := getLancamentoFilter(c)
	pagination := getPaginationParams(c)

	result, err := h.svc.FindLancamentosResumo(filter, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetLancamentoByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("codigo"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "código inválido"})
		return
	}

	lancamento, err := h.svc.FindLancamentoByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lançamento não encontrado"})
		return
	}
	c.JSON(http.StatusOK, lancamento)
}

func (h *Handler) CreateLancamento(c *gin.Context) {
	var req dto.CreateLancamentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lancamento, err := h.svc.CreateLancamento(req)
	if err != nil {
		if err == service.ErrPessoaInexistenteInativa {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, lancamento)
}

func (h *Handler) UpdateLancamento(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("codigo"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "código inválido"})
		return
	}

	var req dto.UpdateLancamentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lancamento, err := h.svc.UpdateLancamento(id, req)
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "lançamento não encontrado"})
			return
		}
		if err == service.ErrPessoaInexistenteInativa {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lancamento)
}

func (h *Handler) DeleteLancamento(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("codigo"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "código inválido"})
		return
	}

	if err := h.svc.DeleteLancamento(id); err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "lançamento não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) EstatisticasPorCategoria(c *gin.Context) {
	mes := getMesReferencia(c)
	result, err := h.svc.EstatisticasPorCategoria(mes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) EstatisticasPorDia(c *gin.Context) {
	mes := getMesReferencia(c)
	result, err := h.svc.EstatisticasPorDia(mes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Helper functions
func getPaginationParams(c *gin.Context) dto.PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	return dto.PaginationParams{Page: page, Size: size}
}

func getLancamentoFilter(c *gin.Context) dto.LancamentoFilter {
	return dto.LancamentoFilter{
		Descricao:         c.Query("descricao"),
		DataVencimentoDe:  c.Query("dataVencimentoDe"),
		DataVencimentoAte: c.Query("dataVencimentoAte"),
	}
}

func getMesReferencia(c *gin.Context) time.Time {
	mesParam := c.Query("mes")
	if mesParam != "" {
		if mes, err := time.Parse("2006-01", mesParam); err == nil {
			return mes
		}
	}
	return time.Now()
}
