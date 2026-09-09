package theme

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

var (
	scssCache      sync.Map
	sassLookupOnce sync.Once
	cachedSassBin  string
	cachedSassType string // "native", "npx", or "none"
)

func detectSassCompiler() {
	if bin, err := exec.LookPath("sass"); err == nil {
		cachedSassBin = bin
		cachedSassType = "native"
		return
	}
	if bin, err := exec.LookPath("dart-sass"); err == nil {
		cachedSassBin = bin
		cachedSassType = "native"
		return
	}
	if bin, err := exec.LookPath("npx"); err == nil {
		cachedSassBin = bin
		cachedSassType = "npx"
		return
	}
	cachedSassType = "none"
}

// CompileSCSS attempts to transpile SCSS content to standard CSS using local sass tools.
// If no sass compiler is available on the system, it passes the raw content through,
// allowing modern browser CSS native nesting and features to execute directly.
// Results are cached in memory and on disk (~/.cache/deckforge/scss/) by SHA-256 hash.
func CompileSCSS(content string) (string, error) {
	if content == "" {
		return "", nil
	}

	hash := sha256.Sum256([]byte(content))
	key := hex.EncodeToString(hash[:])

	// 1. In-memory cache hit
	if val, ok := scssCache.Load(key); ok {
		return val.(string), nil
	}

	// 2. On-disk cache hit
	var diskCacheFile string
	if cacheDir, err := os.UserCacheDir(); err == nil {
		diskDir := filepath.Join(cacheDir, "deckforge", "scss")
		diskCacheFile = filepath.Join(diskDir, key+".css")
		if data, err := os.ReadFile(diskCacheFile); err == nil {
			compiled := string(data)
			scssCache.Store(key, compiled)
			return compiled, nil
		}
	}

	sassLookupOnce.Do(detectSassCompiler)

	var compiled string
	var compileErr error

	switch cachedSassType {
	case "native":
		cmd := exec.Command(cachedSassBin, "--stdin", "--no-source-map", "--style=expanded")
		cmd.Stdin = bytes.NewBufferString(content)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err == nil {
			compiled = out.String()
		} else {
			compileErr = fmt.Errorf("native sass compilation error: %s (%w)", stderr.String(), err)
		}
	case "npx":
		cmd := exec.Command(cachedSassBin, "-y", "sass", "--stdin", "--no-source-map", "--style=expanded")
		cmd.Stdin = bytes.NewBufferString(content)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err == nil {
			compiled = out.String()
		} else {
			compileErr = fmt.Errorf("npx sass compilation error: %s (%w)", stderr.String(), err)
		}
	default:
		compiled = content
		compileErr = fmt.Errorf("no sass compiler found; using native modern CSS")
	}

	if compiled != "" && compileErr == nil {
		scssCache.Store(key, compiled)
		if diskCacheFile != "" {
			_ = os.MkdirAll(filepath.Dir(diskCacheFile), 0755)
			_ = os.WriteFile(diskCacheFile, []byte(compiled), 0644)
		}
		return compiled, nil
	}

	if compiled != "" {
		return compiled, compileErr
	}

	return content, compileErr
}

// ClearSCSSCache flushes the in-memory SCSS compilation cache
func ClearSCSSCache() {
	scssCache.Range(func(key, _ interface{}) bool {
		scssCache.Delete(key)
		return true
	})
}
