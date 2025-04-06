package cfg_dotenv

import (
	"io/fs"

	"github.com/joho/godotenv"
)

// Decoder of DotENV files for cfg.
type Decoder struct {
	fsys fs.FS
}

// New .ENV decoder for cfg.
func New() *Decoder { return &Decoder{} }

// Format of the decoder.
func (d *Decoder) Format() string {
	return "env"
}

// DecodeFile implements cfg.FileDecoder.
func (d *Decoder) DecodeFile(filename string) (map[string]interface{}, error) {
	file, err := d.fsys.Open(filename)
	if err != nil {
		return nil, err
	}

	raw, err := godotenv.Parse(file)
	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{}, len(raw))
	for key, value := range raw {
		res[key] = value
	}
	return res, nil
}

// DecodeFile implements cfg.FileDecoder.
func (d *Decoder) Init(fsys fs.FS) {
	d.fsys = fsys
}
