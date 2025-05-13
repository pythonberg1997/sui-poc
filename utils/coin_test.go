package utils

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/stretchr/testify/require"
)

func TestGetCoinObjectId(t *testing.T) {
	cli := sui.NewSuiClient("https://sui-fullnode-qa.internal.nodereal.io")

	suiSigner, err := NewSignerFromSecretKey(os.Getenv("PRIVATE_KEY"))
	require.NoError(t, err)

	ids, err := GetCoinObjectId(cli, suiSigner.Address, "0xdba34672e30cb065b1f93e3ab55318768fd6fef66c15942c9f7cb846e2f900e7::usdc::USDC")
	require.NoError(t, err)
	for _, id := range ids {
		fmt.Println("Coin Object ID:", id.CoinObjectId)
	}
}

func TestMergeCoin1(t *testing.T) {
	ctx := context.Background()
	cli := sui.NewSuiClient("https://sui-fullnode-qa.internal.nodereal.io")

	suiSigner, err := NewSignerFromSecretKey(os.Getenv("PRIVATE_KEY"))
	require.NoError(t, err)

	txnBytes, err := cli.MergeCoins(ctx, models.MergeCoinsRequest{
		Signer:      suiSigner.Address,
		PrimaryCoin: "0x903160997126fbcc591924dd948af3c5ba9082ae5a4025dcd12d89e76ad98e20",
		CoinToMerge: "0xe8f4183b75a1442366f9fe7fb11e1911267ec08387bb8ccde3b5d99f779d0035",
		GasBudget:   "50000000",
	})
	require.NoError(t, err)

	resp, err := cli.SignAndExecuteTransactionBlock(ctx, models.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: txnBytes,
		PriKey:      suiSigner.PriKey,
		Options: models.SuiTransactionBlockOptions{
			ShowInput:          true,
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowEvents:         true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
		},
		RequestType: "WaitForLocalExecution",
	})
	require.NoError(t, err)
	fmt.Println("Transaction response:", resp)
}

func TestMergeCoin2(t *testing.T) {
	ctx := context.Background()
	cli := sui.NewSuiClient("https://sui-fullnode-qa.internal.nodereal.io")
	suiClient, ok := cli.(*sui.Client)
	if !ok {
		panic("not sui client")
	}

	suiSigner, err := NewSignerFromSecretKey(os.Getenv("PRIVATE_KEY"))
	require.NoError(t, err)

	gasCoinObjectId := "0xd93b7b103bb7488320d13337400c18fca9e3e8eac20a53e2661d2ef1fe32e949"
	gasCoinObj, err := suiClient.SuiGetObject(ctx, models.SuiGetObjectRequest{ObjectId: gasCoinObjectId})
	if err != nil {
		fmt.Printf("Error getting gas coin object: %v\n", err)
		return
	}
	gasCoin, err := transaction.NewSuiObjectRef(
		models.SuiAddress(gasCoinObjectId),
		gasCoinObj.Data.Version,
		models.ObjectDigest(gasCoinObj.Data.Digest),
	)
	if err != nil {
		fmt.Printf("Error creating gas coin reference: %v\n", err)
		return
	}

	tx := transaction.NewTransaction()
	tx.SetSuiClient(suiClient).
		SetSigner(suiSigner).
		SetSender(models.SuiAddress(suiSigner.Address)).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasPayment([]transaction.SuiObjectRef{*gasCoin}).
		SetGasOwner(models.SuiAddress(suiSigner.Address))

	dstCoinObj, err := NewOwnedObjectRefFromObjectId(cli, "0x903160997126fbcc591924dd948af3c5ba9082ae5a4025dcd12d89e76ad98e20")
	require.NoError(t, err)

	sourceCoinObj, err := NewOwnedObjectRefFromObjectId(cli, "0xa65f83d5804dc1efd8248591cf7ad8e441a12193dbf8eb967638aac771d618e4")
	require.NoError(t, err)

	tx.MergeCoins(
		tx.Object(transaction.CallArg{
			Object: &transaction.ObjectArg{
				ImmOrOwnedObject: dstCoinObj,
			},
		},
		),
		[]transaction.Argument{
			tx.Object(transaction.CallArg{
				Object: &transaction.ObjectArg{
					ImmOrOwnedObject: sourceCoinObj,
				},
			},
			),
		},
	)

	resp, err := tx.Execute(
		ctx,
		models.SuiTransactionBlockOptions{
			ShowInput:    true,
			ShowRawInput: true,
			ShowEffects:  true,
		},
		"WaitForLocalExecution",
	)
	require.NoError(t, err)
	fmt.Printf("Transaction success. digest: %v\n", resp.Digest)
}
