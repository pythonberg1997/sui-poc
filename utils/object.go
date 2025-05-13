package utils

import (
	"context"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
)

func NewSharedObjectRefFromObjectId(cli sui.ISuiAPI, objectId string, mutable bool) (*transaction.SharedObjectRef, error) {
	ctx := context.Background()

	objectIdBytes, err := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(objectId))
	if err != nil {
		return nil, fmt.Errorf("error converting object ID to bytes: %w", err)
	}

	obj, err := cli.SuiGetObject(ctx, models.SuiGetObjectRequest{
		ObjectId: objectId,
		Options: models.SuiObjectDataOptions{
			ShowOwner: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error getting object: %w", err)
	}

	var initialSharedVersion uint64
	if owner, ok := obj.Data.Owner.(map[string]interface{}); ok {
		initialSharedVersion = uint64(owner["Shared"].(map[string]interface{})["initial_shared_version"].(float64))
	} else {
		return nil, fmt.Errorf("error parsing object owner: %w", err)
	}

	sharedObj := transaction.SharedObjectRef{
		ObjectId:             *objectIdBytes,
		InitialSharedVersion: initialSharedVersion,
		Mutable:              mutable,
	}

	return &sharedObj, nil
}

func NewOwnedObjectRefFromObjectId(cli sui.ISuiAPI, objectId string) (*transaction.SuiObjectRef, error) {
	ctx := context.Background()

	obj, err := cli.SuiGetObject(ctx, models.SuiGetObjectRequest{ObjectId: objectId})
	if err != nil {
		return nil, fmt.Errorf("error getting object: %w", err)
	}

	return transaction.NewSuiObjectRef(
		models.SuiAddress(objectId),
		obj.Data.Version,
		models.ObjectDigest(obj.Data.Digest),
	)
}
