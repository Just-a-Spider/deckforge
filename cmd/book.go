package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func runBook(args []string) {
	if len(args) < 1 {
		printBookHelp()
		return
	}

	sub := args[0]
	switch sub {
	case "status":
		showBookStatus()
	case "log":
		if len(args) < 2 {
			fmt.Println("Usage: deckforge book log \"<achievement description>\"")
			return
		}
		logAchievement(strings.Join(args[1:], " "))
	case "audit":
		auditBook()
	default:
		fmt.Printf("Unknown book command: %s\n", sub)
		printBookHelp()
	}
}

func printBookHelp() {
	fmt.Println(`DeckForge Development Book CLI

USAGE:
  deckforge book status             Show milestone progress & active goals
  deckforge book log "<note>"       Record verified achievement in archives
  deckforge book audit              Validate integrity of book documents & links`)
}

func findBookDir() string {
	candidates := []string{"book", "../book", "../../book"}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && fi.IsDir() {
			abs, _ := filepath.Abs(c)
			return abs
		}
	}
	return "book"
}

func showBookStatus() {
	bookDir := findBookDir()
	trackerPath := filepath.Join(bookDir, "02_roadmap", "02_milestone_tracker.md")
	data, err := os.ReadFile(trackerPath)
	if err != nil {
		fmt.Printf("Could not read milestone tracker at %s: %v\n", trackerPath, err)
		return
	}

	lines := strings.Split(string(data), "\n")
	var total, done int
	var currentPhase string

	fmt.Println("=== DeckForge Development Book: Status ===")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## Phase") {
			currentPhase = strings.TrimPrefix(trimmed, "## ")
		}
		if strings.HasPrefix(trimmed, "- [x]") {
			total++
			done++
		} else if strings.HasPrefix(trimmed, "- [ ]") {
			total++
		}
	}

	pct := 0
	if total > 0 {
		pct = (done * 100) / total
	}

	fmt.Printf("Current Phase: %s\n", currentPhase)
	fmt.Printf("Tasks Completed: %d / %d (%d%%)\n", done, total, pct)

	barWidth := 30
	filled := (pct * barWidth) / 100
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	fmt.Printf("Progress: [%s] %d%%\n\n", bar, pct)

	fmt.Println("Recent Active Milestones:")
	phaseCount := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## Phase") {
			phaseCount++
			if phaseCount <= 3 {
				fmt.Printf("\n%s\n", trimmed)
			}
		} else if (strings.HasPrefix(trimmed, "- [ ]") || strings.HasPrefix(trimmed, "- [x]")) && phaseCount <= 3 {
			fmt.Printf("  %s\n", trimmed)
		}
	}
}

func logAchievement(desc string) {
	bookDir := findBookDir()
	logPath := filepath.Join(bookDir, "06_archives", "01_achievements_log.md")

	gitHash := "n/a"
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	if out, err := cmd.Output(); err == nil {
		gitHash = strings.TrimSpace(string(out))
	}

	today := time.Now().Format("2006-01-02")
	author := os.Getenv("USER")
	if author == "" {
		author = "dev"
	}

	entry := fmt.Sprintf("| %s | %s | `%s` | %s | %s |\n", today, "Progress", gitHash, desc, author)

	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Failed to open achievements log: %v\n", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(entry); err != nil {
		fmt.Printf("Failed to write entry: %v\n", err)
		return
	}

	fmt.Printf("Logged achievement in %s:\n  -> %s\n", logPath, strings.TrimSpace(entry))
}

func auditBook() {
	bookDir := findBookDir()
	fmt.Printf("Auditing Development Book at: %s\n", bookDir)

	requiredDirs := []string{
		"01_foundations",
		"02_roadmap",
		"03_adrs",
		"04_specifications",
		"05_agent_playbook",
		"06_archives",
	}

	allPassed := true
	for _, d := range requiredDirs {
		dirPath := filepath.Join(bookDir, d)
		fi, err := os.Stat(dirPath)
		if err != nil || !fi.IsDir() {
			fmt.Printf("  [FAIL] Missing directory: %s\n", d)
			allPassed = false
		} else {
			fmt.Printf("  [PASS] Directory exists: %s\n", d)
		}
	}

	summaryPath := filepath.Join(bookDir, "SUMMARY.md")
	data, err := os.ReadFile(summaryPath)
	if err != nil {
		fmt.Printf("  [FAIL] Missing SUMMARY.md\n")
		return
	}

	linkRegex := regexp.MustCompile(`\[.*?\]\((file://.*?\.md)\)`)
	matches := linkRegex.FindAllStringSubmatch(string(data), -1)

	checkedFiles := 0
	for _, m := range matches {
		if len(m) > 1 {
			target := strings.TrimPrefix(m[1], "file://")
			if _, err := os.Stat(target); err != nil {
				fmt.Printf("  [FAIL] Broken document link: %s\n", target)
				allPassed = false
			} else {
				checkedFiles++
			}
		}
	}

	fmt.Printf("  [PASS] Verified %d linked documents\n", checkedFiles)
	if allPassed {
		fmt.Println("\nDevBook Audit: 100% HEALTHY")
	} else {
		fmt.Println("\nDevBook Audit: ISSUES FOUND")
	}
}
