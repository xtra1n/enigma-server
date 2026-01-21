package enigma

import (
	"fmt"
	"strings"
	"github.com/xtra1n/enigma-server/internal/config"
)

type Enigma struct {
	rotors    [3]*Rotor
	reflector *Reflector
	plugboard Plugboard
}

func NewEnigma() *Enigma {
	return &Enigma{
		rotors: [3]*Rotor{
			NewRotor(config.RotorI, config.RotorNotches[config.RotorI]),
			NewRotor(config.RotorII, config.RotorNotches[config.RotorII]),
			NewRotor(config.RotorIII, config.RotorNotches[config.RotorIII]),
		},
		reflector: NewReflector(config.Reflector),
		plugboard: NewPlugboard(),
	}
}

func (e *Enigma) SetRotors(names [3]string) error {
	rotorMap := map[string]string{
		"I":   config.RotorI,
		"II":  config.RotorII,
		"III": config.RotorIII,
	}

	notchMap := map[string]int{
		"I":   config.RotorNotches[config.RotorI],
		"II":  config.RotorNotches[config.RotorII],
		"III": config.RotorNotches[config.RotorIII],
	}

	for i, name := range names {
		if writing, ok := rotorMap[name]; ok {
			e.rotors[i] = NewRotor(writing, notchMap[name])
		} else {
			return fmt.Errorf("unknown rotor: %s", name)
		}
	}

	return nil
}

func (e *Enigma) SetPositions(positions string) error {
	positions = strings.ToUpper(positions)

	if len(positions) != 3 {
		return fmt.Errorf("position must be 3 characteres, got %d", len(positions))
	}

	for i, ch := range positions {
		if ch < 'A' || ch > 'Z' {
			return fmt.Errorf("invalid position")
		}
		e.rotors[i].SetPosition(int(ch - 'A'))
	}

	return nil
}

func (e *Enigma) SetPlugboard(pairs string) error {
	e.plugboard = NewPlugboard()

	if pairs == "" {
		return nil
	}

	pairs = strings.ToUpper(pairs)
	pairList := strings.Fields(pairs)

	for _, pair := range pairList {
		if len(pair) != 2 {
			return fmt.Errorf("invalid pair format: %s (must be 2 characters)", pair)
		}

		a := rune(pair[0])
		b := rune(pair[1])

		if a < 'A' || a > 'Z' || b < 'A' || b > 'Z' {
			return fmt.Errorf("invalid plugboard pair: %c%c (must be A-Z)", a, b)
		}

		e.plugboard.SetPair(a, b)
	}

	return nil
}

func (e *Enigma) rotateRotors() {
	if e.rotors[1].AtNotch() {
		e.rotors[1].Rotate()
		e.rotors[2].Rotate()
	}

	if e.rotors[0].AtNotch() {
		e.rotors[1].Rotate()
	}

	e.rotors[0].Rotate()
}

func (e *Enigma) Translate(ch rune) rune {
	if ch < 'A' || ch > 'Z' {
		return ch
	}

	e.rotateRotors()

	signal := byte(e.plugboard.Substitute(ch) - 'A')

	signal = e.rotors[0].Forward(signal)
	signal = e.rotors[1].Forward(signal)
	signal = e.rotors[2].Forward(signal)

	signal = e.reflector.Reflect(signal)

	signal = e.rotors[2].Backward(signal)
	signal = e.rotors[1].Backward(signal)
	signal = e.rotors[0].Backward(signal)

	result := e.plugboard.Substitute(rune(signal) + 'A')

	return result
}

func (e *Enigma) TranlateString(text string) string {
	text = strings.ToUpper(text)

	var result string

	for _, ch := range text {
		result += string(e.Translate(ch))
	}

	return  result
}

func (e *Enigma) GetPosition() string {
	pos := ""

	for _, r := range e.rotors {
		pos += string(rune('A' + r.pos))
	}

	return  pos
}
