package services

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	apiKeyLength = 32 // 32 bytes = 256 bits
	apiKeyFile   = ".api-key"
)

// APIKeyManager gère la génération et validation des clés API
type APIKeyManager struct {
	apiKey    string
	keyPath   string
}

// NewAPIKeyManager crée un nouveau gestionnaire de clé API
func NewAPIKeyManager(basePath string) (*APIKeyManager, error) {
	keyPath := filepath.Join(basePath, apiKeyFile)
	manager := &APIKeyManager{
		keyPath: keyPath,
	}

	// Essaie de charger une clé existante
	err := manager.loadKey()
	if err != nil {
		// Génère et sauvegarde une nouvelle clé si aucune n'existe
		err = manager.generateAndSaveKey()
		if err != nil {
			return nil, err
		}
	}

	return manager, nil
}

// loadKey charge la clé API depuis le fichier
func (a *APIKeyManager) loadKey() error {
	data, err := os.ReadFile(a.keyPath)
	if err != nil {
		return err
	}
	a.apiKey = strings.TrimSpace(string(data))
	return nil
}

// generateAndSaveKey génère une nouvelle clé API et la sauvegarde dans le fichier
func (a *APIKeyManager) generateAndSaveKey() error {
	// Génère des octets aléatoires
	randBytes := make([]byte, apiKeyLength)
	_, err := rand.Read(randBytes)
	if err != nil {
		return err
	}

	// Encode en base64 pour une clé API lisible
	a.apiKey = base64.StdEncoding.EncodeToString(randBytes)

	// Sauvegarde dans le fichier
	return os.WriteFile(a.keyPath, []byte(a.apiKey), 0600)
}

// GetAPIKey retourne la clé API courante
func (a *APIKeyManager) GetAPIKey() string {
	return a.apiKey
}

// ValidateKey vérifie si la clé fournie correspond à celle stockée
func (a *APIKeyManager) ValidateKey(key string) bool {
	log.Printf("Validating API key: %s\n", a.apiKey)
	log.Printf("Provided API key: %s\n", key)
	return key == a.apiKey
}
