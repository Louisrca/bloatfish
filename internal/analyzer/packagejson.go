package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
)

type packageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// LoadDependenciesSeparately returns dependencies and devDependencies from package.json
func LoadDependenciesSeparately(path string) (deps []string, devDeps []string, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot read package.json: %w", err)
	}

	var pkg packageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, nil, fmt.Errorf("cannot parse package.json: %w", err)
	}

	deps = make([]string, 0, len(pkg.Dependencies))
	for name := range pkg.Dependencies {
		deps = append(deps, name)
	}

	devDeps = make([]string, 0, len(pkg.DevDependencies))
	for name := range pkg.DevDependencies {
		devDeps = append(devDeps, name)
	}

	return deps, devDeps, nil
}
