package serverCalculix

import (
	"fmt"
	"math/rand"
)

// ServerName - name of server
type ServerName struct {
	A string
}

// GetServerName - get server name
func (c *Calculix) GetServerName(empty string, name *ServerName) error {
	_ = empty
	if len(c.serverName) == 0 {
		c.mutex.Lock()
		c.serverName = fmt.Sprintf("%v", rand.Int())
		c.mutex.Unlock()
	}
	name.A = c.serverName
	return nil
}
