package new

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Confbase/cfg/cmdrunner"
	"github.com/Confbase/cfg/dotcfg"
	"github.com/Confbase/cfg/tag"
)

func MustNew(cfg *Config) {
	if err := New(cfg.TemplName, cfg.Targets); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func New(templName string, targets []string) error {
	baseDir, err := dotcfg.GetBaseDir()
	if err != nil {
		return err
	}
	cfgFile, err := dotcfg.LoadCfg(baseDir)
	if err != nil {
		return err
	}

	var templObj *dotcfg.Template
	for _, templ := range cfgFile.Templates {
		if templ.Name == templName {
			templObj = &templ
			break
		}
	}
	if templObj == nil {
		return fmt.Errorf("a template with the name '%v' does not exist", templName)
	}

	var templSchemaPath string
	if templObj.Schema.FilePath != "" {
		templSchemaPath = templObj.Schema.FilePath
	} else {
		templSchemaPath = filepath.Join(dotcfg.SchemasDirName, templObj.FilePath)
	}

	args := append([]string{"init", "-s", templSchemaPath}, targets...)
	cmd := exec.Command("schema", args...)
	if err := cmdrunner.PipeTo(cmd, os.Stdout, os.Stderr); err != nil {
		return err
	}

	for _, target := range targets {
		if strings.HasPrefix(target, "--") || strings.HasPrefix(target, "-") {
			continue
		}
		if err := tag.Tag(target, templName); err != nil {
			return err
		}
	}
	return nil
}
