package logger

type CommandLog struct {
	Username     string   `json:"username"`
	Timestamp    string   `json:"timestamp"`
	Command      string   `json:"command"`
	Args         []string `json:"args"`
	Status       string   `json:"status"`
	StatusReason string   `json:"statusReason"`
	ExecutionID  string   `json:"executionId"`
}
