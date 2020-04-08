package cli

import "fmt"

type Status struct {
	spinner       *Spinner
	status        string
	successFormat string
	failureFormat string
}

func NewStatus() *Status {
	s := &Status{
		spinner:       NewSpinner(),
		successFormat: "✔︎ %s\n",
		failureFormat: "x %s\n",
	}
	return s
}

func (s *Status) End(success bool) {
	if s.status == "" {
		return
	}

	if s.spinner != nil {
		s.spinner.Stop()
		fmt.Print("\r")
	}
	if success {
		fmt.Printf(s.successFormat, s.status)
	} else {
		fmt.Printf(s.failureFormat, s.status)
	}
	s.status = ""
}

func (s *Status) Start(status string) {
	s.End(true)

	s.status = status
	if s.spinner != nil {
		s.spinner.SetSuffix(fmt.Sprintf(" %s ", s.status))
		s.spinner.Start()
	}
}
