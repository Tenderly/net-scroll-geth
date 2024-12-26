package rawdb

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/tenderly/net-scroll-geth/common"
	"github.com/tenderly/net-scroll-geth/core/types"
)

func TestReadBlockRowConsumption(t *testing.T) {
	l2BlockHash := common.BigToHash(big.NewInt(10))
	rc := types.RowConsumption{
		types.SubCircuitRowUsage{Name: "aa", RowNumber: 12},
		types.SubCircuitRowUsage{Name: "bb", RowNumber: 100},
	}
	db := NewMemoryDatabase()
	WriteBlockRowConsumption(db, l2BlockHash, &rc)
	got := ReadBlockRowConsumption(db, l2BlockHash)
	if got == nil || !reflect.DeepEqual(rc, *got) {
		t.Fatal("RowConsumption mismatch", "expected", rc, "got", got)
	}
}
