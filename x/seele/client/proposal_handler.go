package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/Seele-N/Seele/x/seele/client/cli"
	"github.com/Seele-N/Seele/x/seele/client/rest"
)

// ProposalHandler is the token mapping change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewSubmitTokenMappingChangeProposalTxCmd, rest.ProposalRESTHandler)
