package dcrlibwallet

import (
	"context"
	"fmt"
	"time"

	pwww "github.com/planetdecred/dcrlibwallet/politeiawww"
)

//proposal categories
const (
	PreVote = iota
	Active
	Approved
	Rejected
	Abandoned
)

var invTypesToHasVoteSummary map[int]bool = map[int]bool{
	0: false,
	1: true,
	2: true,
	3: true,
	4: false,
}

//Politeia handles the loading of proposials on
//proposals.decred.org via a politeiawww client
//This politeia struct handles the management of
//token inventory for lazy loading.
type Politeia struct {
	inventory *pwww.TokenInventory
	client    pwww.Client
	ctx       context.Context
	lcount    [5]int //loaded count for each token category
	version   *pwww.ServerVersion
	policy    *pwww.ServerPolicy
}

//NewPoliteia returns a new Politeia type
func NewPoliteia(timeoutSeconds int64) *Politeia {
	return &Politeia{
		client: pwww.NewClient(time.Duration(timeoutSeconds * int64(time.Second))),
		ctx:    context.Background(),
	}
}

//PoliteiaProposal contains a specific proposals's details
//and vote summary if applicable. If a proposal is Pre-vote
//or abandoned, VoteSummary will be nil
type PoliteiaProposal struct {
	Details     *pwww.Proposal
	VoteSummary *pwww.VoteSummary
}

type Proposal struct {
	Name             string
	State            int
	Status           int
	Timestamp        int64
	Username         string
	PublicKey        string
	Signature        string
	NumComments      int
	Version          string
	PublishedAt      int64
	Files            *ProposalFiles
	MetaData         *ProposalMetadataIterator
	CensorshipRecord *pwww.ProposalCensorshipRecord
	VoteStatus       *pwww.VoteStatus
}

type ProposalFiles struct {
	files        []pwww.ProposalFile
	currentIndex int
}

func (fs *ProposalFiles) Next() *pwww.ProposalFile {
	if fs.currentIndex < len(fs.files) {
		pf := fs.files[fs.currentIndex]
		fs.currentIndex++
		return &pf
	}

	return nil
}

//Reset resets the current index to 0
func (fs *ProposalFiles) Reset() {
	fs.currentIndex = 0
}

type ProposalMetadataIterator struct {
	metadata     []pwww.ProposalMetaData
	currentIndex int
}

type PoliteiaProposalMobile struct {
	Details     *Proposal
	VoteSummary *pwww.VoteSummary
}

func (pm *ProposalMetadataIterator) Next() *pwww.ProposalMetaData {
	if pm.currentIndex < len(pm.metadata) {
		md := pm.metadata[pm.currentIndex]
		pm.currentIndex++
		return &md
	}

	return nil
}

//Reset resets the current index to 0
func (pm *ProposalMetadataIterator) Reset() {
	pm.currentIndex = 0
}

func (p *Politeia) getVersion() error {
	if !p.client.GotVersion() {
		fmt.Println("setting version")
		v, err := p.client.GetVersion(p.ctx)
		if err != nil {
			return err
		}
		p.version = v
	}
	return nil
}

func (p *Politeia) getInventory() (*pwww.TokenInventory, error) {
	if p.inventory != nil {
		return p.inventory, nil
	}

	if err := p.getVersion(); err != nil {
		return nil, err
	}

	fmt.Println("setting inventory")
	inv, err := p.client.GetTokenInventory(p.ctx)
	if err != nil {
		return nil, err
	}
	p.inventory = inv
	return p.inventory, nil
}

//PreVoteCount returns the number of proposals in the pre-vote
//stage in the pi inventory
func (p *Politeia) PreVoteCount() (int, error) {
	inv, err := p.getInventory()
	if err != nil {
		return -1, err
	}
	return len(inv.Pre), nil
}

//ActiveCount returns the number of proposals in the active
//voting stage in the pi inventory
func (p *Politeia) ActiveCount() (int, error) {
	inv, err := p.getInventory()
	if err != nil {
		return -1, err
	}
	return len(inv.Active), nil
}

