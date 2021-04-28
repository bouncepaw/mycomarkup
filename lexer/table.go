package lexer

// If the input starts on any of the prefices, call the function.
type tableEntry struct {
	prefices []string
	λ        func(s *SourceText, tw *TokenWriter)
}

func executeTable(table []tableEntry, s *SourceText, tw *TokenWriter) bool {
	for _, rule := range table {
		for _, prefix := range rule.prefices {
			if startsWithStr(s.b, prefix) {
				rule.λ(s, tw)
				return true
			}
		}
	}
	return false
}
