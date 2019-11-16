package pointer

// String ...
func String(str string) *string {
	return &str
}

// Bool ...
func Bool(val bool) *bool {
	return &val
}
