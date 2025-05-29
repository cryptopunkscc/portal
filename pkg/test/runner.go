package test

import (
	"github.com/cryptopunkscc/astrald/sig"
	"sync"
	"testing"
)

type Runner struct {
	Tasks []Task
	sig.Map[string, *Task]
}

func (r *Runner) Run(tasks []Task, task Task) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		if task.Name == "" {
			t.Helper()
			t.SkipNow()
		}
		r.Tasks = tasks
		r.init()
		tt, _ := r.Get(task.Test.name)
		r.run(t, tt)
	}
}

func (r *Runner) init() {
	if r.Len() > 0 {
		return
	}
	for _, t := range r.Tasks {
		t.mu = &sync.Mutex{}
		r.Map.Set(t.Test.name, &t)
	}
	for _, t := range r.Map.Values() {
		t.tasks = make([]*Task, len(t.Require))
		for i, s := range t.Require {
			t.tasks[i], _ = r.Get(s.name)
		}
	}
}

func (r *Runner) run(t *testing.T, task *Task) {
	if task.status != Initial {
		return
	}
	t.Run(task.Test.name, func(t *testing.T) {
		for _, tt := range task.tasks {
			r.run(t, tt)
			if tt.status == Failure {
				t.FailNow()
			}
		}
		task.run(t)
	})
}

type Test struct {
	name string
	run  func(t *testing.T)
}

func (test Test) Run(t *testing.T) {
	t.Run(test.name, test.run)
}

func New(name string, run func(t *testing.T)) Test {
	return Test{name, run}
}

type Task struct {
	Name    string
	Require Tests
	tasks   []*Task
	Test    Test
	mu      *sync.Mutex
	status  status
}
type Tests []Test

type status int

const (
	Initial status = iota
	Success
	Failure
)

func (t *Task) run(tt *testing.T) {
	t.mu.Lock()
	if t.status != Initial {
		t.mu.Unlock()
		if t.status == Failure {
			tt.FailNow()
		}
		tt.Log(t.Name + " already done")
		return
	}
	tt.Cleanup(func() {
		tt.Log(t.Name + " done")
		if tt.Failed() || tt.Skipped() {
			t.status = Failure
		} else {
			t.status = Success
		}
		t.mu.Unlock()
	})
	t.Test.run(tt)
}
