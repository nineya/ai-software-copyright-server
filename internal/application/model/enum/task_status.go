package enum

import "errors"

// WAITING：等待中
// INITIATE：已发起
// EXECUTION：执行中
// ABORTED：已终止
// COMPLETE：已完成
var TASK_STATUS = [...]string{"", "WAITING", "INITIATE", "EXECUTION", "ABORTED", "COMPLETE"}

type TaskStatus uint

// JsonDate反序列化
func (t *TaskStatus) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range TASK_STATUS {
		if status == value {
			*t = TaskStatus(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t TaskStatus) MarshalJSON() ([]byte, error) {
	return []byte("\"" + TASK_STATUS[t] + "\""), nil
}

func TaskStatusValue(value string) (TaskStatus, error) {
	for i, status := range TASK_STATUS {
		if status == value {
			return TaskStatus(i), nil
		}
	}
	return TaskStatus(0), errors.New("未找到状态码：" + value)
}
