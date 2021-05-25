package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/oklog/run"

	"github.com/LasTshaMAN/txstore"
	"github.com/LasTshaMAN/txstore/internal/inmemory"
)

func main() {
	var g run.Group
	{
		ctx, cancel := context.WithCancel(context.Background())

		g.Add(func() error {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
			select {
			case <-ctx.Done():
				return nil
			case s := <-sig:
				fmt.Printf("termination signal received: %v\n", s)
				return nil
			}
		}, func(_ error) {
			cancel()
		})
	}
	{
		ctx, cancel := context.WithCancel(context.Background())

		g.Add(func() error {
			reader := bufio.NewReader(os.Stdin)
			store := inmemory.NewStore()
			for {
				select {
				case <-ctx.Done():
					return nil
				default:
				}

				handleConsoleInput(reader, store)
			}
		}, func(_ error) {
			cancel()
		})
	}

	err := g.Run()
	if err != nil {
		fmt.Printf("finished with error: %v\n", err)
	}
}

func handleConsoleInput(reader *bufio.Reader, store txstore.Store) {
	fmt.Print(": ")
	input, _ := reader.ReadString('\n')
	cmd := strings.Fields(input)
	if len(cmd) == 0 {
		return
	}
	switch cmd[0] {
	case "SET":
		store.Set(cmd[1], cmd[2])
	case "GET":
		value := store.Get(cmd[1])
		fmt.Println(value)
	case "DELETE":
		store.Delete(cmd[1])
	case "COUNT":
		value := store.Count(cmd[1])
		fmt.Println(value)
	case "BEGIN":
		store.Begin()
	case "COMMIT":
		store.Commit()
	case "ROLLBACK":
		store.Rollback()
	default:
		fmt.Println("Unknown command", cmd[0])
	}
}
