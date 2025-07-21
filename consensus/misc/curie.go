package misc

import (
	"github.com/tenderly/net-scroll-geth/common"
	"github.com/tenderly/net-scroll-geth/core/state"
	"github.com/tenderly/net-scroll-geth/log"
	"github.com/tenderly/net-scroll-geth/rollup/rcfg"
)

// ApplyCurieHardFork modifies the state database according to the Curie hard-fork rules,
// updating the bytecode and storage of the L1GasPriceOracle contract.
func ApplyCurieHardFork(statedb *state.StateDB) {
	log.Info("Applying Curie hard fork")

	// update contract byte code
	statedb.SetCode(rcfg.L1GasPriceOracleAddress, rcfg.CurieL1GasPriceOracleBytecode)

	// initialize new storage slots
	statedb.SetState(rcfg.L1GasPriceOracleAddress, rcfg.IsCurieSlot, common.BytesToHash([]byte{1}))
	statedb.SetState(rcfg.L1GasPriceOracleAddress, rcfg.L1BlobBaseFeeSlot, common.BytesToHash([]byte{1}))
	statedb.SetState(rcfg.L1GasPriceOracleAddress, rcfg.CommitScalarSlot, common.BigToHash(rcfg.InitialCommitScalar))
	statedb.SetState(rcfg.L1GasPriceOracleAddress, rcfg.BlobScalarSlot, common.BigToHash(rcfg.InitialBlobScalar))
}
