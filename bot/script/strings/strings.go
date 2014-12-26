package strings

// First returns the first character of a string or an empty string.
// If the string is blank it will return "".
func First(str string) string {
	if len(str) > 0 {
		return string(str[0])
	} else {
		return ""
	}
}

// Rest returns all but the first character of a string. It returns
// an empty string if the string is blank.
func Rest(str string) string {
	if len(str) > 0 {
		return str[1:]
	} else {
		return ""
	}
}

// Shuck removes the first and last character of a string, analogous to
// shucking off the husk of an ear of corn.
func Shuck(victim string) string {
	return victim[1 : len(victim)-1]
}
