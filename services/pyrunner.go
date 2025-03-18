package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// PythonRunner provides utilities for running Python scripts
type PythonRunner struct {
	PythonPath string // Path to python interpreter
}

// NewPythonRunner creates a new PythonRunner with the given Python path
func NewPythonRunner(pythonPath string) *PythonRunner {
	if pythonPath == "" {
		// Default to "python3" if not specified
		pythonPath = "python3"
	}
	return &PythonRunner{PythonPath: pythonPath}
}

// RunForge executes a forge.py script with the given lines count and unique ID
func (p *PythonRunner) RunForge(scriptPath string, linesCount int, uniqueID string) ([]string, error) {
	// Get the directory where the script is located
	// scriptDir := filepath.Dir(scriptPath)

	// Create a Python script that directly returns the lines
	tempScript := `import sys
import importlib.util

# Import the forge module
spec = importlib.util.spec_from_file_location("forge", sys.argv[1])
forge_module = importlib.util.module_from_spec(spec)
sys.modules["forge"] = forge_module
spec.loader.exec_module(forge_module)

# Initialize and run
lines_count = int(sys.argv[2])
unique_id = sys.argv[3]
forge = forge_module.Forge(lines_count, unique_id)
lines = forge.run()

# Print lines directly to stdout
for line in lines:
    print(line)
`

	// Create a temporary Python file
	tempFile, err := os.CreateTemp("", "forge_runner_*.py")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(tempScript); err != nil {
		return nil, fmt.Errorf("failed to write temp script: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temp file: %v", err)
	}

	// Run the script that imports forge.py and outputs lines
	cmd := exec.Command(p.PythonPath, tempFile.Name(), scriptPath, fmt.Sprintf("%d", linesCount), uniqueID)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run forge.py: %v, stderr: %s", err, stderr.String())
	}

	// Split output into lines
	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}, nil // Return empty slice if no output
	}
	return lines, nil
}

// RunDecrypt executes a decrypt.py script with the given input lines
func (p *PythonRunner) RunDecrypt(scriptPath string, inputLines []string) (string, error) {
	// Similar approach using importlib to avoid writing to files
	tempScript := `import sys
import importlib.util

# Import the decrypt module
spec = importlib.util.spec_from_file_location("decrypt", sys.argv[1])
decrypt_module = importlib.util.module_from_spec(spec)
sys.modules["decrypt"] = decrypt_module
spec.loader.exec_module(decrypt_module)

# Read lines from stdin
lines = []
for line in sys.stdin:
    lines.append(line.rstrip('\n'))

# Initialize and run
decrypt = decrypt_module.Decrypt(lines)
solution = decrypt.run()
print(solution)
`

	return p.runPythonWithImport(tempScript, scriptPath, inputLines)
}

// RunUnveil executes an unveil.py script with the given input lines
func (p *PythonRunner) RunUnveil(scriptPath string, inputLines []string) (string, error) {
	tempScript := `import sys
import importlib.util

# Import the unveil module
spec = importlib.util.spec_from_file_location("unveil", sys.argv[1])
unveil_module = importlib.util.module_from_spec(spec)
sys.modules["unveil"] = unveil_module
spec.loader.exec_module(unveil_module)

# Read lines from stdin
lines = []
for line in sys.stdin:
    lines.append(line.rstrip('\n'))

# Initialize and run
unveil = unveil_module.Unveil(lines)
solution = unveil.run()
print(solution)
`

	return p.runPythonWithImport(tempScript, scriptPath, inputLines)
}

// Helper function to run a Python script with dynamic imports
func (p *PythonRunner) runPythonWithImport(scriptContent string, targetScriptPath string, inputLines []string) (string, error) {
	// Create a temporary Python file
	tempFile, err := os.CreateTemp("", "python_runner_*.py")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(scriptContent); err != nil {
		return "", fmt.Errorf("failed to write temp script: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file: %v", err)
	}

	// Create a process with stdin pipe
	cmd := exec.Command(p.PythonPath, tempFile.Name(), targetScriptPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Start the process
	err = cmd.Start()
	if err != nil {
		return "", fmt.Errorf("failed to start Python script: %v", err)
	}

	// Write input lines to stdin
	for _, line := range inputLines {
		fmt.Fprintln(stdin, line)
	}
	stdin.Close()

	// Wait for the process to complete
	err = cmd.Wait()
	if err != nil {
		return "", fmt.Errorf("python script execution failed: %v, stderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}
