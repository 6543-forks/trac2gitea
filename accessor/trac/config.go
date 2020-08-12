package trac

// GetStringConfig retrieves a value from the Trac config as a string.
func (accessor *Accessor) GetStringConfig(sectionName string, configName string) string {
	configValue, err := accessor.config.Section(sectionName).GetKey(configName)
	if err != nil {
		return ""
	}

	return configValue.String()
}