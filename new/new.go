package new

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Confbase/cfg/cmdrunner"
	"github.com/Confbase/cfg/dotcfg"
	"github.com/Confbase/cfg/tag"
)

func MustNew(cfg *Config) {
	if err := New(cfg.TemplName, cfg.Targets, cfg.DoUseDefaults); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func New(templName string, targets []string, doUseDefaults bool) error {
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

	if doUseDefaults {
		args := append([]string{"init", "-s", templSchemaPath}, targets...)
		cmd := exec.Command("schema", args...)
		if err := cmdrunner.PipeTo(cmd, os.Stdout, os.Stderr); err != nil {
			return err
		}
	}

	templFile, err := os.Open(templObj.FilePath)

	if err != nil {
		return err
	}

	for _, target := range targets {
		if !doUseDefaults {
			f, err := os.Create(target)
			if _, err = io.Copy(f, templFile); err != nil {
				return err
			}
			if _, err = templFile.Seek(0, 0); err != nil {
				return err
			}
		}
		if strings.HasPrefix(target, "--") || strings.HasPrefix(target, "-") {
			continue
		}
		if err := tag.Tag(target, templName); err != nil {
			return err
		}
	}
	return nil
}
