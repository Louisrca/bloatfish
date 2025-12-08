package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type LockFile interface {
	AllInstalledPackages() []string
	GetPackageSizes() map[string]int64
}

type NPMPackageEntry struct {
	Version      string            `json:"version"`
	Dev          bool              `json:"dev"`
	Dependencies map[string]string `json:"dependencies"`
}

type NPMLockFile struct {
	Packages map[string]NPMPackageEntry `json:"packages"`
}

func (n *NPMLockFile) AllInstalledPackages() []string {
	result := []string{}
	for path := range n.Packages {
		if path == "" {
			continue // root package
		}

		// Extract package name from path
		// node_modules/@angular/core → @angular/core
		// node_modules/lodash → lodash

		name := extractPackageNameFromPath(path)
		if name != "" {
			result = append(result, name)
		}
	}
	return result
}

func (n *NPMLockFile) GetPackageSizes() map[string]int64 {
	sizes := make(map[string]int64)

	for path := range n.Packages {
		if path == "" {
			continue // root package
		}

		name := extractPackageNameFromPath(path)
		if name == "" {
			continue
		}

		// Calculate package size by scanning node_modules
		size := calculatePackageSize(path)

		// If the package already exists, keep the largest size
		if existing, ok := sizes[name]; !ok || size > existing {
			sizes[name] = size
		}
	}

	return sizes
}

func calculatePackageSize(packagePath string) int64 {

	// Check if the directory exists
	info, err := os.Stat(packagePath)
	if err != nil || !info.IsDir() {
		return 0
	}

	var totalSize int64

	err = filepath.WalkDir(packagePath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		totalSize += info.Size()
		return nil
	})

	return totalSize
}

func extractPackageNameFromPath(path string) string {
	// Remove "node_modules/" prefix
	path = strings.TrimPrefix(path, "node_modules/")

	// Handle scoped packages: @angular/core
	if strings.HasPrefix(path, "@") {
		parts := strings.SplitN(path, "/", 3)
		if len(parts) >= 2 {
			return parts[0] + "/" + parts[1]
		}
		return ""
	}

	// Handle normal packages: lodash/fp → lodash
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}

// LoadLockFile loads the appropriate lock file
func LoadLockFile(path string) (LockFile, error) {
	if strings.HasSuffix(path, "package-lock.json") {
		return loadNPMLockFile(path)
	}

	// TODO: Add support for yarn.lock and pnpm-lock.yaml
	return nil, fmt.Errorf("unsupported lock file format: %s", path)
}

func loadNPMLockFile(path string) (*NPMLockFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read package-lock.json: %w", err)
	}

	var lock NPMLockFile
	if err := json.Unmarshal(data, &lock); err != nil {
		return nil, fmt.Errorf("cannot parse package-lock.json: %w", err)
	}

	return &lock, nil
}
