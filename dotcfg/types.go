package dotcfg

const (
	DirName        = ".cfg"         // this dir resides in ./
	FileName       = ".cfg.json"    // this file resides in ./
	KeyfileName    = "key.json"     // this file resides in ./.cfg/
	SnapsFileName  = "snaps.json"   // this file resides in ./.cfg/
	SchemasDirName = ".cfg_schemas" // this file resides in ./
)

type JSONSchema struct {
	FilePath string `json:"filePath"`
}

type Template struct {
	Name     string     `json:"name"`
	FilePath string     `json:"filePath"`
	Schema   JSONSchema `json:"schema"`
}

// .cfg.json is tracked by git
type File struct {
	Templates  []Template  `json:"templates"`
	Instances  []Instance  `json:"instances"`
	Singletons []Singleton `json:"singletons"`
	NoGit      bool        `json:"noGit"`
}

type Singleton struct {
	FilePath string     `json:"filePath"`
	Schema   JSONSchema `json:"schema"`
}

type Instance struct {
	FilePath   string     `json:"filePath"`
	TemplNames []string   `json:"templateName"`
	Schema     JSONSchema `json:"schema"`
}

// everything in .cfg/ (including .cfg/key.json) is not tracked by git
type Key struct {
	Email    string            `json:"email"`
	Remotes  map[string]string `json:"remotes"`
	BaseName string            `json:"baseName"`
}

type Snapshot struct {
	Name string `json:"name"`
}

// everything in .cfg/ (including .cfg/snaps) is not tracked by git
// however, snaps are pushed to Confbase servers
type Snaps struct {
	Current   Snapshot   `json:"current"`
	Snapshots []Snapshot `json:"snapshots"`
}
