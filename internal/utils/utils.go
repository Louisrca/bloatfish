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

type Meta struct {
	Tool    string `json:"tool"`
	Version string `json:"version,omitempty"`
}

type ChromeDPReport struct {
	URL           string  `json:"URL"`
	Requests      int     `json:"Requests"`
	TransferredKB float64 `json:"TransferredKB"`
	DomNodes      int     `json:"DomNodes"`
	HTMLSizeKB    int     `json:"HTMLSizeKB"`
	ExternalReqs  int     `json:"ExternalReqs"`
	Score         int     `json:"Score"`
	Grade         string  `json:"Grade"`
	Deep          bool    `json:"Deep"`
}

type DependencyReport struct {
	Declared     []PackageInfo `json:"declared"`
	Used         []PackageInfo `json:"used"`
	Indirect     []PackageInfo `json:"indirect"`
	Unused       []PackageInfo `json:"unused"`
	UnusedDev    []PackageInfo `json:"unused_dev"`
	Framework    string        `json:"framework"`
	TotalSize    int64         `json:"total_size"`
	TotalSizeStr string        `json:"total_size_str"`
	Errors       []string      `json:"errors"`
}

type ReportFile struct {
	Meta    Meta         `json:"meta"`
	Reports []JSONReport `json:"reports"`
}

type JSONReport struct {
	DependencyReport *DependencyReport `json:"dependency_report,omitempty"`
	ChromeDPReport   []*ChromeDPReport `json:"chrome_dp_report,omitempty"`
}

func WriteJSONReport(newReport JSONReport, filename string) error {
	var fileData ReportFile

	if data, err := os.ReadFile(filename); err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &fileData); err != nil {
			return err
		}
	} else {
		fileData = ReportFile{
			Meta: Meta{
				Tool: "bloatfish",
			},
			Reports: []JSONReport{},
		}
	}

	fileData.Reports = append(fileData.Reports, newReport)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(fileData)
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
