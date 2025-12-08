package analyzer

type PackageInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`     // Taille en octets
	SizeStr string `json:"size_str"` // Taille formatée (ex: "1.2 MB")
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

// AnalyzeDependencies analyzes all dependencies of the project
func AnalyzeDependencies() (*DependencyReport, error) {
	report := &DependencyReport{}

	//Load Dependencies Separately
	deps, devDeps, err := LoadDependenciesSeparately("package.json")
	if err != nil {
		report.Errors = append(report.Errors, "Could not load package.json: "+err.Error())
	}

	allDeclared := append([]string{}, deps...)
	allDeclared = append(allDeclared, devDeps...)

	// Detect framework
	framework := DetectFramework(allDeclared)
	report.Framework = string(framework)

	// Get framework core and dev tools packages
	frameworkCore := GetFrameworkCorePackages(framework)
	frameworkDevTools := GetFrameworkDevTools(framework)

	// Find lockfile and load sizes
	lockfile := findLockFile()
	packageSizes := make(map[string]int64)

	var installed []string
	if lockfile == "" {
		report.Errors = append(report.Errors, "No lockfile found (npm/yarn/pnpm). Cannot resolve indirect dependencies.")
	} else {
		lock, err := LoadLockFile(lockfile)
		if err != nil {
			report.Errors = append(report.Errors, "Could not read lockfile: "+err.Error())
		} else if lock != nil {
			installed = lock.AllInstalledPackages()
			packageSizes = lock.GetPackageSizes()
		}
	}

	// Scan used imports in the code
	used, err := ScanImportsInProject(".")
	if err != nil {
		report.Errors = append(report.Errors, "Error scanning imports: "+err.Error())
	}

	// Create sets for efficient comparison
	declaredSet := buildSet(allDeclared)
	usedSet := buildSet(used)
	installedSet := buildSet(installed)

	// Calculate unused dependencies (ONLY dependencies, not devDependencies)
	unused := []string{}
	for _, pkg := range deps {
		// Ignore framework core packages
		if frameworkCore[pkg] {
			continue
		}

		// Check if used in the code
		if !usedSet[pkg] {
			unused = append(unused, pkg)
		}
	}

	// Calculate unused devDependencies (optional, for info)
	unusedDev := []string{}
	for _, pkg := range devDeps {
		// Ignore framework dev tools
		if frameworkDevTools[pkg] {
			continue
		}

		// Ignore common patterns
		if isCommonDevPattern(pkg) {
			continue
		}

		// Check if used in the code (rare for devDeps)
		if !usedSet[pkg] {
			unusedDev = append(unusedDev, pkg)
		}
	}

	// Calculate indirect dependencies (installed but not declared)
	indirect := []string{}
	for pkg := range installedSet {
		if !declaredSet[pkg] && !shouldIgnoreFromIndirect(pkg) {
			indirect = append(indirect, pkg)
		}
	}

	// Add framework core packages to "used" even if not detected
	for pkg := range frameworkCore {
		if declaredSet[pkg] {
			usedSet[pkg] = true
		}
	}

	// Get all used packages
	usedList := []string{}
	for pkg := range usedSet {
		usedList = append(usedList, pkg)
	}

	// Convert to PackageInfo with sizes
	report.Declared = toPackageInfoList(allDeclared, packageSizes)
	report.Used = toPackageInfoList(usedList, packageSizes)
	report.Indirect = toPackageInfoList(indirect, packageSizes)
	report.Unused = toPackageInfoList(unused, packageSizes)
	report.UnusedDev = toPackageInfoList(unusedDev, packageSizes)

	// Calculate total size
	var totalSize int64
	for _, pkg := range report.Used {
		totalSize += pkg.Size
	}
	report.TotalSize = totalSize
	report.TotalSizeStr = formatSize(totalSize)

	return report, nil
}
