package scheduler

import (
	"fmt"
	"testing"
	"time"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/clock"
	"github.com/benjohns1/scheduled-tasks/services/internal/core/schedule"
	"github.com/benjohns1/scheduled-tasks/services/internal/data/transient"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

type loggerStub struct{}

func (l *loggerStub) Printf(format string, v ...interface{}) {
	if testing.Verbose() {
		fmt.Printf(fmt.Sprintf("    LOG: %v\n", format), v...)
	}
}

func closeNonBlocking(ch chan<- bool) {
	select {
	case ch <- true:
	default:
	}
}

func TestRun(t *testing.T) {

	timeout := 10 * time.Millisecond

	now := time.Now()
	prevClock := clock.Get()
	clock.Set(clock.NewStaticMock(now))
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
				clock.Set(clock.NewStaticMock(testNow))

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
				clock.Set(clock.NewStaticMock(testNow))

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
				f, err := schedule.NewHourFrequency([]int{5, 10, 15})
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
			defer closeNonBlocking(close)
			tt.assert(t, args, resp{close, check, closed})
		})
	}
}

func TestHourFrequencyIntervalOffsets(t *testing.T) {

	timeout := 10 * time.Millisecond

	now := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	prevClock := clock.Get()
	clock.Set(clock.NewStaticMock(now))
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
			name: "hour frequency with default interval and offset should be scheduled to run at +0:05",
			arrange: func(t *testing.T) args {
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewHourFrequency([]int{5})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)

				return args{
					l:            &loggerStub{},
					taskRepo:     transient.NewTaskRepo(),
					scheduleRepo: sr,
					nextRun:      make(chan time.Time),
				}
			},
			assert: func(t *testing.T, a args, r resp) {
				select {
				case <-r.closed:
					t.Errorf("scheduler.Run() should not have closed")
					return
				case next := <-a.nextRun:
					want := now.Add(5 * time.Minute).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() unexpected next run time, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
				}
			},
		},
		{
			name: "hour frequency with default interval and offset of 1 should be scheduled to run at +1:05",
			arrange: func(t *testing.T) args {
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewHourFrequency([]int{5})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				f.SetOffset(1)
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)

				return args{
					l:            &loggerStub{},
					taskRepo:     transient.NewTaskRepo(),
					scheduleRepo: sr,
					nextRun:      make(chan time.Time),
				}
			},
			assert: func(t *testing.T, a args, r resp) {
				select {
				case <-r.closed:
					t.Errorf("scheduler.Run() should not have closed")
					return
				case next := <-a.nextRun:
					want := now.Add(1*time.Hour + 5*time.Minute).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() unexpected next run time, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
				}
			},
		},
		{
			name: "hour frequency with interval of 2 and default offset should be scheduled to run at +2:05",
			arrange: func(t *testing.T) args {
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewHourFrequency([]int{5})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				f.SetInterval(2)
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)

				prevClock := clock.Get()
				checkTime := time.Date(2000, time.January, 1, 0, 10, 0, 0, time.UTC)
				clock.Set(clock.NewStaticMock(checkTime))

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
					t.Errorf("scheduler.Run() should not have closed")
					return
				case next := <-a.nextRun:
					want := now.Add(2*time.Hour + 5*time.Minute).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() unexpected next run time, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
				}
			},
		},
		{
			name: "hour frequency with interval of 2 and offset 1 should be scheduled to run at +3:05",
			arrange: func(t *testing.T) args {
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewHourFrequency([]int{5})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				f.SetInterval(2)
				f.SetOffset(1)
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)

				prevClock := clock.Get()
				checkTime := time.Date(2000, time.January, 1, 1, 10, 0, 0, time.UTC)
				clock.Set(clock.NewStaticMock(checkTime))

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
					t.Errorf("scheduler.Run() should not have closed")
					return
				case next := <-a.nextRun:
					want := now.Add(3*time.Hour + 5*time.Minute).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() unexpected next run time, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
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
			defer closeNonBlocking(close)
			tt.assert(t, args, resp{close, check, closed})
		})
	}
}

