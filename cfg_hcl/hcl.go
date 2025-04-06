package cfg_hcl

import (
	"io/fs"

	"github.com/hashicorp/hcl"
)

// Decoder of HCL files for cfg.
type Decoder struct {
	fsys fs.FS
}

// New HCL decoder for cfg.
func New() *Decoder { return &Decoder{} }

// Format of the decoder.
func (d *Decoder) Format() string {
	return "hcl"
}

// DecodeFile implements cfg.FileDecoder.
func (d *Decoder) DecodeFile(filename string) (map[string]interface{}, error) {
	b, err := fs.ReadFile(d.fsys, filename)
	if err != nil {
		return nil, err
	}

	f, err := hcl.ParseBytes(b)
	if err != nil {
		return nil, err
	}

	var raw map[string]interface{}
	if err := hcl.DecodeObject(&raw, f); err != nil {
		return nil, err
	}
	return raw, nil
}

// DecodeFile implements cfg.FileDecoder.
func (d *Decoder) Init(fsys fs.FS) {
	d.fsys = fsys
}
