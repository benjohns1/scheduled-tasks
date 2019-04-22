package postgres

// BoolVal returns a boolean value representation for queries
func BoolVal(val bool) string {
	if val {
		return "1::bit"
	} else {
		return "0::bit"
	}
}
