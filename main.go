package main

// Code from https://www.youtube.com/watch?v=kaZOXRqFPCw

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	start := time.Now()
	ctx := context.Background()

	val, err := fetchUserData(ctx, 10)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("result", val)
	fmt.Println("took: ", time.Since(start))
}

type Response struct {
	value int
	err   error
}

func fetchUserData(ctx context.Context, userID int) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	respch := make(chan Response)

	go func() {
		val, err := fetchDataFromThirdPartyThatCanBeSlow()
		respch <- Response{
			val,
			err,
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return 0, fmt.Errorf("fetching took too long")
		case resp := <-respch:
			return resp.value, resp.err
		}
	}
}

func fetchDataFromThirdPartyThatCanBeSlow() (int, error) {
	time.Sleep(1500 * time.Millisecond)

	return 15, nil
}
