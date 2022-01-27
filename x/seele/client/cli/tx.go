package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/Seele-N/Seele/x/seele/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// this line is used by starport scaffolding # 1

	cmd.AddCommand(CmdUpdateTokenMapping())

	return cmd
}

// NewSubmitTokenMappingChangeProposalTxCmd returns a CLI command handler for creating
// a token mapping change proposal governance transaction.
func NewSubmitTokenMappingChangeProposalTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-mapping-change [denom] [contract]",
		Args:  cobra.ExactArgs(2),
		Short: "Submit a token mapping change proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a token mapping change proposal.

Example:
$ %s tx gov submit-proposal token-mapping-change gravity0x0000...0000 0x0000...0000 --from=<key_or_address>
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(govcli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(govcli.FlagDescription)
			if err != nil {
				return err
			}

			var contract *common.Address
			if len(args[1]) > 0 {
				addr := common.HexToAddress(args[1])
				contract = &addr
			}

			content := types.NewTokenMappingChangeProposal(
				title, description, args[0], contract,
			)

			from := clientCtx.GetFromAddress()

			strDeposit, err := cmd.Flags().GetString(govcli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(strDeposit)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String(govcli.FlagTitle, "", "The proposal title")
	cmd.Flags().String(govcli.FlagDescription, "", "The proposal description")
	cmd.Flags().String(govcli.FlagDeposit, "", "The proposal deposit")

	return cmd
}

// CmdUpdateTokenMapping returns a CLI command handler for update token mapping
func CmdUpdateTokenMapping() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-token-mapping [denom] [contract]",
		Short: "Update token mapping",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateTokenMapping(clientCtx.GetFromAddress().String(), args[0], args[1])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
