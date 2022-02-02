package util

import (
	"github.com/GoncharovMikhail/go-sql/errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var (
	separator = "/"
)

func init() {
	envSeparator, exists := os.LookupEnv("SEPARATOR")
	if exists {
		separator = envSeparator
	}
}

func ValidateDirExistence(dir string) error {
	if _, err := os.Open(dir); err != nil {
		log.Printf("Dir: <%s> doesn't exist", dir)
		return err
	}
	return nil
}

func ListAllDirsFromStatedDir(dir string) ([]string, error) {
	if err := ValidateDirExistence(dir); err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for dir = filepath.Dir(dir); dir != separator; {
		result = append(result, dir)
		dir = filepath.Dir(dir)
	}
	return result, nil
}

type PredicateFileInfoFunc func(os.FileInfo) bool

const (
	Conjunction = iota
	Disjunction
)

func ListAllFilesMatchingPatternsAllOverOsFromSpecifiedDir(dir string,
	predicate PredicateFileInfoFunc,
	/* Conjunction, Disjunction */
	matchMode int,
	restrictingPatterns ...string) ([]string, errors.Errors) {
	rootDirs, err := ListAllDirsFromStatedDir(dir)
	if err != nil {
		return nil,
			errors.NewErrors(
				errors.BuildSimpleErrMsg("err", err),
				err,
				nil,
			)
	}
	var result = make([]string, 0)
	// the root of all dirs
	theRootDir := rootDirs[len(rootDirs)-1]
	err = filepath.Walk(theRootDir, func(currentPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		var conjunctionFlag bool = true
		for _, restrictingPattern := range restrictingPatterns {
			matchString, errRegExpMatch := regexp.MatchString(restrictingPattern, currentPath)
			if errRegExpMatch != nil {
				return errRegExpMatch
			}
			if matchString && predicate(info) {
				if matchMode == Disjunction {
					result = append(result, currentPath)
				}
			} else {
				if matchMode == Conjunction {
					conjunctionFlag = false
				}
			}
		}
		if matchMode == Conjunction && conjunctionFlag {
			result = append(result, currentPath)
		}
		return nil
	})
	if err != nil {
		return nil,
			errors.NewErrors(
				errors.BuildSimpleErrMsg("err", err),
				err,
				nil,
			)
	}
	return result, nil
}

func GetGoModDir() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	rootDirs, err := ListAllDirsFromStatedDir(pwd)
	if err != nil {
		return "", err
	}
	var goModDir string = ""
	for _, rootDir := range rootDirs {
		err = filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
			var match bool
			match, err = regexp.MatchString("go.mod", path)
			if err != nil {
				return err
			}
			if match && !info.IsDir() {
				goModDir = rootDir
			}
			return nil
		})
		if err != nil {
			return "", err
		}
		if goModDir != "" {
			break
		}
	}
	return goModDir, nil
}
