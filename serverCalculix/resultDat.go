package serverCalculix

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
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

	var mutex = &sync.Mutex{}
	mutex.Lock()

	c.amountTasks--
	defer func() { c.amountTasks++ }()
	if c.amountTasks < 0 {
		return fmt.Errorf(ErrorServerBusy)
	}

	mutex.Unlock()

	// create temp folder
	dir, err := c.createNewTempDir()
	if err != nil {
		return err
	}
	// remove temp folder
	defer func() {
		err2 := os.RemoveAll(dir)
		if err2 != nil {
			err = fmt.Errorf("Cannot remove folder: %v. Other: %v", err2, err)
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

	// try all posibile to execute by any ccx
	for _, ccx := range ccxExecutionLocation {
		// execute
		_, err := exec.Command(ccx, "-i", file).Output()
		if err != nil {
			continue
			//return fmt.Errorf("Try install from https://pkgs.org/download/calculix-ccx\nError in calculix execution: %v\n%v", err, out)
		}
		b, err := c.getDatFileBody(dir)
		if err != nil {
			return fmt.Errorf("Cannot take .dat file: %v", err)
		}

		var buffer string
		for _, s := range b {
			buffer += s + "\n"
		}

		datFileBody.A = buffer
		return nil
	}
	return fmt.Errorf("Cannot found ccx")
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
