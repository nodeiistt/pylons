package handlers

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Pylons-tech/pylons/x/pylons/keep"
	"github.com/Pylons-tech/pylons/x/pylons/msgs"
	"github.com/Pylons-tech/pylons/x/pylons/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestHandlerMsgDisableRecipe(t *testing.T) {
	tci := keep.SetupTestCoinInput()
	sender1, sender2 := keep.SetupTestAccounts(t, tci, types.NewPylon(1000000))

	// mock cookbook
	cbData := MockCookbook(tci, sender1)

	// mock recipe
	rcpData := MockPopularRecipe(RcpDefault, tci, "existing recipe", cbData.CookbookID, sender1)

	cases := map[string]struct {
		rcpID        string
		sender       sdk.AccAddress
		desiredError string
		showError    bool
	}{
		"wrong recipe check": {
			rcpID:        "invalidRecipeID",
			sender:       sender1,
			desiredError: "The recipe doesn't exist",
			showError:    true,
		},
		"owner of recipe check": {
			rcpID:        rcpData.RecipeID,
			sender:       sender2,
			desiredError: "msg sender is not the owner of the recipe",
			showError:    true,
		},
		"successful update check": {
			rcpID:        rcpData.RecipeID,
			sender:       sender1,
			desiredError: "",
			showError:    false,
		},
	}
	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			msg := msgs.NewMsgDisableRecipe(tc.rcpID, tc.sender)
			result, err := HandlerMsgDisableRecipe(tci.Ctx, tci.PlnK, msg)

			if tc.showError == false {
				disableRcpResponse := DisableRecipeResponse{}
				err := json.Unmarshal(result.Data, &disableRcpResponse)

				require.True(t, err == nil)
				require.True(t, disableRcpResponse.Status == "Success")
				require.True(t, disableRcpResponse.Message == "successfully disabled the recipe")

				uRcp, err := tci.PlnK.GetRecipe(tci.Ctx, tc.rcpID)
				// t.Errorf("DisableRecipeTEST LOG:: %+v", uRcp)
				require.True(t, err == nil)
				require.True(t, uRcp.Disabled == true)
			} else {
				require.True(t, strings.Contains(err.Error(), tc.desiredError))
			}
		})
	}
}
