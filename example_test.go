package sisyphus_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/soroushj/sisyphus"
)

func Example() {
	f := func() error {
		fmt.Println("trying")
		if rand.Intn(5) != 0 {
			return errors.New("failure")
		}
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := sisyphus.Do(ctx, f)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("success")
}

func Example_shouldRetry() {
	unrecoverable := errors.New("unrecoverable failure")
	f := func() error {
		fmt.Println("trying")
		if r := rand.Intn(5); r == 0 {
			return unrecoverable
		} else if r != 1 {
			return errors.New("failure")
		}
		return nil
	}
	shouldRetry := func(err error) bool {
		return err != unrecoverable
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := sisyphus.DoIf(ctx, f, shouldRetry)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("success")
}
