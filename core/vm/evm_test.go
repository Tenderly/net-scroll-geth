package vm

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/tenderly/net-scroll-geth/common"
	"github.com/tenderly/net-scroll-geth/core/rawdb"
	"github.com/tenderly/net-scroll-geth/core/state"
	"github.com/tenderly/net-scroll-geth/params"
)

// echoPrecompile charges a fixed gas and returns its input after a fixed prefix.
type echoPrecompile struct{}

func (echoPrecompile) RequiredGas([]byte) uint64 { return 100 }

func (echoPrecompile) Run(input []byte) ([]byte, error) {
	return append([]byte("echo:"), input...), nil
}

// AI-DEV-TEST: written by an agent, not reviewed by a human.
func TestOverridePrecompiles(t *testing.T) {
	var (
		target = common.HexToAddress("0xC100000000000000000000000000000000000000")
		caller = common.HexToAddress("0x00000000000000000000000000000000000000ca")
		input  = []byte("input")
		gas    = uint64(10_000)
	)
	newEVM := func() *EVM {
		statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
		vmctx := BlockContext{
			CanTransfer: func(StateDB, common.Address, *big.Int) bool { return true },
			Transfer:    func(StateDB, common.Address, common.Address, *big.Int) {},
			BlockNumber: big.NewInt(35_527_145),
			Time:        new(big.Int).SetUint64(*params.ScrollMainnetChainConfig.FeynmanTime),
		}
		return NewEVM(vmctx, TxContext{}, statedb, params.ScrollMainnetChainConfig, Config{})
	}
	t.Run("unset, the lookup is the fork's map", func(t *testing.T) {
		evm := newEVM()
		if !evm.chainRules.IsFeynman || evm.chainRules.IsGalileo {
			t.Fatalf("test EVM is not on Feynman rules")
		}
		for addr, want := range PrecompiledContractsFeynman {
			if got, ok := evm.precompile(addr); !ok || got != want {
				t.Errorf("precompile %x: have %v (%v), want %v", addr, got, ok, want)
			}
		}
		if p, ok := evm.precompile(target); ok {
			t.Errorf("precompile %x: have %v, want none", target, p)
		}

		ret, left, err := evm.Call(AccountRef(caller), target, input, gas, new(big.Int), nil)
		if err != nil || ret != nil || left != gas {
			t.Errorf("call %x: have (%x, %d, %v), want (nil, %d, nil)", target, ret, left, err, gas)
		}
	})

	t.Run("set, a CALL to the address runs the override", func(t *testing.T) {
		evm := newEVM()
		evm.OverridePrecompiles(map[common.Address]PrecompiledContract{target: echoPrecompile{}})

		ret, left, err := evm.Call(AccountRef(caller), target, input, gas, new(big.Int), nil)
		if err != nil {
			t.Fatalf("call %x: %v", target, err)
		}
		if want := []byte("echo:input"); !bytes.Equal(ret, want) {
			t.Errorf("call %x: have output %q, want %q", target, ret, want)
		}
		if used := gas - left; used != 100 {
			t.Errorf("call %x: have gas used %d, want 100", target, used)
		}
		if _, ok := newEVM().precompile(target); ok {
			t.Errorf("override leaked into another EVM")
		}
	})

	t.Run("a nil override removes the fork's precompile", func(t *testing.T) {
		ecrecoverAddr := common.BytesToAddress([]byte{1})
		evm := newEVM()
		evm.OverridePrecompiles(map[common.Address]PrecompiledContract{ecrecoverAddr: nil})

		if p, ok := evm.precompile(ecrecoverAddr); ok {
			t.Errorf("precompile %x: have %v, want none", ecrecoverAddr, p)
		}
	})
}
