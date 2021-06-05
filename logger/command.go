package logger

type CommandLog struct {
	Username     string   `json:"username"`
	Timestamp    string   `json:"timestamp"`
	Command      string   `json:"command"`
	Status       string   `json:"status"`
	StatusReason string   `json:"statusReason"`
	ExecutionID  string   `json:"executionId"`
	Args         []string `json:"args"`
}

type Logger interface {
	LogCommand(command *CommandLog) error
}
