package syncx

import (
	"fmt"
	"log"
	"testing"
	"testing/synctest"
	"time"

	"github.com/luvx21/coding-go/coding-common/common_x"
)

func Test_a_00(t *testing.T) {
	scheduler := NewTaskRunner(3, 100)
	scheduler.Start()
	defer scheduler.Stop()

	tasks := []string{"task1", "task2", "task3", "task4", "task5", "task6", "task7"}

	for i, task := range tasks {
		i, task := i, task
		err := scheduler.AddTask(func() {
			fmt.Printf("执行任务 %d: %s\n", i, task)
			time.Sleep(time.Second * time.Duration(common_x.IfThen(i%2 == 0, 2, 3)))
			fmt.Printf("完成任务 %d: %s\n", i, task)
		})

		if err != nil {
			log.Printf("添加任务失败: %v", err)
		}
	}

	time.Sleep(time.Second * 10) // 等待任务完成
	fmt.Println("程序结束")
}

func Test_a_01(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		Test_a_00(t)
		// synctest.Wait()
	})
}
