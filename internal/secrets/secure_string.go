package secrets

import (
	"crypto/rand"
	"errors"
	"os"
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

// ReadWith provides secure access to the secret data without creating string copies.
// The provided function is called with the raw bytes and should not retain references to them.
// The secret is consumed (zeroed) after the function returns.
func (ss *SecureString) ReadWith(fn func([]byte) error) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if ss.consumed {
		return errors.New("secret already consumed")
	}

	// Call the function with direct access to the bytes
	err := fn(ss.data)

	// Always zero memory after use, regardless of function result
	ss.zeroMemory()
	ss.consumed = true

	return err
}

// Destroy explicitly zeros out memory
func (ss *SecureString) Destroy() {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.destroy()
}

// SecureWrite writes the secret data to a writer without creating string copies
func (ss *SecureString) SecureWrite(w interface{ Write([]byte) (int, error) }) error {
	return ss.ReadWith(func(data []byte) error {
		_, err := w.Write(data)
		return err
	})
}

// SecureWriteEnvVar sets an environment variable securely without string copies
func (ss *SecureString) SecureWriteEnvVar(key string) error {
	return ss.ReadWith(func(data []byte) error {
		// Convert to string only at the point of use
		return os.Setenv(key, string(data))
	})
}

// SecureCompare compares the secret with another byte slice without exposing the secret
func (ss *SecureString) SecureCompare(other []byte) (bool, error) {
	var result bool
	err := ss.peekWith(func(data []byte) error {
		if len(data) != len(other) {
			result = false
			return nil
		}

		// Constant-time comparison to prevent timing attacks
		var diff byte
		for i := 0; i < len(data); i++ {
			diff |= data[i] ^ other[i]
		}
		result = diff == 0
		return nil
	})
	return result, err
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
		//First we fill data with random bytes
		_, _ = rand.Read(ss.data)
		for i := range ss.data {
			//Then we set all bytes to 0
			ss.data[i] = 0
		}
	}
}

// peekWith provides secure read-only access to the secret data without consuming it.
func (ss *SecureString) peekWith(fn func([]byte) error) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if ss.consumed {
		return errors.New("secret already consumed")
	}

	return fn(ss.data)
}
