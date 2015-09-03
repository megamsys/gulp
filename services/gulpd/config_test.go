package gulpd

import (
	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
)

// Ensure the configuration can be parsed.
func (s *S) TestConfig_Parse(c *check.C) {
	// Parse configuration.
	var c activity.Config
	if _, err := toml.Decode(`
		assembly_id = "ASM000"
	
`, &c); err != nil {
		t.Fatal(err)
	}

	c.Assert(c.AssemblyID, check.Equals, "ASM000")	
}
