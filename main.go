package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ahafin/appteam/internal/config"
	"github.com/ahafin/appteam/internal/generator"
	"github.com/ahafin/appteam/internal/wizard"
)

const version = "0.18.0"

func main() {
	var (
		showHelp   bool
		showVer    bool
		regen      bool
		targetDir  string
		dirFlagSet bool
	)

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-h", "--help":
			showHelp = true
		case "-v", "--version":
			showVer = true
		case "-r", "--regenerate":
			regen = true
		case "-d", "--dir":
			dirFlagSet = true
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
				fmt.Fprintln(os.Stderr, "Error: -d/--dir requires a directory argument")
				os.Exit(1)
			}
			i++
			targetDir = args[i]
			if targetDir == "" {
				fmt.Fprintln(os.Stderr, "Error: -d/--dir requires a non-empty directory argument")
				os.Exit(1)
			}
		default:
			fmt.Fprintf(os.Stderr, "Unknown option: %s\n", args[i])
			fmt.Fprintln(os.Stderr, "Run 'appteam --help' for usage.")
			os.Exit(1)
		}
	}

	if showHelp {
		fmt.Println("appteam — Claude Code Agent Team Generator")
		fmt.Println()
		fmt.Println("Usage: appteam [options]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -h, --help         Show this help message")
		fmt.Println("  -v, --version      Show version")
		fmt.Println("  -r, --regenerate   Regenerate from saved .appteam/settings.json")
		fmt.Println("  -d, --dir <folder> Target directory (created if missing)")
		fmt.Println()
		fmt.Println("Run without arguments to start the interactive wizard.")
		return
	}

	if showVer {
		fmt.Printf("appteam v%s\n", version)
		return
	}

	// Create target directory if -d was provided
	if dirFlagSet {
		if err := os.MkdirAll(targetDir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "Error: creating directory %q: %v\n", targetDir, err)
			os.Exit(1)
		}
	}

	if regen {
		color := wizard.IsTTY()
		if err := regenerate(color, targetDir, dirFlagSet); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	color := wizard.IsTTY()

	// Determine the base directory for checking existing settings
	baseDir := targetDir
	if baseDir == "" {
		var err error
		baseDir, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	// Check for existing settings
	if config.SettingsExist(baseDir) {
		cfg, loadErr := config.LoadSettings(baseDir)
		if loadErr == nil {
			s := wizard.NewStyler(color)
			fmt.Printf("%s Found existing .appteam/settings.json\n", s.Bold("●"))
			fmt.Printf("  Project: %s  |  SWEs: %d  |  Model: %s\n",
				cfg.ProjectName, len(cfg.SWEs), cfg.ModelName)
			fmt.Println()

			scanner := bufio.NewScanner(os.Stdin)
			fmt.Printf("  %s Use saved config? %s: ", s.Green("▸"), s.Dim("(Y/n)"))
			if scanner.Scan() {
				answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
				if answer == "" || answer == "y" || answer == "yes" {
					cfg.TargetDir = baseDir
					if err := generator.Generate(cfg, color); err != nil {
						fmt.Fprintf(os.Stderr, "Error: %v\n", err)
						os.Exit(1)
					}
					return
				}
			}
			fmt.Println()
		}
	}

	cfg, err := wizard.Run(os.Stdin, os.Stdout, color, targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := generator.Generate(cfg, color); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// regenerate loads settings from .appteam/settings.json and regenerates all files.
// If dirFlagSet is true, targetDir overrides the working directory.
func regenerate(color bool, targetDir string, dirFlagSet bool) error {
	dir := targetDir
	if !dirFlagSet {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting working directory: %w", err)
		}
		dir = cwd
	}
	if !config.SettingsExist(dir) {
		return fmt.Errorf("No settings.json found. Run appteam first to create one.")
	}
	cfg, err := config.LoadSettings(dir)
	if err != nil {
		return err
	}
	cfg.TargetDir = dir
	return generator.Generate(cfg, color)
}
