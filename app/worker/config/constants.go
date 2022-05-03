package config

const (
	// JOB_SAVE_DIR 任务保存目录
	JOB_SAVE_DIR = "cron/jobs/"

	// JOB_LOCK_DIR 任务锁目录
	JOB_LOCK_DIR = "cron/lock/"

	// JOB_KILL_DIR 任务杀死目录
	JOB_KILLER_DIR = "cron/killer/"

	// JOB_EVENT_SAVE 保存任务事件
	JOB_EVENT_SAVE = 1

	// JOB_EVENT_DELETE 删除任务事件
	JOB_EVENT_DELETE = 2

	// JOB_EVENT_KILLER 删除任务事件
	JOB_EVENT_KILLER = 3

	// JOB_EVENT_LOCK 锁任务事件
	JOB_EVENT_LOCK = 4
)
