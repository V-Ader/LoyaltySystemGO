package common

type Entity interface {
	GetHash() string
}

type RequestError struct {
	StatusCode int
	Err        error
}
