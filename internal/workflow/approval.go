package workflow

import "errors"

type Action string

const (
	SubmitApplication Action = "submit_application"
	SendOutreach Action = "send_outreach"
)

type Approval struct {
	Action Action
	PayloadHash string
	ApprovedBy string
}

func ValidateApproval(approval Approval) error {
	if approval.Action != SubmitApplication && approval.Action != SendOutreach { return errors.New("unsupported external action") }
	if approval.PayloadHash == "" { return errors.New("approval must bind to an exact payload") }
	if approval.ApprovedBy == "" { return errors.New("approval actor is required") }
	return nil
}
