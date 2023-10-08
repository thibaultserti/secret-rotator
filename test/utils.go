package test

func IgnoreFields(data map[string]interface{}, fields ...string) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range data {
		if !contains(fields, key) {
			result[key] = value
		}
	}
	return result
}

// Fonction pour vérifier si une valeur existe dans un slice de chaînes
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
