package politeiawww

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

//GetTokenInventory retrieves the censorship record tokens of all proposals in the inventory.
//The tokens are categorized by stage of the voting process and sorted according to the rules listed below.
//Unvetted proposal tokens are only returned to admins.
//Unvetted proposals include unvreviewed and censored proposals.
//Sorted by record timestamp in descending order: Pre, Abandonded, Unreviewed, Censored
//Sorted by voting period end block height in descending order: Active, Approved, Rejected
func (c *Client) GetTokenInventory(ctx context.Context) (*TokenInventory, error) {
	inv := new(TokenInventory)
	if _, err := c.makeRequest(ctx, http.MethodGet, tokenInventoryPath, nil, nil, inv); err != nil {
		return nil, err
	}

	return inv, nil
}

func marshalBatchRequest(ctokens []string) ([]byte, error) {
	b, err := json.Marshal(&BatchProposalsRequest{ctokens})
	if err != nil {
		return nil, fmt.Errorf("error marshaling batch request: err")
	}
	return b, nil
}

//GetProposalBatch retrieves the proposal details for a list of proposals. This route wil not return the proposal files.
//The number of proposals that may be requested is limited by the ProposalListPageSize property,
//which is provided via Policy
func (c *Client) GetProposalBatch(ctx context.Context, ctokens []string) ([]Proposal, error) {
	b, err := marshalBatchRequest(ctokens)
	if err != nil {
		return nil, err
	}

	props := new(Proposals)
	if _, err := c.makeRequest(ctx, http.MethodPost, batchProposalsPath, nil, b, props); err != nil {
		return nil, err
	}

	return props.Proposals, err
}

//GetVoteSummaryBatch retrieves the vote summaries for a list of proposals. The number of vote summaries that may be
//requested is limited by the ProposalListPageSize property, which is provided via Policy.
func (c *Client) GetVoteSummaryBatch(ctx context.Context, ctokens []string) (*BatchVoteSummaryResponse, error) {
	b, err := marshalBatchRequest(ctokens)
	if err != nil {
		return nil, err
	}

	vsums := new(BatchVoteSummaryResponse)
	if _, err := c.makeRequest(ctx, http.MethodPost, batchProposalsPath, nil, b, vsums); err != nil {
		return nil, err
	}

	return vsums, err
}

//GetVetted retrieves a page of vetted proposals; the number of proposals returned in the page is limited
//by the ProposalListPageSize property, which is provided via Policy.
//Before/After are censorship tokens.
//If before is provided, the page of proposals returned will end right before the proposal whose token is provided,
//when sorted in reverse chronological order.
//This parameter should be and empty string if after is set.
//If after is provided,  the page of proposals returned will begin right after the proposal whose token is provided,
//when sorted in reverse chronological order. This parameter should not be specified if before is set.
//This parameter should be and empty string if before is set.
func (c *Client) GetVetted(ctx context.Context, before, after string) ([]Proposal, error) {
	qs := make(map[string]string)
	if len(after) > 0 {
		qs["after"] = after
	}
	if len(before) > 0 {
		qs["before"] = after
	}

	props := new(Proposals)
	if _, err := c.makeRequest(ctx, http.MethodGet, vettedProposalsPath, qs, nil, props); err != nil {
		return nil, err
	}

	return props.Proposals, nil
}

// GetProposalDetails Retrieve proposal and its details. This request can be made with the full censorship token
//or its 7 character prefix.
func (c *Client) GetProposalDetails(ctx context.Context, censorshipToken, version string) (*Proposal, error) {
	var qs map[string]string
	if version != "" {
		qs = map[string]string{
			"version": version,
		}
	}

	prop := new(Proposal)
	if _, err := c.makeRequest(ctx, http.MethodGet, fmt.Sprintf(proposalDetailsPath, censorshipToken), qs, nil, prop); err != nil {
		return nil, err
	}
	return prop, nil
}

//GetVoteStatus returns the vote status for a single public proposal.
//This route deprecated by Batch Vote Status.
func (c *Client) GetVoteStatus(ctx context.Context, ctoken string) (*VoteStatus, error) {
	vs := new(VoteStatus)
	if _, err := c.makeRequest(ctx, http.MethodGet, fmt.Sprintf(voteStatusPath, ctoken), nil, nil, vs); err != nil {
		return nil, err
	}
	return vs, nil
}
