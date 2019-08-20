package tests

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"

	"github.com/gebv/mgraph"
)

func TestMain(m *testing.M) {
	log.SetPrefix("testmain: ")
	log.SetFlags(0)

	flag.Parse()

	var cancel context.CancelFunc
	Ctx, cancel = context.WithCancel(context.Background())

	DB = mgraph.NewPostgtresConnect("postgres://app:app@127.0.0.1:5432/app?sslmode=disable")
	MGP = mgraph.NewMGraphPostgres(DB)

	var exitCode int
	defer func() {
		if p := recover(); p != nil {
			panic(p)
		}
		log.Printf("stoped with exit code %d\n", exitCode)
		os.Exit(exitCode)
	}()

	exitCode = m.Run()
	cancel()
	log.Println("bye bye.")
}
