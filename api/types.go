package api

type Message struct {
	Action    string    `json:"action"`
	Rotors    [3]string `json:"rotors,omitempty"`
	Positions string    `json:"positions,omitempty"`
	Plugboard string    `json:"plugboard,omitempty"`
	Text      string    `json:"text,omitempty"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Result string `json:"result,omitempty"`
	Pos    string `json:"pos,omitempty"`
}
