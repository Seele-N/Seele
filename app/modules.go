package app

import (
	"encoding/json"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"

	"github.com/cosmos/cosmos-sdk/x/distribution"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/Seele-N/Seele/x/mintx"
	mintxtypes "github.com/Seele-N/Seele/x/mintx/types"

	"github.com/tharsis/ethermint/x/evm"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/peggyjv/gravity-bridge/module/x/gravity"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/x/gravity/types"
)

const (

	// Day time of day
	Day = 24 * time.Hour
	// DefaultBondDenom Default Bond Denom
	DefaultBondDenom = "snp"
	// DefaultMintDenom Default Mint Denom
	DefaultMintDenom = "seele"
)

// staking
const (
	// DefaultUnbondingTime reflects 30 days in seconds as the default
	// unbonding time.
	DefaultUnbondingTime = 30 * Day

	// DefaultMaxValidators Default maximum number of bonded validators
	DefaultMaxValidators uint32 = 21
)

// slashinging
const (
	// DefaultSignedBlocksWindow reflects signed blocks window
	DefaultSignedBlocksWindow = 10000
)

// Crisis
var (
	// DefaultCrisisConstantFee Default Crisis Constant Fee
	DefaultCrisisConstantFee = sdk.NewInt(1000)
)

// gov
var (
	// DefaultGovMinDepositAmount Default Gov Min Deposit Amount
	DefaultGovMinDepositAmount = sdk.NewInt(1000000000000000000)
)

// gov
const (
	// DefaultPeriod Default period for deposits & voting
	DefaultPeriod = 14 * Day
	// VotingPeriod Voting Period
	VotingPeriod = 7 * Day
)

// gravity
const (
	// DefaultPeriod Default period for deposits & voting
	DefaultGravityId             = "seele_bridge"
	DefaultBridgeEthereumAddress = "0xCad5A42d74F66d96650fdf1a1b1d738DeDB7d876"
	DefaultWindowTime            = 3144960
)

// StakingModuleBasic staking module basic replace
type StakingModuleBasic struct {
	staking.AppModuleBasic
}

// DefaultGenesis defaut genesis for extend params
func (StakingModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := stakingtypes.DefaultGenesisState()
	genesisState.Params.UnbondingTime = DefaultUnbondingTime
	genesisState.Params.MaxValidators = DefaultMaxValidators
	genesisState.Params.BondDenom = DefaultBondDenom
	return cdc.MustMarshalJSON(genesisState)
}

// SlashingModuleBasic slashing module basic replace
type SlashingModuleBasic struct {
	slashing.AppModuleBasic
}

// DefaultGenesis defaut genesis for extend params
func (SlashingModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := slashingtypes.DefaultGenesisState()
	genesisState.Params.SignedBlocksWindow = DefaultSignedBlocksWindow
	//genesisState.Params.SlashFractionDowntime = sdk.NewDecWithPrec(1, 2)
	return cdc.MustMarshalJSON(genesisState)
}

// CrisisModuleBasic Crisis Module Basic
type CrisisModuleBasic struct {
	crisis.AppModuleBasic
}

// DefaultGenesis defaut genesis for extend params
func (CrisisModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := crisistypes.DefaultGenesisState()
	genesisState.ConstantFee.Denom = DefaultBondDenom
	genesisState.ConstantFee.Amount = DefaultCrisisConstantFee
	return cdc.MustMarshalJSON(genesisState)
}

// DistributionModuleBasic Crisis Module Basic
type DistributionModuleBasic struct {
	distribution.AppModuleBasic
}

// DefaultGenesis defaut genesis for extend params
func (DistributionModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := distributiontypes.DefaultGenesisState()
	genesisState.Params.CommunityTax = sdk.NewDecWithPrec(2, 2)
	genesisState.Params.BaseProposerReward = sdk.NewDecWithPrec(1, 2)
	genesisState.Params.BonusProposerReward = sdk.NewDecWithPrec(4, 2)
	return cdc.MustMarshalJSON(genesisState)
}

// GovModuleBasic Gov Module Basic
type GovModuleBasic struct {
	gov.AppModuleBasic
}

