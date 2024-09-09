package cache

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/V-Ader/Loyality_GO/api/resource/common"
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

func (tc *TokenCache) RemoveToken(token string) *common.RequestError {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	if !tc.cache.Contains(token) {
		return &common.RequestError{StatusCode: http.StatusNotFound, Err: errors.New("token not found")}
	}
	tc.cache.Remove(token)
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
