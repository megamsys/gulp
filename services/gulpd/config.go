package gulpd

const (
	// DefaultAssemblyID.
	DefaultAssemblyID = "ASM00"
)

type Config struct {
	AssemblyID    string 			`toml:"assembly_id"`	
}

func (c Config) String() string {
	table := NewTable()
	table.AddRow(Row{Colorfy("Config:", "white", "", "bold"), Colorfy("Activity", "green", "", "")})
	table.AddRow(Row{"AssemblyID", c.AssemblyID})	
	table.AddRow(Row{"", ""})
	return table.String() 
}

func NewConfig() *Config {
	c := &Config{}
	c.AssemblyID = DefaultAssemblyID
	return c

	return &Config{
		AssemblyID:    DefaultAssemblyID,	
	}
}
