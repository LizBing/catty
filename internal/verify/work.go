package verify

import "fmt"

// Work is the deterministic verification-work budget for one admitted JAR.
// A nil Context.Work disables accounting for normal builds and Gate 2 tests.
type Work struct {
	Limit int64
	Used  int64
}

// WorkError reports the first operation that would exceed the budget.
type WorkError struct {
	Used      int64
	Limit     int64
	Operation string
}

func (e *WorkError) Error() string {
	return fmt.Sprintf("verification work limit %d exceeded before %s (already used %d)", e.Limit, e.Operation, e.Used)
}

// Charge reserves n units before the associated work or allocation. A rejected
// charge does not consume the budget.
func (w *Work) Charge(n int64, operation string) error {
	if w == nil {
		return nil
	}
	if n < 0 {
		return fmt.Errorf("negative verification work charge %d for %s", n, operation)
	}
	if n > w.Limit-w.Used {
		return &WorkError{Used: w.Used, Limit: w.Limit, Operation: operation}
	}
	w.Used += n
	return nil
}

// charge reserves work and records the first failure on the context.
func (c *Context) charge(n int64, operation string) error {
	if c == nil || c.Work == nil {
		return nil
	}
	if err := c.Work.Charge(n, operation); err != nil {
		if c.workErr == nil {
			c.workErr = err
		}
		return err
	}
	return nil
}

// chargeSoft records a budget failure but leaves the caller free to finish a
// boolean query conservatively; the analysis loop rechecks WorkError.
func (c *Context) chargeSoft(n int64, operation string) bool {
	return c.charge(n, operation) == nil
}

// WorkError returns the first recorded budget failure, if any.
func (c *Context) WorkError() error {
	if c == nil {
		return nil
	}
	return c.workErr
}
