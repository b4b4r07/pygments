package pygments

import (
	"bytes"
	"os/exec"

	"github.com/oblitum/config"
)

type Pygments struct {
	Executable string
	Formatter  string
	Style      string
}

func (pygments *Pygments) Configure(appliers ...config.Applier) (rollback config.Applier, err error) {
	return config.Configure(pygments, appliers...)
}

var defaultPygments Pygments = Pygments{
	Executable: "pygmentize",
	Formatter:  "html",
	Style:      "default",
}

func New(appliers ...config.Applier) (*Pygments, error) {
	pygments := defaultPygments
	if _, err := pygments.Configure(appliers...); err != nil {
		return nil, err
	}
	return &pygments, nil
}

func Executable(path ...string) config.Applier {
	return func(configurable interface{}) (config.Applier, error) {
		if len(path) == 0 {
			path = []string{"pygmentize"}
		}
		if _, err := exec.LookPath(path[0]); err != nil {
			return nil, err
		}
		pygments := configurable.(*Pygments)
		previous := pygments.Executable
		pygments.Executable = path[0]
		return Executable(previous), nil
	}
}

func Formatter(formatter string) config.Applier {
	return func(configurable interface{}) (config.Applier, error) {
		pygments := configurable.(*Pygments)
		previous := pygments.Formatter
		pygments.Formatter = formatter
		return Formatter(previous), nil
	}
}

func Style(style string) config.Applier {
	return func(configurable interface{}) (config.Applier, error) {
		pygments := configurable.(*Pygments)
		previous := pygments.Style
		pygments.Style = style
		return Style(previous), nil
	}
}

func (pygments *Pygments) Pygmentize(path string) (string, error) {
	cmd := exec.Command(
		pygments.Executable,
		"-f"+pygments.Formatter,
		"-O", "style="+pygments.Style,
		"-g",
		path,
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return stderr.String(), err
	}

	return out.String(), nil
}

func (pygments *Pygments) Stylesheet() (string, error) {
	cmd := exec.Command(
		pygments.Executable,
		"-S"+pygments.Style,
		"-f"+pygments.Formatter,
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return stderr.String(), err
	}

	return out.String(), nil
}
