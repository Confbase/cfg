package mark

type Config struct {
	Targets    []string
	Template   string
	InstanceOf string
	Singleton  bool
	Force      bool
	UnMark     bool
}
