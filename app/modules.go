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
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

// MintModuleBasic Mint Module Basic
type MintModuleBasic struct {
	mint.AppModuleBasic
}

// DefaultGenesis defaut genesis for extend params
func (MintModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	genesisState := minttypes.DefaultGenesisState()
	genesisState.Minter.Inflation = sdk.NewDecWithPrec(13, 2)       // current annual inflation rate
	genesisState.Minter.AnnualProvisions = sdk.NewDecWithPrec(0, 1) // current annual expected provisions
	genesisState.Params.MintDenom = appparams.DefaultMintDenom
	genesisState.Params.InflationRateChange = sdk.NewDecWithPrec(13, 2) // maximum annual change in inflation rate
	genesisState.Params.InflationMax = sdk.NewDecWithPrec(2, 1)         // maximum inflation rate
	genesisState.Params.InflationMin = sdk.NewDecWithPrec(7, 2)         // minimum inflation rate
	genesisState.Params.GoalBonded = sdk.NewDecWithPrec(67, 2)          // goal of percent bonded snps
	genesisState.Params.BlocksPerYear = appparams.BlocksPerYear         // expected blocks per year
	return cdc.MustMarshalJSON(genesisState)
}
