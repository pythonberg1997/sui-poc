package main

import (
	"context"
	"fmt"
	"os"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/joho/godotenv"

	"sui-poc/utils"
)

func main() {
	ctx := context.Background()

	cli := sui.NewSuiClient("https://sui-fullnode-qa.internal.nodereal.io")
	suiClient, ok := cli.(*sui.Client)
	if !ok {
		panic("not sui client")
	}

	godotenv.Load(".env")
	suiSigner, err := utils.NewSignerFromSecretKey(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		fmt.Println("Error creating signer:", err)
		return
	}
	signerAddress := suiSigner.Address

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
		SetSender(models.SuiAddress(signerAddress)).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasPayment([]transaction.SuiObjectRef{*gasCoin}).
		SetGasOwner(models.SuiAddress(signerAddress))

	// 1. split coin
	splitCoin := tx.SplitCoins(tx.Gas(), []transaction.Argument{
		tx.Pure(uint64(1000000000 * 0.0001)),
	})

	// 2. move call
	packageId := "0x11451575c775a3e633437b827ecbc1eb51a5964b0302210b28f5b89880be21a2"
	module := "cetus"
	funcName := "swap_b2a"

	coinAAddressBytes, err := transaction.ConvertSuiAddressStringToBytes("0xdba34672e30cb065b1f93e3ab55318768fd6fef66c15942c9f7cb846e2f900e7")
	if err != nil {
		fmt.Printf("Error converting coin A address: %v\n", err)
		return
	}
	coinBAddressBytes, err := transaction.ConvertSuiAddressStringToBytes("0x0000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {
		fmt.Printf("Error converting coin B address: %v\n", err)
		return
	}

	globalConfig, err := utils.NewSharedObjectRefFromObjectId(cli, "0xdaa46292632c3c4d8f31f23ea0f9b36a28ff3677e9684980e4438403a67a3d8f", true)
	if err != nil {
		fmt.Printf("Error creating global config reference: %v\n", err)
		return
	}
	pool, err := utils.NewSharedObjectRefFromObjectId(cli, "0xb8d7d9e66a60c239e7a60110efcf8de6c705580ed924d0dde141f4a0e2c90105", true)
	if err != nil {
		fmt.Printf("Error creating pool reference: %v\n", err)
		return
	}
	partner, err := utils.NewSharedObjectRefFromObjectId(cli, "0x639b5e433da31739e800cd085f356e64cae222966d0f1b11bd9dc76b322ff58b", true)
	if err != nil {
		fmt.Printf("Error creating partner reference: %v\n", err)
		return
	}
	clock, err := utils.NewSharedObjectRefFromObjectId(cli, "0x6", false)
	if err != nil {
		fmt.Printf("Error creating clock reference: %v\n", err)
		return
	}

	outCoin := tx.MoveCall(
		models.SuiAddress(packageId),
		module,
		funcName,
		[]transaction.TypeTag{
			{
				Struct: &transaction.StructTag{
					Address: *coinAAddressBytes,
					Module:  "usdc",
					Name:    "USDC",
				},
			},
			{
				Struct: &transaction.StructTag{
					Address: *coinBAddressBytes,
					Module:  "sui",
					Name:    "SUI",
				},
			},
		},
		[]transaction.Argument{
			tx.Object(transaction.CallArg{
				Object: &transaction.ObjectArg{
					SharedObject: globalConfig,
				},
			},
			),
			tx.Object(transaction.CallArg{
				Object: &transaction.ObjectArg{
					SharedObject: pool,
				},
			},
			),
			tx.Object(transaction.CallArg{
				Object: &transaction.ObjectArg{
					SharedObject: partner,
				},
			},
			),
			splitCoin,
			tx.Object(transaction.CallArg{
				Object: &transaction.ObjectArg{
					SharedObject: clock,
				},
			},
			),
		},
	)

	// 3. transfer coins
	tx.TransferObjects([]transaction.Argument{outCoin}, tx.Pure(signerAddress))

	resp, err := tx.Execute(
		ctx,
		models.SuiTransactionBlockOptions{
			ShowInput:          true,
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowEvents:         true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
		},
		"WaitForLocalExecution",
	)
	if err != nil {
		fmt.Printf("Error executing transaction: %v\n", err)
		return
	}

	fmt.Println(resp.Digest, resp.Effects, resp.Results)
}
