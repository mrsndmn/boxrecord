package boxrecordc

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"testing"

	"github.com/mrsndmn/boxrecord/schema"
	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {

	uidField := schema.Field{
		Name:       "UserID",
		Type:       schema.Types.UInt32,
		FieldNo:    0,
		Size:       4,
		PackFunc:   "tnt.PackInt",
		UnpackFunc: "tnt.UnpackInt",
	}

	primary := schema.Index{
		Name:         "User",
		Fields:       []string{"UserID"},
		FieldsStucts: []schema.Field{uidField},
		Type:         "TREE",
		Uniq:         true,
		IndexNo:      0,
	}

	secondary := []schema.Field{
		{
			Name:       "OauthSource",
			Type:       schema.Types.String,
			Size:       0,
			FieldNo:    1,
			PackFunc:   "[]byte",
			UnpackFunc: "string",
		},
		{
			Name:       "OauthUserID",
			Type:       schema.Types.UInt32,
			Size:       0,
			FieldNo:    2,
			PackFunc:   "tnt.PackInt",
			UnpackFunc: "tnt.UnpackInt",
		},
		{
			Name:       "Flags",
			Type:       schema.Types.UInt32,
			FieldNo:    3,
			Size:       4,
			PackFunc:   "tnt.PackInt",
			UnpackFunc: "tnt.UnpackInt",
		},
	}

	idx2 := schema.Index{
		Name:         "OauthUserID",
		Fields:       []string{"OauthUserID"},
		FieldsStucts: []schema.Field{secondary[1]},
		Type:         "TREE",
		Uniq:         true,
		IndexNo:      1,
	}


	fields := []schema.Field{uidField}
	fields = append(fields, secondary...)

	bs := schema.BoxSchema{
		BoxName: "User",
		Package: "main",
		Space:   0,
		Fields:  fields,
		Indexes: []schema.Index{
			primary,
			idx2,
		},
		PrimaryIndex:    primary,
		SecondaryFields: secondary,
	}

	recordTmpl := template.Must(template.ParseFiles("template/box.tmpl"))
	recordTestTmpl := template.Must(template.ParseFiles("template/box_test.tmpl"))

	dir, err := ioutil.TempDir("", "boxrecord-test")
	if err != nil {
		log.Fatal(err)
	}

	// defer os.RemoveAll(dir) // clean up

	boxrecordfile := filepath.Join(dir, "user.go")
	boxrecordfileTest := filepath.Join(dir, "user_test.go")
	boxrecordfileRunTest := filepath.Join(dir, "runtests.sh")

	br, err := os.Create(boxrecordfile)
	require.NoError(t, err)

	brt, err := os.Create(boxrecordfileTest)
	require.NoError(t, err)

	err = recordTmpl.Execute(br, bs)
	require.NoError(t, err)

	err = recordTestTmpl.Execute(brt, bs)
	require.NoError(t, err)

	message := []byte(fmt.Sprintf(`
#!/usr/bin/bash -ex
cd %s
export GOROOT=""
go mod init test
go test -v .
`, dir))
	err = ioutil.WriteFile(boxrecordfileRunTest, message, 0744)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("bash", boxrecordfileRunTest)
	fmt.Println("\ndir:\n", dir)

	out, err := cmd.CombinedOutput()
	fmt.Printf("out\n\n====\n%s====\n\n", out)
	require.NoError(t, err, dir)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
