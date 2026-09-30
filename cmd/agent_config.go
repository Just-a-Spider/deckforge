package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"deckforge/internal/models"
)

func runAgentConfigGet(args []string) {
	if len(args) < 1 {
		fmt.Println(`{"error": "usage: deckforge agent config-get <key> [deck-or-workspace-path]"}`)
		os.Exit(1)
	}
	key := args[0]
	root := ""
	if len(args) > 1 {
		root = args[1]
	}
	cfg := models.LoadWorkspaceConfigForRoot(root)
	val, err := cfg.Get(key)
	if err != nil {
		fmt.Printf(`{"error": "%s"}`+"\n", err)
		os.Exit(1)
	}
	path := models.ConfigFilePath()
	if root != "" {
		path = models.WorkspaceConfigFilePath(root)
	}
	res := map[string]string{
		"key":   key,
		"value": val,
		"path":  path,
	}
	_ = json.NewEncoder(os.Stdout).Encode(res)
}

func runAgentConfigSet(args []string) {
	if len(args) < 2 {
		fmt.Println(`{"error": "usage: deckforge agent config-set <key> <value> [--local]"}`)
		os.Exit(1)
	}
	key := args[0]
	val := args[1]
	local := false
	for i := 2; i < len(args); i++ {
		if args[i] == "--local" || args[i] == "-l" {
			local = true
		}
	}

	if local {
		cwd, _ := os.Getwd()
		cfg := models.LoadWorkspaceConfigForRoot(cwd)
		if err := cfg.SetLocal(key, val, cwd); err != nil {
			fmt.Printf(`{"error": "%s"}`+"\n", err)
			os.Exit(1)
		}
		res := map[string]string{
			"status": "ok",
			"key":    key,
			"value":  val,
			"scope":  "local",
			"path":   models.WorkspaceConfigFilePath(cwd),
		}
		_ = json.NewEncoder(os.Stdout).Encode(res)
		return
	}

	cfg := models.LoadWorkspaceConfig()
	if err := cfg.Set(key, val); err != nil {
		fmt.Printf(`{"error": "%s"}`+"\n", err)
		os.Exit(1)
	}
	res := map[string]string{
		"status": "ok",
		"key":    key,
		"value":  val,
		"scope":  "global",
		"path":   models.ConfigFilePath(),
	}
	_ = json.NewEncoder(os.Stdout).Encode(res)
}
