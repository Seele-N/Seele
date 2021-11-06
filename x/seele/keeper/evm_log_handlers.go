package keeper

import (
	//"fmt"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Seele-N/Seele/x/seele/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	_ types.EvmLogHandler = SendSnpStakeHandler{}
	_ types.EvmLogHandler = SendUnSnpStakeHandler{}
	_ types.EvmLogHandler = SendSnpClaimRewardHandler{}
	_ types.EvmLogHandler = SendSnpClaimCommissionHandler{}
	_ types.EvmLogHandler = SendReSnpStakeHandler{}
)

const (
	SnpStakingEventName         = "Snp_Staking"
	SnpUnStakingEventName       = "Snp_UnStaking"
	SnpClaimRewardEventName     = "Snp_ClaimReward"
	SnpClaimCommissionEventName = "Snp_ClaimCommission"
	SnpReStakingEventName       = "Snp_ReStaking"
)

var (

	// SnpStakeEvent represent the signature of
	// `event Snp_Staking(address validator, address delegator,uint256 amount)`
	SnpStakeEvent abi.Event

	// SnpUnStakeEvent represent the signature of
	// `event Snp_UnStaking(address validator, address delegator,uint256 amount)`
	SnpUnStakeEvent abi.Event

	// SnpClaimRewardEvent represent the signature of
	// `event Snp_ClaimReward(address validator, address delegator)`
	SnpClaimRewardEvent abi.Event

	// SnpReStakeEvent represent the signature of
	// `event Snp_ReStaking(address srcVal,address destVal,address delegator,uint256 amount)`
	SnpReStakeEvent abi.Event

	// SnpClaimCommissionEvent represent the signature of
	// `event Snp_ClaimCommission(address validator)`
	SnpClaimCommissionEvent abi.Event
)

