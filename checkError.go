package libagent

type Severity string

const (
	Low    Severity = "LOW"
	Medium Severity = "MEDIUM"
	High   Severity = "HIGH"
)

type CheckError struct {
	Sev Severity
	Err error
}

func (r *CheckError) Error() string {
	return r.Err.Error()
}
