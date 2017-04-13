package serverCalculix

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ExecuteForDat - calculute by Calculix and return body of .dat file
func (c *Calculix) ExecuteForDat(inpFileBody string, datFileBody *[]string) error {
	c.amountTasks--
	defer func() { c.amountTasks++ }()
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
		return nil
	}
	return fmt.Errorf("Cannot found ccx")
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
