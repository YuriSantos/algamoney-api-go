package main

import (
	"log"

	"github.com/algamoney/api/internal/config"
	"github.com/algamoney/api/internal/handler"
	"github.com/algamoney/api/internal/middleware"
	"github.com/algamoney/api/internal/repository"
	"github.com/algamoney/api/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set gin mode
	gin.SetMode(cfg.GinMode)

	// Database connection
	db := config.NewDatabase(cfg)

	// Initialize layers
	repo := repository.NewRepository(db)
	svc := service.NewService(repo, cfg)
	h := handler.NewHandler(svc)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	// Create router
	r := gin.Default()

	// CORS
	r.Use(middleware.CORSMiddleware(cfg.CORSOrigin))

	// Public routes
	r.POST("/oauth2/token", h.Login)

	// Protected routes
	api := r.Group("")
	api.Use(authMiddleware.Authenticate())
	{
		// Categorias
		categorias := api.Group("/categorias")
		{
			categorias.GET("", authMiddleware.RequirePermission("ROLE_PESQUISAR_CATEGORIA"), h.GetCategorias)
			categorias.GET("/:codigo", authMiddleware.RequirePermission("ROLE_PESQUISAR_CATEGORIA"), h.GetCategoriaByID)
			categorias.POST("", authMiddleware.RequirePermission("ROLE_CADASTRAR_CATEGORIA"), h.CreateCategoria)
		}

		// Estados
		estados := api.Group("/estados")
		{
			estados.GET("", h.GetEstados)
		}

		// Cidades
		cidades := api.Group("/cidades")
		{
			cidades.GET("", h.GetCidades)
		}

		// Pessoas
		pessoas := api.Group("/pessoas")
		{
			pessoas.GET("", authMiddleware.RequirePermission("ROLE_PESQUISAR_PESSOA"), h.GetPessoas)
			pessoas.GET("/:codigo", authMiddleware.RequirePermission("ROLE_PESQUISAR_PESSOA"), h.GetPessoaByID)
			pessoas.POST("", authMiddleware.RequirePermission("ROLE_CADASTRAR_PESSOA"), h.CreatePessoa)
			pessoas.PUT("/:codigo", authMiddleware.RequireAnyPermission("ROLE_CADASTRAR_PESSOA"), h.UpdatePessoa)
			pessoas.DELETE("/:codigo", authMiddleware.RequirePermission("ROLE_REMOVER_PESSOA"), h.DeletePessoa)
			pessoas.PUT("/:codigo/ativo", authMiddleware.RequirePermission("ROLE_CADASTRAR_PESSOA"), h.UpdatePessoaAtivo)
		}

		// Lancamentos
		lancamentos := api.Group("/lancamentos")
		{
			lancamentos.GET("", authMiddleware.RequirePermission("ROLE_PESQUISAR_LANCAMENTO"), h.GetLancamentos)
			lancamentos.GET("/resumo", authMiddleware.RequirePermission("ROLE_PESQUISAR_LANCAMENTO"), h.GetLancamentosResumo)
			lancamentos.GET("/:codigo", authMiddleware.RequirePermission("ROLE_PESQUISAR_LANCAMENTO"), h.GetLancamentoByID)
			lancamentos.POST("", authMiddleware.RequirePermission("ROLE_CADASTRAR_LANCAMENTO"), h.CreateLancamento)
			lancamentos.PUT("/:codigo", authMiddleware.RequireAnyPermission("ROLE_CADASTRAR_LANCAMENTO"), h.UpdateLancamento)
			lancamentos.DELETE("/:codigo", authMiddleware.RequirePermission("ROLE_REMOVER_LANCAMENTO"), h.DeleteLancamento)
			lancamentos.GET("/estatisticas/por-categoria", authMiddleware.RequirePermission("ROLE_PESQUISAR_LANCAMENTO"), h.EstatisticasPorCategoria)
			lancamentos.GET("/estatisticas/por-dia", authMiddleware.RequirePermission("ROLE_PESQUISAR_LANCAMENTO"), h.EstatisticasPorDia)
		}
	}

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
