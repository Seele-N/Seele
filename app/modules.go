package app

import (
	"encoding/json"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/Seele-N/Seele/x/mintx"
	mintxtypes "github.com/Seele-N/Seele/x/mintx/types"
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
	// DefaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.
	DefaultUnbondingTime = 21 * Day

	// DefaultMaxValidators Default maximum number of bonded validators
	DefaultMaxValidators uint32 = 21
)

// Crisis
var (
	// DefaultCrisisConstantFee Default Crisis Constant Fee
	DefaultCrisisConstantFee = sdk.NewInt(1000)
)

// gov
var (
	// DefaultGovMinDepositAmount Default Gov Min Deposit Amount
	DefaultGovMinDepositAmount = sdk.NewInt(10000000)
)

// gov
const (
	// DefaultPeriod Default period for deposits & voting
	DefaultPeriod = 14 * Day
	// VotingPeriod Voting Period
	VotingPeriod = 7 * Day
)

// mint
const (
	BlocksPerYear = 6311520
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