func TestDayFrequency(t *testing.T) {

	timeout := 10 * time.Millisecond

	now := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	prevClock := clock.Get()
	clock.Set(clock.NewStaticMock(now))
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
			name: "day frequency should be scheduled to run at 2000-01-01 00:05",
			arrange: func(t *testing.T) args {
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewDayFrequency([]int{5}, []int{0})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)

				return args{
					l:            &loggerStub{},
					taskRepo:     transient.NewTaskRepo(),
					scheduleRepo: sr,
					nextRun:      make(chan time.Time),
				}
			},
			assert: func(t *testing.T, a args, r resp) {
				select {
				case <-r.closed:
					t.Errorf("scheduler.Run() should not have closed")
					return
				case next := <-a.nextRun:
					want := now.Add(5 * time.Minute).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() unexpected next run time, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
				}
			},
		},
		{
			name: "day frequency with interval of 2 and offset 1 should be scheduled to run at 2000-01-04 00:05",
			arrange: func(t *testing.T) args {
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewDayFrequency([]int{5}, []int{0})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				f.SetInterval(2)
				f.SetOffset(1)
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)

				prevClock := clock.Get()
				checkTime := time.Date(2000, time.January, 2, 1, 10, 0, 0, time.UTC)
				clock.Set(clock.NewStaticMock(checkTime))

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
					t.Errorf("scheduler.Run() should not have closed")
					return
				case next := <-a.nextRun:
					want := now.AddDate(0, 0, 3).Add(5 * time.Minute).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() unexpected next run time, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
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
			defer closeNonBlocking(close)
			tt.assert(t, args, resp{close, check, closed})
		})
	}
}

func TestWeekFrequency(t *testing.T) {

	timeout := 100 * time.Millisecond

	now := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	prevClock := clock.Get()
	clock.Set(clock.NewStaticMock(now))
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
			name: "sunday week frequency should be scheduled to run at 2000-01-02 00:05",
			arrange: func(t *testing.T) args {
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewWeekFrequency([]int{5}, []int{0}, []time.Weekday{time.Sunday})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)

				return args{
					l:            &loggerStub{},
					taskRepo:     transient.NewTaskRepo(),
					scheduleRepo: sr,
					nextRun:      make(chan time.Time),
				}
			},
			assert: func(t *testing.T, a args, r resp) {
				select {
				case <-r.closed:
					t.Errorf("scheduler.Run() should not have closed")
					return
				case next := <-a.nextRun:
					want := time.Date(2000, time.January, 2, 0, 5, 0, 0, time.UTC).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() unexpected next run time, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
				}
			},
		},
		{
			name: "sunday week frequency with interval of 2 and offset 1 should be scheduled to run at 2000-01-16 00:05",
			arrange: func(t *testing.T) args {
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewWeekFrequency([]int{5}, []int{0}, []time.Weekday{time.Sunday})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				f.SetInterval(2)
				f.SetOffset(1)
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)

				prevClock := clock.Get()
				checkTime := time.Date(2000, time.January, 3, 1, 10, 0, 0, time.UTC)
				clock.Set(clock.NewStaticMock(checkTime))

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
					t.Errorf("scheduler.Run() should not have closed")
					return
				case next := <-a.nextRun:
					want := time.Date(2000, time.January, 16, 0, 5, 0, 0, time.UTC).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() unexpected next run time, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
				}
			},
		},
		{
			name: "monday week frequency with interval of 2 and offset 1 should be scheduled to run at 2000-01-10 00:05",
			arrange: func(t *testing.T) args {
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewWeekFrequency([]int{5}, []int{0}, []time.Weekday{time.Monday})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				f.SetInterval(2)
				f.SetOffset(1)
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)

				prevClock := clock.Get()
				checkTime := time.Date(2000, time.January, 3, 1, 10, 0, 0, time.UTC)
				clock.Set(clock.NewStaticMock(checkTime))

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
					t.Errorf("scheduler.Run() should not have closed")
					return
				case next := <-a.nextRun:
					want := time.Date(2000, time.January, 10, 0, 5, 0, 0, time.UTC).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() unexpected next run time, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
				}
			},
		},
		{
			name: "tuesday week frequency with interval of 2 and offset 1 should be scheduled to run at 2000-01-11 00:05",
			arrange: func(t *testing.T) args {
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewWeekFrequency([]int{5}, []int{0}, []time.Weekday{time.Tuesday})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				f.SetInterval(2)
				f.SetOffset(1)
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)

				prevClock := clock.Get()
				checkTime := time.Date(2000, time.January, 3, 0, 10, 0, 0, time.UTC)
				clock.Set(clock.NewStaticMock(checkTime))

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
					t.Errorf("scheduler.Run() should not have closed")
					return
				case next := <-a.nextRun:
					want := time.Date(2000, time.January, 11, 0, 5, 0, 0, time.UTC).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() unexpected next run time, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
				}
			},
		},
		{
			name: "thursday, friday week frequency with interval of 2 and offset 1 should be scheduled to run at 2000-01-13 00:05",
			arrange: func(t *testing.T) args {
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewWeekFrequency([]int{5}, []int{0}, []time.Weekday{time.Thursday, time.Friday})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				f.SetInterval(2)
				f.SetOffset(1)
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)

				prevClock := clock.Get()
				checkTime := time.Date(2000, time.January, 3, 1, 10, 0, 0, time.UTC)
				clock.Set(clock.NewStaticMock(checkTime))

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
					t.Errorf("scheduler.Run() should not have closed")
					return
				case next := <-a.nextRun:
					want := time.Date(2000, time.January, 13, 0, 5, 0, 0, time.UTC).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() unexpected next run time, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
				}
			},
		},
		{
			name: "saturday week frequency with interval of 2 and offset 1 should be scheduled to run at 2000-01-15 00:05",
			arrange: func(t *testing.T) args {
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewWeekFrequency([]int{5}, []int{0}, []time.Weekday{time.Saturday})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				f.SetInterval(2)
				f.SetOffset(1)
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)

				prevClock := clock.Get()
				checkTime := time.Date(2000, time.January, 3, 1, 10, 0, 0, time.UTC)
				clock.Set(clock.NewStaticMock(checkTime))

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
					t.Errorf("scheduler.Run() should not have closed")
					return
				case next := <-a.nextRun:
					want := time.Date(2000, time.January, 15, 0, 5, 0, 0, time.UTC).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() unexpected next run time, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
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
			defer closeNonBlocking(close)
			tt.assert(t, args, resp{close, check, closed})
		})
	}
}

