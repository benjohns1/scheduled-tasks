package schedule

import (
	"testing"
)

func TestRecurringTask_Equal(t *testing.T) {
	rt1 := NewRecurringTask("task 1", "desc")
	rt1dupe := NewRecurringTask("task 1", "desc")
	rt2 := NewRecurringTask("task 2", "desc")
	rt2b := NewRecurringTask("task 2", "different description")

	type args struct {
		rtc RecurringTask
	}
	tests := []struct {
		name string
		rt   *RecurringTask
		args args
		want bool
	}{
		{
			name: "same recurring task should be equal to itself",
			rt:   &rt1,
			args: args{rtc: rt1},
			want: true,
		},
		{
			name: "duplicate recurring tasks should be equal",
			rt:   &rt1,
			args: args{rtc: rt1dupe},
			want: true,
		},
		{
			name: "recurring tasks with different name should be different",
			rt:   &rt1,
			args: args{rtc: rt2},
			want: false,
		},
		{
			name: "recurring tasks with different description should be different",
			rt:   &rt2,
			args: args{rtc: rt2b},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rt.Equal(tt.args.rtc); got != tt.want {
				t.Errorf("RecurringTask.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
