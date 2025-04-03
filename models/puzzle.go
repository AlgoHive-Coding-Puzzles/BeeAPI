package models

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"plugin"
)

// Puzzle represents a programming challenge
type Puzzle struct {
	Id		    string `json:"id"`
	Path        string `json:"-"`
	Cipher      string `json:"-"`
	Obscure     string `json:"-"`
	ForgePlugin *plugin.Plugin `json:"-"`
	DecryptPlugin *plugin.Plugin `json:"-"`
	UnveilPlugin *plugin.Plugin `json:"-"`
	MetaProps   *MetaProps `json:"-"`
	DescProps   *DescProps `json:"-"`
}

// PuzzleResponse represents puzzle data for API responses
type PuzzleResponse struct {
	Name            string `json:"name"`
	Title 		 	string `json:"title"`
	Index           string `json:"index"`
	Difficulty      string `json:"difficulty"`
	Language        string `json:"language"`
	CompressedSize  int64  `json:"compressedSize"`
	UncompressedSize int64 `json:"uncompressedSize"`
	HivecraftVersion string `json:"hivecraftVersion"`
	Cipher          string `json:"cipher"`
	Obscure         string `json:"obscure"`
	ID              string `json:"id"`
	Author          string `json:"author"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// MetaProps represents metadata XML properties for a puzzle
type MetaProps struct {
	XMLName  xml.Name `xml:"Properties"`
	Author   			string   `xml:"author"`
	Created  			string   `xml:"created"`
	Modified 			string   `xml:"modified"`
	HivecraftVersion  	string   `xml:"hivecraft-version"`
	Title    			string   `xml:"title"`
	ID       			string   `xml:"id"`
}

// DescProps represents description XML properties for a puzzle
type DescProps struct {
	XMLName    xml.Name `xml:"Properties"`
	Difficulty string   `xml:"difficulty"`
	Language   string   `xml:"language"`
	Title      string   `xml:"title"`
	Index      string   `xml:"index"`
}

// GetName returns the name of the puzzle (last part of the path)
func (p *Puzzle) GetName() string {
	return filepath.Base(p.Path)
}

// LoadMetaProps loads metadata properties from an XML file
func (p *Puzzle) LoadMetaProps(xmlContent []byte) error {
	p.MetaProps = &MetaProps{}
	return xml.Unmarshal(xmlContent, p.MetaProps)
}

// LoadDescProps loads description properties from an XML file
func (p *Puzzle) LoadDescProps(xmlContent []byte) error {
	p.DescProps = &DescProps{}
	return xml.Unmarshal(xmlContent, p.DescProps)
}

// ReadFileContent reads the content of a file given its path
func ReadFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// GetForgePath returns the path to the forge.py script
func (p *Puzzle) GetForgePath() string {
	return filepath.Join(p.Path, "forge.py")
}

// GetDecryptPath returns the path to the decrypt.py script
func (p *Puzzle) GetDecryptPath() string {
	return filepath.Join(p.Path, "decrypt.py")
}

// GetUnveilPath returns the path to the unveil.py script
func (p *Puzzle) GetUnveilPath() string {
	return filepath.Join(p.Path, "unveil.py")
}

// GetId returns the ID of the puzzle
func (p *Puzzle) GetId() string {
	return p.MetaProps.ID
}

// GetPuzzleTitle returns the title of the puzzle
func (p *Puzzle) GetPuzzleTitle() string {
	return p.DescProps.Title
}

// GetPuzzleIndex returns the index of the puzzle
func (p *Puzzle) GetPuzzleIndex() string {
	return p.DescProps.Index
}

// GetPuzzleHivecraftVersion returns the Hivecraft version of the puzzle
func (p *Puzzle) GetPuzzleHivecraftVersion() string {
	return p.MetaProps.HivecraftVersion
}