// DefaultGenesis defaut genesis for extend params
func (GovModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := govtypes.DefaultGenesisState()
	genesisState.DepositParams.MinDeposit[0].Denom = DefaultBondDenom
	genesisState.DepositParams.MinDeposit[0].Amount = DefaultGovMinDepositAmount
	genesisState.DepositParams.MaxDepositPeriod = DefaultPeriod
	genesisState.VotingParams.VotingPeriod = VotingPeriod
	genesisState.TallyParams = govtypes.TallyParams{
		Quorum:        sdk.NewDecWithPrec(4, 1),   //  Minimum percentage of total stake needed to vote for a result to be considered valid
		Threshold:     sdk.NewDecWithPrec(5, 1),   // Minimum proportion of Yes votes for proposal to pass. Default value: 0.5.
		VetoThreshold: sdk.NewDecWithPrec(334, 3), //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Default value: 1/3.
	}
	return cdc.MustMarshalJSON(genesisState)
}

// EvmModuleBasic evm module basic replace
type EvmModuleBasic struct {
	evm.AppModuleBasic
}

// DefaultGenesis defaut genesis for extend params
func (EvmModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := evmtypes.DefaultGenesisState()
	genesisState.Params.EvmDenom = DefaultMintDenom
	eips := []int64{2929, 2200, 1884, 1344}
	genesisState.Params.ExtraEIPs = eips
	return cdc.MustMarshalJSON(genesisState)
}

// GravityModuleBasic gravity bridge module basic replace
type GravityModuleBasic struct {
	gravity.AppModuleBasic
}

// DefaultGenesis defaut genesis for extend params
func (GravityModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := gravitytypes.DefaultGenesisState()
	genesisState.Params.GravityId = DefaultGravityId
	genesisState.Params.BridgeEthereumAddress = DefaultBridgeEthereumAddress
	genesisState.Params.SignedBatchesWindow = DefaultWindowTime
	genesisState.Params.UnbondSlashingSignerSetTxsWindow = DefaultWindowTime
	var items []*gravitytypes.ERC20ToDenom
	item := &gravitytypes.ERC20ToDenom{
		Erc20: "0x795dBF627484F8248D3d6c09c309825c1563E873",
		Denom: "snp",
	}
	items = append(items, item)
	item = &gravitytypes.ERC20ToDenom{
		Erc20: "0xB1e93236ab6073fdAC58adA5564897177D4bcC43",
		Denom: "seele",
	}
	items = append(items, item)
	genesisState.Erc20ToDenoms = items
	//genesisState.Params.BridgeChainId
	return cdc.MustMarshalJSON(genesisState)
}

// MintxModuleBasic Mint Module Basic
type MintxModuleBasic struct {
	mintx.AppModuleBasic
}

// DefaultGenesis defaut genesis for extend params
func (MintxModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := mintxtypes.DefaultGenesisState()
	genesisState.Minter.HeightAdjustment = 0
	genesisState.Params.MintDenom = DefaultMintDenom
	genesisState.Params.DefaultRewardPerBlock = sdk.NewDec(20)
	mintplans := []mintxtypes.MintPlan{}
	plan1 := mintxtypes.MintPlan{
		StartHeight:    0,
		EndHeight:      10512000,
		RewardPerBlock: sdk.NewDec(200),
	}
	plan2 := mintxtypes.MintPlan{
		StartHeight:    10512000,
		EndHeight:      21024000,
		RewardPerBlock: sdk.NewDec(160),
	}
	plan3 := mintxtypes.MintPlan{
		StartHeight:    21024000,
		EndHeight:      31536000,
		RewardPerBlock: sdk.NewDec(120),
	}

	plan4 := mintxtypes.MintPlan{
		StartHeight:    31536000,
		EndHeight:      42048000,
		RewardPerBlock: sdk.NewDec(80),
	}

	plan5 := mintxtypes.MintPlan{
		StartHeight:    42048000,
		EndHeight:      52560000,
		RewardPerBlock: sdk.NewDec(40),
	}

	plan6 := mintxtypes.MintPlan{
		StartHeight:    52560000,
		EndHeight:      63072000,
		RewardPerBlock: sdk.NewDec(20),
	}

	plan7 := mintxtypes.MintPlan{
		StartHeight:    63072000,
		EndHeight:      ^uint64(0),
		RewardPerBlock: sdk.NewDec(20),
	}
	mintplans = append(mintplans, plan1)
	mintplans = append(mintplans, plan2)
	mintplans = append(mintplans, plan3)
	mintplans = append(mintplans, plan4)
	mintplans = append(mintplans, plan5)
	mintplans = append(mintplans, plan6)
	mintplans = append(mintplans, plan7)

	genesisState.Params.MintPlans = mintplans

	return cdc.MustMarshalJSON(genesisState)
}
