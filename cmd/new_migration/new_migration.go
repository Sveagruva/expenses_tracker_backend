package main

import (
	"os"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) < 2 {
		panic("provide name for migration")
	}
	name := os.Args[1]

	createCmd := exec.Command("migrate", "create", "-ext", "sql", "-dir", "migrations", "-seq", name)
	createCmd.Stdout = os.Stdout
	createCmd.Stderr = os.Stderr
	createCmd.Run()
}
