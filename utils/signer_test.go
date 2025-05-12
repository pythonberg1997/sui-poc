package utils

import (
	"fmt"
	"testing"

	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/stretchr/testify/require"
)

func TestKey(t *testing.T) {
	seed, err := decodeSuiPrivateKey("")
	require.NoError(t, err)
	require.Equal(t, 32, len(seed))
	fmt.Printf("Private Key: %x\n", seed)

	suiSigner := signer.NewSigner(seed)
	fmt.Printf("address: %s\n", suiSigner.Address)
}
