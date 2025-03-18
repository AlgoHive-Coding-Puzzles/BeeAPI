package models

// Theme represents a collection of puzzles
type Theme struct {
	Name    string   `json:"name"`
	Path    string   `json:"-"`
	Puzzles []Puzzle `json:"puzzles"`
}

// ThemeResponse represents a theme with additional information
type ThemeResponse struct {
	Name         string          `json:"name"`
	EnigmesCount int             `json:"enigmes_count"`
	Puzzles      []PuzzleResponse `json:"puzzles"`
	Size         int64           `json:"size"`
}
