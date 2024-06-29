package mamedattools

func removeFormList(s []Game, i int) []Game {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
