package main

import (
    "log"
    "os"
    "os/exec"
)

func main() {
    migration := exec.Command("/app/career-migrate")
    migration.Stdout = os.Stdout
    migration.Stderr = os.Stderr
    migration.Env = os.Environ()
    if err := migration.Run(); err != nil {
        log.Fatalf("database migration failed: %v", err)
    }

    app := exec.Command("/app/careeros")
    app.Stdout = os.Stdout
    app.Stderr = os.Stderr
    app.Stdin = os.Stdin
    app.Env = os.Environ()
    if err := app.Run(); err != nil {
        log.Fatalf("CareerOS failed: %v", err)
    }
}
