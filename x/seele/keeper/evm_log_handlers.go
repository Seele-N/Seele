package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Seele-N/Seele/x/seele/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	_ types.EvmLogHandler = SendToAccountHandler{}
	_ types.EvmLogHandler = SendToEthereumHandler{}
	_ types.EvmLogHandler = SendToIbcHandler{}
	_ types.EvmLogHandler = SendCroToIbcHandler{}

	_ types.EvmLogHandler = SendSnpStakeHandler{}
)

const (
	SendToAccountEventName  = "__SeeleSendToAccount"
	SendToEthereumEventName = "__SeeleSendToEthereum"
	SendToIbcEventName      = "__SeeleSendToIbc"
	SendCroToIbcEventName   = "__SeeleSendSeeleToIbc"

	SnpStakingEventName = "Snp_Staking"
)

var (
	// SendToAccountEvent represent the signature of
	// `event __SeeleSendToAccount(address recipient, uint256 amount)`
	SendToAccountEvent abi.Event

	// SendToEthereumEvent represent the signature of
	// `event __SeeleSendToEthereum(address recipient, uint256 amount, uint256 bridge_fee)`
	SendToEthereumEvent abi.Event

	// SendToIbcEvent represent the signature of
	// `event __SeeleSendToIbc(string recipient, uint256 amount)`
	SendToIbcEvent abi.Event

	// SnpStakingEvent represent the signature of
	// `event __SeeleSendSeeleToIbc(string recipient, uint256 amount)`
	SendCroToIbcEvent abi.Event

	// SnpStakeEvent represent the signature of
	// `event Snp_Staking(address validator, address delegator,uint256 amount)`
	SnpStakeEvent abi.Event
)

func init() {
	addressType, _ := abi.NewType("address", "", nil)
	uint256Type, _ := abi.NewType("uint256", "", nil)
	stringType, _ := abi.NewType("string", "", nil)
	SendToAccountEvent = abi.NewEvent(
		SendToAccountEventName,
		SendToAccountEventName,
		false,
		abi.Arguments{abi.Argument{
			Name:    "recipient",
			Type:    addressType,
			Indexed: false,
		}, abi.Argument{
			Name:    "amount",
			Type:    uint256Type,
			Indexed: false,
		}},
	)
	SendToEthereumEvent = abi.NewEvent(
		SendToEthereumEventName,
		SendToEthereumEventName,
		false,
		abi.Arguments{abi.Argument{
			Name:    "recipient",
			Type:    addressType,
			Indexed: false,
		}, abi.Argument{
			Name:    "amount",
			Type:    uint256Type,
			Indexed: false,
		}, abi.Argument{
			Name:    "bridge_fee",
			Type:    uint256Type,
			Indexed: false,
		}},
	)
	SendToIbcEvent = abi.NewEvent(
		SendToIbcEventName,
		SendToIbcEventName,
		false,
		abi.Arguments{abi.Argument{
			Name:    "sender",
			Type:    addressType,
			Indexed: false,
		}, abi.Argument{
			Name:    "recipient",
			Type:    stringType,
			Indexed: false,
		}, abi.Argument{
			Name:    "amount",
			Type:    uint256Type,
			Indexed: false,
		}},
	)
	SendCroToIbcEvent = abi.NewEvent(
		SendCroToIbcEventName,
		SendCroToIbcEventName,
		false,
		abi.Arguments{abi.Argument{
			Name:    "sender",
			Type:    addressType,
			Indexed: false,
		}, abi.Argument{
			Name:    "recipient",
			Type:    stringType,
			Indexed: false,
		}, abi.Argument{
			Name:    "amount",
			Type:    uint256Type,
			Indexed: false,
		}},
	)

	SnpStakeEvent = abi.NewEvent(
		SnpStakingEventName,
		SnpStakingEventName,
		false,
		abi.Arguments{abi.Argument{
			Name:    "validator",
			Type:    addressType,
			Indexed: false,
		}, abi.Argument{
			Name:    "delegator",
			Type:    addressType,
			Indexed: false,
		}, abi.Argument{
			Name:    "amount",
			Type:    uint256Type,
			Indexed: false,
		}},
	)
}

// SendSnpStakeHandler handles `Snp_Staking` log
type SendSnpStakeHandler struct {
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
	seeleKeeper   Keeper
}

