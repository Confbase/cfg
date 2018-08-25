package init

type Config struct {
	Dest               string
	AppendGitIgnore    bool
	OverwriteGitIgnore bool
	NoGit              bool
	NoModGitIgnore     bool
	Force              bool
}
