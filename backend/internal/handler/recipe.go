package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yellow78/mini-mes/backend/internal/repository"
)

// RecipeHandler Recipe HTTP 處理器
type RecipeHandler struct {
	repo repository.RecipeRepository
}

func NewRecipeHandler(repo repository.RecipeRepository) *RecipeHandler {
	return &RecipeHandler{repo: repo}
}

// ListRecipes GET /api/v1/recipes
func (h *RecipeHandler) ListRecipes(c *gin.Context) {
	recipes, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": recipes, "total": len(recipes)})
}