func NewSendSnpStakeHandler(bankKeeper types.BankKeeper, stakingKeeper types.StakingKeeper, seeleKeeper Keeper) *SendSnpStakeHandler {
	return &SendSnpStakeHandler{
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		seeleKeeper:   seeleKeeper,
	}
}

func (h SendSnpStakeHandler) EventID() common.Hash {
	return SnpStakeEvent.ID
}

func (h SendSnpStakeHandler) Handle(ctx sdk.Context, contract common.Address, data []byte) error {
	h.seeleKeeper.Logger(ctx).Info("SendSnpStakeHandler")
	unpacked, err := SnpStakeEvent.Inputs.Unpack(data)
	if err != nil {
		// log and ignore
		h.seeleKeeper.Logger(ctx).Error("log signature matches but failed to decode", "error", err)
		return nil
	}
	h.seeleKeeper.Logger(ctx).Info("Event from contract:" + contract.Hex())
	h.seeleKeeper.Logger(ctx).Info("validator:" + unpacked[0].(common.Address).Hex())
	h.seeleKeeper.Logger(ctx).Info("delegator:" + unpacked[1].(common.Address).Hex())
	h.seeleKeeper.Logger(ctx).Info("amount:" + unpacked[2].(*big.Int).String())

	amount := unpacked[2].(*big.Int)
	valAddress := sdk.ValAddress(unpacked[0].(common.Address).Bytes())
	h.seeleKeeper.Logger(ctx).Info("valAddress:" + valAddress.String())
	validator, found := h.stakingKeeper.GetValidator(ctx, valAddress)
	if !found {
		return stakingtypes.ErrNoValidatorFound
	}
	delegator := sdk.AccAddress(unpacked[1].(common.Address).Bytes())
	coin := sdk.NewCoin("snp", sdk.NewIntFromBigInt(amount))
	h.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
	h.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegator, sdk.NewCoins(coin))
	newShares, err := h.stakingKeeper.Delegate(ctx, delegator, sdk.NewIntFromBigInt(amount), stakingtypes.Unbonded, validator, true)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingtypes.EventTypeDelegate,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, valAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyNewShares, newShares.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, delegator.String()),
		),
	})
	//contractAddr := sdk.AccAddress(contract.Bytes())
	//recipient := sdk.AccAddress(unpacked[0].(common.Address).Bytes())
	//coins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewIntFromBigInt(unpacked[1].(*big.Int))))
	//err = h.bankKeeper.SendCoins(ctx, contractAddr, recipient, coins)
	//if err != nil {
	//	return err
	//}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// SendToAccountHandler handles `__SeeleSendToAccount` log
type SendToAccountHandler struct {
	bankKeeper  types.BankKeeper
	seeleKeeper Keeper
}

func NewSendToAccountHandler(bankKeeper types.BankKeeper, seeleKeeper Keeper) *SendToAccountHandler {
	return &SendToAccountHandler{
		bankKeeper:  bankKeeper,
		seeleKeeper: seeleKeeper,
	}
}

func (h SendToAccountHandler) EventID() common.Hash {
	return SendToAccountEvent.ID
}

func (h SendToAccountHandler) Handle(ctx sdk.Context, contract common.Address, data []byte) error {
	unpacked, err := SendToAccountEvent.Inputs.Unpack(data)
	if err != nil {
		// log and ignore
		h.seeleKeeper.Logger(ctx).Error("log signature matches but failed to decode", "error", err)
		return nil
	}

	denom, found := h.seeleKeeper.GetDenomByContract(ctx, contract)
	if !found {
		return fmt.Errorf("contract %s is not connected to native token", contract)
	}

	contractAddr := sdk.AccAddress(contract.Bytes())
	recipient := sdk.AccAddress(unpacked[0].(common.Address).Bytes())
	coins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewIntFromBigInt(unpacked[1].(*big.Int))))
	err = h.bankKeeper.SendCoins(ctx, contractAddr, recipient, coins)
	if err != nil {
		return err
	}

	return nil
}

// SendToEthereumHandler handles `__SeeleSendToEthereum` log
type SendToEthereumHandler struct {
	seeleKeeper Keeper
}

func NewSendToEthereumHandler(seeleKeeper Keeper) *SendToEthereumHandler {
	return &SendToEthereumHandler{
		seeleKeeper: seeleKeeper,
	}
}

