package politeiawww

import (
	"encoding/json"
	"fmt"
)

// Error status codes
const (
	ErrorStatusInvalid = iota
	ErrorStatusInvalidPassword
	ErrorStatusMalformedEmail
	ErrorStatusVerificationTokenInvalid
	ErrorStatusVerificationTokenExpired
	ErrorStatusProposalMissingFiles
	ErrorStatusProposalNotFound
	ErrorStatusProposalDuplicateFilenames
	ErrorStatusProposalInvalidTitle
	ErrorStatusMaxMDsExceededPolicy
	ErrorStatusMaxImagesExceededPolicy
	ErrorStatusMaxMDSizeExceededPolicy
	ErrorStatusMaxImageSizeExceededPolicy
	ErrorStatusMalformedPassword
	ErrorStatusCommentNotFound
	ErrorStatusInvalidFilename
	ErrorStatusInvalidFileDigest
	ErrorStatusInvalidBase64
	ErrorStatusInvalidMIMEType
	ErrorStatusUnsupportedMIMEType
	ErrorStatusInvalidPropStatusTransition
	ErrorStatusInvalidPublicKey
	ErrorStatusNoPublicKey
	ErrorStatusInvalidSignature
	ErrorStatusInvalidInput
	ErrorStatusInvalidSigningKey
	ErrorStatusCommentLengthExceededPolicy
	ErrorStatusUserNotFound
	ErrorStatusWrongStatus
	ErrorStatusNotLoggedIn
	ErrorStatusUserNotPaid
	ErrorStatusReviewerAdminEqualsAuthor
	ErrorStatusMalformedUsername
	ErrorStatusDuplicateUsername
	ErrorStatusVerificationTokenUnexpired
	ErrorStatusCannotVerifyPayment
	ErrorStatusDuplicatePublicKey
	ErrorStatusInvalidPropVoteStatus
	ErrorStatusUserLocked
	ErrorStatusNoProposalCredits
	ErrorStatusInvalidUserManageAction
	ErrorStatusUserActionNotAllowed
	ErrorStatusWrongVoteStatus
	ErrorStatusCannotVoteOnPropComment
	ErrorStatusChangeMessageCannotBeBlank
	ErrorStatusCensorReasonCannotBeBlank
	ErrorStatusCannotCensorComment
	ErrorStatusUserNotAuthor
	ErrorStatusVoteNotAuthorized
	ErrorStatusVoteAlreadyAuthorized
	ErrorStatusInvalidAuthVoteAction
	ErrorStatusUserDeactivated
	ErrorStatusInvalidPropVoteBits
	ErrorStatusInvalidPropVoteParams
	ErrorStatusEmailNotVerified
	ErrorStatusInvalidUUID
	ErrorStatusInvalidLikeCommentAction
	ErrorStatusInvalidCensorshipToken
	ErrorStatusEmailAlreadyVerified
	ErrorStatusNoProposalChanges
	ErrorStatusMaxProposalsExceededPolic
	ErrorStatusDuplicateComment
	ErrorStatusInvalidLogin
	ErrorStatusCommentIsCensored
	ErrorStatusInvalidProposalVersion
)

// Proposal vote status codes
const (
	PropVoteStatusInvalid = iota
	// Vote has not been authorized by author
	PropVoteStatusNotAuthorized
	// Vote has been authorized by author
	PropVoteStatusAuthorized
	// Proposal vote has been started
	PropVoteStatusStarted
	// Proposal vote has been finished
	PropVoteStatusFinished
	// Proposal doesn't exist
	PropVoteStatusDoesntExist
)

