package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

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

	// initialVersion := gasCoinObj.Data.Version
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go monitorGasObjectVersion(ctx, suiClient, gasCoinObjectId, initialVersion, &wg)

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

	globalConfig, err := utils.NewSharedObjectRefFromObjectId(cli, "0xdaa46292632c3c4d8f31f23ea0f9b36a28ff3677e9684980e4438403a67a3d8f", false)
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

	// latestSequenceNumber, err := suiClient.SuiGetLatestCheckpointSequenceNumber(ctx)
	// if err != nil {
	// 	fmt.Printf("Error getting latest checkpoint sequence number: %v\n", err)
	// 	return
	// }
	// fmt.Printf("Sequence number before execute: %d\n", latestSequenceNumber)
	// fmt.Printf("Time before execute: %s\n", time.Now().Format(time.RFC3339Nano))

	gasPrice, err := suiClient.SuiXGetReferenceGasPrice(ctx)
	if err != nil {
		fmt.Printf("Error getting gas price: %v\n", err)
		return
	}
	tx.SetGasPrice(uint64(float64(gasPrice) * 1.05))

	gasBudget, err := utils.EstimateGasBudge(ctx, suiClient, tx)
	if err != nil {
		fmt.Printf("Error estimating gas budge: %v\n", err)
		return
	}
	tx.SetGasBudget(uint64(float64(gasBudget) * 1.2))

	req, err := tx.ToSuiExecuteTransactionBlockRequest(
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
		fmt.Printf("Error converting transaction to request: %v\n", err)
		return
	}
	fmt.Printf("tx bytes: %s\n", req.TxBytes)
	fmt.Printf("tx sig: %s\n", req.Signature[0])

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

func monitorGasObjectVersion(ctx context.Context, client *sui.Client, objectId string, initialVersion string, wg *sync.WaitGroup) {
	defer wg.Done()

	startTime := time.Now()
	fmt.Printf("Starting to monitor gas object with initial version: %s at time: %s\n", initialVersion, startTime.Format(time.RFC3339))

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			obj, err := client.SuiGetObject(ctx, models.SuiGetObjectRequest{ObjectId: objectId})
			if err != nil {
				fmt.Printf("Error checking gas object: %v\n", err)
				continue
			}

			currentVersion := obj.Data.Version

			if currentVersion != initialVersion {
				elapsed := time.Since(startTime)
				fmt.Printf("Gas object version changed from %s to %s\n", initialVersion, currentVersion)
				fmt.Printf("Time elapsed: %s\n", elapsed)
				fmt.Printf("New version detected at: %s\n", time.Now().Format(time.RFC3339Nano))

				latestSequenceNumber, _ := client.SuiGetLatestCheckpointSequenceNumber(ctx)
				fmt.Printf("Sequence number: %d\n", latestSequenceNumber)
				return
			}
		case <-ctx.Done():
			fmt.Println("Context cancelled, stopping monitoring")
			return
		}
	}
}
