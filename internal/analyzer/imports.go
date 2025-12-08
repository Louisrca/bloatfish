package analyzer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// import { Component } from '@angular/core';
	importFromRegex = regexp.MustCompile(`import\s+(?:type\s+)?(?:\{[^}]*\}|\*\s+as\s+\w+|\w+)\s+from\s+['"]([^'"]+)['"]`)

	// import '@angular/platform-browser';
	importBareRegex = regexp.MustCompile(`import\s+['"]([^'"]+)['"]`)

	// const express = require('express');
	requireRegex = regexp.MustCompile(`require\s*\(\s*['"]([^'"]+)['"]\s*\)`)

	// import('@angular/core')
	dynamicImportRegex = regexp.MustCompile(`import\s*\(\s*['"]([^'"]+)['"]\s*\)`)
)

func ScanImportsInProject(root string) ([]string, error) {
	result := map[string]bool{}

	// ensure root exists
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory %s does not exist", root)
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			// Ignore unnecessary folders
			name := d.Name()
			if name == "node_modules" || name == ".git" ||
				name == "dist" || name == "build" || name == ".angular" {
				return filepath.SkipDir
			}
			return nil
		}

		// Only source files
		if !isSourceFile(path) {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		text := string(content)

		// Try all regexes
		allMatches := [][]string{}

		allMatches = append(allMatches, importFromRegex.FindAllStringSubmatch(text, -1)...)
		allMatches = append(allMatches, importBareRegex.FindAllStringSubmatch(text, -1)...)
		allMatches = append(allMatches, requireRegex.FindAllStringSubmatch(text, -1)...)
		allMatches = append(allMatches, dynamicImportRegex.FindAllStringSubmatch(text, -1)...)

		for _, match := range allMatches {
			if len(match) < 2 {
				continue
			}

			importPath := match[1]
			pkg := extractPackageName(importPath)

			if pkg != "" {
				result[pkg] = true
			}
		}

		return nil
	})

	list := []string{}
	for pkg := range result {
		list = append(list, pkg)
	}

	return list, err
}

func isSourceFile(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".ts" || ext == ".js" || ext == ".tsx" || ext == ".jsx"
}

// extractPackageName extracts the correct package name from an import
// Properly handles scoped packages
func extractPackageName(importPath string) string {
	if importPath == "" {
		return ""
	}

	// Ignore relative imports
	if strings.HasPrefix(importPath, ".") || strings.HasPrefix(importPath, "/") {
		return ""
	}

	// Ignore Node.js built-in imports
	nodeBuiltins := map[string]bool{
		"fs": true, "path": true, "http": true, "https": true,
		"crypto": true, "os": true, "util": true, "stream": true,
		"buffer": true, "events": true, "url": true, "querystring": true,
		"child_process": true, "cluster": true, "net": true, "tls": true,
		"dns": true, "dgram": true, "readline": true, "repl": true,
		"vm": true, "zlib": true, "assert": true, "constants": true,
		"node:fs": true, "node:path": true, "node:http": true,
	}

	firstPart := strings.Split(importPath, "/")[0]
	if nodeBuiltins[firstPart] {
		return ""
	}

	// Scoped packages: @angular/core, @angular/common/http
	if strings.HasPrefix(importPath, "@") {
		parts := strings.SplitN(importPath, "/", 3)
		if len(parts) >= 2 {
			return parts[0] + "/" + parts[1]
		}
		return importPath
	}

	// Normal packages: lodash, primeng/button
	idx := strings.Index(importPath, "/")
	if idx > 0 {
		return importPath[:idx]
	}

	return importPath
}
