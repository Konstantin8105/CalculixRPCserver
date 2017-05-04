package serverCalculix

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// DatBody - body of dat file
type DatBody struct {
	A string
}

// ErrorServerBusy - system error of server
const (
	ErrorServerBusy string = "Server is busy"
)

// ExecuteForDat - calculute by Calculix and return body of .dat file
func (c *Calculix) ExecuteForDat(inpFileBody string, datFileBody *DatBody) error {

	if len(inpFileBody) == 0 {
		return fmt.Errorf("Input inp file is empty")
	}

	c.Lock()

	c.amountTasks--
	if c.amountTasks < 0 {
		c.amountTasks++
		c.Unlock()
		return fmt.Errorf(ErrorServerBusy)
	}
	c.Unlock()
	defer func() { c.amountTasks++ }()

	// create temp folder
	dir, err := c.createNewTempDir()
	if err != nil {
		return err
	}
	// remove temp folder
	defer func() {
		// TODO: remove only if calculation haven`t mistakes
		if err == nil {
			err2 := os.RemoveAll(dir)
			if err2 != nil {
				err = fmt.Errorf("Cannot remove folder: %v. Other: %v", err2, err)
			}
		}

	}()

	// create inp file
	inpFilename := modelName + ".inp"
	file := dir + string(filepath.Separator) + inpFilename
	// check file is exist
	if _, err := os.Stat(file); os.IsNotExist(err) {
		// create file
		newFile, err := os.Create(file)
		if err != nil {
			return err
		}
		err = newFile.Close()
		if err != nil {
			return err
		}
	}

	// open file
	f, err := os.OpenFile(file, os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, "%v", inpFileBody)
	err = f.Close()
	if err != nil {
		return err
	}

	// check file is exist
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return fmt.Errorf("File inp is not exist : %v", err)
	}

	// remove .INP in filename
	file = file[:(len(file) - 4)]

	// set amount of processors
	err = os.Setenv("NUMBER_OF_CPUS", string(runtime.GOMAXPROCS(runtime.NumCPU())))
	if err != nil {
		return fmt.Errorf("Cannot set environment for use all CPU by calculix")
	}

	var summError error
	var outs []string

	// try all posibile to execute by any ccx
	for _, ccx := range ccxExecutionLocation {
		// execute
		out, err := exec.Command(ccx, "-i", file).Output()
		if err != nil {
			summError = fmt.Errorf("%v\n%v", err, summError)
			outs = append(outs, string(out))
			continue
		}
		lines, err := c.getDatFileBody(dir)
		if err != nil {
			return fmt.Errorf("Cannot take .dat file: %v", err)
		}
		if len(lines) == 0 {
			return fmt.Errorf("Cannot read dat file in temp folder = %v", dir)
		}

		datFileBody.A = strings.Join(lines, "\n")
		return nil
	}
	return fmt.Errorf("Cannot found ccx %v\nOUT - %v", err, strings.Join(outs, "\n"))
}

func (c *Calculix) getDatFileBody(dir string) (datBody []string, err error) {
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
		datBody = append(datBody, scanner.Text())
	}

	return datBody, nil
}
