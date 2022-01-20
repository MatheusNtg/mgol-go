package lexer

func ContainsState(states []State, element State) bool {
	for _, e := range states {
		if e == element {
			return true
		}
	}
	return false
}

func ContainsSymbol(symbols []Symbol, element Symbol) bool {
	for _, e := range symbols {
		if e == element {
			return true
		}
	}
	return false
}

func ContainsByte(bytes []byte, element byte) bool {
	for _, e := range bytes {
		if e == element {
			return true
		}
	}
	return false
}
