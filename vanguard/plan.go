package vanguard

type rawPlan struct {
	Tasks []map[string]interface{}
}

// A Plan to be executed by a Runner.
// The plan contains a list of tasks.
type Plan struct {
	Tasks []*Task
}
