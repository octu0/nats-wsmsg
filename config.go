package wsmsg

type Config struct {
	DebugMode      bool
	VerboseMode    bool
	Procs          int
	LogDir         string
	NatsLogStdout  bool
	HttpLogStdout  bool
	BindIP         string
	BindPort       int
	MaxMessageSize int
}
