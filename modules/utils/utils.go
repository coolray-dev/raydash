package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"

	"github.com/coolray-dev/raydash/modules/log"
	"github.com/sirupsen/logrus"
)

// Hash return sha256 checksum
func Hash(source string) string {
	hashed := sha256.Sum256([]byte(source))
	return base64.StdEncoding.EncodeToString(hashed[:])
}

// AbsPath add project root path before a relative path
func AbsPath(rel string) string {
	_, projectRoot, _, _ := runtime.Caller(0)                           // get dir of current file setting.go
	projectRoot = filepath.Dir(filepath.Dir(filepath.Dir(projectRoot))) // get project root path
	if !filepath.IsAbs(rel) {
		return filepath.Join(projectRoot, rel)
	}
	return rel
}

// VerifyEmailFormat check given string is a email or not
func VerifyEmailFormat(email string) error {
	pattern := `^[A-Za-z0-9]+([_\.][A-Za-z0-9]+)*@([A-Za-z0-9\-]+\.)+[A-Za-z]{2,6}$`
	reg := regexp.MustCompile(pattern)
	if !reg.MatchString(email) {
		return fmt.Errorf("Invalid email")
	}
	return nil
}

// RemoveLoggerHook remove a hook of given name from logger
func RemoveLoggerHook(l *logrus.Logger, hook interface{}) error {
	newHooks := make(logrus.LevelHooks)
	var flag bool = false
	for k := range l.Hooks {
		for _, v2 := range l.Hooks[k] {
			varType := reflect.TypeOf(v2)
			if varType != reflect.TypeOf(hook) {
				newHooks[k] = append(newHooks[k], v2)
			} else {
				flag = true
			}
		}
	}
	if flag {
		log.Log.Debugf("Removed Logger Hook %s", reflect.TypeOf(hook))
	}

	l.ReplaceHooks(newHooks)
	return nil
}

//  UInt64SliceDeDuplicate remove duplicated item in a int slice
func UInt64SliceDeDuplicate(list []uint64) []uint64 {
	set := make(map[uint64]bool, len(list))
	for _, x := range list {
		set[x] = true
	}
	result := make([]uint64, len(set))
	i := 0
	for x := range set {
		result[i] = x
		i++
	}
	return result
}
