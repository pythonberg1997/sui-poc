package utils

import (
	"fmt"
	"github.com/block-vision/sui-go-sdk/models"
	"strconv"
)

func GetGasUsedFromSimulationResponse(response *models.SuiTransactionBlockResponse) (uint64, error) {
	if response == nil {
		return 0, fmt.Errorf("simulation response is nil")
	}

	computationCost, err := strconv.ParseUint(response.Effects.GasUsed.ComputationCost, 10, 64)
	if err != nil {
		return 0, err
	}
	storageCost, err := strconv.ParseUint(response.Effects.GasUsed.StorageCost, 10, 64)
	if err != nil {
		return 0, err
	}
	storageRebate, err := strconv.ParseUint(response.Effects.GasUsed.StorageRebate, 10, 64)
	if err != nil {
		return 0, err
	}

	return computationCost + storageCost - storageRebate, nil
}
