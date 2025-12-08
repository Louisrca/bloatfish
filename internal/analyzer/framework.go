package analyzer

import (
	"strings"
)

type Framework string

const (
	FrameworkAngular Framework = "angular"
	FrameworkReact   Framework = "react"
	FrameworkVue     Framework = "vue"
	FrameworkSvelte  Framework = "svelte"
	FrameworkNext    Framework = "next"
	FrameworkNuxt    Framework = "nuxt"
	FrameworkUnknown Framework = "unknown"
)

// DetectFramework detects the main framework of the project
func DetectFramework(deps []string) Framework {
	depsSet := buildSet(deps)

	// Next.js (check before React because Next includes React)
	if depsSet["next"] {
		return FrameworkNext
	}

	// Nuxt (check before Vue because Nuxt includes Vue)
	if depsSet["nuxt"] {
		return FrameworkNuxt
	}

	// Angular
	if depsSet["@angular/core"] {
		return FrameworkAngular
	}

	// React
	if depsSet["react"] {
		return FrameworkReact
	}

	// Vue
	if depsSet["vue"] {
		return FrameworkVue
	}

	// Svelte
	if depsSet["svelte"] {
		return FrameworkSvelte
	}

	return FrameworkUnknown
}

// GetFrameworkCorePackages returns core packages according to the framework
func GetFrameworkCorePackages(framework Framework) map[string]bool {
	switch framework {
	case FrameworkAngular:
		return map[string]bool{
			"@angular/core":             true,
			"@angular/common":           true,
			"@angular/compiler":         true,
			"@angular/platform-browser": true,
			"@angular/router":           true,
			"@angular/forms":            true,
			"@angular/animations":       true,
			"@angular/ssr":              true,
			"@angular/platform-server":  true,
			"@angular/pwa":              true,
			"@angular/service-worker":   true,
			"@angular/cdk":              true,
			"@angular/material":         true,
			"rxjs":                      true,
			"tslib":                     true,
			"zone.js":                   true,
		}

	case FrameworkReact, FrameworkNext:
		return map[string]bool{
			"react":            true,
			"react-dom":        true,
			"react-router":     true,
			"react-router-dom": true,
			"next":             true, // Next.js
		}

	case FrameworkVue, FrameworkNuxt:
		return map[string]bool{
			"vue":        true,
			"vue-router": true,
			"vuex":       true,
			"pinia":      true,
			"nuxt":       true, // Nuxt
			"@nuxt/kit":  true,
		}

	case FrameworkSvelte:
		return map[string]bool{
			"svelte":                   true,
			"@sveltejs/kit":            true,
			"@sveltejs/adapter-auto":   true,
			"@sveltejs/adapter-node":   true,
			"@sveltejs/adapter-static": true,
		}

	default:
		return map[string]bool{}
	}
}

// GetFrameworkDevTools returns dev tools according to the framework
func GetFrameworkDevTools(framework Framework) map[string]bool {
	base := map[string]bool{
		// Common build tools
		"typescript": true,
		"vite":       true,
		"esbuild":    true,
		"rollup":     true,
		"webpack":    true,

		// CSS/PostCSS
		"postcss":              true,
		"autoprefixer":         true,
		"tailwindcss":          true,
		"@tailwindcss/postcss": true,
		"sass":                 true,

		// Linters & formatters
		"eslint":                           true,
		"prettier":                         true,
		"@typescript-eslint/parser":        true,
		"@typescript-eslint/eslint-plugin": true,

		// Git hooks
		"husky":       true,
		"lint-staged": true,

		// Common types
		"@types/node": true,
	}

	switch framework {
	case FrameworkAngular:
		base["@angular/cli"] = true
		base["@angular/build"] = true
		base["@angular/compiler-cli"] = true
		base["@angular-devkit/build-angular"] = true
		base["@angular-devkit/architect"] = true
		base["@angular-devkit/core"] = true
		base["@angular-devkit/schematics"] = true
		base["@schematics/angular"] = true
		base["karma"] = true
		base["karma-chrome-launcher"] = true
		base["karma-jasmine"] = true
		base["karma-jasmine-html-reporter"] = true
		base["karma-coverage"] = true
		base["jasmine-core"] = true
		base["@types/jasmine"] = true

	case FrameworkReact, FrameworkNext:
		base["@types/react"] = true
		base["@types/react-dom"] = true
		base["@vitejs/plugin-react"] = true
		base["@vitejs/plugin-react-swc"] = true
		base["eslint-plugin-react"] = true
		base["eslint-plugin-react-hooks"] = true
		base["@testing-library/react"] = true
		base["@testing-library/jest-dom"] = true

	case FrameworkVue, FrameworkNuxt:
		base["@vitejs/plugin-vue"] = true
		base["@vue/compiler-sfc"] = true
		base["@vue/test-utils"] = true
		base["eslint-plugin-vue"] = true
		base["vue-tsc"] = true

	case FrameworkSvelte:
		base["@sveltejs/vite-plugin-svelte"] = true
		base["svelte-check"] = true
		base["@sveltejs/package"] = true
	}

	return base
}

func isCommonDevPattern(pkg string) bool {
	// All @types/*
	if strings.HasPrefix(pkg, "@types/") {
		return true
	}

	// All plugins from frameworks
	patterns := []string{
		"@angular-devkit/",
		"@schematics/",
		"@vitejs/plugin-",
		"eslint-plugin-",
		"@sveltejs/adapter-",
		"vite-plugin-",
		"rollup-plugin-",
		"babel-plugin-",
		"@babel/plugin-",
	}

	for _, pattern := range patterns {
		if strings.HasPrefix(pkg, pattern) {
			return true
		}
	}

	return false
}

func shouldIgnoreFromIndirect(pkg string) bool {
	// Patterns to ignore from indirect dependencies
	patterns := []string{
		"@types/",
		"@angular-devkit/",
		"@schematics/",
		"@esbuild/",
		"@rollup/",
		"@babel/",
		"@jridgewell/",
		"@msgpackr-extract/",
		"@napi-rs/",
	}

	for _, pattern := range patterns {
		if strings.HasPrefix(pkg, pattern) {
			return true
		}
	}

	// Common individual packages to ignore
	common := map[string]bool{
		".bin": true, ".cache": true,
		"esbuild": true, "rollup": true,
		"typescript": true,
	}

	return common[pkg]
}
