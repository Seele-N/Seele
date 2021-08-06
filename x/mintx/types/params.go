package types

import (
	"errors"
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyMintDenom = []byte("MintDenom")
	KeyMintPlans = []byte("MintPlans")
)

// ParamTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewMintPlan create a new MintPlan object
func NewMintPlan(startHeight, endHeight uint64, rewardPerBlock sdk.Dec) MintPlan {
	return MintPlan{
		StartHeight:    startHeight,
		EndHeight:      endHeight,
		RewardPerBlock: rewardPerBlock,
	}
}

// NewParams create a new Params object
func NewParams(
	mintDenom string, mintPlan []MintPlan) Params {
	return Params{
		MintDenom: mintDenom,
		MintPlans: mintPlan,
	}
}

// DefaultParams default minting module parameters
func DefaultParams() Params {
	return Params{
		MintDenom:             sdk.DefaultBondDenom,
		DefaultRewardPerBlock: sdk.OneDec(),
		MintPlans:             []MintPlan{},
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}

	if err := validateMintPlan(p.MintPlans); err != nil {
		return err
	}

	if p.DefaultRewardPerBlock.LTE(sdk.ZeroDec()) {
		return fmt.Errorf("reward per block must be positive: %d", p.DefaultRewardPerBlock)
	}

	return nil

}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		paramtypes.NewParamSetPair(KeyMintPlans, &p.MintPlans, validateMintPlan),
	}
}

// GetRewardByHeight return reward by block height
func (p *Params) GetRewardByHeight(height uint64) sdk.Dec {
	for _, value := range p.MintPlans {
		if value.StartHeight < height && value.EndHeight >= height {
			return value.RewardPerBlock
		}
	}
	return p.DefaultRewardPerBlock
	//return sdk.OneDec()
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateMintPlan(i interface{}) error {
	v, ok := i.([]MintPlan)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, value := range v {
		if value.EndHeight < value.StartHeight {
			return fmt.Errorf("end height:%d must great start height: %d", value.EndHeight, value.StartHeight)
		}
		if value.RewardPerBlock.LTE(sdk.ZeroDec()) {
			return fmt.Errorf("reward per block must be positive: %d", value.RewardPerBlock)
		}
	}

	return nil
}
