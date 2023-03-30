package constant

type ctxKey string

const (
	KeyDBCtx     ctxKey = "DB"
	KeyUserIDCtx ctxKey = "USERID"

	SystemID = string("SYSTEM")
	GuestID  = string("GUEST")
)
