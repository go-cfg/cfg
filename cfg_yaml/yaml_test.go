package cfg_yaml_test

import (
	"embed"
	"os"
	"reflect"
	"testing"

	"gopkg.in/cfg.v0"
	"gopkg.in/cfg.v0/cfg_yaml"
)

//go:embed testdata
var configEmbed embed.FS

func TestYAMLEmbed(t *testing.T) {
	var mcfg struct {
		Foo string
		Bar string
	}
	loader := cfg.LoaderFor(&mcfg, cfg.Config{
		SkipDefaults:       true,
		SkipEnv:            true,
		SkipFlags:          true,
		FailOnFileNotFound: true,
		FileDecoders: map[string]cfg.FileDecoder{
			".yaml": cfg_yaml.New(),
		},
		Files:      []string{"testdata/config.yaml"},
		FileSystem: configEmbed,
	})

	if err := loader.Load(); err != nil {
		t.Fatal(err)
	}

	if mcfg.Foo != "value1" {
		t.Fatalf("have: %v", mcfg.Foo)
	}
	if mcfg.Bar != "value2" {
		t.Fatalf("have: %v", mcfg.Bar)
	}
}

func TestYAML(t *testing.T) {
	filepath := createTestFile(t)

	var mcfg structConfig
	loader := cfg.LoaderFor(&mcfg, cfg.Config{
		SkipDefaults: true,
		SkipEnv:      true,
		SkipFlags:    true,
		FileDecoders: map[string]cfg.FileDecoder{
			".yaml": cfg_yaml.New(),
		},
		Files: []string{filepath},
	})

	if err := loader.Load(); err != nil {
		t.Fatal(err)
	}

	i := int32(42)
	j := int64(420)
	mInterface := make([]interface{}, 2)
	for iI, vI := range []string{"q", "w"} {
		mInterface[iI] = vI
	}
	want := structConfig{
		A: "b",
		C: 10,
		E: 123.456,
		B: []byte("abc"),
		I: &i,
		J: &j,
		Y: structY{
			X: "y",
			Z: []string{"1", "2", "3"},
			A: structD{
				I: true,
			},
		},
		AA: structA{
			X: "y",
			BB: structB{
				CC: structC{
					MM: "n",
					BB: []byte("boo"),
				},
				DD: []string{"x", "y", "z"},
			},
		},
		StructM: StructM{
			M: "n",
		},
		MI: mInterface,
	}

	if got := mcfg; !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func TestLoadResources(t *testing.T) {
	type ResourceA struct {
		Field string `yaml:"field"`
	}
	type ResourceB struct {
		Field int `yaml:"field"`
	}
	type TestConfig struct {
		ResourcesA []ResourceA `yaml:"resources_a"`
		ResourcesB []ResourceB `yaml:"resources_b"`
	}

	var mcfg TestConfig

	resourcesLoader := cfg.LoaderFor(&mcfg,
		cfg.Config{
			SkipFlags:          true,
			Files:              []string{"res.yaml"},
			FailOnFileNotFound: true,
			FileDecoders: map[string]cfg.FileDecoder{
				".yaml": cfg_yaml.New(),
			},
		})
	if err := resourcesLoader.Load(); err != nil {
		t.Errorf("failed to load resources configurations [err=%s]", err)
	}
}

func createTestFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})

	filepath := dir + "/testfile.yaml"

	f, err := os.Create(filepath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, err = f.WriteString(testfileContent)
	if err != nil {
		t.Fatal(err)
	}
	return filepath
}

type structConfig struct {
	A string
	C int
	E float64
	B []byte
	I *int32
	J *int64
	Y structY

	AA structA `yaml:"A"`
	StructM
	MI interface{} `yaml:"MI"`
}

type structY struct {
	X string
	Z []string
	A structD
}

type structA struct {
	X  string  `yaml:"x"`
	BB structB `yaml:"B"`
}

type structB struct {
	CC structC  `yaml:"C"`
	DD []string `yaml:"D"`
}

type structC struct {
	MM string `yaml:"m"`
	BB []byte `yaml:"b"`
}

type structD struct {
	I bool
}

type StructM struct {
	M string
}

const testfileContent = `
a: "b"
c: 10
e: 123.456
b: "abc"
i: 42
j: 420

y:
    x: "y"
    z: ["1", "2", "3"]
    a:
        "i": true

A:
    x: "y"
    B: 
        C:
            m: "n"
            b: "boo"
        D: ["x", "y", "z"]

m: "n"

MI: ["q", "w"]
`
