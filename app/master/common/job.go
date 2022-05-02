package common

const (
	_        int = iota
	Pending      // 等待执行
	Running      // 执行中
	Stopped      // 任务正常完成
	Kill         // 停止当前任务
	Canceled     // 任务被强制取消
)

type Job struct {
	Name     string `json:"name"`     // 任务名称
	Command  string `json:"command"`  // 任务命令
	CronExpr string `json:"cronExpr"` // cron表达式
}
