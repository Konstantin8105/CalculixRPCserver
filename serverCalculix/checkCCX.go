package serverCalculix

import (
	"os/exec"
)

// ChechCCXResult - result of CheckCCX
type ChechCCXResult struct {
	A bool
}

const (
	ccxCorrectAnswer = "Usage: CalculiX.exe -i jobname\n"
)

// CheckCCX - check ccx is allowable
func (c *Calculix) CheckCCX(empty string, result *ChechCCXResult) error {
	for _, ccx := range ccxExecutionLocation {
		out, err := exec.Command(ccx).Output()
		if err.Error() == "exit status 201" && string(out) == ccxCorrectAnswer {
			result.A = true
			break
		}
	}
	return nil
}
