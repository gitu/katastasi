package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"strings"
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
	// Conditions for this component to be other than healthy
	Conditions []Condition
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

func (c Conditional) UnmarshalYAML(value *yaml.Node) error {
	c = StringToConditional[strings.ToLower(value.Value)]
	if c == "" {
		return errors.New("invalid conditional")
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
func (s Severity) UnmarshalYAML(value *yaml.Node) error {
	s = StringToSeverity[strings.ToLower(value.Value)]
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
	Components []ServiceComponent
}

type Environment struct {
	// ID of the environment
	ID string
	// Services in this environment by id
	Services map[string]Service
	// Status pages in this environment by id
	StatusPages map[string]StatusPage
}

type StatusPage struct {
	// Status page ID
	ID string
	// Name of the status page
	Name string
	// Environment of this status page
	Environment string
	// URL of the page
	URL string
	// Contact for this service - can be a URL or email
	Contact string
	// Owner of this service - can be a team name or email
	Owner string
	// Services in this status page
	Services []string
}

type Config struct {
	Environments map[string]Environment
	Queries      map[string]string
}

func (c *Config) AddService(s Service) error {
	if _, f := c.Environments[s.Environment]; !f {
		c.Environments[s.Environment] = newEnv(s.Environment)
	}
	if _, f := c.Environments[s.Environment].Services[s.ID]; f {
		return errors.New("service " + s.ID + " already exists in environment " + s.Environment)
	}
	c.Environments[s.Environment].Services[s.ID] = s
	return nil
}

func (c *Config) AddStatusPage(s StatusPage) error {
	if _, f := c.Environments[s.Environment]; !f {
		c.Environments[s.Environment] = newEnv(s.Environment)
	}
	if _, f := c.Environments[s.Environment].StatusPages[s.ID]; f {
		return errors.New("status page " + s.ID + " already exists in environment " + s.Environment)
	}
	c.Environments[s.Environment].StatusPages[s.ID] = s
	return nil
}

func newEnv(env string) Environment {
	return Environment{ID: env,
		StatusPages: map[string]StatusPage{},
		Services:    map[string]Service{},
	}
}

func (c *Config) GetService(environment string, service string) Service {
	if _, f := c.Environments[environment]; !f {
		return Service{ID: service, Name: service, Environment: environment, Components: []ServiceComponent{}}
	}
	if _, f := c.Environments[environment].Services[service]; !f {
		return Service{ID: service, Name: service, Environment: environment, Components: []ServiceComponent{}}
	}
	return c.Environments[environment].Services[service]

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
	Components []ComponentStatus
}

type PageStatus struct {
	// ID of this page
	ID string
	// Status of this page overall
	Status Severity
	// Last update of this page
	LastUpdate time.Time
	// Services of this page
	Services map[string]ServiceStatus
}
