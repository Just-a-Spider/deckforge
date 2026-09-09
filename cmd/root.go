package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"deckforge/internal/compiler"
	"deckforge/internal/exporter"
	"deckforge/internal/models"
	"deckforge/internal/scaffold"
	"deckforge/internal/server"
	"deckforge/internal/theme"
	"deckforge/internal/tui"
	"deckforge/internal/workspace"

	tea "github.com/charmbracelet/bubbletea"
)

// Execute handles CLI routing or launches interactive TUI
func Execute() {
	if len(os.Args) < 2 {
		cwd, err := os.Getwd()
		if err == nil {
			runTUI(cwd)
		} else {
			runTUI("")
		}
		return
	}

	command := os.Args[1]

	switch command {
	case "tui":
		path := ""
		if len(os.Args) > 2 {
			path = os.Args[2]
		} else {
			path, _ = os.Getwd()
		}
		runTUI(path)

	case "init", "new", "create":
		runInit(os.Args[2:])

	case "build":
		runBuild(os.Args[2:])

	case "serve":
		runServe(os.Args[2:])

	case "export":
		runExport(os.Args[2:])

	case "list":
		runList(os.Args[2:])

	case "slide":
		runSlide(os.Args[2:])

	case "theme":
		runTheme(os.Args[2:])

	case "book":
		runBook(os.Args[2:])

	case "agent":
		runAgent(os.Args[2:])

	case "mcp":
		deckPath := "."
		if len(os.Args) > 2 {
			deckPath = os.Args[2]
		}
		if err := server.RunMCPServer(deckPath); err != nil {
			fmt.Printf("MCP Server error: %v\n", err)
			os.Exit(1)
		}

	case "-h", "--help", "help":
		printHelp()

	default:
		// If unknown command, treat as path or print help
		if fi, err := os.Stat(command); err == nil && fi.IsDir() {
			runTUI(command)
		} else {
			fmt.Printf("Unknown command: %s\n\n", command)
			printHelp()
			os.Exit(1)
		}
	}
}

func runTUI(initialPath string) {
	if initialPath == "" {
		if cwd, err := os.Getwd(); err == nil {
			initialPath = cwd
		}
	}
	if abs, err := filepath.Abs(initialPath); err == nil {
		initialPath = abs
	}
	app := tui.NewAppModel(initialPath)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

func normalizeArgs(args []string) []string {
	var flags []string
	var pos []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
			if !strings.Contains(arg, "=") && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				flags = append(flags, args[i+1])
				i++
			}
		} else {
			pos = append(pos, arg)
		}
	}
	return append(flags, pos...)
}

