package share

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ipfs/go-blockservice"

	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/pkg/da"
	"github.com/celestiaorg/celestia-app/pkg/shares"
	"github.com/celestiaorg/rsmt2d"
)

var (
	emptyRoot *Root
	emptyEDS  *rsmt2d.ExtendedDataSquare
)

func init() {
	// compute empty block EDS and DAH for it
	shares := emptyDataSquare()
	eds, err := da.ExtendShares(shares)
	if err != nil {
		panic(fmt.Errorf("failed to create empty EDS: %w", err))
	}
	emptyEDS = eds

	dah := da.NewDataAvailabilityHeader(eds)
	minDAH := da.MinDataAvailabilityHeader()
	if !bytes.Equal(minDAH.Hash(), dah.Hash()) {
		panic(fmt.Sprintf("mismatch in calculated minimum DAH and minimum DAH from celestia-app, "+
			"expected %s, got %s", minDAH.String(), dah.String()))
	}
	emptyRoot = &dah

	// precompute Hash, so it's cached internally to avoid potential races
	emptyRoot.Hash()
}

// EmptyRoot returns Root of an empty EDS.
func EmptyRoot() *Root {
	return emptyRoot
}

// EnsureEmptySquareExists checks if the given DAG contains an empty block data square.
// If it does not, it stores an empty block. This optimization exists to prevent
// redundant storing of empty block data so that it is only stored once and returned
// upon request for a block with an empty data square. Ref: header/constructors.go#L56
func EnsureEmptySquareExists(ctx context.Context, bServ blockservice.BlockService) (*rsmt2d.ExtendedDataSquare, error) {
	shares := emptyDataSquare()
	return AddShares(ctx, shares, bServ)
}

// EmptyExtendedDataSquare returns the EDS of the empty block data square.
func EmptyExtendedDataSquare() *rsmt2d.ExtendedDataSquare {
	return emptyEDS
}

// emptyDataSquare returns the minimum size data square filled with tail padding.
func emptyDataSquare() [][]byte {
	result := shares.TailPaddingShares(appconsts.MinShareCount)
	return shares.ToBytes(result)
}
