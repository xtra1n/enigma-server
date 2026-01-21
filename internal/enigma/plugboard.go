package enigma

type Plugboard map[rune]rune

func NewPlugboard() Plugboard {
	return make(Plugboard)
}

func (p Plugboard) Substitute(c rune) rune {
	var result rune

	if value, ok := p[c]; ok {
		result = value
	} else {
		result = c
	}

	return  result
}

func (p Plugboard) SetPair(a, b rune) {
	p[a] = b
	p[b] = a
}