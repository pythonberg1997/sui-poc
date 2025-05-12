package utils

import (
	"crypto/ed25519"
	"fmt"

	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/btcsuite/btcutil/bech32"
)

const SuiPrivateKeyPrefix = "suiprivkey"

func NewSignerFromSecretKey(bech32Key string) (*signer.Signer, error) {
	seed, err := decodeSuiPrivateKey(bech32Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Sui private key: %w", err)
	}

	if len(seed) != ed25519.SeedSize {
		return nil, fmt.Errorf("invalid private key length: expected %d, got %d",
			ed25519.SeedSize, len(seed))
	}

	return signer.NewSigner(seed), nil
}

func decodeSuiPrivateKey(bech32Key string) (ed25519.PrivateKey, error) {
	hrp, data, err := bech32.Decode(bech32Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode bech32: %w", err)
	}

	if hrp != SuiPrivateKeyPrefix {
		return nil, fmt.Errorf("unexpected HRP: got %s, want %s", hrp, SuiPrivateKeyPrefix)
	}

	converted, err := bech32.ConvertBits(data, 5, 8, false)
	if err != nil {
		return nil, fmt.Errorf("failed to convert bits: %w", err)
	}

	return converted[1:], nil
}
