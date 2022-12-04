package session

type elementResult struct {
	Element    string `json:"ELEMENT"`
	W3CElement string `json:"element-6066-11e4-a52e-4f735466cecf"`
}

func (er elementResult) ID() string {
	if er.Element != "" {
		return er.Element
	}
	return er.W3CElement
}
