package services

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/algohive/beeapi/models"
)

const PuzzlesDir = "puzzles"

// PuzzlesLoader handles loading/unloading puzzles from the filesystem
type PuzzlesLoader struct {
	Themes []models.Theme
	mu     sync.RWMutex
}

// NewPuzzlesLoader creates a new puzzle loader
func NewPuzzlesLoader() *PuzzlesLoader {
	return &PuzzlesLoader{
		Themes: []models.Theme{},
		mu:     sync.RWMutex{},
	}
}

// Load loads all themes and puzzles
func (p *PuzzlesLoader) Load() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.Themes = []models.Theme{} // Reset themes
	
	// Iterate through themes directory
	themeDirs, err := os.ReadDir(PuzzlesDir)
	if err != nil {
		return err
	}
	
	for _, themeDir := range themeDirs {
		if themeDir.IsDir() {
			theme := models.Theme{
				Name: themeDir.Name(),
				Path: filepath.Join(PuzzlesDir, themeDir.Name()),
				Puzzles: []models.Puzzle{},
			}
			
			// Load puzzles for this theme
			puzzleDirs, err := os.ReadDir(theme.Path)
			if err != nil {
				continue
			}
			
			for _, puzzleDir := range puzzleDirs {
				if puzzleDir.IsDir() {
					puzzlePath := filepath.Join(theme.Path, puzzleDir.Name())
					puzzle, err := p.loadPuzzle(theme.Name, puzzleDir.Name(), puzzlePath)
					if err != nil {
						continue
					}
					theme.Puzzles = append(theme.Puzzles, puzzle)
				}
			}
			
			p.Themes = append(p.Themes, theme)
		}
	}
	
	return nil
}

// Extract extracts all .alghive files within themes
func (p *PuzzlesLoader) Extract() error {
	// Iterate through themes directory
	themeDirs, err := os.ReadDir(PuzzlesDir)
	if err != nil {
		return err
	}
	
	for _, themeDir := range themeDirs {
		if themeDir.IsDir() {
			themePath := filepath.Join(PuzzlesDir, themeDir.Name())
			
			// Find .alghive files
			files, err := os.ReadDir(themePath)
			if err != nil {
				continue
			}
			
			for _, file := range files {
				if !file.IsDir() && filepath.Ext(file.Name()) == ".alghive" {
					alghiveFile := filepath.Join(themePath, file.Name())
					extractDir := filepath.Join(themePath, file.Name()[:len(file.Name())-8])
					
					// Skip if directory already exists
					if _, err := os.Stat(extractDir); err == nil {
						continue
					}
					
					// Extract the .alghive file
					err = unzip(alghiveFile, extractDir)
					if err != nil {
						continue
					}
				}
			}
		}
	}
	
	return nil
}

// Unload deletes extracted puzzle directories
func (p *PuzzlesLoader) Unload() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	for _, theme := range p.Themes {
		for _, puzzle := range theme.Puzzles {
			err := os.RemoveAll(puzzle.Path)
			if err != nil {
				return err
			}
		}
	}
	
	p.Themes = []models.Theme{} // Reset themes
	return nil
}

// Reload reloads all themes and puzzles
func (p *PuzzlesLoader) Reload() error {
	err := p.Unload()
	if err != nil {
		return err
	}
	
	err = p.Extract()
	if err != nil {
		return err
	}
	
	return p.Load()
}

// GetTheme returns a theme by name
func (p *PuzzlesLoader) GetTheme(name string) *models.Theme {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	for i, theme := range p.Themes {
		if theme.Name == name {
			return &p.Themes[i]
		}
	}
	
	return nil
}

// HasTheme checks if a theme exists
func (p *PuzzlesLoader) HasTheme(name string) bool {
	return p.GetTheme(name) != nil
}

func (p *PuzzlesLoader) CreateTheme(name string) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    // Avoid calling HasTheme here to prevent a lock conflict
    for _, theme := range p.Themes {
        if theme.Name == name {
            return os.ErrExist
        }
    }

    themePath := filepath.Join(PuzzlesDir, name)
    err := os.MkdirAll(themePath, 0755)
    if err != nil {
        return err
    }

    p.Themes = append(p.Themes, models.Theme{
        Name:    name,
        Path:    themePath,
        Puzzles: []models.Puzzle{},
    })

    return nil
}

func (p *PuzzlesLoader) DeleteTheme(name string) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    idx := -1
    var themePath string

    for i, theme := range p.Themes {
        if theme.Name == name {
            idx = i
            themePath = theme.Path
            break
        }
    }

    if idx == -1 {
        return os.ErrNotExist
    }

    err := os.RemoveAll(themePath)
    if err != nil {
        return err
    }

    // Remove theme from slice
    p.Themes = append(p.Themes[:idx], p.Themes[idx+1:]...)

    return nil
}

// GetPuzzleSizes returns the compressed and uncompressed sizes of a puzzle
func (p *PuzzlesLoader) GetPuzzleSizes(themeName, puzzleName string) (int64, int64, error) {
	theme := p.GetTheme(themeName)
	if theme == nil {
		return 0, 0, os.ErrNotExist
	}
	
	alghiveFile := filepath.Join(theme.Path, puzzleName+".alghive")
	puzzleDir := filepath.Join(theme.Path, puzzleName)
	
	alghiveInfo, err := os.Stat(alghiveFile)
	if err != nil {
		return 0, 0, err
	}
	
	dirSize, err := getDirSize(puzzleDir)
	if err != nil {
		return 0, 0, err
	}
	
	return alghiveInfo.Size(), dirSize, nil
}

