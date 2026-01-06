package repository

import (
	"time"

	"github.com/algamoney/api/internal/dto"
	"github.com/algamoney/api/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Categoria
func (r *Repository) FindAllCategorias() ([]models.Categoria, error) {
	var categorias []models.Categoria
	err := r.db.Find(&categorias).Error
	return categorias, err
}

func (r *Repository) FindCategoriaByID(id int64) (*models.Categoria, error) {
	var categoria models.Categoria
	err := r.db.First(&categoria, id).Error
	if err != nil {
		return nil, err
	}
	return &categoria, nil
}

func (r *Repository) CreateCategoria(categoria *models.Categoria) error {
	return r.db.Create(categoria).Error
}

// Estado
func (r *Repository) FindAllEstados() ([]models.Estado, error) {
	var estados []models.Estado
	err := r.db.Order("nome").Find(&estados).Error
	return estados, err
}

// Cidade
func (r *Repository) FindCidadesByEstado(codigoEstado *int64) ([]models.Cidade, error) {
	var cidades []models.Cidade
	query := r.db.Preload("Estado").Order("nome")
	if codigoEstado != nil {
		query = query.Where("codigo_estado = ?", *codigoEstado)
	}
	err := query.Find(&cidades).Error
	return cidades, err
}

// Pessoa
func (r *Repository) FindAllPessoas(nome string, page, size int) ([]models.Pessoa, int64, error) {
	var pessoas []models.Pessoa
	var total int64

	query := r.db.Model(&models.Pessoa{}).Preload("Contatos").Preload("Cidade.Estado")
	if nome != "" {
		query = query.Where("nome LIKE ?", "%"+nome+"%")
	}

	query.Count(&total)
	err := query.Order("nome").Offset(page * size).Limit(size).Find(&pessoas).Error
	return pessoas, total, err
}

func (r *Repository) FindPessoaByID(id int64) (*models.Pessoa, error) {
	var pessoa models.Pessoa
	err := r.db.Preload("Contatos").Preload("Cidade.Estado").First(&pessoa, id).Error
	if err != nil {
		return nil, err
	}
	return &pessoa, nil
}

func (r *Repository) CreatePessoa(pessoa *models.Pessoa) error {
	return r.db.Create(pessoa).Error
}

func (r *Repository) UpdatePessoa(pessoa *models.Pessoa) error {
	return r.db.Save(pessoa).Error
}

func (r *Repository) DeletePessoa(id int64) error {
	return r.db.Delete(&models.Pessoa{}, id).Error
}

func (r *Repository) DeleteContatosByPessoa(pessoaID int64) error {
	return r.db.Where("codigo_pessoa = ?", pessoaID).Delete(&models.Contato{}).Error
}

func (r *Repository) CreateContatos(contatos []models.Contato) error {
	if len(contatos) == 0 {
		return nil
	}
	return r.db.Create(&contatos).Error
}

// Lancamento
func (r *Repository) FindAllLancamentos(filter dto.LancamentoFilter, page, size int) ([]models.Lancamento, int64, error) {
	var lancamentos []models.Lancamento
	var total int64

	query := r.db.Model(&models.Lancamento{}).Preload("Categoria").Preload("Pessoa")

	if filter.Descricao != "" {
		query = query.Where("descricao LIKE ?", "%"+filter.Descricao+"%")
	}
	if filter.DataVencimentoDe != "" {
		query = query.Where("data_vencimento >= ?", filter.DataVencimentoDe)
	}
	if filter.DataVencimentoAte != "" {
		query = query.Where("data_vencimento <= ?", filter.DataVencimentoAte)
	}

	query.Count(&total)
	err := query.Order("data_vencimento").Offset(page * size).Limit(size).Find(&lancamentos).Error
	return lancamentos, total, err
}

func (r *Repository) FindLancamentoByID(id int64) (*models.Lancamento, error) {
	var lancamento models.Lancamento
	err := r.db.Preload("Categoria").Preload("Pessoa").First(&lancamento, id).Error
	if err != nil {
		return nil, err
	}
	return &lancamento, nil
}

func (r *Repository) CreateLancamento(lancamento *models.Lancamento) error {
	return r.db.Create(lancamento).Error
}

func (r *Repository) UpdateLancamento(lancamento *models.Lancamento) error {
	return r.db.Save(lancamento).Error
}

func (r *Repository) DeleteLancamento(id int64) error {
	return r.db.Delete(&models.Lancamento{}, id).Error
}

func (r *Repository) EstatisticasPorCategoria(mes time.Time) ([]dto.EstatisticaCategoria, error) {
	primeiroDia := time.Date(mes.Year(), mes.Month(), 1, 0, 0, 0, 0, mes.Location())
	ultimoDia := primeiroDia.AddDate(0, 1, -1)

	var result []dto.EstatisticaCategoria
	err := r.db.Model(&models.Lancamento{}).
		Select("categoria.nome as categoria, SUM(lancamento.valor) as total").
		Joins("JOIN categoria ON categoria.codigo = lancamento.codigo_categoria").
		Where("lancamento.data_vencimento BETWEEN ? AND ?", primeiroDia, ultimoDia).
		Group("categoria.codigo").
		Scan(&result).Error

	return result, err
}

func (r *Repository) EstatisticasPorDia(mes time.Time) ([]dto.EstatisticaDia, error) {
	primeiroDia := time.Date(mes.Year(), mes.Month(), 1, 0, 0, 0, 0, mes.Location())
	ultimoDia := primeiroDia.AddDate(0, 1, -1)

	var result []dto.EstatisticaDia
	err := r.db.Model(&models.Lancamento{}).
		Select("tipo, data_vencimento as dia, SUM(valor) as total").
		Where("data_vencimento BETWEEN ? AND ?", primeiroDia, ultimoDia).
		Group("tipo, data_vencimento").
		Order("data_vencimento").
		Scan(&result).Error

	return result, err
}

// Usuario
func (r *Repository) FindUsuarioByEmail(email string) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.Preload("Permissoes").Where("email = ?", email).First(&usuario).Error
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}
