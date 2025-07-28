package test

import (
	"github.com/cryptopunkscc/astrald/sig"
	"runtime"
	"strings"
	"sync"
	"testing"
)

type Runner struct {
	Tasks []Task
	sig.Map[string, *Task]
	done []sync.WaitGroup
}

func (r *Runner) Run(tasks []Task, task Task) func(t *testing.T) {
	return func(t *testing.T) {
		if task.Name == "" {
			t.Helper()
			t.SkipNow()
		}
		r.Tasks = tasks
		r.init()
		t.Parallel()
		tt, _ := r.Get(task.Test.name)
		r.run(t, tt)
	}
}

func (r *Runner) init() {
	if r.Len() > 0 {
		return
	}
	groups := 0
	for _, t := range r.Tasks {
		t.mu = &sync.Mutex{}
		r.Map.Set(t.Test.name, &t)
		if t.Group > groups {
			groups = t.Group
		}
	}
	for _, t := range r.Map.Values() {
		r.initTask(t)
	}
	r.done = make([]sync.WaitGroup, groups+1)
	for _, t := range r.Map.Values() {
		t.tasks = make([]*Task, len(t.require()))
		for i, s := range t.require() {
			ok := false
			if t.tasks[i], ok = r.Map.Get(s.name); !ok {
				panic("WTF!!!")
			}
		}
	}
}

func (r *Runner) initTask(t *Task) {
	for _, test := range t.require() {
		tt := &Task{Test: test}
		r.initTask(tt)
	}
	if _, ok := r.Map.Get(t.Test.name); !ok {
		t.mu = &sync.Mutex{}
		r.Map.Set(t.Test.name, t)
	}
}

func (r *Runner) run(t *testing.T, task *Task) {
	run := func(t *testing.T) {
		r.done[task.Group].Add(1)
		defer r.done[task.Group].Done()

		for _, tt := range task.tasks {
			r.run(t, tt)
			if tt.status == Failure {
				t.FailNow()
			}
		}
		for i := 0; i < task.Group; i++ {
			r.done[i].Wait()
		}
		task.run(t, &r.done[task.Group])
	}
	if task.Test.name == "" {
		run(t)
	} else {
		t.Run(task.Test.name, run)
	}
}

type Test struct {
	name    string
	run     func(t *testing.T)
	Require Tests
}

func (test Test) Run(t *testing.T) {
	t.Run(test.name, test.run)
}

func New(name string, run func(t *testing.T), require ...Test) Test {
	return Test{name, run, require}
}

type Task struct {
	Name    string
	Require Tests
	tasks   []*Task
	Test    Test
	mu      *sync.Mutex
	status  status
	Group   int
}

type Tests []Test

type status int

const (
	Initial status = iota
	Success
	Failure
)

func (t *Task) require() []Test {
	return append(t.Test.Require, t.Require...)
}

func (t *Task) run(tt *testing.T, wg *sync.WaitGroup) {
	t.mu.Lock()
	if t.status != Initial {
		t.mu.Unlock()
		if t.status == Failure {
			tt.FailNow()
		}
		return
	}
	wg.Add(1)
	tt.Cleanup(func() {
		if tt.Failed() || tt.Skipped() {
			t.status = Failure
		} else {
			t.status = Success
		}
		wg.Done()
		t.mu.Unlock()
	})
	if t.Test.run != nil {
		t.Test.run(tt)
	}
}

func CallerName(depth ...int) (name string) {
	d := 1
	if len(depth) > 0 {
		d = depth[0]
	}
	if pc, _, _, ok := runtime.Caller(d); ok {
		if funcObj := runtime.FuncForPC(pc); funcObj != nil {
			name = funcObj.Name()
			c := strings.Split(name, ".")
			name = c[len(c)-1]
		}
	}
	return
}
