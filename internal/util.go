package cli

func validValue(value string, validValues []string) bool {
	for _,v := range validValues {
		if value == v {
			return true
		}
	}
	return false
}
