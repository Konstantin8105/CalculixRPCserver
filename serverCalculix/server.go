package serverCalculix

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	modelName string = "model"
)

// allowable locations of ccx
var ccxExecutionLocation []string

func init() {
	ccxExecutionLocation = []string{
		"ccx",
		"E:\\CalculiX\\bin\\ccx.bat",
	}
}

// Calculix - general type
type Calculix struct {
	// amount tasks calculated at one moment is equal amount of real processors
	amountTasks int
}

// NewCalculix - constructor for Calculix
func NewCalculix() *Calculix {
	calculix := new(Calculix)
	calculix.amountTasks = runtime.GOMAXPROCS(runtime.NumCPU())
	return calculix
}

func (c *Calculix) createNewTempDir() (dir string, err error) {
	for i := 0; i < 100000; i++ {
		dir = string(".") + string(filepath.Separator) + fmt.Sprintf("Temp(%v)", i)
		if _, err = os.Stat(dir); os.IsNotExist(err) {
			err = os.Mkdir(dir, 0777)
			if err != nil {
				continue
			}
			return dir, nil
		}
	}
	return "", fmt.Errorf("Cannot create temp folder: %v", err)
}
