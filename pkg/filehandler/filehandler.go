package filehandler

import (
	"fmt"
	"path/filepath"

	"github.com/leonidasdeim/goconfig/internal/files"
)

const (
	marshalIndent   = "	"
	emptySpace      = ""
	filePermissions = 0664

	defaultConfig = "%s.default.%s"
	activeConfig  = "%s.%s"
)

type FileHandler struct {
	file   string
	fileIO fileIO
}

type Optional struct {
	Name string
	Path string
	Type FileType
}

type Option func(f *Optional)

// Add custom filename. By default it is set to "app".
func WithName(n string) Option {
	return func(o *Optional) {
		o.Name = n
	}
}

// Add custom config file path. By default library uses work directory.
func WithPath(p string) Option {
	return func(o *Optional) {
		o.Path = p
	}
}

// Specify handler type.
// - filehandler.JSON (default)
// - filehandler.YAML
func WithType(t FileType) Option {
	return func(o *Optional) {
		o.Type = t
	}
}

func New(opts ...Option) (*FileHandler, error) {

	// Set defaults
	o := &Optional{
		Name: "app",
		Path: files.GetWorkDir(),
		Type: JSON,
	}

	for _, opt := range opts {
		opt(o)
	}

	h := FileHandler{}
	h.fileIO = fileIOFactory(o.Type)
	if h.fileIO == nil {
		return nil, fmt.Errorf("bad file handler type: %s", string(o.Type))
	}

	e := h.fileIO.GetExtension()
	h.file = filepath.Join(o.Path, fmt.Sprintf(activeConfig, o.Name, e))
	defaultFile := filepath.Join(o.Path, fmt.Sprintf(defaultConfig, o.Name, e))

	if err := h.initFrom(defaultFile); err != nil {
		return nil, err
	}

	return &h, nil
}

func (h *FileHandler) Load(data any) error {
	return h.fileIO.Read(data, h.file)
}

func (h *FileHandler) Save(data any) error {
	return h.fileIO.Write(data, h.file)
}

func (h *FileHandler) initFrom(def string) error {
	if files.Exists(h.file) {
		return nil
	}

	if !files.Exists(def) {
		return nil
	}

	var t interface{}

	if err := h.fileIO.Read(&t, def); err != nil {
		return err
	}

	if err := h.fileIO.Write(t, h.file); err != nil {
		return err
	}

	return nil
}