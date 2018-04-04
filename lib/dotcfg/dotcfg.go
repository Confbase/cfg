package dotcfg

const (
	Dirname     = ".cfg"      // this dir resides in ./
	FileName    = ".cfg.json" // this file resides in ./
	KeyfileName = "key.json"  // this file resides in ./.cfg/
)

// .cfg.json is tracked by git
type File struct {
	Templates []string          `json:"templates"`
	Instances map[string]string `json:"instances"`
}

// .cfg/ (including .cfg/key.json) is not tracked by git
type Key struct {
	Email      string `json:"email"`
	Key        string `json:"key"`
	EntryPoint string `json:"entryPoint"` // Confbase API base URL
}
