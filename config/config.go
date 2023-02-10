package config

// configurations exported
type Configurations struct {
	Schedule ScheConfigurations
	Monitor  MonitorConfigurations
}

// ScheConfigurations exported
type ScheConfigurations struct {
	Cron int
}

// MonitorConfigurations exported
type MonitorConfigurations struct {
	Monitor []string
}
