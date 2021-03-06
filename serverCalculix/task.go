package serverCalculix

// Amount - amount
type Amount struct {
	A int
}

// MaxAllowableTasks - amount allowable tasks to sending for calculation
func (c *Calculix) MaxAllowableTasks(empty string, amountFreeTasks *Amount) error {
	_ = empty
	amountFreeTasks.A = c.maxAmountTask
	return nil
}

// AmountTasks - amount allowable tasks to sending for calculation
func (c *Calculix) AmountTasks(empty string, amountFreeTasks *Amount) error {
	_ = empty
	amountFreeTasks.A = c.amountTasks
	return nil
}
