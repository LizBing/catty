package build

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Run is the catty command entry point.
func Run(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: catty build --classpath PATH --main CLASS [--mode auto|aot|interp] [--policy FILE] --report FILE --output FILE [--generated-source DIR]")
		return 2
	}
	switch args[0] {
	case "build":
		return runBuild(args[1:])
	case "-h", "--help", "help":
		fmt.Println("usage: catty build --classpath PATH --main CLASS [--mode auto|aot|interp] [--policy FILE] --report FILE --output FILE [--generated-source DIR]")
		return 0
	default:
		fmt.Fprintf(os.Stderr, "catty: unknown command %q\n", args[0])
		return 2
	}
}

func runBuild(args []string) int {
	fs := flag.NewFlagSet("catty build", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	classpath := fs.String("classpath", "", "ordered class directories and JARs")
	main := fs.String("main", "", "Java binary name of the entry class")
	mode := fs.String("mode", "auto", "engine selection mode: auto, aot or interp")
	policy := fs.String("policy", "", "optional method policy JSON")
	report := fs.String("report", "", "build report output path")
	output := fs.String("output", "", "standalone executable output path")
	generated := fs.String("generated-source", "", "optional directory for generated Go source")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "catty: unexpected arguments: %s\n", strings.Join(fs.Args(), " "))
		return 2
	}
	opts := Options{
		Main:            *main,
		Mode:            *mode,
		PolicyPath:      *policy,
		ReportPath:      *report,
		OutputPath:      *output,
		GeneratedSource: *generated,
	}
	if *classpath != "" {
		opts.Classpath = filepath.SplitList(*classpath)
	}
	moduleRoot, err := findModuleRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err := executeBuild(opts, moduleRoot); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func findModuleRoot() (string, error) {
	if env := os.Getenv("CATTY_MODULE_ROOT"); env != "" {
		if isCattyRoot(env) {
			return filepath.Abs(env)
		}
		return "", fmt.Errorf("catty: CATTY_MODULE_ROOT=%s is not a catty module root", env)
	}
	var candidates []string
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, wd)
	}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Dir(exe))
	}
	for _, start := range candidates {
		dir := start
		for {
			if isCattyRoot(dir) {
				return dir, nil
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	return "", fmt.Errorf("catty: cannot find the catty module root (run from the repository or set CATTY_MODULE_ROOT)")
}

func isCattyRoot(dir string) bool {
	data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "module catty" {
			return true
		}
	}
	return false
}

func executeBuild(opts Options, moduleRoot string) error {
	if opts.Mode == "" {
		opts.Mode = "auto"
	}
	outAbs, err := filepath.Abs(opts.OutputPath)
	if err != nil {
		return fmt.Errorf("catty: output path: %w", err)
	}
	reportAbs, err := filepath.Abs(opts.ReportPath)
	if err != nil {
		return fmt.Errorf("catty: report path: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(outAbs), 0o755); err != nil {
		return fmt.Errorf("catty: output directory: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(reportAbs), 0o755); err != nil {
		return fmt.Errorf("catty: report directory: %w", err)
	}
	stage, err := os.MkdirTemp(filepath.Dir(outAbs), ".catty-build-*")
	if err != nil {
		return fmt.Errorf("catty: cannot create staging directory: %w", err)
	}
	defer os.RemoveAll(stage)

	b := newBuild(opts)
	if err := b.prepare(); err != nil {
		return err
	}
	b.report.Toolchains.Java = javaVersion()

	mainSource, generatedSource, err := b.generateSource()
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(stage, "main.go"), []byte(mainSource), 0o644); err != nil {
		return fmt.Errorf("catty: staging main.go: %w", err)
	}
	if err := os.WriteFile(filepath.Join(stage, generatedFileName), []byte(generatedSource), 0o644); err != nil {
		return fmt.Errorf("catty: staging %s: %w", generatedFileName, err)
	}
	if err := os.WriteFile(filepath.Join(stage, "go.mod"), []byte(moduleFile(moduleRoot)), 0o644); err != nil {
		return fmt.Errorf("catty: staging go.mod: %w", err)
	}
	stagedBinary := filepath.Join(stage, "app")
	buildCmd := exec.Command("go", "build", "-trimpath", "-buildvcs=false", "-o", stagedBinary, ".")
	buildCmd.Dir = stage
	buildCmd.Env = append(os.Environ(), "GOWORK=off")
	var stdout, stderr bytes.Buffer
	buildCmd.Stdout = &stdout
	buildCmd.Stderr = &stderr
	if err := buildCmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("catty: generated Go build failed: %s", msg)
	}
	if info, err := os.Stat(stagedBinary); err != nil || info.Size() == 0 {
		return fmt.Errorf("catty: generated Go build produced no executable")
	}

	if opts.GeneratedSource != "" {
		dir, err := filepath.Abs(opts.GeneratedSource)
		if err != nil {
			return fmt.Errorf("catty: generated-source path: %w", err)
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("catty: generated-source directory: %w", err)
		}
		files := map[string][]byte{
			"main.go":         []byte(mainSource),
			generatedFileName: []byte(generatedSource),
			"go.mod":          []byte(moduleFile(moduleRoot)),
		}
		for name, data := range files {
			if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
				return fmt.Errorf("catty: writing generated source %s: %w", name, err)
			}
		}
	}

	reportData, err := json.MarshalIndent(b.report, "", "  ")
	if err != nil {
		return fmt.Errorf("catty: cannot encode build report: %w", err)
	}
	reportData = append(reportData, '\n')
	reportTemp, err := os.CreateTemp(filepath.Dir(reportAbs), ".catty-report-*")
	if err != nil {
		return fmt.Errorf("catty: report staging: %w", err)
	}
	reportTempName := reportTemp.Name()
	defer os.Remove(reportTempName)
	if _, err := reportTemp.Write(reportData); err != nil {
		_ = reportTemp.Close()
		return fmt.Errorf("catty: writing report: %w", err)
	}
	if err := reportTemp.Close(); err != nil {
		return fmt.Errorf("catty: writing report: %w", err)
	}
	if err := os.Chmod(reportTempName, 0o644); err != nil {
		return fmt.Errorf("catty: report permissions: %w", err)
	}
	if err := os.Rename(stagedBinary, outAbs); err != nil {
		return fmt.Errorf("catty: publishing executable: %w", err)
	}
	if err := os.Rename(reportTempName, reportAbs); err != nil {
		return fmt.Errorf("catty: publishing report: %w", err)
	}
	return nil
}

func javaVersion() string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "javac", "-version").CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
