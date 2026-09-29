package workflow

import "testing"

func TestApprovalRequiresPayloadBinding(t *testing.T) {
	if err := ValidateApproval(Approval{Action: SubmitApplication, ApprovedBy: "user"}); err == nil { t.Fatal("approval without payload binding must fail") }
}
