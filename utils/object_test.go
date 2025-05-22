package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/require"
	"encoding/hex"
	"github.com/btcsuite/btcutil/base58"
	"encoding/base64"
)

func TestSharedObject(t *testing.T) {
	cli := sui.NewSuiClient("https://sui-fullnode-qa.internal.nodereal.io")

	obj, err := NewSharedObjectRefFromObjectId(cli, "0x6", false)
	require.NoError(t, err)
	fmt.Println(obj)
}

func TestGetAllObjects(t *testing.T) {
	ctx := context.Background()
	cli := sui.NewSuiClient("https://sui-fullnode-qa.internal.nodereal.io")

	bot := "0xce77f41e64a37defb583b1d04248932afab503e6b02533ed85bf36d8b202f77f"
	resp, err := cli.SuiXGetOwnedObjects(ctx, models.SuiXGetOwnedObjectsRequest{
		Address: bot,
		Query: models.SuiObjectResponseQuery{
			// Filter: models.ObjectFilterByPackage{
			// 	Package: "0x2::coin::Coin",
			// },
			// Filter: models.ObjectFilterByStructType{
			// 	StructType: "0x2::coin::Coin<0xdba34672e30cb065b1f93e3ab55318768fd6fef66c15942c9f7cb846e2f900e7::usdc::USDC>",
			// },
			Options: models.SuiObjectDataOptions{
				ShowOwner: true,
				ShowType:  true,
			},
		},
		Cursor: nil,
		Limit:  10,
	})
	require.NoError(t, err)
	fmt.Println("Response:", resp)
}

func TestDecode(t *testing.T) {
	objectId := "3864c7c59a4889fec05d1aae4bc9dba5a0e0940594b424fbed44cb3f6ac4c032"
	bz, err := hex.DecodeString(objectId)
	require.NoError(t, err)
	fmt.Println("objectId:", bz)

	digest := "J6Ki9HpWpf9EmUahr5gGCXqrt3hd6vTYaNweTj1HuqBV"
	bz = base58.Decode(digest)
	fmt.Println("digest:", bz)

	pure := "TtookpwDAAA="
	bz, err = base64.StdEncoding.DecodeString(pure)
	require.NoError(t, err)
	fmt.Println("pure:", bz)

	b64Tx := "AAAIAQBiBnaLHr5TzVJu29SyLYueU6dFUyLeaY9ZeK11c+8IrXPUayEAAAAAIP3zvI2RTxZdaMnH47J4OJQbyHdIb50BtZ2JxNb4gD7AAQDo6vegaubVlLDf0hxZfOtmAkBc4ryrHxXfLLjFo72CF4bUayEAAAAAIECXRV/tLtwVmiAi/xXB6L5v1KFuEZQ0oah1NwYdBcLfAAgAhNcXAAAAAAEB2qRikmMsPE2PMfI+oPmzaij/NnfpaEmA5EOEA6Z6PY8uBRgAAAAAAAABAa2qRWjdYdo2MMQNLRRttp5jmUbSJZnqQxXT0+oBmR8aY8ddIQAAAAABAQFjm15DPaMXOegAzQhfNW5kyuIilm0PGxG9ncdrMi/1i2LvCxMAAAAAAQEBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAYBAAAAAAAAAAAACE7aKJKcAwAABwMBAAABAQEAAgEAAAEBAgACAwEAAAABAQIAADhkx8WaSIn+wF0arkvJ26Wg4JQFlLQk++1Eyz9qxMAyBWNldHVzCHN3YXBfYjJhAgfhtFoOZBuZVaIKoK0cH0rYaq2K+wcpbUCF40mlDpC9ygRibHVlBEJMVUUAB9ujRnLjDLBlsfk+OrVTGHaP1v72bBWULJ98uEbi+QDnBHVzZGMEVVNEQwAFAQMAAQQAAQUAAwIAAAABBgAAJIX+udQsfDvLjs3lVa1A8bBz2ftPrzVPotMKCxg6I84FdXRpbHMYdHJhbnNmZXJfb3JfZGVzdHJveV9jb2luAQfbo0Zy4wywZbH5Pjq1Uxh2j9b+9mwVlCyffLhG4vkA5wR1c2RjBFVTREMAAQMBAAAAACSF/rnULHw7y47N5VWtQPGwc9n7T681T6LTCgsYOiPOBXV0aWxzFGNoZWNrX2NvaW5fdGhyZXNob2xkAQfhtFoOZBuZVaIKoK0cH0rYaq2K+wcpbUCF40mlDpC9ygRibHVlBEJMVUUAAgMDAAAAAQcAACSF/rnULHw7y47N5VWtQPGwc9n7T681T6LTCgsYOiPOBXV0aWxzGHRyYW5zZmVyX29yX2Rlc3Ryb3lfY29pbgEH4bRaDmQbmVWiCqCtHB9K2GqtivsHKW1AheNJpQ6QvcoEYmx1ZQRCTFVFAAEDAwAAAPGLdfcB6yxirRUeIzQYdrBeeCZsC49tvfE1MjQ3HRe+AeFyymxHLoC1zuCDjWceFIVGpQinWoTY3E/nQtJ0PXeLhtRrIQAAAAAg7wnlIFKxbd6S0eZ7pktb3VGeUCezpe3rCIK/6e1qR+7xi3X3AessYq0VHiM0GHawXngmbAuPbb3xNTI0Nx0Xvu4CAAAAAAAAAOH1BQAAAAAA"
	txBz, err := base64.StdEncoding.DecodeString(b64Tx)
	require.NoError(t, err)
	fmt.Println("tx:", txBz)
}
