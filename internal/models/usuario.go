package models

type Usuario struct {
	Codigo     int64       `gorm:"primaryKey;column:codigo" json:"codigo"`
	Nome       string      `gorm:"column:nome;size:100" json:"nome"`
	Email      string      `gorm:"column:email;size:100;uniqueIndex" json:"email"`
	Senha      string      `gorm:"column:senha;size:150" json:"-"`
	Permissoes []Permissao `gorm:"many2many:usuario_permissao;foreignKey:Codigo;joinForeignKey:codigo_usuario;References:Codigo;joinReferences:codigo_permissao" json:"permissoes,omitempty"`
}

func (Usuario) TableName() string {
	return "usuario"
}

func (u *Usuario) GetPermissoes() []string {
	var perms []string
	for _, p := range u.Permissoes {
		perms = append(perms, p.Descricao)
	}
	return perms
}

func (u *Usuario) HasPermissao(permissao string) bool {
	for _, p := range u.Permissoes {
		if p.Descricao == permissao {
			return true
		}
	}
	return false
}

type Permissao struct {
	Codigo    int64  `gorm:"primaryKey;column:codigo" json:"codigo"`
	Descricao string `gorm:"column:descricao;size:50;uniqueIndex" json:"descricao"`
}

func (Permissao) TableName() string {
	return "permissao"
}