var (
	// ErrorStatus converts error status codes to human readable text.
	ErrorStatus = map[int]string{
		ErrorStatusInvalid:                     "invalid error status",
		ErrorStatusInvalidPassword:             "invalid password",
		ErrorStatusMalformedEmail:              "malformed email",
		ErrorStatusVerificationTokenInvalid:    "invalid verification token",
		ErrorStatusVerificationTokenExpired:    "expired verification token",
		ErrorStatusProposalMissingFiles:        "missing proposal files",
		ErrorStatusProposalNotFound:            "proposal not found",
		ErrorStatusProposalDuplicateFilenames:  "duplicate proposal files",
		ErrorStatusProposalInvalidTitle:        "invalid proposal title",
		ErrorStatusMaxMDsExceededPolicy:        "maximum markdown files exceeded",
		ErrorStatusMaxImagesExceededPolicy:     "maximum image files exceeded",
		ErrorStatusMaxMDSizeExceededPolicy:     "maximum markdown file size exceeded",
		ErrorStatusMaxImageSizeExceededPolicy:  "maximum image file size exceeded",
		ErrorStatusMalformedPassword:           "malformed password",
		ErrorStatusCommentNotFound:             "comment not found",
		ErrorStatusInvalidFilename:             "invalid filename",
		ErrorStatusInvalidFileDigest:           "invalid file digest",
		ErrorStatusInvalidBase64:               "invalid base64 file content",
		ErrorStatusInvalidMIMEType:             "invalid MIME type detected for file",
		ErrorStatusUnsupportedMIMEType:         "unsupported MIME type for file",
		ErrorStatusInvalidPropStatusTransition: "invalid proposal status",
		ErrorStatusInvalidPublicKey:            "invalid public key",
		ErrorStatusNoPublicKey:                 "no active public key",
		ErrorStatusInvalidSignature:            "invalid signature",
		ErrorStatusInvalidInput:                "invalid input",
		ErrorStatusInvalidSigningKey:           "invalid signing key",
		ErrorStatusCommentLengthExceededPolicy: "maximum comment length exceeded",
		ErrorStatusUserNotFound:                "user not found",
		ErrorStatusWrongStatus:                 "wrong proposal status",
		ErrorStatusNotLoggedIn:                 "user not logged in",
		ErrorStatusUserNotPaid:                 "user hasn't paid paywall",
		ErrorStatusReviewerAdminEqualsAuthor:   "user cannot change the status of his own proposal",
		ErrorStatusMalformedUsername:           "malformed username",
		ErrorStatusDuplicateUsername:           "duplicate username",
		ErrorStatusVerificationTokenUnexpired:  "verification token not yet expired",
		ErrorStatusCannotVerifyPayment:         "cannot verify payment at this time",
		ErrorStatusDuplicatePublicKey:          "public key already taken by another user",
		ErrorStatusInvalidPropVoteStatus:       "invalid proposal vote status",
		ErrorStatusUserLocked:                  "user locked due to too many login attempts",
		ErrorStatusNoProposalCredits:           "no proposal credits",
		ErrorStatusInvalidUserManageAction:     "invalid user edit action",
		ErrorStatusUserActionNotAllowed:        "user action is not allowed",
		ErrorStatusWrongVoteStatus:             "wrong proposal vote status",
		ErrorStatusCannotVoteOnPropComment:     "cannot vote on proposal comment",
		ErrorStatusChangeMessageCannotBeBlank:  "status change message cannot be blank",
		ErrorStatusCensorReasonCannotBeBlank:   "censor comment reason cannot be blank",
		ErrorStatusCannotCensorComment:         "cannot censor comment",
		ErrorStatusUserNotAuthor:               "user is not the proposal author",
		ErrorStatusVoteNotAuthorized:           "vote has not been authorized",
		ErrorStatusVoteAlreadyAuthorized:       "vote has already been authorized",
		ErrorStatusInvalidAuthVoteAction:       "invalid authorize vote action",
		ErrorStatusUserDeactivated:             "user account is deactivated",
		ErrorStatusInvalidPropVoteBits:         "invalid proposal vote option bits",
		ErrorStatusInvalidPropVoteParams:       "invalid proposal vote parameters",
		ErrorStatusEmailNotVerified:            "email address is not verified",
		ErrorStatusInvalidUUID:                 "invalid user UUID",
		ErrorStatusInvalidLikeCommentAction:    "invalid like comment action",
		ErrorStatusInvalidCensorshipToken:      "invalid proposal censorship token",
		ErrorStatusEmailAlreadyVerified:        "email address is already verified",
		ErrorStatusNoProposalChanges:           "no changes found in proposal",
		ErrorStatusDuplicateComment:            "duplicate comment",
		ErrorStatusInvalidLogin:                "invalid login credentials",
		ErrorStatusCommentIsCensored:           "comment is censored",
		ErrorStatusInvalidProposalVersion:      "invalid proposal version",
	}
)

type errorStatusCode int

func (sc *errorStatusCode) String() string {
	if s, ok := ErrorStatus[int(*sc)]; ok {
		return s
	}
	return "unknown error"
}

