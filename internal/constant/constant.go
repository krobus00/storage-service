package constant

type ctxKey string

const (
	KeyDBCtx     ctxKey = "DB"
	KeyUserIDCtx ctxKey = "USERID"

	SystemID = string("SYSTEM")
	GuestID  = string("GUEST")
)

var RedactedField = map[string]bool{
	"password": true,
	"file":     true,
}
