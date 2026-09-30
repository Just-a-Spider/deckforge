package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"deckforge/internal/compiler"
	"deckforge/internal/components"
	"deckforge/internal/editor"
	"deckforge/internal/exporter"
	"deckforge/internal/mcp"
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

	case "config", "settings":
		runConfig(os.Args[2:])

	case "component", "components":
		runComponent(os.Args[2:])

	case "book":
		runBook(os.Args[2:])

	case "agent":
		runAgent(os.Args[2:])

	case "mcp":
		deckPath := "."
		if len(os.Args) > 2 {
			deckPath = os.Args[2]
		}
		if err := mcp.RunMCPServer(deckPath); err != nil {
			fmt.Printf("MCP Server error: %v\n", err)
			os.Exit(1)
		}

	case "version", "-v", "--version":
		runVersion()

	case "completion":
		runCompletion(os.Args[2:])

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
	cfg := models.LoadWorkspaceConfig()
	defaultTheme := "academic-crimson"
	defaultSlides := 5
	if cfg != nil {
		if cfg.DefaultTheme != "" {
			defaultTheme = cfg.DefaultTheme
		}
		if cfg.DefaultSlideCount > 0 {
			defaultSlides = cfg.DefaultSlideCount
		}
	}

	norm := normalizeArgs(args)
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	targetPath := fs.String("path", "", "Destination directory path")
	themeName := fs.String("theme", defaultTheme, "Base theme preset")
	slideCount := fs.Int("slides", defaultSlides, "Number of initial slides")
	_ = fs.Parse(norm)

	name := "new_deck"
	if fs.NArg() > 0 {
		name = fs.Arg(0)
	}

	path := *targetPath
	if path == "" {
		cwd, _ := os.Getwd()
		dotDir := models.WorkspaceDotDir(cwd)
		if fi, err := os.Stat(dotDir); err == nil && fi.IsDir() {
			path = filepath.Join(dotDir, "decks", name)
		} else {
			path = filepath.Join(cwd, "decks", name)
		}
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

	wsRoot := workspace.DetectWorkspaceRoot(targetPath)
	tm := theme.NewThemeManager(wsRoot)
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

	wsRoot := workspace.DetectWorkspaceRoot(targetPath)
	tm := theme.NewThemeManager(wsRoot)
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
	engine := fs.String("engine", "chromedp", "Export engine: 'chromedp' (CDP library, default), 'node' (Node CDP runner), 'cli' (headless Chrome CLI)")
	timeout := fs.Duration("timeout", 30*time.Second, "Timeout for export rendering and font loading")
	_ = fs.Parse(norm)

	targetPath := "."
	if fs.NArg() > 0 {
		targetPath = fs.Arg(0)
	}

	format := "pdf"
	if !*pdf {
		format = "png"
	}

	wsRoot := workspace.DetectWorkspaceRoot(targetPath)
	tm := theme.NewThemeManager(wsRoot)
	comp := compiler.NewCompiler(tm)

	out, err := exporter.ExportDeckWithOptions(targetPath, format, *engine, *timeout, comp)
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
	showThemes := fs.Bool("themes", false, "Include global theme presets in listing")
	_ = fs.Parse(norm)

	cwd, _ := os.Getwd()
	absTarget, err := filepath.Abs(*path)
	if err != nil {
		absTarget = *path
	}

	displayTarget := *path
	if rel, err := filepath.Rel(cwd, absTarget); err == nil && !strings.HasPrefix(rel, "..") {
		displayTarget = rel
		if displayTarget == "" {
			displayTarget = "."
		}
	}

	if *showThemes {
		wsRoot := workspace.DetectWorkspaceRoot(*path)
		tm := theme.NewThemeManager(wsRoot)
		themes := tm.ListThemes(*path)
		fmt.Println("=== DeckForge Themes ===")
		for _, t := range themes {
			badge := ""
			if t.IsCustom {
				badge = " [CUSTOM]"
			}
			fmt.Printf("  * %-20s | %s%s\n", t.Tokens.Name, t.Tokens.DisplayName, badge)
		}
		fmt.Println()
	}

	decks, err := workspace.ScanDecksInRoot(absTarget)
	if err != nil {
		fmt.Printf("Scan error: %v\n", err)
		return
	}

	fmt.Printf("=== Presentations in %s ===\n", displayTarget)
	if len(decks) == 0 {
		fmt.Println("  (No slide presentations found)")
		fmt.Println("  -> Run 'deckforge new <name>' to create a presentation in this directory")
		return
	}

	for _, d := range decks {
		deckDisplayPath := d.Path
		if rel, err := filepath.Rel(cwd, d.Path); err == nil && !strings.HasPrefix(rel, "..") {
			deckDisplayPath = rel
		}
		fmt.Printf("  * %-22s | %02d slides | Theme: %-16s | %s\n",
			d.Name, len(d.Slides), d.ThemeName, deckDisplayPath)
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
	wsRoot := workspace.DetectWorkspaceRoot(".")
	tm := theme.NewThemeManager(wsRoot)

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
		deckWsRoot := workspace.DetectWorkspaceRoot(deckPath)
		tm := theme.NewThemeManager(deckWsRoot)
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

func runConfig(args []string) {
	if len(args) < 1 {
		showConfigList()
		return
	}

	action := args[0]
	cfg := models.LoadWorkspaceConfig()

	switch action {
	case "list", "show":
		showConfigList()

	case "path":
		for _, a := range args[1:] {
			if a == "--local" || a == "-l" {
				cwd, _ := os.Getwd()
				fmt.Println(models.WorkspaceConfigFilePath(cwd))
				return
			}
		}
		fmt.Println(models.ConfigFilePath())

	case "editors":
		showEditors()

	case "get":
		if len(args) < 2 {
			fmt.Println("Usage: deckforge config get <key> [--local]")
			return
		}
		key := args[1]
		local := false
		for _, a := range args[2:] {
			if a == "--local" || a == "-l" {
				local = true
			}
		}
		var activeCfg *models.WorkspaceConfig
		if local {
			cwd, _ := os.Getwd()
			activeCfg = models.LoadWorkspaceConfigForRoot(cwd)
		} else {
			activeCfg = cfg
		}
		val, err := activeCfg.Get(key)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(val)

	case "set":
		if len(args) < 3 {
			fmt.Println("Usage: deckforge config set <key> <value> [--local]")
			return
		}
		local := false
		var cleanArgs []string
		for _, a := range args[1:] {
			if a == "--local" || a == "-l" {
				local = true
			} else {
				cleanArgs = append(cleanArgs, a)
			}
		}
		if len(cleanArgs) < 2 {
			fmt.Println("Usage: deckforge config set <key> <value> [--local]")
			return
		}
		key := cleanArgs[0]
		val := cleanArgs[1]

		if local {
			cwd, _ := os.Getwd()
			localCfg := models.LoadWorkspaceConfigForRoot(cwd)
			if err := localCfg.SetLocal(key, val, cwd); err != nil {
				fmt.Printf("Error setting local config: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Configuration updated locally: %s = %s\nSaved to %s\n", key, val, models.WorkspaceConfigFilePath(cwd))
			return
		}

		if err := cfg.Set(key, val); err != nil {
			fmt.Printf("Error setting config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Configuration updated: %s = %s\nSaved to %s\n", key, val, models.ConfigFilePath())

	default:
		fmt.Printf("Unknown config action: %s\nUsage: deckforge config [list|get|set|editors|path] [--local]\n", action)
	}
}

func showConfigList() {
	cwd, _ := os.Getwd()
	cfg := models.LoadWorkspaceConfigForRoot(cwd)
	fmt.Printf("=== DeckForge Preferences ===\n")
	fmt.Printf("Global Config:       %s\n", models.ConfigFilePath())
	localPath := models.WorkspaceConfigFilePath(cwd)
	if _, err := os.Stat(localPath); err == nil {
		fmt.Printf("Local Config:        %s\n", localPath)
	}
	fmt.Printf("Active Root:         %s\n", cfg.ActiveRoot)
	fmt.Printf("Preferred Editor:    %s\n", cfg.PreferredEditor)
	fmt.Printf("Default Theme:       %s\n", cfg.DefaultTheme)
	fmt.Printf("Default Slide Count: %d\n", cfg.DefaultSlideCount)
	fmt.Printf("Server Port:         %d\n", cfg.ServerPort)
	fmt.Printf("Live Watcher:        %v\n", cfg.AutoWatch)
	if len(cfg.RecentRoots) > 0 {
		fmt.Printf("\nRecent Workspaces (%d):\n", len(cfg.RecentRoots))
		for _, r := range cfg.RecentRoots {
			fmt.Printf("  - %s\n", r)
		}
	}
}

func showEditors() {
	installed := editor.DetectInstalledEditors()
	installedMap := make(map[string]editor.InstalledEditor)
	for _, inst := range installed {
		installedMap[inst.ID] = inst
	}

	fmt.Println("=== Supported External Editors ===")
	for _, spec := range editor.KnownEditors() {
		status := "[not found]"
		pathStr := ""
		if inst, ok := installedMap[spec.ID]; ok {
			status = "[INSTALLED]"
			pathStr = " -> " + inst.BinaryPath
			if len(inst.WaitFlags) > 0 {
				pathStr += fmt.Sprintf(" (wait: %s)", strings.Join(inst.WaitFlags, " "))
			}
		}
		fmt.Printf("  %-12s %-16s %-12s%s\n", spec.ID, spec.DisplayName, status, pathStr)
	}
}

func runComponent(args []string) {
	if len(args) < 1 {
		printComponentHelp()
		return
	}

	action := args[0]
	deckPath := "."
	wsRoot := workspace.DetectWorkspaceRoot(deckPath)
	cm := components.NewComponentManager(wsRoot)

	switch action {
	case "list":
		if len(args) > 1 && !strings.HasPrefix(args[1], "-") {
			deckPath = args[1]
			wsRoot = workspace.DetectWorkspaceRoot(deckPath)
			cm = components.NewComponentManager(wsRoot)
		}
		comps := cm.ListComponents(deckPath)
		fmt.Printf("=== Available Components (%d registered) ===\n", len(comps))
		for _, c := range comps {
			scopeBadge := fmt.Sprintf("[%s]", strings.ToUpper(c.Scope))
			hasCSS := ""
			if c.Styles != "" {
				hasCSS = " (with custom CSS)"
			}
			fmt.Printf("  %-24s %-12s %-12s %s%s\n", c.Selector, c.Category, scopeBadge, c.Name, hasCSS)
		}

	case "inspect":
		if len(args) < 2 {
			fmt.Println("Usage: deckforge component inspect <selector> [deck-path]")
			return
		}
		selector := args[1]
		if len(args) > 2 {
			deckPath = args[2]
			wsRoot = workspace.DetectWorkspaceRoot(deckPath)
			cm = components.NewComponentManager(wsRoot)
		}
		comp, err := cm.GetComponent(selector, deckPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		data, _ := json.MarshalIndent(comp, "", "  ")
		fmt.Println(string(data))

	case "create", "new":
		if len(args) < 2 {
			fmt.Println("Usage: deckforge component create <name> [--category <cat>] [--global] [--styles]")
			return
		}
		name := args[1]
		category := "cards"
		global := false
		withCSS := false

		for i := 2; i < len(args); i++ {
			if (args[i] == "--category" || args[i] == "-c") && i+1 < len(args) {
				category = args[i+1]
				i++
			} else if args[i] == "--global" || args[i] == "-g" {
				global = true
			} else if args[i] == "--styles" || args[i] == "--css" {
				withCSS = true
			}
		}

		opts := components.ScaffoldOptions{
			Name:          name,
			Category:      category,
			Global:        global,
			WithCSS:       withCSS,
			WorkspaceRoot: ".",
		}

		created, err := components.ScaffoldComponent(opts)
		if err != nil {
			fmt.Printf("Error creating component: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created component '%s' (scope: %s)\n", created.Selector, created.Scope)
		fmt.Printf("  -> Location: %s\n", created.SourcePath)
		fmt.Println("Inspect with: deckforge component inspect " + created.Selector)

	default:
		fmt.Printf("Unknown component action: %s\n", action)
		printComponentHelp()
	}
}

func printComponentHelp() {
	fmt.Println(`DeckForge Component Management

USAGE:
  deckforge component list [deck-path]
  deckforge component create <name> [--category <cat>] [--global] [--styles]
  deckforge component inspect <selector> [deck-path]`)
}

// Version holds the current release version, overridable at build time via -ldflags
var Version = "0.3.0"

func runVersion() {
	fmt.Printf("DeckForge v%s (%s/%s, runtime: %s)\n", Version, runtime.GOOS, runtime.GOARCH, runtime.Version())
}

func runCompletion(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: deckforge completion <shell> [--install]")
		fmt.Println("Supported shells: fish")
		return
	}

	shell := args[0]
	install := false
	for _, a := range args[1:] {
		if a == "--install" || a == "-i" {
			install = true
		}
	}

	switch shell {
	case "fish":
		script := GenerateFishCompletion()
		if install {
			home, _ := os.UserHomeDir()
			targetDir := filepath.Join(home, ".config", "fish", "completions")
			_ = os.MkdirAll(targetDir, 0755)
			targetFile := filepath.Join(targetDir, "deckforge.fish")
			if err := os.WriteFile(targetFile, []byte(script), 0644); err != nil {
				fmt.Printf("Error installing fish completion: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Fish completion installed to %s\n", targetFile)
		} else {
			fmt.Print(script)
		}
	default:
		fmt.Printf("Unsupported shell '%s'. Supported shells: fish\n", shell)
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
  deckforge component [list|create|inspect] Multi-tier component library and scaffolding
  deckforge config [list|get|set|editors|path] Global and local (.deckforge/) preferences
  deckforge agent <subcommand>       AI agent interface for autonomous tools (JSON output)
  deckforge mcp [deck-path]          Model Context Protocol (MCP) JSON-RPC 2.0 stdio server
  deckforge book [status|log|audit]  Manage living Development Book
  deckforge list [--path <root>]     List themes and discovered presentations
  deckforge completion <shell>       Generate shell autocompletions (fish) [--install]
  deckforge version                  Display version, OS/architecture, and Go runtime

FLAGS:
  --path <dir>       Target destination directory
  --theme <name>     Select presentation theme (14 presets available)
  --slides <count>   Set number of slides for scaffolding (default: 5)
  --port <port>      HTTP port for serve (default: 8080)
  --watch            Enable live file watcher (default: true)
  --pdf              Export as PDF format`)
}
