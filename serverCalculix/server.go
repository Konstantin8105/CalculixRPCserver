package serverCalculix

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var modelName string

var ccxExecutionLocation []string

func init() {
	modelName = "model"

	// allowable locations of ccx
	ccxExecutionLocation = append(ccxExecutionLocation, "ccx")
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

// AmountTasks - amount allowable tasks to sending for calculation
func (c *Calculix) AmountTasks(empty string, amountFreeTasks *int) error {
	_ = empty
	amountFreeTasks = &c.amountTasks
	return nil
}

// ExecuteForDat - calculute by Calculix and return body of .dat file
func (c *Calculix) ExecuteForDat(inpFileBody string, datFileBody *[]string) error {
	c.amountTasks--
	defer func() { c.amountTasks++ }()
	// create temp folder
	dir, err := c.createNewTempDir()
	if err != nil {
		return err
	}
	// create inp file
	inpFilename := modelName + ".inp"
	file := dir + string(filepath.Separator) + inpFilename

	err = ioutil.WriteFile(file, []byte(inpFileBody), 0777)
	if err != nil {
		return fmt.Errorf("Cannot write to inp file: %v", err)
	}
	// check file is exist
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return fmt.Errorf("File inp is not exist : %v", err)
	}
	// try to execute by any ccx
	for _, ccx := range ccxExecutionLocation {
		// remove .INP in filename
		file = file[:(len(file) - 4)]
		// execute
		out, err := exec.Command(ccx, "-i", file).Output()
		if err != nil {
			return fmt.Errorf("Try install from https://pkgs.org/download/calculix-ccx\nError in calculix execution: %v\n%v", err, out)
		}
		datFileBody, err = c.getDatFileBody(dir)
		if err != nil {
			return fmt.Errorf("Cannot take .dat file: %v", err)
		}
		// remove temp folder
		err = os.RemoveAll(dir)
		if err != nil {
			return fmt.Errorf("Cannot remove folder: %v", err)
		}
		return nil
	}
	return fmt.Errorf("Cannot found ccx")
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

func (c *Calculix) getDatFileBody(dir string) (datBody *[]string, err error) {
	file := dir + string(filepath.Separator) + modelName + ".dat"
	if strings.ToUpper(file[len(file)-3:]) != "DAT" {
		return datBody, fmt.Errorf("Wrong filename : %v", file)
	}
	// check file is exist
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return datBody, fmt.Errorf("Cannot find .dat file : %v", err)
	}
	// open file
	inFile, err := os.Open(file)
	if err != nil {
		return datBody, err
	}
	defer func() {
		errFile := inFile.Close()
		if errFile != nil {
			if err != nil {
				err = fmt.Errorf("%v ; %v", err, errFile)
			} else {
				err = errFile
			}
		}
	}()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		*datBody = append(*datBody, scanner.Text())
	}
	return datBody, nil
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
