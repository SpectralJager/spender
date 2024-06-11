package handlers

type ctxKey string

const (
	ownerIDKey = ctxKey("ownerID")
)

const (
	dateLayout = "2006-01-02"
	orderASC   = "asc"
	orderDESC  = "desc"
)
