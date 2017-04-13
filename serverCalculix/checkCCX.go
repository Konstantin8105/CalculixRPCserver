package serverCalculix

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ChechCCXResult - result of CheckCCX
type ChechCCXResult struct {
	A bool
}

// CheckCCX - check ccx is allowable
func (c *Calculix) CheckCCX(empty string, result *ChechCCXResult) error {
	for _, ccx := range ccxExecutionLocation {
		fmt.Println("ccx = ", ccx)
		out, err := exec.Command(ccx).Output()
		//if err != nil {
		//	return fmt.Errorf("Try install from https://pkgs.org/download/calculix-ccx\nError in calculix execution: %v\n%v", err, out)
		//}
		fmt.Println("out = ", out)
		fmt.Println("err = ", err)

		cmd := exec.Command(ccx)
		cmd.Stdin = strings.NewReader("")
		var out2 bytes.Buffer
		cmd.Stdout = &out2
		err = cmd.Run()
		if err != nil {
			fmt.Println("err ccx ---> ", err)
		}
		fmt.Printf("in all caps: %q\n", out2.String())

	}
	result.A = true
	return nil
}
