package tokens

import (
	"crypto/rand"
	"fmt"
)

func New() string {
	b := make([]byte, 40)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
