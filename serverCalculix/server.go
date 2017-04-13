package serverCalculix

import (
	"fmt"
	"os"
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
		dir = fmt.Sprintf("Temp (%v)", i)
		if _, err = os.Stat(dir); os.IsExist(err) {
			err = os.Mkdir(dir, 0777)
			if err != nil {
				continue
			}
			return dir, nil
		}
	}
	return "", fmt.Errorf("Cannot create temp folder: %v", err)
}

/*
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	bucklingHeader := "B U C K L I N G   F A C T O R   O U T P U T"
	var found bool
	var numberLine int
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if !found {
			// empty line
			if len(line) == 0 {
				continue
			}

			if len(line) != len(bucklingHeader) {
				continue
			}

			if line == bucklingHeader {
				found = true
			}
		} else {
			numberLine++
			if numberLine >= 5+amountBuckling {
				break
			}
			if numberLine >= 5 {
				m, f, err := parseBucklingFactor(line)
				if err != nil {
					return bucklingFactor, err
				}
				if m != numberLine-4 {
					return bucklingFactor, fmt.Errorf("Wrong MODE NO: %v (%v) in line: %v", m, numberLine-4, line)
				}
				bucklingFactor = append(bucklingFactor, f)
			}
		}
	}
	if len(bucklingFactor) != amountBuckling {
		return bucklingFactor, fmt.Errorf("Wrong lenght of buckling lines in DAT file")
	}
	return bucklingFactor, nil
}

// Example:
//      4   0.4067088E+03
func parseBucklingFactor(line string) (mode int, factor float64, err error) {
	s := strings.Split(line, "   ")
	for i := range s {
		s[i] = strings.TrimSpace(s[i])
	}
	i, err := strconv.ParseInt(s[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("Error: string parts - %v, error - %v, in line - %v", s, err, line)
	}
	mode = int(i)

	factor, err = strconv.ParseFloat(s[1], 64)
	if err != nil {

		return 0, 0, fmt.Errorf("Error: string parts - %v, error - %v, in line - %v", s, err, line)
	}
	return mode, factor, nil
}
*/
