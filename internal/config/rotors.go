package config

const (
	RotorI    = "EKMFLGDQVZNTOWYHXUSPAIBRCJ"
	RotorII   = "AJDKSIRUXBLHWTMCQGZNPYFVOE"
	RotorIII  = "BDFHJLCPRTXVZNYEIWGAKMUSQO"
	Reflector = "YRUHQSLDPXNGOKMIEBFZCWVJAT"
)

var RotorNotches = map[string]int{
	RotorI:   16,
	RotorII:  4,
	RotorIII: 21,
}
