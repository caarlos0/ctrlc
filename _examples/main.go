package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/ctrlc"
)

var timeout = flag.Duration("timeout", time.Minute, "timeout for the whole thing")

func main() {
	flag.Parse()
	log.Println("creating a new context with", *timeout, "timeout")
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	if err := ctrlc.Default.Run(ctx, func() error {
		log.Println("doing some heavy work for 5s")
		time.Sleep(5 * time.Second)
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	if err := ctrlc.Default.Run(ctx, func() error {
		log.Println("doing some heavy work for 1s and then failing")
		time.Sleep(time.Second)
		return fmt.Errorf("failed for some reason")
	}); err != nil {
		log.Fatal(err)
	}

}
