package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	codeErrIbcCroDenomEmpty = uint32(iota) + 2 // NOTE: code 1 is reserved for internal errors
	codeErrIbcCroDenomInvalid
	codeErrContractAddressInvalid
)

// x/seele module sentinel errors
var (
	ErrIbcCroDenomEmpty       = sdkerrors.Register(ModuleName, codeErrIbcCroDenomEmpty, "ibc seele denom is not set")
	ErrIbcCroDenomInvalid     = sdkerrors.Register(ModuleName, codeErrIbcCroDenomInvalid, "ibc seele denom is invalid")
	ErrContractAddressInvalid = sdkerrors.Register(ModuleName, codeErrContractAddressInvalid, "contract address invalid")
	// this line is used by starport scaffolding # ibc/errors
)
