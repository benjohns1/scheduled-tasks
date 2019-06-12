package scheduler

import (
	"testing"
	"time"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/clock"
	"github.com/benjohns1/scheduled-tasks/services/internal/core/schedule"
	"github.com/benjohns1/scheduled-tasks/services/internal/data/transient"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

type loggerStub struct{}

func (l *loggerStub) Printf(format string, v ...interface{}) {}

func TestRun(t *testing.T) {

	timeout := 10 * time.Millisecond

	now := time.Now()
	prevClock := clock.Get()
	clockMock := clock.NewStaticMock(now)
	clock.Set(clockMock)
	defer clock.Set(prevClock)

	type args struct {
		l            Logger
		taskRepo     usecase.TaskRepo
		scheduleRepo usecase.ScheduleRepo
		nextRun      chan time.Time
		prevClock    clock.Time
	}
	type resp struct {
		close  chan<- bool
		check  chan<- bool
		closed <-chan bool
	}
	tests := []struct {
		name    string
		arrange func(*testing.T) args
		assert  func(*testing.T, args, resp)
	}{
		{
			name: "empty schedule should be scheduled to run again after default wait time",
			arrange: func(t *testing.T) args {
				return args{
					l:            &loggerStub{},
					taskRepo:     transient.NewTaskRepo(),
					scheduleRepo: transient.NewScheduleRepo(),
					nextRun:      make(chan time.Time),
				}
			},
			assert: func(t *testing.T, a args, r resp) {
				select {
				case <-r.closed:
					t.Errorf("scheduler.Run() should not have closed")
					return
				case next := <-a.nextRun:
					want := now.Add(DefaultWait).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() next run time should be the default wait time plus scheduler.Offset, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
				}
			},
		},
		{
			name: "schedule with no upcoming recurrences should be scheduled to run again after default wait time",
			arrange: func(t *testing.T) args {
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
					taskRepo:     transient.NewTaskRepo(),
					scheduleRepo: transient.NewScheduleRepo(),
					nextRun:      make(chan time.Time),
				}
			},
			assert: func(t *testing.T, a args, r resp) {
				select {
				case <-r.closed:
					t.Errorf("scheduler.Run() should not have closed")
					return
				case next := <-a.nextRun:
					want := now.Add(DefaultWait).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() next run time should be the default wait time plus scheduler.Offset, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
				}
			},
		},
		{
			name: "unchecked schedule with a past recurring time should not create a task",
			arrange: func(t *testing.T) args {
				testNow := time.Date(2000, time.January, 1, 12, 30, 0, 0, time.UTC)

				sr := transient.NewScheduleRepo()
				f, err := schedule.NewHourFrequency([]int{25})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)

				prevClock := clock.Get()
				clockMock := clock.NewStaticMock(testNow)
				clock.Set(clockMock)

				return args{
					l:            &loggerStub{},
					taskRepo:     transient.NewTaskRepo(),
					scheduleRepo: sr,
					nextRun:      make(chan time.Time),
					prevClock:    prevClock,
				}
			},
			assert: func(t *testing.T, a args, r resp) {
				select {
				case <-r.closed:
					t.Errorf("scheduler.Run() not have closed yet")
					return
				case next := <-a.nextRun:
					want := time.Date(2000, time.January, 1, 13, 25, 0, 0, time.UTC).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() next run time should be the next scheduled time plus scheduler.Offset, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
					return
				}

				tasks, err := a.taskRepo.GetAll()
				if err != nil {
					t.Errorf("scheduler.Run() error retrieving tasks: %v", err)
					return
				}
				if len(tasks) != 0 {
					t.Errorf("scheduler.Run() should not have created any tasks, but there were %v tasks in repo", len(tasks))
					return
				}
			},
		},
		{
			name: "unchecked schedule should be checked",
			arrange: func(t *testing.T) args {
				testNow := time.Date(2000, time.January, 1, 12, 30, 0, 0, time.UTC)
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewHourFrequency([]int{25})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)

				prevClock := clock.Get()
				clockMock := clock.NewStaticMock(testNow)
				clock.Set(clockMock)

				return args{
					l:            &loggerStub{},
					taskRepo:     transient.NewTaskRepo(),
					scheduleRepo: sr,
					nextRun:      make(chan time.Time),
					prevClock:    prevClock,
				}
			},
			assert: func(t *testing.T, a args, r resp) {
				select {
				case <-r.closed:
					t.Errorf("scheduler.Run() should not have closed yet")
					return
				case next := <-a.nextRun:
					want := time.Date(2000, time.January, 1, 13, 25, 0, 0, time.UTC).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() next run time should be the next scheduled time plus scheduler.Offset, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
					return
				}

				schedules, err := a.scheduleRepo.GetAll()
				if err != nil {
					t.Errorf("scheduler.Run() error retrieving schedules: %v", err)
					return
				}
				if len(schedules) != 1 {
					t.Errorf("scheduler.Run() expected 1 schedule")
					return
				}
				for _, s := range schedules {
					if s.LastChecked() != clock.Now() {
						t.Errorf("scheduler.Run() should have set checked time for schedule")
						return
					}
				}
			},
		},
		{
			name: "after unpausing a paused schedule, recurring tasks should not be created during paused time period",
			arrange: func(t *testing.T) args {

				prevClock := clock.Get()
				
				// setup schedule repo
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewHourFrequency([]int{5,10,15})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)

				firstCheckTime := time.Date(2000, time.January, 1, 12, 1, 0, 0, time.UTC)
				s.Check(firstCheckTime)
				s.Pause()

				// should NOT create a task at 12:35
				unpauseTime := time.Date(2000, time.January, 1, 12, 7, 0, 0, time.UTC)
				clock.Set(clock.NewStaticMock(unpauseTime))
				s.Unpause()

				// should create task at 12:40
				checkNow := time.Date(2000, time.January, 1, 12, 11, 0, 0, time.UTC)
				clock.Set(clock.NewStaticMock(checkNow))

				return args{
					l:            &loggerStub{},
					taskRepo:     transient.NewTaskRepo(),
					scheduleRepo: sr,
					nextRun:      make(chan time.Time),
					prevClock:    prevClock,
				}
			},
			assert: func(t *testing.T, a args, r resp) {
				select {
				case <-r.closed:
					t.Errorf("scheduler.Run() should not have closed yet")
					return
				case next := <-a.nextRun:
					want := time.Date(2000, time.January, 1, 12, 15, 0, 0, time.UTC).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() next run time should be the next scheduled time plus scheduler.Offset, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
					return
				}

				tasks, err := a.taskRepo.GetAll()
				if err != nil {
					t.Errorf("scheduler.Run() there was an error retrieving tasks: %v", err)
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.arrange(t)
			if args.prevClock != nil {
				defer clock.Set(args.prevClock)
			}
			close, check, closed := Run(args.l, args.taskRepo, args.scheduleRepo, args.nextRun)
			defer func() {
				close <- true
			}()
			tt.assert(t, args, resp{close, check, closed})
		})
	}
}