func TestMonthFrequency(t *testing.T) {

	timeout := 10 * time.Millisecond

	now := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	prevClock := clock.Get()
	clock.Set(clock.NewStaticMock(now))
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
			name: "month frequency should be scheduled to run at 2000-01-01 00:05",
			arrange: func(t *testing.T) args {
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewMonthFrequency([]int{5}, []int{0}, []int{1})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)

				return args{
					l:            &loggerStub{},
					taskRepo:     transient.NewTaskRepo(),
					scheduleRepo: sr,
					nextRun:      make(chan time.Time),
				}
			},
			assert: func(t *testing.T, a args, r resp) {
				select {
				case <-r.closed:
					t.Errorf("scheduler.Run() should not have closed")
					return
				case next := <-a.nextRun:
					want := now.Add(5 * time.Minute).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() unexpected next run time, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
				}
			},
		},
		{
			name: "month frequency with interval of 2 and offset 1 should be scheduled to run at 2000-04-01 00:05",
			arrange: func(t *testing.T) args {
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewMonthFrequency([]int{5}, []int{0}, []int{1})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				f.SetInterval(2)
				f.SetOffset(1)
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)

				prevClock := clock.Get()
				checkTime := time.Date(2000, time.February, 1, 0, 10, 0, 0, time.UTC)
				clock.Set(clock.NewStaticMock(checkTime))

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
					t.Errorf("scheduler.Run() should not have closed")
					return
				case next := <-a.nextRun:
					want := time.Date(2000, time.April, 1, 0, 5, 0, 0, time.UTC).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() unexpected next run time, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
				}
			},
		},
		{
			name: "month frequency with offset 1 and February month day overrun should be scheduled to run at 2000-03-02 00:05",
			arrange: func(t *testing.T) args {
				sr := transient.NewScheduleRepo()
				f, err := schedule.NewMonthFrequency([]int{5}, []int{0}, []int{31})
				if err != nil {
					t.Fatalf("error creating frequency: %v", err)
				}
				f.SetInterval(2)
				f.SetOffset(1)
				s := schedule.New(f)
				s.AddTask(schedule.NewRecurringTask("t1", "t1desc"))
				sr.Add(s)

				return args{
					l:            &loggerStub{},
					taskRepo:     transient.NewTaskRepo(),
					scheduleRepo: sr,
					nextRun:      make(chan time.Time),
				}
			},
			assert: func(t *testing.T, a args, r resp) {
				select {
				case <-r.closed:
					t.Errorf("scheduler.Run() should not have closed")
					return
				case next := <-a.nextRun:
					want := time.Date(2000, time.March, 2, 0, 5, 0, 0, time.UTC).Add(Offset)
					if !next.Equal(want) {
						t.Errorf("scheduler.Run() unexpected next run time, got = %v, want = %v", next, want)
						return
					}
				case <-time.After(timeout):
					t.Errorf("scheduler.Run() should have scheduled next run before %v timeout", timeout)
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
			defer closeNonBlocking(close)
			tt.assert(t, args, resp{close, check, closed})
		})
	}
}