// HotSwap replaces a puzzle with another one with the same ID
func (p *PuzzlesLoader) HotSwap(themeName string, puzzleID string, newPuzzleFile string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Find theme directly without using GetTheme to avoid deadlock
	var theme *models.Theme
	for i, t := range p.Themes {
		if t.Name == themeName {
			theme = &p.Themes[i]
			break
		}
	}

	if theme == nil {
		return errors.New("theme not found")
	}

	// Find puzzle with the given ID
	var foundPuzzle *models.Puzzle
	var puzzleName string
	var puzzleIndex int
	for i, puzzle := range theme.Puzzles {
		if puzzle.GetId() == puzzleID {
			foundPuzzle = &theme.Puzzles[i]
			puzzleName = puzzle.GetName()
			puzzleIndex = i
			break
		}
	}

	if foundPuzzle == nil {
		return errors.New("puzzle not found")
	}

	// Create temporary directory for extraction
	tempDir := filepath.Join(os.TempDir(), "puzzle_hotswap_"+puzzleID)
	defer os.RemoveAll(tempDir)
	
	// Extract new puzzle to temp directory
	err := unzip(newPuzzleFile, tempDir)
	if err != nil {
		return fmt.Errorf("failed to extract new puzzle: %w", err)
	}

	// Load new puzzle to verify ID
	newPuzzle, err := p.loadPuzzle(themeName, filepath.Base(tempDir), tempDir)
	if err != nil {
		return fmt.Errorf("failed to load new puzzle: %w", err)
	}

	// Verify that the new puzzle has the same ID
	if newPuzzle.GetId() != puzzleID {
		return fmt.Errorf("new puzzle ID (%s) does not match expected ID (%s)", newPuzzle.GetId(), puzzleID)
	}

	// Get paths
	oldPuzzlePath := foundPuzzle.Path
	oldAlghiveFile := filepath.Join(theme.Path, puzzleName+".alghive")
	newAlghiveFile := filepath.Join(theme.Path, puzzleName+".alghive")

	// Backup old puzzle file
	backupFile := oldAlghiveFile + ".backup"
	if err := copyFile(oldAlghiveFile, backupFile); err != nil {
		return fmt.Errorf("failed to backup old puzzle file: %w", err)
	}

	// Replace old .alghive file with new one
	if err := copyFile(newPuzzleFile, newAlghiveFile); err != nil {
		// Restore backup if copy fails
		os.Rename(backupFile, oldAlghiveFile)
		return fmt.Errorf("failed to replace puzzle file: %w", err)
	}

	// Remove backup file
	os.Remove(backupFile)

	// Remove old extracted puzzle directory
	if err := os.RemoveAll(oldPuzzlePath); err != nil {
		return fmt.Errorf("failed to remove old puzzle directory: %w", err)
	}

	// Extract new puzzle
	if err := unzip(newAlghiveFile, oldPuzzlePath); err != nil {
		return fmt.Errorf("failed to extract new puzzle: %w", err)
	}

	// Reload puzzle
	newLoadedPuzzle, err := p.loadPuzzle(themeName, puzzleName, oldPuzzlePath)
	if err != nil {
		return fmt.Errorf("failed to reload puzzle: %w", err)
	}

	// Update puzzle in memory directly
	theme.Puzzles[puzzleIndex] = newLoadedPuzzle

	return nil
}

// Helper function to copy a file
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return dstFile.Sync()
}

// Helper functions

// loadPuzzle loads a puzzle from the filesystem
func (p *PuzzlesLoader) loadPuzzle(themeName, puzzleName, puzzlePath string) (models.Puzzle, error) {
	puzzle := models.Puzzle{
		Path: puzzlePath,
	}
	
	// Read cipher.html
	cipherContent, err := models.ReadFileContent(filepath.Join(puzzlePath, "cipher.html"))
	if err != nil {
		return puzzle, err
	}
	puzzle.Cipher = cipherContent
	
	// Read obscure.html (unveil.html)
	obscureContent, err := models.ReadFileContent(filepath.Join(puzzlePath, "obscure.html"))
	if err != nil {
		// Try unveil.html as alternative
		obscureContent, err = models.ReadFileContent(filepath.Join(puzzlePath, "unveil.html"))
		if err != nil {
			return puzzle, err
		}
	}
	puzzle.Obscure = obscureContent
	
	// Read XML properties
	metaXML, err := os.ReadFile(filepath.Join(puzzlePath, "props", "meta.xml"))
	if err != nil {
		return puzzle, err
	}
	
	descXML, err := os.ReadFile(filepath.Join(puzzlePath, "props", "desc.xml"))
	if err != nil {
		return puzzle, err
	}
	
	err = puzzle.LoadMetaProps(metaXML)
	if err != nil {
		return puzzle, err
	}
	
	err = puzzle.LoadDescProps(descXML)
	if err != nil {
		return puzzle, err
	}
	
	return puzzle, nil
}

// getDirSize calculates the size of a directory in bytes
func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// unzip extracts a zip file to a destination directory
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// Create destination directory
	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}

	// Extract each file
	for _, f := range r.File {
		extractFile(f, dest)
	}

	return nil
}

// extractFile extracts a single file from a zip archive
func extractFile(f *zip.File, dest string) error {
	filePath := filepath.Join(dest, f.Name)
	
	// Create directory if needed
	if f.FileInfo().IsDir() {
		return os.MkdirAll(filePath, 0755)
	}
	
	// Ensure parent directory exists
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}
	
	// Extract file
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	
	_, err = io.Copy(outFile, rc)
	return err
}
