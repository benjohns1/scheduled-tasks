package scheduler

import (
	"testing"
	"time"

	"github.com/benjohns1/scheduled-tasks/internal/core/schedule"
	"github.com/benjohns1/scheduled-tasks/internal/data/transient"
	"github.com/benjohns1/scheduled-tasks/internal/usecase"
)

type loggerStub struct{}

func (l *loggerStub) Printf(format string, v ...interface{}) {}

type ClockMock struct {
	mockNow func() time.Time
}

func NewClockMock(now func() time.Time) *ClockMock {
	return &ClockMock{mockNow: now}
}

func NewStaticClockMock(now time.Time) *ClockMock {
	return &ClockMock{mockNow: func() time.Time { return now }}
}

// Now implementes the standard time function
func (c *ClockMock) Now() time.Time {
	return c.mockNow()
}

// After implementes the standard time function
func (c *ClockMock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

// Sleep implementes the standard time function
func (c *ClockMock) Sleep(d time.Duration) {
	time.Sleep(d)
}

// Tick implementes the standard time function
func (c *ClockMock) Tick(d time.Duration) <-chan time.Time {
	return time.Tick(d)
}

// Since implementes the standard time function
func (c *ClockMock) Since(t time.Time) time.Duration {
	return c.Now().Sub(t)
}

// Until implementes the standard time function
func (c *ClockMock) Until(t time.Time) time.Duration {
	return t.Sub(c.Now())
}

func TestRun(t *testing.T) {
	now := time.Now()
	immediateTimeout := 10 * time.Nanosecond

	type args struct {
		l            Logger
		c            usecase.Clock
		taskRepo     usecase.TaskRepo
		scheduleRepo usecase.ScheduleRepo
	}
	type resp struct {
		close  chan<- bool
		closed <-chan bool
		next   <-chan time.Time
	}
	tests := []struct {
		name    string
		arrange func() args
		assert  func(args, resp)
	}{
		{
			name: "empty schedule should close scheduler immediately",
			arrange: func() args {
				return args{
					l:            &loggerStub{},
					c:            NewStaticClockMock(now),
					taskRepo:     transient.NewTaskRepo(),
					scheduleRepo: transient.NewScheduleRepo(),
				}
			},
			assert: func(_ args, r resp) {
				select {
				case val := <-r.closed:
					if val {
						return
					}
					t.Errorf("scheduler.Run() closed channel should be true")
				case <-time.After(immediateTimeout):
					t.Errorf("scheduler.Run() should be closed immediately")
				}
			},
		},
		{
			name: "schedule with no upcoming recurrences should close scheduler immediately",
			arrange: func() args {
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewHourFrequency([]int{})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)
				return args{
					l:            &loggerStub{},
					c:            NewStaticClockMock(now),
					taskRepo:     transient.NewTaskRepo(),
					scheduleRepo: transient.NewScheduleRepo(),
				}
			},
			assert: func(_ args, r resp) {
				select {
				case val := <-r.closed:
					if val {
						return
					}
					t.Errorf("scheduler.Run() closed channel should be true")
				case <-time.After(immediateTimeout):
					t.Errorf("scheduler.Run() should be closed immediately")
				}
			},
		},
		{
			name: "schedule with a past recurring time should create a new task",
			arrange: func() args {
				now := time.Date(2000, time.January, 1, 12, 30, 0, 0, time.UTC)
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewHourFrequency([]int{25})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)
				return args{
					l:            &loggerStub{},
					c:            NewStaticClockMock(now),
					taskRepo:     transient.NewTaskRepo(),
					scheduleRepo: sr,
				}
			},
			assert: func(a args, r resp) {
				select {
				case <-r.closed:
					t.Errorf("scheduler.Run() should have created a task before closing")
					return
				case next := <-r.next:
					want := time.Date(2000, time.January, 1, 13, 25, 0, 0, time.UTC).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() next run time should be the next scheduled time plus scheduler.Offset, got = %v, want = %v", next, want)
						return
					}
					break
				}

				tasks, err := a.taskRepo.GetAll()
				if err != nil {
					t.Errorf("scheduler.Run() should have created 1 task, but there was an error retrieving task: %v", err)
					return
				}
				if len(tasks) != 1 {
					t.Errorf("scheduler.Run() should have created 1 task, but there were %v tasks in repo", len(tasks))
					return
				}
				task, ok := tasks[1]
				if !ok {
					t.Errorf("scheduler.Run() should have created 1 task of ID 1, but there was no task with ID 1 in repo")
					return
				}
				if task.Name() != "t1" || task.Description() != "t1desc" {
					t.Errorf("scheduler.Run() should have created 1 task, but task values were not correct")
					return
				}
				r.close <- true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.arrange()
			close, closed, next := Run(args.l, args.c, args.taskRepo, args.scheduleRepo)
			tt.assert(args, resp{close, closed, next})
		})
	}
}
