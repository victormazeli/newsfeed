package utils

func IsTokenInSlice(slice []map[string]string, token string) bool {

	for _, value := range slice {
		if value["token"] == token {
			return true
		}
	}
	return false
}
