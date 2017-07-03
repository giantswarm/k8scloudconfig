package random

import (
	"crypto/rand"
	"io"
	"math/big"
	"testing"
	"time"
)

func Test_Service_CreateNMax_Error_RandFactory(t *testing.T) {
	var err error

	var randomService Service
	{
		randomConfig := DefaultServiceConfig()
		randomConfig.RandFactory = func(randReader io.Reader, max *big.Int) (n *big.Int, err error) {
			return nil, maskAny(timeoutError)
		}
		randomConfig.Timeout = 10 * time.Millisecond
		randomService, err = NewService(randomConfig)
		if err != nil {
			panic(err)
		}
	}

	n := 5
	max := 10

	_, err = randomService.CreateNMax(n, max)
	if !IsTimeout(err) {
		t.Fatal("expected", nil, "got", err)
	}
}

func Test_Service_CreateNMax_Error_Timeout(t *testing.T) {
	var err error

	var randomService Service
	{
		randomConfig := DefaultServiceConfig()
		randomConfig.RandFactory = func(randReader io.Reader, max *big.Int) (n *big.Int, err error) {
			time.Sleep(200 * time.Millisecond)
			return rand.Int(randReader, max)
		}
		randomConfig.Timeout = 20 * time.Millisecond
		randomService, err = NewService(randomConfig)
		if err != nil {
			panic(err)
		}
	}

	n := 5
	max := 10

	_, err = randomService.CreateNMax(n, max)
	if !IsTimeout(err) {
		t.Fatal("expected", nil, "got", err)
	}
}

func Test_Service_CreateNMax_GenerateNNumbers(t *testing.T) {
	var err error

	var randomService Service
	{
		randomConfig := DefaultServiceConfig()
		randomService, err = NewService(randomConfig)
		if err != nil {
			panic(err)
		}
	}

	n := 5
	max := 10

	newRandomNumbers, err := randomService.CreateNMax(n, max)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	if len(newRandomNumbers) != 5 {
		t.Fatal("expected", 5, "got", len(newRandomNumbers))
	}
}

func Test_Service_CreateNMax_GenerateRandomNumbers(t *testing.T) {
	var err error

	var randomService Service
	{
		randomConfig := DefaultServiceConfig()
		randomService, err = NewService(randomConfig)
		if err != nil {
			panic(err)
		}
	}

	n := 100
	max := 10

	newRandomNumbers, err := randomService.CreateNMax(n, max)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	alreadySeen := map[int]struct{}{}

	for _, r := range newRandomNumbers {
		alreadySeen[r] = struct{}{}
	}

	l := len(alreadySeen)
	if l != 10 {
		t.Fatal("expected", 10, "got", l)
	}
}
