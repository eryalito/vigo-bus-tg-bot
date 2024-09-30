package follow

import (
	"sync"
	"time"
)

var (
	FollowTaskManagerInstance = &FollowTaskManager{}
)

// FollowTaskManager manages a list of FollowTasks
type FollowTaskManager struct {
	tasks []FollowTask
	mu    sync.Mutex
	wg    sync.WaitGroup
}

// AddTask adds a new FollowTask to the manager
func (manager *FollowTaskManager) AddTask(task FollowTask, interval time.Duration) {
	manager.mu.Lock()
	manager.tasks = append(manager.tasks, task)
	manager.mu.Unlock()

	manager.wg.Add(1)
	go task.Run(&manager.wg, interval)
}

// RemoveTask removes a FollowTask from the manager
func (manager *FollowTaskManager) RemoveTask(user int64, stop, line int) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	for i, task := range manager.tasks {
		if task.ChatID == user && task.StopNumber == stop && task.LineID == line {
			manager.tasks = append(manager.tasks[:i], manager.tasks[i+1:]...)
			break
		}
	}
}

// Wait waits for all tasks to complete
func (manager *FollowTaskManager) Wait() {
	manager.wg.Wait()
}
