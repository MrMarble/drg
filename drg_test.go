package drg_test

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/mrmarble/drg"
	"github.com/stretchr/testify/require"
)

func readJSON(t *testing.T, filename string) []byte {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	jsonData, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	return jsonData
}

func readSave(t *testing.T, filename string) *os.File {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	return f
}

func TestDecodeMetadata(t *testing.T) {
	f := readSave(t, "testdata/save.sav")
	defer f.Close()

	jsonData := readJSON(t, "testdata/save_metadata.json")

	var expected drg.Metadata
	json.Unmarshal(jsonData, &expected)

	meta, err := drg.DecodeMetadata(f)

	require.NoError(t, err)
	require.NotNil(t, meta)
	require.Equal(t, expected, *meta)
}

func TestDecode(t *testing.T) {
	f := readSave(t, "testdata/save.sav")
	defer f.Close()

	fields, err := drg.Decode(f)

	require.NoError(t, err)
	require.NotNil(t, fields)

	expected := readJSON(t, "testdata/save_fields.json")
	got, err := json.MarshalIndent(fields, "", "  ")

	require.NoError(t, err)
	require.Equal(t, string(expected), string(got))
}

func TestDecodeWithMetadata(t *testing.T) {
	f := readSave(t, "testdata/save.sav")
	defer f.Close()

	drg.DecodeMetadata(f)
	fields, err := drg.Decode(f)

	require.NoError(t, err)
	require.NotNil(t, fields)
}
