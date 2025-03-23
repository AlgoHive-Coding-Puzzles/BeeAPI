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

// GetPuzzlesIds
// @Summary Get puzzle IDs
// @Description Returns IDs of all puzzles for a specific theme
// @Tags Puzzles
// @Produce json
// @Param theme query string true "Theme name"
// @Success 200 {array} string
// @Failure 404 {object} map[string]string
// @Router /puzzles/ids [get]
func (p *PuzzleController) GetPuzzlesIds(c *gin.Context) {
	themeName := c.Query("theme")
	
	theme := p.loader.GetTheme(themeName)
	if theme == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Theme not found"})
		return
	}
	
	var puzzleIds []string
	
	for _, puzzle := range theme.Puzzles {
		puzzleIds = append(puzzleIds, puzzle.MetaProps.ID)
	}
	
	c.JSON(http.StatusOK, puzzleIds)
}

// GetPuzzle godoc
// @Summary Get puzzle details
// @Description Returns details about a specific puzzle
// @Tags Puzzles
// @Produce json
// @Param theme query string true "Theme name"
// @Param puzzle query string true "Puzzle Id"
// @Success 200 {object} models.PuzzleResponse
// @Failure 404 {object} map[string]string
// @Router /puzzle [get]
func (p *PuzzleController) GetPuzzle(c *gin.Context) {
	themeName := c.Query("theme")
	puzzleId := c.Query("puzzle")
	
	theme := p.loader.GetTheme(themeName)
	if theme == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Theme not found"})
		return
	}
	
	var foundPuzzle *models.Puzzle
	for i, puzzle := range theme.Puzzles {
		if puzzle.GetId() == puzzleId {
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
// @Param puzzle query string true "Puzzle Id"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /puzzle [delete]
// @Security Bearer
func (p *PuzzleController) DeletePuzzle(c *gin.Context) {
	themeName := c.Query("theme")
	puzzleId := c.Query("puzzle")
	
	theme := p.loader.GetTheme(themeName)
	if theme == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Theme not found"})
		return
	}
	
	// Check if puzzle exists
	var foundPuzzle *models.Puzzle
	for i, puzzle := range theme.Puzzles {
		if puzzle.GetId() == puzzleId {
			foundPuzzle = &theme.Puzzles[i]
			break
		}
	}
	
	if foundPuzzle == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Puzzle not found"})
		return
	}
	
	// Delete the puzzle opened directory if it exists, and the .alghive file
	puzzleDir := filepath.Join(services.PuzzlesDir, themeName, foundPuzzle.GetName())
	if err := services.RemoveAll(puzzleDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete puzzle directory"})
		return
	}
	alghiveFile := filepath.Join(services.PuzzlesDir, themeName, foundPuzzle.GetName() + ".alghive")
	if err := services.RemoveAll(alghiveFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete puzzle file"})
		return
	}
	
	// Remove puzzle from theme
	for i, puzzle := range theme.Puzzles {
		if puzzle.GetName() == foundPuzzle.GetName() {
			theme.Puzzles = append(theme.Puzzles[:i], theme.Puzzles[i+1:]...)
			break
		}
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Puzzle deleted"})
}

// GeneratePuzzleInput godoc
// @Summary Generate puzzle input
// @Description Generates puzzle input for a given puzzle
// @Tags Puzzles
// @Produce json
// @Param theme query string true "Theme name"
// @Param puzzle query string true "Puzzle Id"
// @Param unique_id query string true "Unique ID for generation"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /puzzle/generate/input [get]
func (p *PuzzleController) GeneratePuzzleInput(c *gin.Context) {
    themeName := c.Query("theme")
    puzzleId := c.Query("puzzle")
    uniqueID := c.Query("unique_id")

    theme := p.loader.GetTheme(themeName)
    if theme == nil {
        c.JSON(http.StatusNotFound, gin.H{"message": "Theme not found"})
        return
    }

    var foundPuzzle *models.Puzzle
    for i, puzzle := range theme.Puzzles {
        if puzzle.GetId() == puzzleId {
            foundPuzzle = &theme.Puzzles[i]
            break
        }
    }

    if foundPuzzle == nil {
        c.JSON(http.StatusNotFound, gin.H{"message": "Puzzle not found"})
        return
    }

    linesCount := 400 // Default value, could be made configurable
    inputLines, err := p.pythonRunner.RunForge(foundPuzzle.GetForgePath(), linesCount, uniqueID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate puzzle input: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "input_lines": inputLines,
    })
}

// CheckFirstSolution godoc
// @Summary Check first solution
// @Description Checks if the first solution matches the provided value
// @Tags Puzzles
// @Produce json
// @Param theme query string true "Theme name"
// @Param puzzle query string true "Puzzle Id"
// @Param unique_id query string true "Unique ID for generation"
// @Param solution query string true "Solution to check"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /puzzle/check/first [get]
func (p *PuzzleController) CheckFirstSolution(c *gin.Context) {
    themeName := c.Query("theme")
    puzzleId := c.Query("puzzle")
    uniqueID := c.Query("unique_id")
    solution := c.Query("solution")

    theme := p.loader.GetTheme(themeName)
    if theme == nil {
        c.JSON(http.StatusNotFound, gin.H{"message": "Theme not found"})
        return
    }

    var foundPuzzle *models.Puzzle
    for i, puzzle := range theme.Puzzles {
        if puzzle.GetId() == puzzleId {
            foundPuzzle = &theme.Puzzles[i]
            break
        }
    }

    if foundPuzzle == nil {
        c.JSON(http.StatusNotFound, gin.H{"message": "Puzzle not found"})
        return
    }

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

    if firstSolution == solution {
        c.JSON(http.StatusOK, gin.H{"message": "First solution matches"})
    } else {
        c.JSON(http.StatusOK, gin.H{"message": "First solution does not match"})
    }
}

// CheckSecondSolution godoc
// @Summary Check second solution
// @Description Checks if the second solution matches the provided value
// @Tags Puzzles
// @Produce json
// @Param theme query string true "Theme name"
// @Param puzzle query string true "Puzzle Id"
// @Param unique_id query string true "Unique ID for generation"
// @Param solution query string true "Solution to check"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /puzzle/check/second [get]
func (p *PuzzleController) CheckSecondSolution(c *gin.Context) {
    themeName := c.Query("theme")
    puzzleId := c.Query("puzzle")
    uniqueID := c.Query("unique_id")
    solution := c.Query("solution")

    theme := p.loader.GetTheme(themeName)
    if theme == nil {
        c.JSON(http.StatusNotFound, gin.H{"message": "Theme not found"})
        return
    }

    var foundPuzzle *models.Puzzle
    for i, puzzle := range theme.Puzzles {
        if puzzle.GetId() == puzzleId {
            foundPuzzle = &theme.Puzzles[i]
            break
        }
    }

    if foundPuzzle == nil {
        c.JSON(http.StatusNotFound, gin.H{"message": "Puzzle not found"})
        return
    }

    linesCount := 400 // Default value, could be made configurable
    inputLines, err := p.pythonRunner.RunForge(foundPuzzle.GetForgePath(), linesCount, uniqueID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate puzzle input: " + err.Error()})
        return
    }

    secondSolution, err := p.pythonRunner.RunUnveil(foundPuzzle.GetUnveilPath(), inputLines)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to solve second part: " + err.Error()})
        return
    }

    if secondSolution == solution {
        c.JSON(http.StatusOK, gin.H{"message": "Second solution matches"})
    } else {
        c.JSON(http.StatusOK, gin.H{"message": "Second solution does not match"})
    }
}