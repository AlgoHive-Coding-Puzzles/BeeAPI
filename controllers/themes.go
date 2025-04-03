package controllers

import (
	"net/http"
	"time"

	"github.com/algohive/beeapi/models"
	"github.com/algohive/beeapi/services"
	"github.com/gin-gonic/gin"
)

// ThemeController handles theme-related endpoints
type ThemeController struct {
	loader *services.PuzzlesLoader
	lastReloadTime map[string]time.Time
	cooldownPeriod time.Duration
}

// NewThemeController creates a new theme controller
func NewThemeController(loader *services.PuzzlesLoader) *ThemeController {
	return &ThemeController{
		loader: loader,
		lastReloadTime: make(map[string]time.Time),
		cooldownPeriod: 10 * time.Second, // 10 seconds cooldown
	}
}

// GetThemes godoc
// @Summary Get all themes
// @Description Returns a list of all available themes
// @Tags Themes
// @Produce json
// @Success 200 {array} models.ThemeResponse
// @Router /themes [get]
func (t *ThemeController) GetThemes(c *gin.Context) {
	var themeResponses []models.ThemeResponse
	
	for _, theme := range t.loader.Themes {
		var puzzleResponses []models.PuzzleResponse
		
		for _, puzzle := range theme.Puzzles {
			compressedSize, uncompressedSize, _ := t.loader.GetPuzzleSizes(theme.Name, puzzle.GetName())
			
			puzzleResponse := models.PuzzleResponse{
				Name:            puzzle.GetName(),
				Title:           puzzle.DescProps.Title,
				Index:           puzzle.DescProps.Index,
				Difficulty:      puzzle.DescProps.Difficulty,
				Language:        puzzle.DescProps.Language,
				CompressedSize:  compressedSize,
				UncompressedSize: uncompressedSize,
				HivecraftVersion: puzzle.MetaProps.HivecraftVersion,
				Cipher:          puzzle.Cipher,
				Obscure:         puzzle.Obscure,
				ID:              puzzle.MetaProps.ID,
				Author:          puzzle.MetaProps.Author,
				CreatedAt:       puzzle.MetaProps.Created,
				UpdatedAt:       puzzle.MetaProps.Modified,
			}
			
			puzzleResponses = append(puzzleResponses, puzzleResponse)
		}
		
		themeSize, _ := services.GetDirSize(theme.Path)
		
		themeResponse := models.ThemeResponse{
			Name:         theme.Name,
			EnigmesCount: len(theme.Puzzles),
			Puzzles:      puzzleResponses,
			Size:         themeSize,
		}
		
		themeResponses = append(themeResponses, themeResponse)
	}
	
	c.JSON(http.StatusOK, themeResponses)
}

// GetThemeNames godoc
// @Summary Get theme names
// @Description Returns a list of theme names
// @Tags Themes
// @Produce json
// @Success 200 {array} string
// @Router /themes/names [get]
func (t *ThemeController) GetThemeNames(c *gin.Context) {
	var themeNames []string
	
	for _, theme := range t.loader.Themes {
		themeNames = append(themeNames, theme.Name)
	}
	
	c.JSON(http.StatusOK, themeNames)
}

// GetTheme godoc
// @Summary Get a specific theme
// @Description Returns details of a specific theme by name
// @Tags Themes
// @Produce json
// @Param name query string true "Theme name"
// @Success 200 {object} models.ThemeResponse
// @Failure 404 {object} map[string]string
// @Router /theme [get]
func (t *ThemeController) GetTheme(c *gin.Context) {
	name := c.Query("name")
	theme := t.loader.GetTheme(name)
	
	if theme == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Theme not found"})
		return
	}
	
	var puzzleResponses []models.PuzzleResponse
	
	for _, puzzle := range theme.Puzzles {
		compressedSize, uncompressedSize, _ := t.loader.GetPuzzleSizes(theme.Name, puzzle.GetName())
		
		puzzleResponse := models.PuzzleResponse{
			Name:            puzzle.GetName(),
			Title:           puzzle.DescProps.Title,
			Index:           puzzle.DescProps.Index,
			Difficulty:      puzzle.DescProps.Difficulty,
			Language:        puzzle.DescProps.Language,
			CompressedSize:  compressedSize,
			UncompressedSize: uncompressedSize,
			HivecraftVersion: puzzle.MetaProps.HivecraftVersion,
			Cipher:          puzzle.Cipher,
			Obscure:         puzzle.Obscure,
			ID:              puzzle.MetaProps.ID,
			Author:          puzzle.MetaProps.Author,
			CreatedAt:       puzzle.MetaProps.Created,
			UpdatedAt:       puzzle.MetaProps.Modified,
		}
		
		puzzleResponses = append(puzzleResponses, puzzleResponse)
	}
	
	themeSize, _ := services.GetDirSize(theme.Path)
	
	themeResponse := models.ThemeResponse{
		Name:         theme.Name,
		EnigmesCount: len(theme.Puzzles),
		Puzzles:      puzzleResponses,
		Size:         themeSize,
	}
	
	c.JSON(http.StatusOK, themeResponse)
}

// CreateTheme godoc
// @Summary Create a new theme
// @Description Creates a new theme with the given name
// @Tags Themes
// @Produce json
// @Param name query string true "Theme name"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /theme [post]
// @Security Bearer
func (t *ThemeController) CreateTheme(c *gin.Context) {
	name := c.Query("name")
	
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Theme name is required"})
		return
	}
	
	if t.loader.HasTheme(name) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Theme already exists"})
		return
	}
	
	err := t.loader.CreateTheme(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create theme"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Theme created"})
}

// DeleteTheme godoc
// @Summary Delete a theme
// @Description Deletes a theme with the given name
// @Tags Themes
// @Produce json
// @Param name query string true "Theme name"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /theme [delete]
// @Security Bearer
func (t *ThemeController) DeleteTheme(c *gin.Context) {
	name := c.Query("name")
	
	if !t.loader.HasTheme(name) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Theme not found"})
		return
	}
	
	err := t.loader.DeleteTheme(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete theme"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Theme deleted"})
}

// ReloadThemes godoc
// @Summary Reload themes
// @Description Reloads all themes and puzzles
// @Tags Themes
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Router /theme/reload [post]
// @Security Bearer
func (t *ThemeController) ReloadThemes(c *gin.Context) {
	userIP := c.ClientIP()
	currentTime := time.Now()
	
	lastReload, exists := t.lastReloadTime[userIP]
	if exists {
		elapsedTime := currentTime.Sub(lastReload)
		if elapsedTime < t.cooldownPeriod {
			remainingSeconds := int(t.cooldownPeriod.Seconds() - elapsedTime.Seconds())
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Cooldown period in effect",
				"wait": remainingSeconds,
			})
			return
		}
	}
	
	t.lastReloadTime[userIP] = currentTime
	
	err := t.loader.Reload()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload themes"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Themes reloaded"})
}
