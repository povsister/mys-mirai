package meta

type DeletePostOptions struct {
	Reason DeleteReason
}

const (
	NoReason              DeleteReason = 0
	AttackOthers          DeleteReason = 1
	SPAM                  DeleteReason = 2
	UnauthorizedRepublish DeleteReason = 3
	AccountTrade          DeleteReason = 4
	OuterLink             DeleteReason = 5
	UnrelatedContent      DeleteReason = 6
	FakeNews              DeleteReason = 7
	// TODO
)

type DeleteReason uint8

type MovePostOptions struct {
	To Forum
	// TODO
}
