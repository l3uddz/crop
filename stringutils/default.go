package stringutils

func NewOrExisting(new string, existing string) string {
	if new == "" {
		return existing
	}

	return new
}
