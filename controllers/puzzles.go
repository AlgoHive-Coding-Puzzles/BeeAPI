package controllers

import (
	"net/http"
	"path/filepath"

	"github.com/algohive/beeapi/models"
	"github.com/algohive/beeapi/services"
	"github.com/gin-gonic/gin"
)

// PuzzleController handles puzzle-related endpoints
type PuzzleController struct {
	loader      *services.PuzzlesLoader
	pythonRunner *services.PythonRunner
}

// NewPuzzleController creates a new puzzle controller
func NewPuzzleController(loader *services.PuzzlesLoader, pythonRunner *services.PythonRunner) *PuzzleController {
	return &PuzzleController{
		loader:      loader,
		pythonRunner: pythonRunner,
	}
}

// GetPuzzles godoc
// @Summary Get puzzles for a theme
// @Description Returns all puzzles for a specific theme
// @Tags Puzzles
// @Produce json
// @Param theme query string true "Theme name"
// @Success 200 {array} models.PuzzleResponse
// @Failure 404 {object} map[string]string
// @Router /puzzles [get]
func (p *PuzzleController) GetPuzzles(c *gin.Context) {
	themeName := c.Query("theme")
	
	theme := p.loader.GetTheme(themeName)
	if theme == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Theme not found"})
		return
	}
	
	var puzzleResponses []models.PuzzleResponse
	
	for _, puzzle := range theme.Puzzles {
		compressedSize, uncompressedSize, _ := p.loader.GetPuzzleSizes(theme.Name, puzzle.GetName())
		
		puzzleResponse := models.PuzzleResponse{
			Name:            puzzle.GetName(),
			Difficulty:      puzzle.DescProps.Difficulty,
			Language:        puzzle.DescProps.Language,
			CompressedSize:  compressedSize,
			UncompressedSize: uncompressedSize,
			Cipher:          puzzle.Cipher,
			Obscure:         puzzle.Obscure,
			ID:              puzzle.MetaProps.ID,
			Author:          puzzle.MetaProps.Author,
			CreatedAt:       puzzle.MetaProps.Created,
			UpdatedAt:       puzzle.MetaProps.Modified,
		}
		
		puzzleResponses = append(puzzleResponses, puzzleResponse)
	}
	
	c.JSON(http.StatusOK, puzzleResponses)
}

// GetPuzzleNames godoc
// @Summary Get puzzle names
// @Description Returns names of all puzzles for a specific theme
// @Tags Puzzles
// @Produce json
// @Param theme query string true "Theme name"
// @Success 200 {array} string
// @Failure 404 {object} map[string]string
// @Router /puzzles/names [get]
func (p *PuzzleController) GetPuzzleNames(c *gin.Context) {
	themeName := c.Query("theme")
	
	theme := p.loader.GetTheme(themeName)
	if theme == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Theme not found"})
		return
	}
	
	var puzzleNames []string
	
	for _, puzzle := range theme.Puzzles {
		puzzleNames = append(puzzleNames, puzzle.GetName())
	}
	
	c.JSON(http.StatusOK, puzzleNames)
}

// GetPuzzle godoc
// @Summary Get puzzle details
// @Description Returns details about a specific puzzle
// @Tags Puzzles
// @Produce json
// @Param theme query string true "Theme name"
// @Param puzzle query string true "Puzzle name"
// @Success 200 {object} models.PuzzleResponse
// @Failure 404 {object} map[string]string
// @Router /puzzle [get]
func (p *PuzzleController) GetPuzzle(c *gin.Context) {
	themeName := c.Query("theme")
	puzzleName := c.Query("puzzle")
	
	theme := p.loader.GetTheme(themeName)
	if theme == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Theme not found"})
		return
	}
	
	var foundPuzzle *models.Puzzle
	for i, puzzle := range theme.Puzzles {
		if puzzle.GetName() == puzzleName {
			foundPuzzle = &theme.Puzzles[i]
			break
		}
	}
	
	if foundPuzzle == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Puzzle not found"})
		return
	}
	
	compressedSize, uncompressedSize, _ := p.loader.GetPuzzleSizes(theme.Name, foundPuzzle.GetName())
	
	puzzleResponse := models.PuzzleResponse{
		Name:            foundPuzzle.GetName(),
		Difficulty:      foundPuzzle.DescProps.Difficulty,
		Language:        foundPuzzle.DescProps.Language,
		CompressedSize:  compressedSize,
		UncompressedSize: uncompressedSize,
		Cipher:          foundPuzzle.Cipher,
		Obscure:         foundPuzzle.Obscure,
		ID:              foundPuzzle.MetaProps.ID,
		Author:          foundPuzzle.MetaProps.Author,
		CreatedAt:       foundPuzzle.MetaProps.Created,
		UpdatedAt:       foundPuzzle.MetaProps.Modified,
	}
	
	c.JSON(http.StatusOK, puzzleResponse)
}

