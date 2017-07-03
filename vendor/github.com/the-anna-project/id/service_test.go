package id

import (
	"io"
	"math/big"
	"sync"
	"testing"

	"github.com/the-anna-project/random"
)

func Test_Service_WithType_Error(t *testing.T) {
	var err error

	var randomService random.Service
	{
		randomConfig := random.DefaultServiceConfig()
		randomConfig.RandFactory = func(randReader io.Reader, max *big.Int) (n *big.Int, err error) {
			return nil, maskAny(invalidConfigError)
		}
		randomService, err = random.NewService(randomConfig)
		if err != nil {
			panic(err)
		}
	}

	var idService Service
	{
		idConfig := DefaultServiceConfig()
		idConfig.RandomService = randomService
		idService, err = NewService(idConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	_, err = idService.WithType(Hex128)
	if !IsInvalidConfig(err) {
		t.Fatal("expected", nil, "got", err)
	}
}

// Test_Service_New checks that a generated ID is still unique after a certain
// number of concurrent generations.
func Test_Service_New(t *testing.T) {
	var err error

	var idService Service
	{
		idConfig := DefaultServiceConfig()
		idService, err = NewService(idConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	alreadySeen := map[string]struct{}{}

	var mutex sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			newObjectID, err := idService.New()
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}

			mutex.Lock()
			defer mutex.Unlock()
			if _, ok := alreadySeen[newObjectID]; ok {
				t.Fatal("idService.New returned the same ID twice")
			}
			alreadySeen[newObjectID] = struct{}{}
		}()
	}
	wg.Wait()
}

// Test_Service_WithType checks that a generated ID is still unique after a
// certain number of concurrent generations.
func Test_Service_WithType(t *testing.T) {
	var err error

	var idService Service
	{
		idConfig := DefaultServiceConfig()
		idService, err = NewService(idConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	alreadySeen := map[string]struct{}{}

	var mutex sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			newObjectID, err := idService.WithType(Hex128)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}

			mutex.Lock()
			defer mutex.Unlock()
			if _, ok := alreadySeen[newObjectID]; ok {
				t.Fatal("idService.New returned the same ID twice")
			}
			alreadySeen[newObjectID] = struct{}{}
		}()
	}
	wg.Wait()
}
