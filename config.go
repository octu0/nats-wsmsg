package wsmsg

type Config struct {
  DebugMode      bool
  VerboseMode    bool
  Procs          int
  LogDir         string
  BindIP         string
  BindPort       int
}