// UploadPuzzle godoc
// @Summary Upload a puzzle
// @Description Uploads a new puzzle to a theme
// @Tags Puzzles
// @Accept multipart/form-data
// @Produce json
// @Param theme query string true "Theme name"
// @Param file formData file true "Puzzle file (.alghive)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /puzzle/upload [post]
// @Security Bearer
func (p *PuzzleController) UploadPuzzle(c *gin.Context) {
	themeName := c.Query("theme")
	
	theme := p.loader.GetTheme(themeName)
	if theme == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Theme not found"})
		return
	}
	
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	
	// Check file extension
	if filepath.Ext(file.Filename) != ".alghive" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only .alghive files are allowed"})
		return
	}
	
	// Save the file
	dst := filepath.Join(services.PuzzlesDir, themeName, file.Filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Puzzle uploaded"})
}

// DeletePuzzle godoc
// @Summary Delete a puzzle
// @Description Deletes a puzzle from a theme
// @Tags Puzzles
// @Produce json
// @Param theme query string true "Theme name"
// @Param puzzle query string true "Puzzle name"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /puzzle [delete]
// @Security Bearer
func (p *PuzzleController) DeletePuzzle(c *gin.Context) {
	themeName := c.Query("theme")
	puzzleName := c.Query("puzzle")
	
	theme := p.loader.GetTheme(themeName)
	if theme == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Theme not found"})
		return
	}
	
	// Check if puzzle exists
	var foundPuzzle *models.Puzzle
	for i, puzzle := range theme.Puzzles {
		if puzzle.GetName() == puzzleName {
			foundPuzzle = &theme.Puzzles[i]
			break
		}
	}
	
	if foundPuzzle == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Puzzle not found"})
		return
	}
	
	// Delete the puzzle opened directory if it exists, and the .alghive file
	puzzleDir := filepath.Join(services.PuzzlesDir, themeName, puzzleName)
	if err := services.RemoveAll(puzzleDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete puzzle directory"})
		return
	}
	alghiveFile := filepath.Join(services.PuzzlesDir, themeName, puzzleName+".alghive")
	if err := services.RemoveAll(alghiveFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete puzzle file"})
		return
	}
	
	// Remove puzzle from theme
	for i, puzzle := range theme.Puzzles {
		if puzzle.GetName() == puzzleName {
			theme.Puzzles = append(theme.Puzzles[:i], theme.Puzzles[i+1:]...)
			break
		}
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Puzzle deleted"})
}

// GeneratePuzzle godoc
// @Summary Generate puzzle input and solutions
// @Description Generates puzzle input and calculates solutions for a given puzzle
// @Tags Puzzles
// @Produce json
// @Param theme query string true "Theme name"
// @Param puzzle query string true "Puzzle name"
// @Param unique_id query string true "Unique ID for generation"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /puzzle/generate [get]
func (p *PuzzleController) GeneratePuzzle(c *gin.Context) {
	themeName := c.Query("theme")
	puzzleName := c.Query("puzzle")
	uniqueID := c.Query("unique_id")
	
	theme := p.loader.GetTheme(themeName)
	if theme == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Theme not found"})
		return
	}
	
	var foundPuzzle *models.Puzzle
	for i, puzzle := range theme.Puzzles {
		if puzzle.GetName() == puzzleName {
			foundPuzzle = &theme.Puzzles[i]
			break
		}
	}
	
	if foundPuzzle == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Puzzle not found"})
		return
	}
	
	// Use pythonRunner to generate input and solutions
	linesCount := 400 // Default value, could be made configurable
	inputLines, err := p.pythonRunner.RunForge(foundPuzzle.GetForgePath(), linesCount, uniqueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate puzzle input: " + err.Error()})
		return
	}
	
	firstSolution, err := p.pythonRunner.RunDecrypt(foundPuzzle.GetDecryptPath(), inputLines)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to solve first part: " + err.Error()})
		return
	}
	
	secondSolution, err := p.pythonRunner.RunUnveil(foundPuzzle.GetUnveilPath(), inputLines)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to solve second part: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"first_solution":  firstSolution,
		"second_solution": secondSolution,
		"input_lines":     inputLines,
	})
}
