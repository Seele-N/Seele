package app

import (
	"encoding/json"

	appparams "github.com/seele-n/seele/app/params"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/seele-n/seele/x/mintx"
	mintxtypes "github.com/seele-n/seele/x/mintx/types"
)

// StakingModuleBasic staking module basic replace
type StakingModuleBasic struct {
	staking.AppModuleBasic
}

// DefaultGenesis defaut genesis for extend params
func (StakingModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	genesisState := stakingtypes.DefaultGenesisState()
	genesisState.Params.UnbondingTime = appparams.DefaultUnbondingTime
	genesisState.Params.MaxValidators = appparams.DefaultMaxValidators
	genesisState.Params.BondDenom = appparams.DefaultBondDenom
	return cdc.MustMarshalJSON(genesisState)
}

// CrisisModuleBasic Crisis Module Basic
type CrisisModuleBasic struct {
	crisis.AppModuleBasic
}

// DefaultGenesis defaut genesis for extend params
func (CrisisModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	genesisState := crisistypes.DefaultGenesisState()
	genesisState.ConstantFee.Denom = appparams.DefaultBondDenom
	genesisState.ConstantFee.Amount = appparams.DefaultCrisisConstantFee
	return cdc.MustMarshalJSON(genesisState)
}

// GovModuleBasic Gov Module Basic
type GovModuleBasic struct {
	gov.AppModuleBasic
}

// DefaultGenesis defaut genesis for extend params
func (GovModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	genesisState := govtypes.DefaultGenesisState()
	genesisState.DepositParams.MinDeposit[0].Denom = appparams.DefaultBondDenom
	genesisState.DepositParams.MinDeposit[0].Amount = appparams.DefaultGovMinDepositAmount
	genesisState.DepositParams.MaxDepositPeriod = appparams.DefaultPeriod
	genesisState.VotingParams.VotingPeriod = appparams.VotingPeriod
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
func (MintxModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	genesisState := mintxtypes.DefaultGenesisState()
	genesisState.Minter.HeightAdjustment = 0
	genesisState.Params.MintDenom = appparams.DefaultMintDenom
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
