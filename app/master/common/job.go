package common

const (
	_        int32 = iota
	Pending        // 1 等待执行
	Running        // 2 执行中
	Success        // 3 执行成功
	Kill           // 4 停止当前任务
	Canceled       // 5 任务被强制取消
	// JOB_SAVE_DIR 任务保存目录
	JOB_SAVE_DIR = "cron/jobs/"

	// JOB_LOCK_DIR 任务锁目录
	JOB_LOCK_DIR = "cron/lock/"

	// JOB_KILLER_DIR 停止掉本次任务
	JOB_KILLER_DIR = "cron/killer/"
)

type Job struct {
	Id       int    `json:"id"`       // 任务 ID
	Name     string `json:"name"`     // 任务名称
	Command  string `json:"command"`  // 任务命令
	CronExpr string `json:"cronExpr"` // cron表达式
}
