package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

type runner struct {
	name    string
	detect  func(dir string) bool
	command string
	args    []string
	// passSep: insert "--" before user-passed flags so the underlying tool
	// forwards them to the script/binary instead of consuming them itself.
	passSep bool
}

func hasFile(dir, name string) bool {
	_, err := os.Stat(filepath.Join(dir, name))
	return err == nil
}

func hasGlob(dir, pattern string) bool {
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	return err == nil && len(matches) > 0
}

func makeHasDevTarget(dir string) bool {
	f, err := os.Open(filepath.Join(dir, "Makefile"))
	if err != nil {
		return false
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "dev:") || strings.HasPrefix(line, "dev :") {
			return true
		}
	}
	return false
}

func hasLuaEntry(dir string) (string, bool) {
	for _, name := range []string{"main.lua", "init.lua"} {
		if hasFile(dir, name) {
			return name, true
		}
	}
	matches, err := filepath.Glob(filepath.Join(dir, "*.lua"))
	if err == nil && len(matches) > 0 {
		return filepath.Base(matches[0]), true
	}
	return "", false
}

func detectRunner(dir string) (*runner, error) {
	if hasFile(dir, "bun.lockb") || hasFile(dir, "bun.lock") {
		return &runner{name: "Bun", command: "bun", args: []string{"run", "dev"}, passSep: true}, nil
	}
	if hasFile(dir, "pnpm-lock.yaml") {
		return &runner{name: "pnpm", command: "pnpm", args: []string{"dev"}, passSep: true}, nil
	}
	if hasFile(dir, "yarn.lock") {
		return &runner{name: "Yarn", command: "yarn", args: []string{"dev"}, passSep: true}, nil
	}
	if hasFile(dir, "deno.json") || hasFile(dir, "deno.jsonc") || hasFile(dir, "deno.lock") {
		return &runner{name: "Deno", command: "deno", args: []string{"task", "dev"}}, nil
	}
	if hasFile(dir, "package.json") {
		return &runner{name: "NodeJS", command: "npm", args: []string{"run", "dev"}, passSep: true}, nil
	}
	if hasFile(dir, "go.mod") {
		return &runner{name: "Go", command: "go", args: []string{"run", "."}}, nil
	}
	if hasFile(dir, "Cargo.toml") {
		return &runner{name: "Rust", command: "cargo", args: []string{"run"}, passSep: true}, nil
	}
	if hasFile(dir, "uv.lock") {
		return &runner{name: "Python (uv)", command: "uv", args: []string{"run", "python", "main.py"}}, nil
	}
	if hasFile(dir, "poetry.lock") {
		return &runner{name: "Python (poetry)", command: "poetry", args: []string{"run", "python", "main.py"}}, nil
	}
	if hasFile(dir, "manage.py") {
		return &runner{name: "Django", command: "python", args: []string{"manage.py", "runserver"}}, nil
	}
	if hasGlob(dir, "*.csproj") || hasGlob(dir, "*.sln") || hasGlob(dir, "*.fsproj") {
		return &runner{name: "Dotnet", command: "dotnet", args: []string{"run"}, passSep: true}, nil
	}
	if entry, ok := hasLuaEntry(dir); ok {
		return &runner{name: "Lua", command: "lua", args: []string{entry}}, nil
	}
	if makeHasDevTarget(dir) {
		return &runner{name: "Make", command: "make", args: []string{"dev"}}, nil
	}
	return nil, fmt.Errorf("no supported project detected in %s", dir)
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "dcrun: cannot get working directory:", err)
		os.Exit(1)
	}

	r, err := detectRunner(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dcrun:", err)
		os.Exit(1)
	}

	if _, err := exec.LookPath(r.command); err != nil {
		fmt.Fprintf(os.Stderr, "dcrun: %s detected but %q not found in PATH\n", r.name, r.command)
		os.Exit(1)
	}

	userArgs := os.Args[1:]
	if len(userArgs) > 0 && userArgs[0] == "--" {
		userArgs = userArgs[1:]
	}

	finalArgs := append([]string{}, r.args...)
	if len(userArgs) > 0 {
		if r.passSep {
			finalArgs = append(finalArgs, "--")
		}
		finalArgs = append(finalArgs, userArgs...)
	}

	fmt.Printf("dcrun: %s detected, running: %s %v\n", r.name, r.command, finalArgs)

	cmd := exec.Command(r.command, finalArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "dcrun: failed to start:", err)
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for s := range sigCh {
			if cmd.Process != nil {
				_ = cmd.Process.Signal(s)
			}
		}
	}()

	err = cmd.Wait()
	signal.Stop(sigCh)
	close(sigCh)

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintln(os.Stderr, "dcrun:", err)
		os.Exit(1)
	}
}
