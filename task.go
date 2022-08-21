package restart

import "time"

type Task struct {
	isGo bool
	task func()
}

var beforeTask []Task

var afterTask []func()

func AddBeforeTask(task func(), isGo bool) {
	beforeTask = append(beforeTask, Task{task: task, isGo: isGo})
}

func AddAfterTask(task func()) {
	afterTask = append(afterTask, task)
}

func StartAt(start time.Time) {
	beforeTask = append(beforeTask, Task{
		task: func() {
			time.Sleep(time.Until(start))
		},
		isGo: false,
	})
}

func WaitFor(internal time.Duration) {
	beforeTask = append(beforeTask, Task{
		task: func() {
			time.Sleep(internal)
		},
		isGo: false,
	})
}
