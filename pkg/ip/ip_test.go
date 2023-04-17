package ip

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultAddressService(t *testing.T) {
	as := NewDefaultAddressService()
	ip := "128.199.146.208"
	res := as.Query(ip)
	assert.NotEmpty(t, res)
	fmt.Printf("%s: %s\n", ip, res)
}
