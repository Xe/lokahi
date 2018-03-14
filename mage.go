// +build mage

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// the version of Go to use in docker images.
	goVersion = "1.9.4"
)

// Generate runs all relevant code generation tasks.
func Generate() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shouldWork(ctx, nil, wd, "statik", "-src", "./public", "-f")
	shouldWork(ctx, nil, filepath.Join(wd, "rpc", "lokahi"), "sh", "./regen.sh")
	shouldWork(ctx, nil, filepath.Join(wd, "rpc", "lokahiadmin"), "sh", "./regen.sh")
	shouldWork(ctx, nil, filepath.Join(wd, "internal", "database", "migrations"), "go-bindata", "-pkg=dmigrations", "-o=../dmigrations/bindata.go", ".")

	fmt.Println("reran code generation")
}

// Travis runs initial setup needed for travis, then a full build and test cycle.
func Travis() {
	os.Setenv("DATABASE_URL", "postgres://postgres:hunter2@127.0.0.1/postgres?sslmode=disable")
	os.Setenv("PATH", os.Getenv("PATH")+":"+os.Getenv("HOME")+"/.local/bin")

	fmt.Println("[-] building lokahi...")
	Generate()
	Build()

	fmt.Println("[-] testing lokahi...")
	Test()
}

// Test runs lokahi's test suite.
func Test() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shouldWork(ctx, nil, wd, "go", "test", "-v", "./...")
}

// Build builds the command code into binaries in ./bin.
func Build() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	os.Mkdir("bin", 0777)

	outd := filepath.Join(wd, "bin")
	cmds := []string{
		"lokahid",
		"lokahictl",
		"sample_hook",
		"duke-of-york",
		"webhookworker",
		"healthworker",
		"slack_hook",
		"discord_hook",
	}

	for _, c := range cmds {
		shouldWork(ctx, nil, outd, "go", "build", "../cmd/"+c)
		fmt.Println("built ./bin/" + c)
	}
}

// Dep reruns `dep`.
func Dep() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shouldWork(ctx, nil, wd, "dep", "ensure", "-update")
	shouldWork(ctx, nil, wd, "dep", "prune")
}

// Docker creates the docker image xena/lokahi with the lokahi server.
func Docker() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shouldWork(ctx, nil, wd, "docker", "pull", "xena/alpine")
	shouldWork(ctx, nil, wd, "docker", "pull", "xena/go:"+goVersion)
	shouldWork(ctx, nil, wd, "docker", "build", "-t", "xena/lokahi", ".")
}

// Run starts an instance of lokahid with default configuration and no
// authentication for debugging.
func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("building docker images")
	Docker()

	fmt.Println("Starting docker compose")

	defer shouldWork(ctx, nil, wd, "docker-compose", "down")
	shouldWork(ctx, nil, wd, "docker-compose", "up")
}
