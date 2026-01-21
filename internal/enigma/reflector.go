package enigma

type Reflector struct {
	writing string
}

func NewReflector(wringStr string) *Reflector {
	return &Reflector{
		writing: wringStr,
	}
}

func (r *Reflector) Reflect(in byte) byte {
	return byte(r.writing[in] - 'A')
}
