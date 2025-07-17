package secrets

import (
	"crypto/rand"
	"errors"
	"runtime"
	"sync"
)

type SecureString struct {
	data     []byte
	consumed bool
	mu       sync.Mutex
}

// NewSecureString creates a new SecureString from a regular string
func NewSecureString(s string) *SecureString {
	data := make([]byte, len(s))
	copy(data, []byte(s))

	ss := &SecureString{
		data: data,
	}

	// Set finalizer to ensure cleanup if not explicitly called
	runtime.SetFinalizer(ss, (*SecureString).destroy)
	return ss
}

// Read returns the secret value and zeros it out (one-time use)
func (ss *SecureString) Read() (string, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if ss.consumed {
		return "", errors.New("secret already consumed")
	}

	value := string(ss.data)
	ss.zeroMemory()
	ss.consumed = true

	return value, nil
}

func (ss *SecureString) MustRead() string {
	value, err := ss.Read()
	if err != nil {
		panic(err)
	}
	return value
}

// Peek returns the value without consuming it (for cases where you need multiple reads)
func (ss *SecureString) Peek() (string, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if ss.consumed {
		return "", errors.New("secret already consumed")
	}

	return string(ss.data), nil
}

// Destroy explicitly zeros out memory
func (ss *SecureString) Destroy() {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.destroy()
}

func (ss *SecureString) destroy() {
	if !ss.consumed {
		ss.zeroMemory()
		ss.consumed = true
	}
	runtime.SetFinalizer(ss, nil)
}

func (ss *SecureString) zeroMemory() {
	if len(ss.data) > 0 {
		_, _ = rand.Read(ss.data)
		for i := range ss.data {
			ss.data[i] = 0
		}
	}
}
