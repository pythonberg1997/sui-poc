package utils

import (
	"context"
	"fmt"
	"strconv"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
)

func EstimateGasBudge(ctx context.Context, cli sui.ISuiAPI, tx *transaction.Transaction) (uint64, error) {
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
		return 0, fmt.Errorf("error converting transaction to request: %w", err)
	}
	simResp, err := cli.SuiDryRunTransactionBlock(ctx, models.SuiDryRunTransactionBlockRequest{TxBytes: req.TxBytes})
	if err != nil {
		return 0, fmt.Errorf("error simulating transaction: %w", err)
	}

	computationCost, err := strconv.ParseUint(simResp.Effects.GasUsed.ComputationCost, 10, 64)
	if err != nil {
		return 0, err
	}
	storageCost, err := strconv.ParseUint(simResp.Effects.GasUsed.StorageCost, 10, 64)
	if err != nil {
		return 0, err
	}
	storageRebate, err := strconv.ParseUint(simResp.Effects.GasUsed.StorageRebate, 10, 64)
	if err != nil {
		return 0, err
	}

	return computationCost + storageCost - storageRebate, nil
}
