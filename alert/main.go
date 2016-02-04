package alert

type Message struct {
	Service string            `json:"source"`
	Meta    map[string]string `json:"event"`
}
