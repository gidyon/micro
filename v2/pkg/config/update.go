package config

func (cfg *config) updateConfigWith(newCfg *config) {
	// Set config only if zero value
	cfg.ServiceName = setStringIfEmpty(cfg.ServiceName, newCfg.ServiceName)
	cfg.ServiceType = setStringIfEmpty(cfg.ServiceType, newCfg.ServiceType)
	cfg.HTTPort = setIntIfZero(cfg.HTTPort, newCfg.HTTPort)
	cfg.GRPCPort = setIntIfZero(cfg.GRPCPort, newCfg.GRPCPort)
	cfg.StartupSleepSeconds = setIntIfZero(cfg.StartupSleepSeconds, newCfg.StartupSleepSeconds)

	// Service log
	cfg.LogLevel = setLogLevl(cfg.LogLevel, newCfg.LogLevel)

	// Service security
	if newCfg.Security != nil {
		cfg.Security.TLSCertFile = setStringIfEmpty(cfg.Security.TLSCertFile, newCfg.Security.TLSCertFile)
		cfg.Security.TLSKeyFile = setStringIfEmpty(cfg.Security.TLSKeyFile, newCfg.Security.TLSKeyFile)
		cfg.Security.ServerName = setStringIfEmpty(cfg.Security.ServerName, newCfg.Security.ServerName)
		cfg.Security.Insecure = setBoolIfEmpty(cfg.Security.Insecure, newCfg.Security.Insecure)
	}

	if newCfg.HttpOtions != nil {
		cfg.HttpOtions.CorsEnabled = setBoolIfEmpty(cfg.HttpOtions.CorsEnabled, newCfg.HttpOtions.CorsEnabled)
	}

	// Update databases options
	if len(newCfg.Databases) > 0 {
		cfg.Databases = newCfg.Databases
	}

	// External services
	if len(newCfg.ExternalServices) != 0 {
		// cfg.ExternalServices
		cfg.ExternalServices = newCfg.ExternalServices
	}
}

func setStringIfEmpty(def, val string) string {
	if val == "" {
		return def
	}
	return val
}

func setBoolIfEmpty(def, val bool) bool {
	if val {
		return val
	}
	return def
}

func setIntIfZero(def, val int) int {
	if val == 0 {
		return def
	}
	return val
}

func setLogLevl(def, val int) int {
	if val == unknownLevel {
		return def
	}
	return val
}
