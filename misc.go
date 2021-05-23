package main

// to get the moves still needs optimization
func filterMoves(mv string, mvs []string) []string {
	res := []string{}
	if mv == "" {
		return res
	}
	for _, str := range mvs {
		if mv == string(str[0:2]) {
			res = append(res, str[2:])
		}
	}
	return res
}