//ApprovedCount returns the number of approved proposals
//in the pi inventory
func (p *Politeia) ApprovedCount() (int, error) {
	inv, err := p.getInventory()
	if err != nil {
		return -1, err
	}
	return len(inv.Approved), nil
}

//RejectedCount returns the number of rejected proposals
//in the pi inventory
func (p *Politeia) RejectedCount() (int, error) {
	inv, err := p.getInventory()
	if err != nil {
		return -1, err
	}
	return len(inv.Rejected), nil
}

//AbandonedCount returns the number of approved proposals
//in the pi inventory
func (p *Politeia) AbandonedCount() (int, error) {
	inv, err := p.getInventory()
	if err != nil {
		return -1, err
	}
	return len(inv.Abandoned), nil
}

type propResponse struct {
	Prop *pwww.Proposal
	Err  error
}

type propsResponse struct {
	Props []*pwww.Proposal
	Err   error
}

type voteSummaryResponse struct {
	Vsum *pwww.BatchVoteSummaryResponse
	Err  error
}

func (p *Politeia) getVoteSummaryResp(ctx context.Context, ctokens []string) voteSummaryResponse {
	vs, err := p.client.GetVoteSummaryBatch(ctx, ctokens)
	return voteSummaryResponse{vs, err}
}

func (p *Politeia) loadProposal(ctx context.Context, ctoken string, voteSummary bool) (*PoliteiaProposal, error) {
	pchan := make(chan propResponse, 1)
	vschan := make(chan voteSummaryResponse, 1)
	go func() {
		p, err := p.client.GetProposalDetails(ctx, ctoken, "")
		pchan <- propResponse{Prop: p, Err: err}
	}()
	go func() {
		var vs voteSummaryResponse
		if voteSummary {
			vs = p.getVoteSummaryResp(ctx, []string{ctoken})
		}
		vschan <- vs
	}()

	prop := new(PoliteiaProposal)
	propresp := <-pchan
	if propresp.Err != nil {
		return nil, propresp.Err
	}
	prop.Details = propresp.Prop

	vsresp := <-vschan
	if vsresp.Err != nil {
		return nil, vsresp.Err
	}

	prop.VoteSummary = vsresp.Vsum.Summaries[prop.Details.CensorshipRecord.Token]
	return prop, nil
}

func (p *Politeia) loadProposals(ctx context.Context, ctokens []string, voteSummary bool) ([]PoliteiaProposal, error) {
	pchan := make(chan propsResponse, 1)
	vschan := make(chan voteSummaryResponse, 1)
	go func() {
		p, err := p.client.GetProposalBatch(ctx, ctokens)
		pchan <- propsResponse{Props: p, Err: err}
	}()
	go func() {
		var vs voteSummaryResponse
		if voteSummary {
			vs = p.getVoteSummaryResp(ctx, ctokens)
		}
		vschan <- vs
	}()

	propresp := <-pchan
	if propresp.Err != nil {
		return nil, propresp.Err
	}
	vsresp := <-vschan
	if vsresp.Err != nil {
		return nil, vsresp.Err
	}
	pprops := make([]PoliteiaProposal, 0)
	for _, prop := range propresp.Props {
		pp := PoliteiaProposal{Details: prop}
		if voteSummary {
			if v, ok := vsresp.Vsum.Summaries[prop.CensorshipRecord.Token]; ok {
				pp.VoteSummary = v
			}
		}
		pprops = append(pprops, pp)
	}
	return pprops, nil
}

func mapProposalToMobile(prop *PoliteiaProposal) *PoliteiaProposalMobile {
	if prop == nil {
		return nil
	}
	d := prop.Details
	return &PoliteiaProposalMobile{
		Details: &Proposal{
			Name:             d.Name,
			State:            d.State,
			Status:           d.Status,
			Timestamp:        d.Timestamp,
			Username:         d.Username,
			PublicKey:        d.PublicKey,
			Signature:        d.Signature,
			NumComments:      d.NumComments,
			Version:          d.Version,
			PublishedAt:      d.PublishedAt,
			Files:            &ProposalFiles{d.Files, 0},
			MetaData:         &ProposalMetadataIterator{d.MetaData, 0},
			CensorshipRecord: d.CensorshipRecord,
			VoteStatus:       d.VoteStatus,
		},
		VoteSummary: prop.VoteSummary,
	}
}

