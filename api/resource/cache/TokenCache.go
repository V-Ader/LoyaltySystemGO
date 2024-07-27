package cache

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set"
)

type TokenCache struct {
	cache mapset.Set
	mutex sync.Mutex
}

func NewTokenCache(defaultExpiration time.Duration, cleanupInterval time.Duration) *TokenCache {
	return &TokenCache{
		cache: mapset.NewSet(),
	}
}

func (tc *TokenCache) CreateToken() (string, error) {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	token, err := generateUniqueToken(tc)
	if err != nil {
		return "", err
	}

	tc.cache.Add(token)
	fmt.Printf("Token %s added\n", token)
	return token, nil
}

func generateUniqueToken(tc *TokenCache) (string, error) {
	var token string
	var err error
	for {
		token, err = generateRandomToken(16)
		if err != nil {
			return "", err
		}
		if !tc.cache.Contains(token) {
			break
		}
	}
	return token, nil
}

func (tc *TokenCache) RemoveToken(token string) error {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	if !tc.cache.Contains(token) {
		return errors.New("token not found in cache")
	}

	tc.cache.Remove(token)
	fmt.Printf("Token %s removed\n", token)

	return nil
}

func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
