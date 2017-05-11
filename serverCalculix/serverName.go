package serverCalculix

import (
	"math/rand"
	"time"
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
		const charset = "abcdefghijklmnopqrstuvwxyz" +
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		var seededRand *rand.Rand = rand.New(
			rand.NewSource(time.Now().UnixNano()))
		b := make([]byte, 20)
		for i := range b {
			b[i] = charset[seededRand.Intn(len(charset))]
		}
		c.serverName = string(b)
		c.mutex.Unlock()
	}
	name.A = c.serverName
	return nil
}