func runInit(args []string) {
	norm := normalizeArgs(args)
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	targetPath := fs.String("path", "", "Destination directory path")
	themeName := fs.String("theme", "academic-crimson", "Base theme preset")
	slideCount := fs.Int("slides", 5, "Number of initial slides")
	_ = fs.Parse(norm)

	name := "new_deck"
	if fs.NArg() > 0 {
		name = fs.Arg(0)
	}

	path := *targetPath
	if path == "" {
		cwd, _ := os.Getwd()
		path = filepath.Join(cwd, "decks", name)
	}

	deck, err := scaffold.ScaffoldDeck(path, name, *themeName, *slideCount)
	if err != nil {
		fmt.Printf("Scaffold error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Initialized presentation '%s' in %s\n", deck.Name, deck.Path)
	fmt.Printf("  -> Theme: %s\n", deck.ThemeName)
	fmt.Printf("  -> Generated: %d slides\n", len(deck.Slides))
	fmt.Println("Run 'deckforge build' or 'deckforge serve' to compile and view.")
}

func runBuild(args []string) {
	norm := normalizeArgs(args)
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	themeOverride := fs.String("theme", "", "Override presentation theme")
	_ = fs.Parse(norm)

	targetPath := "."
	if fs.NArg() > 0 {
		targetPath = fs.Arg(0)
	}

	tm := theme.NewThemeManager(targetPath)
	comp := compiler.NewCompiler(tm)

	outFile, err := comp.Build(targetPath, *themeOverride)
	if err != nil {
		fmt.Printf("Build error: %v\n", err)
		os.Exit(1)
	}

	fi, _ := os.Stat(outFile)
	sizeKB := float64(fi.Size()) / 1024.0
	fmt.Printf("Compilation Complete: %s\n", outFile)
	fmt.Printf("  -> Size: %.1f KB (Self-contained standalone 1080p HTML)\n", sizeKB)
}

func runServe(args []string) {
	norm := normalizeArgs(args)
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.Int("port", 0, "HTTP server port (defaults to configured port or 8080)")
	watch := fs.Bool("watch", true, "Auto-recompile on slide edits")
	_ = fs.Parse(norm)

	targetPath := "."
	if fs.NArg() > 0 {
		targetPath = fs.Arg(0)
	}

	cfg := models.LoadWorkspaceConfig()
	serverPort := *port
	if serverPort <= 0 {
		if cfg != nil && cfg.ServerPort > 0 {
			serverPort = cfg.ServerPort
		} else {
			serverPort = 8080
		}
	}

	tm := theme.NewThemeManager(targetPath)
	comp := compiler.NewCompiler(tm)

	srv := server.NewDeckServer(targetPath, serverPort, comp)
	if err := srv.ServeSync(*watch); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}

func runExport(args []string) {
	norm := normalizeArgs(args)
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	pdf := fs.Bool("pdf", true, "Export as single PDF")
	_ = fs.Parse(norm)

	targetPath := "."
	if fs.NArg() > 0 {
		targetPath = fs.Arg(0)
	}

	format := "pdf"
	if !*pdf {
		format = "png"
	}

	tm := theme.NewThemeManager(targetPath)
	comp := compiler.NewCompiler(tm)

	out, err := exporter.ExportDeck(targetPath, format, comp)
	if err != nil {
		fmt.Printf("Export error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Exported deck to: %s\n", out)
}

func runList(args []string) {
	norm := normalizeArgs(args)
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	path := fs.String("path", ".", "Root path to inspect")
	_ = fs.Parse(norm)

	tm := theme.NewThemeManager(*path)
	themes := tm.ListThemes(*path)

	fmt.Println("=== DeckForge Themes ===")
	for _, t := range themes {
		badge := ""
		if t.IsCustom {
			badge = " [CUSTOM]"
		}
		fmt.Printf("  * %-20s | %s%s\n", t.Tokens.Name, t.Tokens.DisplayName, badge)
	}

	fmt.Printf("\n=== Discovered Presentations in %s ===\n", *path)
	decks, err := workspace.ScanDecksInRoot(*path)
	if err != nil {
		fmt.Printf("Scan error: %v\n", err)
		return
	}

	if len(decks) == 0 {
		fmt.Println("  (No slide presentations found)")
		return
	}

	for _, d := range decks {
		fmt.Printf("  * %-22s | %02d slides | Theme: %-16s | %s\n",
			d.Name, len(d.Slides), d.ThemeName, d.Path)
	}
}

func runSlide(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: deckforge slide add <deck-path> <title> [--layout 2col|3col|hero|code]")
		return
	}

	action := args[0]
	deckPath := args[1]

	if action == "add" {
		title := "New Slide"
		if len(args) > 2 {
			title = args[2]
		}
		layout := "2col"
		if len(args) > 4 && args[3] == "--layout" {
			layout = args[4]
		}

		s, err := scaffold.AddSlide(deckPath, title, layout)
		if err != nil {
			fmt.Printf("Error adding slide: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created Slide %02d: %s\n", s.Index, s.Path)
	}
}

func runTheme(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: deckforge theme [list | clone <preset> <new-name> [target-dir]]")
		return
	}

	action := args[0]
	tm := theme.NewThemeManager(".")

	switch action {
	case "list":
		for _, p := range theme.GetBuiltinPresets() {
			fmt.Printf("  - %-18s (%s) - %s\n", p.Tokens.Name, p.Tokens.DisplayName, p.Tokens.Vibe)
		}
	case "set":
		if len(args) < 3 {
			fmt.Println("Usage: deckforge theme set <deck-path> <theme-name>")
			return
		}
		deckPath := args[1]
		themeName := args[2]
		if err := workspace.SetDeckTheme(deckPath, themeName); err != nil {
			fmt.Printf("Error setting theme: %v\n", err)
			os.Exit(1)
		}
		tm := theme.NewThemeManager(deckPath)
		comp := compiler.NewCompiler(tm)
		outFile, err := comp.Build(deckPath, themeName)
		if err != nil {
			fmt.Printf("Error recompiling deck: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Theme updated to '%s' in %s\n  -> Compiled: %s\n", themeName, deckPath, outFile)
	case "create", "new":
		if len(args) < 2 {
			fmt.Println("Usage: deckforge theme create <name> [--base <preset>] [--global] [--scss]")
			return
		}
		themeName := args[1]
		basePreset := ""
		global := false
		scss := false
		for i := 2; i < len(args); i++ {
			if (args[i] == "--base" || args[i] == "-b") && i+1 < len(args) {
				basePreset = args[i+1]
				i++
			} else if args[i] == "--global" || args[i] == "-g" {
				global = true
			} else if args[i] == "--scss" {
				scss = true
			}
		}
		created, err := tm.CreateTheme(themeName, basePreset, global, scss)
		if err != nil {
			fmt.Printf("Error creating theme: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created theme '%s' (scope: %s) in %s\n", created.Tokens.Name, created.Scope, created.Dir)
		fmt.Println("Segment files: tokens.json, typography, surfaces, backdrop, components")
	case "clone":
		if len(args) < 3 {
			fmt.Println("Usage: deckforge theme clone <preset-name> <new-name> [target-dir] [--global] [--scss]")
			return
		}
		preset := args[1]
		newName := args[2]
		targetDir := ""
		global := false
		scss := false
		for i := 3; i < len(args); i++ {
			if args[i] == "--global" || args[i] == "-g" {
				global = true
			} else if args[i] == "--scss" {
				scss = true
			} else if targetDir == "" && !strings.HasPrefix(args[i], "-") {
				targetDir = args[i]
			}
		}
		if targetDir == "" {
			if global {
				targetDir = filepath.Join(tm.GlobalDir, newName)
			} else {
				targetDir = "themes/" + newName
			}
		}
		cloned, err := tm.ClonePresetSCSS(preset, newName, targetDir, scss)
		if err != nil {
			fmt.Printf("Clone error: %v\n", err)
			os.Exit(1)
		}
		if global {
			cloned.Scope = "global"
		} else {
			cloned.Scope = "workspace"
		}
		fmt.Printf("Cloned theme '%s' (scope: %s) into %s\n", cloned.Tokens.Name, cloned.Scope, targetDir)
		fmt.Println("Segment files: tokens.json, typography, surfaces, backdrop, components")
	case "seed":
		if len(args) < 2 {
			fmt.Println("Usage: deckforge theme seed <preset-name> [--name <custom-name>] [--global] [--scss]")
			return
		}
		preset := args[1]
		customName := preset
		global := false
		scss := false
		for i := 2; i < len(args); i++ {
			if (args[i] == "--name" || args[i] == "-n") && i+1 < len(args) {
				customName = args[i+1]
				i++
			} else if args[i] == "--global" || args[i] == "-g" {
				global = true
			} else if args[i] == "--scss" {
				scss = true
			}
		}
		seeded, err := tm.SeedThemeAs(preset, customName, global, scss)
		if err != nil {
			fmt.Printf("Seed error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Seeded theme '%s' (scope: %s) into %s\n", seeded.Tokens.Name, seeded.Scope, seeded.Dir)
	}
}

func printHelp() {
	fmt.Println(`DeckForge - 100% Pure Go Modular 1080p Presentation Platform

USAGE:
  deckforge                          Launch interactive TUI
  deckforge [path]                   Launch interactive TUI in specific root directory
  deckforge tui [path]               Explicit TUI launch
  deckforge new <name> [flags]       Initialize new presentation with N slides (alias: init, create)
  deckforge build [path] [flags]     Compile presentation to standalone HTML
  deckforge serve [path] [flags]     Start local web server with live reload watcher
  deckforge export [path] [flags]    Headless export to PDF (via Chrome/Chromium)
  deckforge slide add <path> <title> Append a slide with selected layout
  deckforge theme [list|set|create|seed|clone] Theme management, assignment, and scaffolding
  deckforge book [status|log|audit]  Manage living Development Book
  deckforge list [--path <root>]     List themes and discovered presentations

FLAGS:
  --path <dir>       Target destination directory
  --theme <name>     Select presentation theme (14 presets available)
  --slides <count>   Set number of slides for scaffolding (default: 5)
  --port <port>      HTTP port for serve (default: 8080)
  --watch            Enable live file watcher (default: true)
  --pdf              Export as PDF format`)
}
