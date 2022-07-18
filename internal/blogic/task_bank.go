package blogic

import (
	"OnlineSchool/internal/structs"
	"context"
	"fmt"
)

//возврат словарь задач где ключ = task_id, успешность операции true/false
func (b *BLogic) getTasks(tasksId []int, handed bool) (map[int]structs.Task, bool) {
	masTask, err := b.DBTaskBank.GetTasks(context.TODO(), tasksId)
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	if len(masTask) != len(tasksId) {
		fmt.Println("len(masTask) != len(tasksId) ")
		return nil, false
	}
	mapTasks := make(map[int]structs.Task)
	for _, val := range masTask {
		mapTasks[val.TaskId] = val
	} /*

		}*/
	return mapTasks, true
}
