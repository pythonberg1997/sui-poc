package utils

import (
	"context"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
)

func GetCoinObjectId(client sui.ISuiAPI, signerAddress, coinType string) ([]models.CoinData, error) {
	ctx := context.Background()
	ownedCoins, err := client.SuiXGetCoins(ctx, models.SuiXGetCoinsRequest{
		Owner:    signerAddress,
		CoinType: coinType,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching available coins: %w", err)
	}

	if len(ownedCoins.Data) == 0 {
		return nil, fmt.Errorf("no coins found of type %s", coinType)
	}

	return ownedCoins.Data, nil
}
