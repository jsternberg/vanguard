package vanguard

type rawPlan struct {
	Tasks []map[string]interface{}
}

type Plan struct {
	Tasks []*Task
}
