package merkle

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

func RegisterCodec(cdc *wire.Codec) {
	cdc.RegisterInterface((*KeyProof)(nil), nil)
	cdc.RegisterConcrete(ExistsProof{}, "cosmos-sdk/ExistsProof", nil)
	cdc.RegisterConcrete(AbsentProof{}, "cosmos-sdk/AbsentProof", nil)
	cdc.RegisterConcrete(RangeProof{}, "cosmos-sdk/RangeProof", nil)
}