//UnmarshalJSON satisfies the json unmarshaler interface
func (sc errorStatusCode) UnmarshalJSON(b []byte) error {
	var c int
	if err := json.Unmarshal(b, &c); err != nil {
		return err
	}
	sc = errorStatusCode(c)
	return nil
}

type PoliteiawwwError struct {
	code    errorStatusCode `json:"errorcode"`
	context []string        `json:"errorcontext,omitempty"`
}

//Code returns the error code as an int
func (e *PoliteiawwwError) Code() int {
	return int(e.code)
}

func (e *PoliteiawwwError) String() string {
	return fmt.Sprintf("%d - %s%s", e.Code, e.code.String(), e.contextstr())
}

func (e *PoliteiawwwError) contextstr() string {
	var context string
	if e.context != nil {
		context = fmt.Sprintf(" context: %s", e.context)
	}
	return context
}

//Error satisfies the Error interface
func (e *PoliteiawwwError) Error() string {

	return fmt.Sprintf("Politeiawww server responded with error code: %s", e.String())
}

type ServerVersion struct {
	Version int `json:"version"`
}

type ServerPolicy struct {
	ProposalListPageSize int `json:"proposallistpagesize"`
}

type ProposalFile struct {
	Name    string `json:"name"`
	Mime    string `json:"mime"`
	Digest  string `json:"digest"`
	Payload string `json:"payload"`
}

type ProposalMetaData struct {
	Name   string `json:"name"`
	LinkTo string `json:"linkto"`
	LinkBy int64  `json:"linkby"`
}

type ProposalCensorshipRecord struct {
	Token     string `json:"token"`
	Merkle    string `json:"merkle"`
	Signature string `json:"signature"`
}

type Proposal struct {
	Name             string                   `json:"name"`
	State            int                      `json:"state"`
	Status           int                      `json:"status"`
	Timestamp        int64                    `json:"timestamp"`
	UserID           string                   `json:"userid"`
	Username         string                   `json:"username"`
	PublicKey        string                   `json:"publickey"`
	Signature        string                   `json:"signature"`
	NumComments      int                      `json:"numcomments"`
	Version          string                   `json:"version"`
	PublishedAt      int64                    `json:"publishedat"`
	Files            []ProposalFile           `json:"files"`
	MetaData         []ProposalMetaData       `json:"metadata"`
	CensorshipRecord ProposalCensorshipRecord `json:"censorshiprecord"`
	VoteStatus       VoteStatus               `json:"votestatus"`
}

type Proposals struct {
	Proposals []Proposal `json:"proposals"`
}

type ProposalResult struct {
	Proposal Proposal `json:"proposal"`
}

type VoteOption struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Bits        int    `json:"bits"`
}

type VoteOptionResult struct {
	Option        VoteOption `json:"option"`
	VotesReceived int64      `json:"votesreceived"`
}

type VoteStatus struct {
	Token              string             `json:"token"`
	Status             int                `json:"status"`
	TotalVotes         int                `json:"totalvotes"`
	OptionsResult      []VoteOptionResult `json:"optionsresult"`
	EndHeight          string             `json:"endheight"`
	BestBlock          string             `json:"bestblock"`
	NumOfEligibleVotes int                `json:"numofeligiblevotes"`
	QuorumPercentage   int                `json:"quorumpercentage"`
	PassPercentage     int                `json:"passpercentage"`
}

type VotesStatus struct {
	VotesStatus []VoteStatus `json:"votesstatus"`
}

type TokenInventory struct {
	Pre       []string `json:"pre"`
	Active    []string `json:"active"`
	Approved  []string `json:"approved"`
	Rejected  []string `json:"rejected"`
	Abandoned []string `json:"abandoned"`
}

type Tokens struct {
	Tokens []string `json:"tokens"`
}

type VoteSummary struct {
	Status           int                `json:"status"`
	Approved         bool               `json:"approved,omitempty"`
	Type             VoteType           `json:"type,omitempty"`
	EligibleTickets  int                `json:"eligibletickets"`
	Duration         int64              `json:"duration,omitempty"`
	EndHeight        int64              `json:"endheight,omitempty"`
	QuorumPercentage int                `json:"quorumpercentage,omitempty"`
	PassPercentage   int                `json:"passpercentage,omitempty"`
	OptionsResult    []VoteOptionResult `json:"optionsresult,omitempty"`
}

type VoteType struct {
	VoteTypeInvalid  int
	VoteTypeStandard int
	VoteType         int
}
