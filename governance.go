package dcrlibwallet

import (
	"context"
	"fmt"
	"time"

	pwww "github.com/planetdecred/dcrlibwallet/politeiawww"
)

//proposal categories
const (
	PropCategoryPreVote = iota
	PropCategoryActive
	PropCategoryApproved
	PropCategoryRejected
	PropCategoryAbandoned
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

//PoliteiaProposalMobile is a go mobile compatible wrapper around
//PoliteiaProposal
type PoliteiaProposalMobile struct {
	pp *PoliteiaProposal
}

func (pm *PoliteiaProposalMobile) Name() string {
	return pm.pp.Details.Name
}

func (pm *PoliteiaProposalMobile) State() int {
	return pm.pp.Details.State
}

func (pm *PoliteiaProposalMobile) Status() int {
	return pm.pp.Details.Status
}

func (pm *PoliteiaProposalMobile) Timestamp() int64 {
	return pm.pp.Details.Timestamp
}

func (pm *PoliteiaProposalMobile) UserID() string {
	return pm.pp.Details.Username
}

func (pm *PoliteiaProposalMobile) Username() string {
	return pm.pp.Details.Username
}

func (pm *PoliteiaProposalMobile) PublicKey() string {
	return pm.pp.Details.PublicKey
}

func (pm *PoliteiaProposalMobile) Signature() string {
	return pm.pp.Details.Signature
}

func (pm *PoliteiaProposalMobile) NumComments() int {
	return pm.pp.Details.NumComments
}

func (pm *PoliteiaProposalMobile) Version() string {
	return pm.pp.Details.Version
}

func (pm *PoliteiaProposalMobile) PublishedAt() int64 {
	return pm.pp.Details.Timestamp
}

func (pm *PoliteiaProposalMobile) CensorshipRecord() *pwww.ProposalCensorshipRecord {
	return pm.pp.Details.CensorshipRecord
}

func (pm *PoliteiaProposalMobile) Files() *ProposalFiles {
	return &ProposalFiles{files: pm.pp.Details.Files}
}

func (pm *PoliteiaProposalMobile) Metadata() *ProposalMetadataIterator {
	return &ProposalMetadataIterator{metadata: pm.pp.Details.MetaData}
}

func (pm *PoliteiaProposalMobile) VoteStatus() int {
	return pm.pp.VoteSummary.Status
}

func (pm *PoliteiaProposalMobile) VoteApproved() bool {
	return pm.pp.VoteSummary.Approved
}

func (pm *PoliteiaProposalMobile) VoteType() int {
	return pm.pp.VoteSummary.Type
}

func (pm *PoliteiaProposalMobile) VoteEligibleTickets() int {
	return pm.pp.VoteSummary.EligibleTickets
}

func (pm *PoliteiaProposalMobile) VoteDuration() int64 {
	return pm.pp.VoteSummary.Duration
}

func (pm *PoliteiaProposalMobile) VoteEndHeight() int64 {
	return pm.pp.VoteSummary.EndHeight
}

func (pm *PoliteiaProposalMobile) VoteQuorumPercentage() int {
	return pm.pp.VoteSummary.QuorumPercentage
}

func (pm *PoliteiaProposalMobile) VotePassPercentage() int {
	return pm.pp.VoteSummary.PassPercentage
}

func (pm *PoliteiaProposalMobile) VoteOptionsresult() *VoteOptionResultIterator {
	return &VoteOptionResultIterator{opts: pm.pp.VoteSummary.OptionsResult}
}

type ProposalFiles struct {
	files        []*pwww.ProposalFile
	currentIndex int
}

func (fs *ProposalFiles) Next() *pwww.ProposalFile {
	if fs.currentIndex < len(fs.files) {
		pf := fs.files[fs.currentIndex]
		fs.currentIndex++
		return pf
	}

	return nil
}

//Reset resets the current index to 0
func (fs *ProposalFiles) Reset() {
	fs.currentIndex = 0
}

type ProposalMetadataIterator struct {
	metadata     []*pwww.ProposalMetaData
	currentIndex int
}

func (pm *ProposalMetadataIterator) Next() *pwww.ProposalMetaData {
	if pm.currentIndex < len(pm.metadata) {
		md := pm.metadata[pm.currentIndex]
		pm.currentIndex++
		return md
	}

	return nil
}

//Reset resets the current index to 0
func (pm *ProposalMetadataIterator) Reset() {
	pm.currentIndex = 0
}

type VoteOptionResultIterator struct {
	opts         []*pwww.VoteOptionResult
	currentIndex int
}

func (vor *VoteOptionResultIterator) Next() *pwww.VoteOptionResult {
	if vor.currentIndex < len(vor.opts) {
		r := vor.opts[vor.currentIndex]
		vor.currentIndex++
		return r
	}
	return nil
}

func (vor *VoteOptionResultIterator) Reset() {
	vor.currentIndex = 0
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

func (p *Politeia) loadProposals(ctx context.Context, ctokens []string, voteSummary bool) ([]*PoliteiaProposal, error) {
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
	pprops := make([]*PoliteiaProposal, 0)
	for _, prop := range propresp.Props {
		pp := PoliteiaProposal{Details: prop}
		if voteSummary {
			if v, ok := vsresp.Vsum.Summaries[prop.CensorshipRecord.Token]; ok {
				pp.VoteSummary = v
			}
		}
		pprops = append(pprops, &pp)
	}
	return pprops, nil
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

//CategoryCount returns the number of proposals
//in the given proposal category in the loaded the pi inventory
func (p *Politeia) CategoryCount(category int) (int, error) {
	inv, err := p.getInventory()
	if err != nil {
		return -1, err
	}
	return len(inv.Abandoned), nil
}

//LoadProposalsInCategory loads n proposals in a given category.
//category should be one of the PropCategory constants
func (p *Politeia) LoadProposalsInCategory(n, category int) ([]*PoliteiaProposal, error) {
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

//LoadProposalsInCategoryMobile loads n proposals in a given category wrapped
//in a mobile compatibly ProposalsIterator type.
//category should be one of the PropCategory constants
func (p *Politeia) LoadProposalsInCategoryMobile(n, category int) (*ProposalsIterator, error) {
	props, err := p.LoadProposalsInCategory(n, category)
	if err != nil {
		return nil, err
	}
	mprops := make([]*PoliteiaProposalMobile, len(props))
	for i, p := range props {
		mprops[i] = &PoliteiaProposalMobile{p}
	}
	return &ProposalsIterator{proposals: mprops}, nil
}
