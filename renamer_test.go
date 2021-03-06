package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func touch(fileName string) error {
	return ioutil.WriteFile(fileName, []byte{}, os.ModePerm)
}

func contains(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

func all(list []bool) bool {
	for _, item := range list {
		if !item {
			return false
		}
	}
	return true
}

func containsAll(haystack []string, needles []string) bool {
	var results []bool
	for _, item := range needles {
		results = append(results, contains(haystack, item))
	}
	return all(results)
}

func TestMoveFiles(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "sftp_downloader_test")
	if err != nil {
		t.Fail()
	}
	defer os.RemoveAll(tempDir)

	dirA := filepath.Join(tempDir, "foo", "bar", "baz")
	dirB := filepath.Join(tempDir, "foo", "bar", "baz")

	os.MkdirAll(dirA, os.ModePerm)
	os.MkdirAll(dirB, os.ModePerm)
	t.Run("changeme", func(t *testing.T) {
		touch(filepath.Join(dirA, "foo.csV"))
		touch(filepath.Join(dirA, "foo.sAv"))
		touch(filepath.Join(dirA, "foo.txt"))
		touch(filepath.Join(dirB, "bar.csv"))
		touch(filepath.Join(dirB, "BAR.SAV"))
		touch(filepath.Join(dirB, "nothing.txt"))
		err = moveFiles(tempDir)
		if err != nil {
			t.Errorf("Expected no error, got '%s'", err.Error())
		}
		fh, err := os.Open(tempDir)
		if err != nil {
			t.Fail()
		}
		defer fh.Close()
		infos, err := fh.Readdir(-1)
		if err != nil {
			t.Fail()
		}
		var names []string
		for _, info := range infos {
			names = append(names, info.Name())
			t.Logf("name is %s", info.Name())
		}
		want := []string{"foo.csV", "foo.sAv", "bar.csv", "BAR.SAV"}
		if !containsAll(names, want) {
			t.Errorf("Not all items from %s were in directory which contains:\n%s", want, names)
		}

		if len(want) != len(infos) {
			t.Errorf("Directory should contain %d items, got %d", len(want), len(infos))
		}

	})
}

func TestRemoveSuffix(t *testing.T) {
	t.Run("basic-test", func(t *testing.T) {
		input := "STI-EX.MERLIN-29012018220001.csv"
		expected := "STI-EX.MERLIN.csv"
		actual := removeSuffix(input)
		if actual != expected {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	})
	t.Run("without-dashes", func(t *testing.T) {
		input := "bar.csv"
		expected := "bar.csv"
		actual := removeSuffix(input)
		if actual != expected {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	})

	t.Run("sav-extension-no-dashes", func(t *testing.T) {
		input := "bar.sav"
		expected := "bar.sav"
		actual := removeSuffix(input)
		if actual != expected {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	})

	t.Run("sav-extension-with-dashes", func(t *testing.T) {
		input := "bar-foo-quux.sav"
		expected := "bar-foo.sav"
		actual := removeSuffix(input)
		if actual != expected {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	})

	t.Run("bug-from-lili", func(t *testing.T) {
		input := "AE.MERLIN-29012018220001.csv"
		expected := "AE.MERLIN.csv"
		actual := removeSuffix(input)
		if actual != expected {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	})

	// TODO FIXME add text with .sav suffix - make sure func changes that too
}
