package test

import (
	"fmt"
	"github.com/cryptopunkscc/astrald/sig"
	"reflect"
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

// Initializes Runner to prepare it for Runner.run.
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
		r.initMissingTasks(t)
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

// Converts missing Test into Task and sets them to Runner.Map.
func (r *Runner) initMissingTasks(t *Task) {
	for _, test := range t.require() {
		tt := &Task{Test: test}
		r.initMissingTasks(tt)
	}
	if _, ok := r.Map.Get(t.Test.name); !ok {
		t.mu = &sync.Mutex{}
		r.Map.Set(t.Test.name, t)
	}
}

// Runs given Task and Task.Require recursively.
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
	name     string
	run      func(t *testing.T)
	requires Tests
}

func (test Test) Args(args ...any) Test {
	if len(args) == 1 {
		switch reflect.ValueOf(args[0]).Kind() {
		case reflect.Slice, reflect.Map, reflect.Struct:
			test.name = fmt.Sprintf("%s%v", test.name, args[0])
		default:
			test.name = fmt.Sprintf("%s(%v)", test.name, args[0])
		}
	} else if len(args) > 1 {
		test.name = fmt.Sprintf("%s%v", test.name, args)
	}
	return test
}

func (test Test) Requires(requires ...Test) Test {
	test.requires = append(test.requires, requires...)
	return test
}

func (test Test) Func(run func(t *testing.T), requires ...Test) Test {
	test.run = run
	test.requires = append(test.requires, requires...)
	return test
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
	return append(t.Test.requires, t.Require...)
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
