package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type PackageInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`     // Taille en octets
	SizeStr string `json:"size_str"` // Taille formatée (ex: "1.2 MB")
}

func ToPackageInfoList(names []string, sizes map[string]int64) []PackageInfo {
	seen := make(map[string]bool)
	result := []PackageInfo{}

	for _, name := range names {
		if !seen[name] {
			seen[name] = true
			size := sizes[name]
			result = append(result, PackageInfo{
				Name:    name,
				Size:    size,
				SizeStr: FormatSize(size),
			})
		}
	}

	return result
}

func WriteJSONReport(report interface{}, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Failed to create report file: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		fmt.Printf("Failed to write report: %v\n", err)
	}
}

func FindLockFile() string {
	if FileExists("package-lock.json") {
		return "package-lock.json"
	}
	if FileExists("yarn.lock") {
		return "yarn.lock"
	}
	if FileExists("pnpm-lock.yaml") {
		return "pnpm-lock.yaml"
	}
	return ""
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func BuildSet(list []string) map[string]bool {
	set := make(map[string]bool)
	for _, item := range list {
		set[item] = true
	}
	return set
}

func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
