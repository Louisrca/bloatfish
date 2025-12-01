package analyser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type UnusedPackageAnalyser struct {
	PackagesDeclared []string `json:"packages_declared"`
	PackagesUsed     []string `json:"packages_used"`
	PackageUnused    []string `json:"package_unused"`
	Errors           []string `json:"errors,omitempty"`
}

func (upa *UnusedPackageAnalyser) Analyze(directionPath string) {

	dirPath, err := os.ReadDir(directionPath)

	if err != nil {
		upa.Errors = append(upa.Errors, fmt.Sprintf("Failed to read directory: %v", err))
		return
	}

	for _, entry := range dirPath {

		fullPath := filepath.Join(directionPath, entry.Name())
		dir, err := DirSize(fullPath)
		if err != nil {
			upa.Errors = append(upa.Errors, fmt.Sprintf("Failed to get size for %s: %v", entry.Name(), err))
			continue
		}

		fmt.Printf("Size of %s: %s\n", entry.Name(), BytesConversion(dir))

		// check if package is used
		used, err := CheckIfPackageIsInSrcFile(entry.Name())
		if err != nil {
			upa.Errors = append(upa.Errors, fmt.Sprintf("Failure checking usage for %s: %v", entry.Name(), err))
			continue
		}

		if used {
			upa.PackagesUsed = append(upa.PackagesUsed, entry.Name())
		} else {
			upa.PackageUnused = append(upa.PackageUnused, entry.Name())
		}

	}
	WriteJSONReport(upa)
}

func DirSize(path string) (int64, error) {
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

func BytesConversion(size int64) string {
	const (
		KB = 1 << (10 * 1)
		MB = 1 << (10 * 2)
		GB = 1 << (10 * 3)
	)

	red := "\033[31m"
	yellow := "\033[33m"
	blue := "\033[34m"
	grey := "\033[90m"
	reset := "\033[0m"

	switch {
	case size >= GB:
		return fmt.Sprintf("%s%.2f GB%s", red, float64(size)/GB, reset)
	case size >= MB:
		return fmt.Sprintf("%s%.2f MB%s", yellow, float64(size)/MB, reset)
	case size >= KB:
		return fmt.Sprintf("%s%.2f KB%s", blue, float64(size)/KB, reset)
	default:
		return fmt.Sprintf("%s%d bytes%s", grey, size, reset)
	}
}

func CheckIfPackageIsInSrcFile(packageName string) (bool, error) {
	return searchInDirectory("./src", packageName)
}

func searchInDirectory(dir string, packageName string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			// recurse into sub-directories
			found, err := searchInDirectory(path, packageName)
			if err != nil {
				return false, err
			}
			if found {
				return true, nil
			}
			continue
		}

		// read only .js, .ts, .tsx, .jsx
		if !strings.HasSuffix(path, ".js") &&
			!strings.HasSuffix(path, ".ts") &&
			!strings.HasSuffix(path, ".tsx") &&
			!strings.HasSuffix(path, ".jsx") {
			continue
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return false, err
		}

		importLine := fmt.Sprintf("from \"%s\"", packageName)
		if strings.Contains(string(content), importLine) {
			return true, nil
		}
	}

	return false, nil
}

func WriteJSONReport(upa *UnusedPackageAnalyser) {
	file, err := os.Create("unused_packages_report.json")
	if err != nil {
		fmt.Printf("Failed to create report file: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(upa); err != nil {
		fmt.Printf("Failed to write report: %v\n", err)
	}
}
