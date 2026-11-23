package frame

// ComponentConfig holds one component's settings and the metadata used by its loader.
type ComponentConfig struct {
	// Name is the component's configured name and archive filename.
	Name string `json:"名称"`
	// Description is the description shown to users.
	Description string `json:"描述"`
	// Disabled controls whether the host loads the component.
	Disabled bool `json:"是否禁用"`
	// Version is the component configuration version, used for migration.
	Version string `json:"版本"`
	// Source selects the host loader that owns this configuration.
	Source string `json:"来源"`
	// Configs contains the component-specific settings object.
	Configs       map[string]any `json:"配置"`
	attachedField map[string]string
	upgradeFn     func(*ComponentConfig) error
}

// Upgrade passes the current configuration to the host's write-back function. Call it after migrating
// settings; the host must have installed the function first.
func (c *ComponentConfig) Upgrade() error {
	return c.upgradeFn(c)
}

// SetAttachedField associates string metadata with this configuration. It is not written to JSON.
func (c *ComponentConfig) SetAttachedField(d map[string]string) {
	c.attachedField = d
}

// GetAttachedField returns the attached metadata map, or nil if none was set.
func (c *ComponentConfig) GetAttachedField() (d map[string]string) {
	return c.attachedField
}

// SetUpgradeFn installs the write-back function. It returns false if a function is already installed.
func (c *ComponentConfig) SetUpgradeFn(fn func(*ComponentConfig) error) (ok bool) {
	if c.upgradeFn == nil {
		c.upgradeFn = fn
		return true
	} else {
		return false
	}
}
