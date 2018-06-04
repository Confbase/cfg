package ls

type Config struct {
	NoColors      bool
	NoTty         bool
	DoLsTempls    bool
	DoLsInsts     bool
	DoLsSingles   bool
	DoLsUntracked bool
}
