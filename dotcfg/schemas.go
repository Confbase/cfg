package dotcfg

type Schemas struct {
	BaseDir string
}

func NewSchemas(baseDir string) *Schemas {
	return &Schemas{BaseDir: baseDir}
}
