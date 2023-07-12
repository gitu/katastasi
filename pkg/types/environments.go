package types

type Environment struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`

	// Services names in this environment by id
	Services map[string]string `json:"services,omitempty"`
	// Status pages names in this environment by id
	StatusPages map[string]string `json:"statusPages,omitempty"`
}
