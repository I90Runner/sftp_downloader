package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// FileExists tells you if fileName exists
func FileExists(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("Other error Stat()ing local download directory: %s",
			err.Error())
	}
	return true, nil
}

// AllFilesExist tells you if all files in a list exist
func AllFilesExist(fileNames ...string) (bool, error) {
	for _, fileName := range fileNames {
		ok, err := FileExists(fileName)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

// IsDir tells you if fileName is a directory
func IsDir(fileName string) (bool, error) {
	stat, err := os.Stat(fileName)
	if err != nil {
		return false, fmt.Errorf("Could not stat '%s'", fileName)
	}
	return stat.IsDir(), nil
}

func renameDownloadDir(config Config, fileDate string, phase Phase) (string, error) {
	downloadFolder := getDownloadFolder(phase, config)
	t, err := time.Parse("2006-01-02", fileDate)
	if err != nil {
		return "", fmt.Errorf("Could not convert %s to Time object: %s", fileDate, err.Error())
	}
	newName := t.Format("2006-01-02")
	err = os.Rename(filepath.Join(downloadFolder, fileDate),
		filepath.Join(downloadFolder, newName))
	if err != nil {
		return "", fmt.Errorf("Could not rename download directory: %s", err.Error())
	}
	return newName, nil

}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func convertAccentedToPlain(accented string) (string, error) {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, err := transform.String(t, accented)
	if err != nil {
		return "", err
	}
	return result, nil
}

// change date from dd/mm/yyyy to mm/dd/yyyy
func changeDate(input string) string {
	// fmt.Println(input)
	seps := []string{"/", "-", "."}
	for _, sep := range seps {
		re := regexp.MustCompile(fmt.Sprintf("([0-9]{2})%s([0-9]{2})%s([.]*)", sep, sep))
		if re.MatchString(input) {
			// return re.ReplaceAllString(input, fmt.Sprintf("${2}%s${1}%s${3}", sep, sep))
			out := re.ReplaceAllString(input, "${2}/${1}/${3}")
			re2 := regexp.MustCompile("([0-9]{2})/([0-9]{2})/([0-9]{2})$")
			if re2.MatchString(out) {
				return re2.ReplaceAllString(out, "${1}/${2}/20${3}")
			}
			return out
		}
	}
	return input
}
