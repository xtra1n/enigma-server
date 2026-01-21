package enigma

type Rotor struct {
	writing string
	pos     int
	notch   int
	ring    int
}

func NewRotor(writing string, notchPos int) *Rotor {
	return &Rotor{
		writing: writing,
		pos:     0,
		notch:   notchPos,
		ring:    0,
	}
}

func (r *Rotor) SetPosition(pos int) {
	r.pos = pos % 26
}

func (r *Rotor) AtNotch() bool {
	return r.pos == r.notch
}

func (r *Rotor) Rotate() {
	r.pos = (r.pos + 1) % 26
}

func (r *Rotor) Forward(in byte) byte {
	index := (int(in) + r.pos - r.ring + 26) % 26
	out := byte(r.writing[index] - 'A')

	result := (int(out) - r.pos + r.ring + 26) % 26

	return byte(result)
}

func (r *Rotor) Backward(in byte) byte {
	index := (int(in) + r.pos - r.ring + 26) % 26

	result := byte(0)
	for i, ch := range r.writing {
		if int(byte(ch-'A')) == index {
			result = byte(i)
			break
		}
	}

	final := (int(result) - r.pos + r.ring + 26) % 26

	return byte(final)
}