func init() {
	addressType, _ := abi.NewType("address", "", nil)
	uint256Type, _ := abi.NewType("uint256", "", nil)
	//stringType, _ := abi.NewType("string", "", nil)

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

	SnpUnStakeEvent = abi.NewEvent(
		SnpUnStakingEventName,
		SnpUnStakingEventName,
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

	SnpClaimRewardEvent = abi.NewEvent(
		SnpClaimRewardEventName,
		SnpClaimRewardEventName,
		false,
		abi.Arguments{abi.Argument{
			Name:    "validator",
			Type:    addressType,
			Indexed: false,
		}, abi.Argument{
			Name:    "delegator",
			Type:    addressType,
			Indexed: false,
		}},
	)

	SnpClaimCommissionEvent = abi.NewEvent(
		SnpClaimCommissionEventName,
		SnpClaimCommissionEventName,
		false,
		abi.Arguments{abi.Argument{
			Name:    "validator",
			Type:    addressType,
			Indexed: false,
		}},
	)

	SnpReStakeEvent = abi.NewEvent(
		SnpReStakingEventName,
		SnpReStakingEventName,
		false,
		abi.Arguments{abi.Argument{
			Name:    "srcVal",
			Type:    addressType,
			Indexed: false,
		}, abi.Argument{
			Name:    "destVal",
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
	addr, _ := h.seeleKeeper.getContractByname(ctx, types.SnpDelegateContract.ContractName)
	if addr != contract {
		h.seeleKeeper.Logger(ctx).Error("contract address is not SnpDelegateContract", "error")
		return types.ErrContractAddressInvalid
	}
	unpacked, err := SnpStakeEvent.Inputs.Unpack(data)
	if err != nil {
		// log and ignore
		h.seeleKeeper.Logger(ctx).Error("log signature matches but failed to decode", "error", err)
		return err
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

// SendUnSnpStakeHandler handles `Snp_Staking` log
type SendUnSnpStakeHandler struct {
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
	seeleKeeper   Keeper
}

func NewSendUnSnpStakeHandler(bankKeeper types.BankKeeper, stakingKeeper types.StakingKeeper, seeleKeeper Keeper) *SendUnSnpStakeHandler {
	return &SendUnSnpStakeHandler{
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		seeleKeeper:   seeleKeeper,
	}
}

func (h SendUnSnpStakeHandler) EventID() common.Hash {
	return SnpUnStakeEvent.ID
}

func (h SendUnSnpStakeHandler) Handle(ctx sdk.Context, contract common.Address, data []byte) error {
	h.seeleKeeper.Logger(ctx).Info("SendUnSnpStakeHandler")
	addr, _ := h.seeleKeeper.getContractByname(ctx, types.SnpDelegateContract.ContractName)
	if addr != contract {
		h.seeleKeeper.Logger(ctx).Error("contract address is not SnpDelegateContract", "error")
		return types.ErrContractAddressInvalid
	}
	unpacked, err := SnpUnStakeEvent.Inputs.Unpack(data)
	if err != nil {
		// log and ignore
		h.seeleKeeper.Logger(ctx).Error("log signature matches but failed to decode", "error", err)
		return err
	}
	h.seeleKeeper.Logger(ctx).Info("Event from contract:" + contract.Hex())
	h.seeleKeeper.Logger(ctx).Info("validator:" + unpacked[0].(common.Address).Hex())
	h.seeleKeeper.Logger(ctx).Info("delegator:" + unpacked[1].(common.Address).Hex())
	h.seeleKeeper.Logger(ctx).Info("amount:" + unpacked[2].(*big.Int).String())

	amount := unpacked[2].(*big.Int)
	valAddress := sdk.ValAddress(unpacked[0].(common.Address).Bytes())
	delegator := sdk.AccAddress(unpacked[1].(common.Address).Bytes())
	h.seeleKeeper.Logger(ctx).Info("valAddress:" + valAddress.String())
	shares, err := h.stakingKeeper.ValidateUnbondAmount(ctx, delegator, valAddress, sdk.NewIntFromBigInt(amount))
	if err != nil {
		return err
	}

	completionTime, err := h.stakingKeeper.Undelegate(ctx, delegator, valAddress, shares)
	if err != nil {
		return err
	}

	//coin := sdk.NewCoin("snp", sdk.NewIntFromBigInt(amount))
	//h.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
	//h.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegator, sdk.NewCoins(coin))

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingtypes.EventTypeUnbond,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, valAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
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

// SendSnpClaimRewardHandler handles `Snp_ClaimReward` log
type SendSnpClaimRewardHandler struct {
	bankKeeper         types.BankKeeper
	distributionKeeper types.DistributionKeeper
	seeleKeeper        Keeper
}

func NewSendSnpClaimRewardHandler(bankKeeper types.BankKeeper, distributionKeeper types.DistributionKeeper, seeleKeeper Keeper) *SendSnpClaimRewardHandler {
	return &SendSnpClaimRewardHandler{
		bankKeeper:         bankKeeper,
		distributionKeeper: distributionKeeper,
		seeleKeeper:        seeleKeeper,
	}
}

func (h SendSnpClaimRewardHandler) EventID() common.Hash {
	return SnpClaimRewardEvent.ID
}

func (h SendSnpClaimRewardHandler) Handle(ctx sdk.Context, contract common.Address, data []byte) error {
	h.seeleKeeper.Logger(ctx).Info("SendSnpClaimRewardHandler")
	addr, _ := h.seeleKeeper.getContractByname(ctx, types.SnpDelegateContract.ContractName)
	if addr != contract {
		h.seeleKeeper.Logger(ctx).Error("contract address is not SnpDelegateContract", "error")
		return types.ErrContractAddressInvalid
	}
	unpacked, err := SnpClaimRewardEvent.Inputs.Unpack(data)
	if err != nil {
		// log and ignore
		h.seeleKeeper.Logger(ctx).Error("log signature matches but failed to decode", "error", err)
		return err
	}
	h.seeleKeeper.Logger(ctx).Info("Event from contract:" + contract.Hex())
	h.seeleKeeper.Logger(ctx).Info("validator:" + unpacked[0].(common.Address).Hex())
	h.seeleKeeper.Logger(ctx).Info("delegator:" + unpacked[1].(common.Address).Hex())

	valAddress := sdk.ValAddress(unpacked[0].(common.Address).Bytes())
	delegator := sdk.AccAddress(unpacked[1].(common.Address).Bytes())
	h.seeleKeeper.Logger(ctx).Info("valAddress:" + valAddress.String())
	_, err = h.distributionKeeper.WithdrawDelegationRewards(ctx, delegator, valAddress)

	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, delegator.String()),
		),
	)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// SendSnpClaimCommissionHandler handles `Snp_ClaimCommission` log
type SendSnpClaimCommissionHandler struct {
	bankKeeper         types.BankKeeper
	distributionKeeper types.DistributionKeeper
	seeleKeeper        Keeper
}

func NewSendSnpClaimCommissionHandler(bankKeeper types.BankKeeper, distributionKeeper types.DistributionKeeper, seeleKeeper Keeper) *SendSnpClaimCommissionHandler {
	return &SendSnpClaimCommissionHandler{
		bankKeeper:         bankKeeper,
		distributionKeeper: distributionKeeper,
		seeleKeeper:        seeleKeeper,
	}
}

func (h SendSnpClaimCommissionHandler) EventID() common.Hash {
	return SnpClaimCommissionEvent.ID
}

func (h SendSnpClaimCommissionHandler) Handle(ctx sdk.Context, contract common.Address, data []byte) error {
	h.seeleKeeper.Logger(ctx).Info("SendSnpClaimCommissionHandler")
	addr, _ := h.seeleKeeper.getContractByname(ctx, types.SnpDelegateContract.ContractName)
	if addr != contract {
		h.seeleKeeper.Logger(ctx).Error("contract address is not SnpDelegateContract", "error")
		return types.ErrContractAddressInvalid
	}
	unpacked, err := SnpClaimCommissionEvent.Inputs.Unpack(data)
	if err != nil {
		// log and ignore
		h.seeleKeeper.Logger(ctx).Error("log signature matches but failed to decode", "error", err)
		return err
	}
	h.seeleKeeper.Logger(ctx).Info("Event from contract:" + contract.Hex())
	h.seeleKeeper.Logger(ctx).Info("validator:" + unpacked[0].(common.Address).Hex())

	valAddress := sdk.ValAddress(unpacked[0].(common.Address).Bytes())
	h.seeleKeeper.Logger(ctx).Info("valAddress:" + valAddress.String())
	_, err = h.distributionKeeper.WithdrawValidatorCommission(ctx, valAddress)

	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, valAddress.String()),
		),
	)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// SendReSnpStakeHandler handles `Snp_ReStaking` log
type SendReSnpStakeHandler struct {
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
	seeleKeeper   Keeper
}

func NewSendReSnpStakeHandler(bankKeeper types.BankKeeper, stakingKeeper types.StakingKeeper, seeleKeeper Keeper) *SendReSnpStakeHandler {
	return &SendReSnpStakeHandler{
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		seeleKeeper:   seeleKeeper,
	}
}

func (h SendReSnpStakeHandler) EventID() common.Hash {
	return SnpReStakeEvent.ID
}

func (h SendReSnpStakeHandler) Handle(ctx sdk.Context, contract common.Address, data []byte) error {
	h.seeleKeeper.Logger(ctx).Info("SendReSnpStakeHandler")
	addr, _ := h.seeleKeeper.getContractByname(ctx, types.SnpDelegateContract.ContractName)
	if addr != contract {
		h.seeleKeeper.Logger(ctx).Error("contract address is not SnpDelegateContract", "error")
		return types.ErrContractAddressInvalid
	}
	unpacked, err := SnpReStakeEvent.Inputs.Unpack(data)
	if err != nil {
		// log and ignore
		h.seeleKeeper.Logger(ctx).Error("log signature matches but failed to decode", "error", err)
		return err
	}
	h.seeleKeeper.Logger(ctx).Info("Event from contract:" + contract.Hex())
	h.seeleKeeper.Logger(ctx).Info("src validator:" + unpacked[0].(common.Address).Hex())
	h.seeleKeeper.Logger(ctx).Info("dest validator:" + unpacked[1].(common.Address).Hex())
	h.seeleKeeper.Logger(ctx).Info("delegator:" + unpacked[2].(common.Address).Hex())
	h.seeleKeeper.Logger(ctx).Info("amount:" + unpacked[3].(*big.Int).String())

	amount := unpacked[3].(*big.Int)
	srcvalAddress := sdk.ValAddress(unpacked[0].(common.Address).Bytes())
	destvalAddress := sdk.ValAddress(unpacked[1].(common.Address).Bytes())
	delegator := sdk.AccAddress(unpacked[2].(common.Address).Bytes())
	//h.seeleKeeper.Logger(ctx).Info("valAddress:" + valAddress.String())
	shares, err := h.stakingKeeper.ValidateUnbondAmount(ctx, delegator, srcvalAddress, sdk.NewIntFromBigInt(amount))
	if err != nil {
		return err
	}

	completionTime, err := h.stakingKeeper.BeginRedelegation(ctx, delegator, srcvalAddress, destvalAddress, shares)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingtypes.EventTypeRedelegate,
			sdk.NewAttribute(stakingtypes.AttributeKeySrcValidator, srcvalAddress.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyDstValidator, destvalAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, delegator.String()),
		),
	})

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
