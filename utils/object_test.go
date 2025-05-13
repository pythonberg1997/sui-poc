package utils

import (
	"fmt"
	"testing"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/require"
)

func TestSharedObject(t *testing.T) {
	cli := sui.NewSuiClient("https://sui-fullnode-qa.internal.nodereal.io")

	obj, err := NewSharedObjectRefFromObjectId(cli, "0x6", false)
	require.NoError(t, err)
	fmt.Println(obj)
}