func (h SendToEthereumHandler) EventID() common.Hash {
	return SendToEthereumEvent.ID
}

// Handle returns error unconditionally.
// Since gravity bridge is removed and could be added later,
// we keep this event handler, but returns error unconditionally to prevent accidental access.
func (h SendToEthereumHandler) Handle(ctx sdk.Context, contract common.Address, data []byte) error {
	return fmt.Errorf("native action %s is not implemented", SendToEthereumEventName)
}

// SendToIbcHandler handles `__SeeleSendToIbc` log
type SendToIbcHandler struct {
	bankKeeper  types.BankKeeper
	seeleKeeper Keeper
}

func NewSendToIbcHandler(bankKeeper types.BankKeeper, seeleKeeper Keeper) *SendToIbcHandler {
	return &SendToIbcHandler{
		bankKeeper:  bankKeeper,
		seeleKeeper: seeleKeeper,
	}
}

func (h SendToIbcHandler) EventID() common.Hash {
	return SendToIbcEvent.ID
}

func (h SendToIbcHandler) Handle(ctx sdk.Context, contract common.Address, data []byte) error {
	unpacked, err := SendToIbcEvent.Inputs.Unpack(data)
	if err != nil {
		// log and ignore
		h.seeleKeeper.Logger(ctx).Info("log signature matches but failed to decode")
		return nil
	}

	denom, found := h.seeleKeeper.GetDenomByContract(ctx, contract)
	if !found {
		return fmt.Errorf("contract %s is not connected to native token", contract)
	}

	if !types.IsValidIBCDenom(denom) {
		return fmt.Errorf("the native token associated with the contract %s is not an ibc voucher", contract)
	}

	contractAddr := sdk.AccAddress(contract.Bytes())
	sender := sdk.AccAddress(unpacked[0].(common.Address).Bytes())
	recipient := unpacked[1].(string)
	amount := sdk.NewIntFromBigInt(unpacked[2].(*big.Int))
	coins := sdk.NewCoins(sdk.NewCoin(denom, amount))

	// First, transfer IBC coin to user so that he will be the refunded address if transfer fails
	if err = h.bankKeeper.SendCoins(ctx, contractAddr, sender, coins); err != nil {
		return err
	}
	// Initiate IBC transfer from sender account
	if err = h.seeleKeeper.IbcTransferCoins(ctx, sender.String(), recipient, coins); err != nil {
		return err
	}
	return nil
}

// SendCroToIbcHandler handles `__SeeleSendSeeleToIbc` log
type SendCroToIbcHandler struct {
	bankKeeper  types.BankKeeper
	seeleKeeper Keeper
}

func NewSendCroToIbcHandler(bankKeeper types.BankKeeper, seeleKeeper Keeper) *SendCroToIbcHandler {
	return &SendCroToIbcHandler{
		bankKeeper:  bankKeeper,
		seeleKeeper: seeleKeeper,
	}
}

func (h SendCroToIbcHandler) EventID() common.Hash {
	return SendCroToIbcEvent.ID
}

func (h SendCroToIbcHandler) Handle(ctx sdk.Context, contract common.Address, data []byte) error {
	unpacked, err := SendCroToIbcEvent.Inputs.Unpack(data)
	if err != nil {
		// log and ignore
		h.seeleKeeper.Logger(ctx).Info("log signature matches but failed to decode")
		return nil
	}

	contractAddr := sdk.AccAddress(contract.Bytes())
	sender := sdk.AccAddress(unpacked[0].(common.Address).Bytes())
	recipient := unpacked[1].(string)
	amount := sdk.NewIntFromBigInt(unpacked[2].(*big.Int))
	evmDenom := h.seeleKeeper.GetEvmParams(ctx).EvmDenom
	coins := sdk.NewCoins(sdk.NewCoin(evmDenom, amount))
	// First, transfer IBC coin to user so that he will be the refunded address if transfer fails
	if err = h.bankKeeper.SendCoins(ctx, contractAddr, sender, coins); err != nil {
		return err
	}
	// Initiate IBC transfer from sender account
	if err = h.seeleKeeper.IbcTransferCoins(ctx, sender.String(), recipient, coins); err != nil {
		return err
	}
	return nil
}
