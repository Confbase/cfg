package rm

// Config represents the configuration for the `cfg rm` command
type Config struct {
	ToRemove  []string
	Recursive bool
}
