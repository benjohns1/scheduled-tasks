package sqlite3

// BoolVal returns a boolean value representation for queries
func BoolVal(val bool) string {
	if val {
		return "1"
	} else {
		return "0"
	}
}
