package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"log"
	"strings"
	"text/template"
	"time"
)

type Condition struct {
	// Severity of this condition
	Severity Severity
	// Operator for this condition to be true
	Condition Conditional
	// Threshold value for his condition to be true
	Threshold string
	// Only raising after duration
	Duration string
}

type ServiceComponent struct {
	Name        string
	Description string
	// Parameters for this components config
	Parameters map[string]string
	// Query to execute for this component
	Query string
	// Conditions for this component to be other than healthy
	Conditions []*Condition
}

type Conditional string

const (
	// Condition is true if the value is greater than the threshold
	GreaterThan Conditional = "greater_than"
	// Condition is true if the value is less than the threshold
	LessThan Conditional = "less_than"
	// Condition is true if the value is equal to the threshold
	Equal Conditional = "equal"
	// Condition is true if the value is not equal to the threshold
	NotEqual Conditional = "not_equal"
	// Always true
	Always Conditional = "always"
)

var StringToConditional = map[string]Conditional{
	"greater_than": GreaterThan,
	"less_than":    LessThan,
	"equal":        Equal,
	"not_equal":    NotEqual,
	"gt":           GreaterThan,
	"lt":           LessThan,
	"eq":           Equal,
	"ne":           NotEqual,
	"always":       Always,
	"true":         Always,
}

func (c Conditional) String() string {
	return string(c)
}

func (c *Conditional) UnmarshalYAML(value *yaml.Node) error {
	var cc Conditional
	cc = StringToConditional[strings.ToLower(value.Value)]
	if cc == "" {
		return errors.New("invalid conditional")
	} else {
		*c = *&cc
	}
	return nil
}

type Severity string

const (
	Critical Severity = "critical"
	Warning  Severity = "warning"
	Info     Severity = "info"
	OK       Severity = "ok"
	Unknown  Severity = "unknown"
)

var StringToSeverity = map[string]Severity{
	"critical": Critical,
	"warning":  Warning,
	"info":     Info,
	"ok":       OK,
	"unknown":  Unknown,
}

func (s Severity) IsHigherThan(other Severity) bool {
	return s.GetPriority() > other.GetPriority()
}

func (s Severity) GetPriority() int {
	switch s {
	case Critical:
		return 4
	case Warning:
		return 3
	case Info:
		return 2
	case OK:
		return 1
	case Unknown:
		return 0
	}
	return 0
}

func (s Severity) String() string {
	return string(s)
}
func (s *Severity) UnmarshalYAML(value *yaml.Node) error {
	var ss Severity
	ss = StringToSeverity[strings.ToLower(value.Value)]
	if ss == "" {
		return errors.New("invalid severity")
	} else {
		*s = *&ss
	}
	return nil
}

type Service struct {
	// Reference of this service - is referenced by status pages
	ID string
	// Name of the service
	Name string
	// Contact for this service - can be a URL or email
	Contact string
	// Owner of this service - can be a team name or email
	Owner string
	// URL of the service
	URL string
	// Environment of this service
	Environment string
	// Components of this service
	Components []*ServiceComponent
}

type Environment struct {
	// ID of the environment
	ID string
	// Services in this environment by id
	Services map[string]*Service
	// Status pages in this environment by id
	StatusPages map[string]*StatusPage
}

type StatusPage struct {
	// Status page ID
	ID string `yaml:"id"`
	// Name of the status page
	Name string `yaml:"name"`
	// Environment of this status page
	Environment string `yaml:"environment"`
	// URL of the page
	URL string `yaml:"url"`
	// Contact for this service - can be a URL or email
	Contact string `yaml:"contact"`
	// Owner of this service - can be a team name or email
	Owner string `yaml:"owner"`
	// Services in this status page
	Services []string `yaml:"services"`
}

type Config struct {
	Environments map[string]*Environment
	Queries      map[string]*template.Template
}

func (c *Config) AddQuery(name, expr, source string) {
	if _, f := c.Queries[name]; f {
		log.Printf("Duplicate query name %s in %s", name, source)
		return
	}

	t, err := template.New(name).Parse(expr)
	if err != nil {
		log.Printf("Error parsing query %s in %s: %s", name, source, err.Error())
		return
	}
	c.Queries[name] = t
}

func newEnv(env string) *Environment {
	return &Environment{ID: env,
		StatusPages: map[string]*StatusPage{},
		Services:    map[string]*Service{},
	}
}

func (c *Config) GetService(environment string, service string) *Service {
	if _, f := c.Environments[environment]; !f {
		return &Service{ID: service, Name: service, Environment: environment, Components: []*ServiceComponent{}}
	}
	if _, f := c.Environments[environment].Services[service]; !f {
		return &Service{ID: service, Name: service, Environment: environment, Components: []*ServiceComponent{}}
	}
	return c.Environments[environment].Services[service]

}

func (c *Config) SetService(s *Service) {
	if _, f := c.Environments[s.Environment]; !f {
		c.Environments[s.Environment] = newEnv(s.Environment)
	}
	c.Environments[s.Environment].Services[s.ID] = s
	log.Printf("Set service %s to environment %s", s.ID, s.Environment)
}

func (c *Config) RemoveService(s *Service) {
	if _, f := c.Environments[s.Environment]; !f {
		return
	}
	delete(c.Environments[s.Environment].Services, s.ID)
}

func (c *Config) SetStatusPage(sp *StatusPage) {
	if _, f := c.Environments[sp.Environment]; !f {
		c.Environments[sp.Environment] = newEnv(sp.Environment)
	}
	c.Environments[sp.Environment].StatusPages[sp.ID] = sp
	log.Printf("Set status page %s to environment %s", sp.ID, sp.Environment)
}

func (c *Config) RemoveStatusPage(sp *StatusPage) {
	if _, f := c.Environments[sp.Environment]; !f {
		return
	}
	delete(c.Environments[sp.Environment].StatusPages, sp.ID)
}

type ComponentStatus struct {
	// Name of this component
	Name string
	// Status of this component
	Status Severity
	// Description of this component status
	StatusString string
}

type ServiceStatus struct {
	// ID of this service
	ID string
	// Status of this service overall
	Status Severity
	// Last update of this service
	LastUpdate time.Time
	// Components of this service
	Components []*ComponentStatus
}

type PageStatus struct {
	// ID of this page
	ID string
	// Status of this page overall
	Status Severity
	// Last update of this page
	LastUpdate time.Time
	// Services of this page
	Services map[string]*ServiceStatus
}

func ParseServiceComponents(data []byte) ([]*ServiceComponent, error) {
	var components []*ServiceComponent
	err := yaml.Unmarshal(data, &components)
	if err != nil {
		return nil, err
	}
	return components, nil
}
