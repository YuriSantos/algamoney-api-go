package service

import (
	"errors"
	"time"

	"github.com/algamoney/api/internal/config"
	"github.com/algamoney/api/internal/dto"
	"github.com/algamoney/api/internal/models"
	"github.com/algamoney/api/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials       = errors.New("credenciais inválidas")
	ErrPessoaInexistenteInativa = errors.New("pessoa inexistente ou inativa")
	ErrNotFound                 = errors.New("registro não encontrado")
)

type Service struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

// Auth
func (s *Service) Login(req dto.LoginRequest) (*dto.TokenResponse, error) {
	usuario, err := s.repo.FindUsuarioByEmail(req.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Senha), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(usuario)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken: token,
		TokenType:   "bearer",
		ExpiresIn:   s.cfg.JWTExpiresHours * 3600,
	}, nil
}

func (s *Service) generateToken(usuario *models.Usuario) (string, error) {
	claims := jwt.MapClaims{
		"sub":         usuario.Email,
		"nome":        usuario.Nome,
		"authorities": usuario.GetPermissoes(),
		"exp":         time.Now().Add(time.Duration(s.cfg.JWTExpiresHours) * time.Hour).Unix(),
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

// Categoria
func (s *Service) FindAllCategorias() ([]models.Categoria, error) {
	return s.repo.FindAllCategorias()
}

func (s *Service) FindCategoriaByID(id int64) (*models.Categoria, error) {
	return s.repo.FindCategoriaByID(id)
}

func (s *Service) CreateCategoria(req dto.CreateCategoriaRequest) (*models.Categoria, error) {
	categoria := &models.Categoria{Nome: req.Nome}
	if err := s.repo.CreateCategoria(categoria); err != nil {
		return nil, err
	}
	return categoria, nil
}

// Estado
func (s *Service) FindAllEstados() ([]models.Estado, error) {
	return s.repo.FindAllEstados()
}

// Cidade
func (s *Service) FindCidadesByEstado(codigoEstado *int64) ([]models.Cidade, error) {
	return s.repo.FindCidadesByEstado(codigoEstado)
}

// Pessoa
func (s *Service) FindAllPessoas(nome string, pagination dto.PaginationParams) (*dto.PaginatedResponse, error) {
	pessoas, total, err := s.repo.FindAllPessoas(nome, pagination.Page, pagination.GetSize())
	if err != nil {
		return nil, err
	}

	size := pagination.GetSize()
	totalPages := int(total) / size
	if int(total)%size > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse{
		Content:       pessoas,
		TotalElements: total,
		TotalPages:    totalPages,
		Size:          size,
		Number:        pagination.Page,
	}, nil
}

func (s *Service) FindPessoaByID(id int64) (*models.Pessoa, error) {
	return s.repo.FindPessoaByID(id)
}

func (s *Service) CreatePessoa(req dto.CreatePessoaRequest) (*models.Pessoa, error) {
	ativo := true
	if req.Ativo != nil {
		ativo = *req.Ativo
	}

	pessoa := &models.Pessoa{
		Nome:         req.Nome,
		Ativo:        ativo,
		Logradouro:   req.Logradouro,
		Numero:       req.Numero,
		Complemento:  req.Complemento,
		Bairro:       req.Bairro,
		Cep:          req.Cep,
		CodigoCidade: req.CodigoCidade,
	}

	if err := s.repo.CreatePessoa(pessoa); err != nil {
		return nil, err
	}

	if len(req.Contatos) > 0 {
		var contatos []models.Contato
		for _, c := range req.Contatos {
			contatos = append(contatos, models.Contato{
				Nome:         c.Nome,
				Email:        c.Email,
				Telefone:     c.Telefone,
				CodigoPessoa: pessoa.Codigo,
			})
		}
		if err := s.repo.CreateContatos(contatos); err != nil {
			return nil, err
		}
	}

	return s.repo.FindPessoaByID(pessoa.Codigo)
}

func (s *Service) UpdatePessoa(id int64, req dto.UpdatePessoaRequest) (*models.Pessoa, error) {
	pessoa, err := s.repo.FindPessoaByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	pessoa.Nome = req.Nome
	if req.Ativo != nil {
		pessoa.Ativo = *req.Ativo
	}
	pessoa.Logradouro = req.Logradouro
	pessoa.Numero = req.Numero
	pessoa.Complemento = req.Complemento
	pessoa.Bairro = req.Bairro
	pessoa.Cep = req.Cep
	pessoa.CodigoCidade = req.CodigoCidade

	if err := s.repo.UpdatePessoa(pessoa); err != nil {
		return nil, err
	}

	// Update contatos
	s.repo.DeleteContatosByPessoa(id)
	if len(req.Contatos) > 0 {
		var contatos []models.Contato
		for _, c := range req.Contatos {
			contatos = append(contatos, models.Contato{
				Nome:         c.Nome,
				Email:        c.Email,
				Telefone:     c.Telefone,
				CodigoPessoa: id,
			})
		}
		s.repo.CreateContatos(contatos)
	}

	return s.repo.FindPessoaByID(id)
}

func (s *Service) DeletePessoa(id int64) error {
	if _, err := s.repo.FindPessoaByID(id); err != nil {
		return ErrNotFound
	}
	return s.repo.DeletePessoa(id)
}

func (s *Service) UpdatePessoaAtivo(id int64, ativo bool) error {
	pessoa, err := s.repo.FindPessoaByID(id)
	if err != nil {
		return ErrNotFound
	}
	pessoa.Ativo = ativo
	return s.repo.UpdatePessoa(pessoa)
}

// Lancamento
func (s *Service) FindAllLancamentos(filter dto.LancamentoFilter, pagination dto.PaginationParams) (*dto.PaginatedResponse, error) {
	lancamentos, total, err := s.repo.FindAllLancamentos(filter, pagination.Page, pagination.GetSize())
	if err != nil {
		return nil, err
	}

	size := pagination.GetSize()
	totalPages := int(total) / size
	if int(total)%size > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse{
		Content:       lancamentos,
		TotalElements: total,
		TotalPages:    totalPages,
		Size:          size,
		Number:        pagination.Page,
	}, nil
}

func (s *Service) FindLancamentosResumo(filter dto.LancamentoFilter, pagination dto.PaginationParams) (*dto.PaginatedResponse, error) {
	lancamentos, total, err := s.repo.FindAllLancamentos(filter, pagination.Page, pagination.GetSize())
	if err != nil {
		return nil, err
	}

	var resumos []dto.ResumoLancamento
	for _, l := range lancamentos {
		resumo := dto.ResumoLancamento{
			Codigo:         l.Codigo,
			Descricao:      l.Descricao,
			DataVencimento: l.DataVencimento.Format("2006-01-02"),
			Valor:          l.Valor,
			Tipo:           string(l.Tipo),
		}
		if l.DataPagamento != nil {
			dp := l.DataPagamento.Format("2006-01-02")
			resumo.DataPagamento = &dp
		}
		if l.Categoria != nil {
			resumo.Categoria = l.Categoria.Nome
		}
		if l.Pessoa != nil {
			resumo.Pessoa = l.Pessoa.Nome
		}
		resumos = append(resumos, resumo)
	}

	size := pagination.GetSize()
	totalPages := int(total) / size
	if int(total)%size > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse{
		Content:       resumos,
		TotalElements: total,
		TotalPages:    totalPages,
		Size:          size,
		Number:        pagination.Page,
	}, nil
}

func (s *Service) FindLancamentoByID(id int64) (*models.Lancamento, error) {
	return s.repo.FindLancamentoByID(id)
}

func (s *Service) CreateLancamento(req dto.CreateLancamentoRequest) (*models.Lancamento, error) {
	if err := s.validarPessoa(req.CodigoPessoa); err != nil {
		return nil, err
	}

	dataVencimento, _ := time.Parse("2006-01-02", req.DataVencimento)
	var dataPagamento *time.Time
	if req.DataPagamento != nil {
		dp, _ := time.Parse("2006-01-02", *req.DataPagamento)
		dataPagamento = &dp
	}

	lancamento := &models.Lancamento{
		Descricao:       req.Descricao,
		DataVencimento:  dataVencimento,
		DataPagamento:   dataPagamento,
		Valor:           req.Valor,
		Observacao:      req.Observacao,
		Tipo:            req.Tipo,
		CodigoCategoria: req.CodigoCategoria,
		CodigoPessoa:    req.CodigoPessoa,
	}

	if err := s.repo.CreateLancamento(lancamento); err != nil {
		return nil, err
	}

	return s.repo.FindLancamentoByID(lancamento.Codigo)
}

func (s *Service) UpdateLancamento(id int64, req dto.UpdateLancamentoRequest) (*models.Lancamento, error) {
	lancamento, err := s.repo.FindLancamentoByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	if err := s.validarPessoa(req.CodigoPessoa); err != nil {
		return nil, err
	}

	dataVencimento, _ := time.Parse("2006-01-02", req.DataVencimento)
	var dataPagamento *time.Time
	if req.DataPagamento != nil {
		dp, _ := time.Parse("2006-01-02", *req.DataPagamento)
		dataPagamento = &dp
	}

	lancamento.Descricao = req.Descricao
	lancamento.DataVencimento = dataVencimento
	lancamento.DataPagamento = dataPagamento
	lancamento.Valor = req.Valor
	lancamento.Observacao = req.Observacao
	lancamento.Tipo = req.Tipo
	lancamento.CodigoCategoria = req.CodigoCategoria
	lancamento.CodigoPessoa = req.CodigoPessoa

	if err := s.repo.UpdateLancamento(lancamento); err != nil {
		return nil, err
	}

	return s.repo.FindLancamentoByID(id)
}

func (s *Service) DeleteLancamento(id int64) error {
	if _, err := s.repo.FindLancamentoByID(id); err != nil {
		return ErrNotFound
	}
	return s.repo.DeleteLancamento(id)
}

func (s *Service) EstatisticasPorCategoria(mesReferencia time.Time) ([]dto.EstatisticaCategoria, error) {
	return s.repo.EstatisticasPorCategoria(mesReferencia)
}

func (s *Service) EstatisticasPorDia(mesReferencia time.Time) ([]dto.EstatisticaDia, error) {
	return s.repo.EstatisticasPorDia(mesReferencia)
}

func (s *Service) validarPessoa(codigoPessoa int64) error {
	pessoa, err := s.repo.FindPessoaByID(codigoPessoa)
	if err != nil || !pessoa.Ativo {
		return ErrPessoaInexistenteInativa
	}
	return nil
}
