package services

import (
	"archive/zip"
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

// CreateTheme creates a new theme
func (p *PuzzlesLoader) CreateTheme(name string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.HasTheme(name) {
		return os.ErrExist
	}
	
	themePath := filepath.Join(PuzzlesDir, name)
	err := os.MkdirAll(themePath, 0755)
	if err != nil {
		return err
	}
	
	p.Themes = append(p.Themes, models.Theme{
		Name: name,
		Path: themePath,
		Puzzles: []models.Puzzle{},
	})
	
	return nil
}

// DeleteTheme deletes a theme
func (p *PuzzlesLoader) DeleteTheme(name string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	theme := p.GetTheme(name)
	if theme == nil {
		return os.ErrNotExist
	}
	
	err := os.RemoveAll(theme.Path)
	if err != nil {
		return err
	}
	
	// Remove theme from slice
	for i, t := range p.Themes {
		if t.Name == name {
			p.Themes = append(p.Themes[:i], p.Themes[i+1:]...)
			break
		}
	}
	
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
	
	// Load plugins (in a real implementation, you would compile Go plugins from Python code)
	// For now, we'll just create placeholder plugins that will be replaced by actual implementation
	
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
