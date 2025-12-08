package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
)

func toPackageInfoList(names []string, sizes map[string]int64) []PackageInfo {
	seen := make(map[string]bool)
	result := []PackageInfo{}

	for _, name := range names {
		if !seen[name] {
			seen[name] = true
			size := sizes[name]
			result = append(result, PackageInfo{
				Name:    name,
				Size:    size,
				SizeStr: formatSize(size),
			})
		}
	}

	return result
}

func WriteJSONReport(report interface{}) {
	file, err := os.Create("unused_packages_report.json")
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

func findLockFile() string {
	if fileExists("package-lock.json") {
		return "package-lock.json"
	}
	if fileExists("yarn.lock") {
		return "yarn.lock"
	}
	if fileExists("pnpm-lock.yaml") {
		return "pnpm-lock.yaml"
	}
	return ""
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func buildSet(list []string) map[string]bool {
	set := make(map[string]bool)
	for _, item := range list {
		set[item] = true
	}
	return set
}

func formatSize(bytes int64) string {
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