//ProposalsIterator allows iterating over a slice of PoliteiaProposals
type ProposalsIterator struct {
	proposals    []*PoliteiaProposalMobile
	currentIndex int
}

//Next returns the next PoliteaiProposal
func (p *ProposalsIterator) Next() *PoliteiaProposalMobile {
	if p.currentIndex < len(p.proposals) {
		prop := p.proposals[p.currentIndex]
		p.currentIndex++
		return prop
	}

	return nil
}

//Reset resets the current index to 0
func (p *ProposalsIterator) Reset() {
	p.currentIndex = 0
}

func (p *Politeia) getTokensToLoad(n int, category int) ([]string, error) {
	var inv []string
	tinv, err := p.getInventory()
	if err != nil {
		return nil, err
	}
	switch category {
	case 0:
		inv = tinv.Pre
	case 1:
		inv = tinv.Active
	case 2:
		inv = tinv.Approved
	case 3:
		inv = tinv.Rejected
	case 4:
		inv = tinv.Abandoned
	default:
		return nil, fmt.Errorf("invalid proposal category: %d", category)
	}

	toload := p.lcount[category] + n
	loaded := p.lcount[category]
	if loaded == len(inv) {
		return nil, nil
	}

	var tokens []string
	if toload > len(inv) {
		tokens = inv[loaded:]
		p.lcount[category] = len(inv)
	} else {
		tokens = inv[loaded:toload]
		p.lcount[category] += len(tokens)
	}
	return tokens, nil
}

//LoadProposalsInCategory loads n proposals in a given category
func (p *Politeia) LoadProposalsInCategory(n int, category int) ([]PoliteiaProposal, error) {
	tokens, err := p.getTokensToLoad(n, category)
	if err != nil {
		return nil, err
	}
	//there are no tokens left to load
	if tokens == nil {
		return nil, nil
	}
	return p.loadProposals(p.ctx, tokens, invTypesToHasVoteSummary[category])
}

func (p *Politeia) getProposalIterator(n, category int) (*ProposalsIterator, error) {
	props, err := p.LoadProposalsInCategory(n, category)
	if err != nil {
		return nil, err
	}
	mprops := make([]*PoliteiaProposalMobile, len(props))
	for i, p := range props {
		mprops[i] = mapProposalToMobile(&p)
	}
	return &ProposalsIterator{proposals: mprops}, nil
}

//LoadPreVoteProposals returns a ProposalIterator after loading the next n
//PreVote Proposals. Returns nil, nil if there are no proposals in this category
//left to load.
func (p *Politeia) LoadPreVoteProposals(n int) (*ProposalsIterator, error) {
	return p.getProposalIterator(n, PreVote)
}

//LoadActiveProposals returns a ProposalIterator after loading the next n
//Active proposals. Returns nil, nil if there are no proposals in this category
//left to load.
func (p *Politeia) LoadActiveProposals(n int) (*ProposalsIterator, error) {
	return p.getProposalIterator(n, Active)
}

//LoadApprovedProposals returns a ProposalIterator after loading the next n
//Approved proposals. Returns nil, nil if there are no proposals in this category
//left to load.
func (p *Politeia) LoadApprovedProposals(n int) (*ProposalsIterator, error) {
	return p.getProposalIterator(n, Approved)
}

//LoadRejectedProposals returns a ProposalIterator after loading the next n
//Rejected proposals. Returns nil, nil if there are no proposals in this category
//left to load.
func (p *Politeia) LoadRejectedProposals(n int) (*ProposalsIterator, error) {
	return p.getProposalIterator(n, Rejected)
}

//LoadAbandonedProposals returns a ProposalIterator after loading the next n
//Abandoned proposals. Returns nil, nil if there are no proposals in this category
//left to load.
func (p *Politeia) LoadAbandonedProposals(n int) (*ProposalsIterator, error) {
	return p.getProposalIterator(n, Abandoned)
}
