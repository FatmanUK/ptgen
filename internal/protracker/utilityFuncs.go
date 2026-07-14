package protracker

import (
	"bufio"
	"bytes"
	"fmt"
	xlha "github.com/FatmanUK/fatgo/xlha"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"text/template"
)

func TemplateFactory(name string, format string) *template.Template {
	t := template.New(name)
	t = template.Must(t.Parse(format))
	return t
}

func MustPrepTemplate(name string, formatString string,
	data any) []byte {
	var wr bytes.Buffer
	err := TemplateFactory(name, formatString).Execute(&wr, data)
	if err != nil {
		panic(err)
	}
	return wr.Bytes()
}

// Row notation must be consistent across pattern files, but I don't
// see a way to enforce it.
func IsHexDetected(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	rgxs := PrepareRegexes()
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !rgxs[2].MatchString(name) {
			continue
		}
		filePath := filepath.Join(path, name)
		detected, err := CheckHex(filePath, rgxs[0], rgxs[1])
		if err != nil {
			return false, err
		}
		if detected {
			return true, nil
		}
	}
	return false, nil
}

// TODO: row number '5x' or '6x' is an instant return false, nil.
func CheckHex(filePath string, reRow *regexp.Regexp,
	reHex *regexp.Regexp) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		mRow := reRow.FindStringSubmatch(scanner.Text())
		if len(mRow) > 1 && reHex.MatchString(mRow[1]) {
			return true, nil
		}
	}
	return false, scanner.Err()
}

func CountPatterns(path string) (uint8, error) {
	var count uint8 = 0
	var err error
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}
	re := regexp.MustCompile(RGX_PTTN_FILE)
	for _, entry := range entries {
		if !entry.IsDir() && re.MatchString(entry.Name()) {
			count++
		}
	}
	return count, nil
}

func PrepareRegexes() []*regexp.Regexp {
	reRow := regexp.MustCompile(RowRegexFactory())
	reHex := regexp.MustCompile(RGX_HEX_EVIDENCE)
	rePttn := regexp.MustCompile(RGX_PTTN_FILE)
	return []*regexp.Regexp{reRow, reHex, rePttn}
}

func uint8FromHexString(s string) (uint8, error) {
	u, err := strconv.ParseUint(s, 16, 8)
	return uint8(u), err
}

// TODO: duplicate. Figure out how to share between packages (symlink won't do it)
// well, another fatgo pkg would do it. Seems a bit much for one function.
func compareSlices[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !(a[i] == b[i]) {
			return false
		}
	}
	return true
}

func CheckFileExistsWithMkdir(path string) (bool, error) {
	cacheDir := filepath.Dir(path)
	err := os.MkdirAll(cacheDir, 0750)
	if err != nil {
		return false, err
	}
	exists, err := fileExists(path)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func fileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		// doesn't exist
		return false, nil
	}
	// permissions issue
	return false, err
}

func prepareArchives() (map[string]string, error) {
	archives := map[string]string{}
	cache, err := os.UserCacheDir()
	if err != nil {
		return archives, err
	}
	archives["st01"] = filepath.Join(cache, "ptgen", FILE_ST01)
	archives["st02"] = filepath.Join(cache, "ptgen", FILE_ST02)
	return archives, nil
}

func extractSample(n string, lr *xlha.Reader,
	data *[]byte) (bool, error) {
	h, err := lr.Next()
	if err == io.EOF {
		return true, nil
	}
	if err != nil {
		return true, fmt.Errorf(ERR_ARCH_HEADER_PARSE, err)
	}
	log.Println(fmt.Sprintf(MSG_ARCH_HEADER_PARSE_OK, h.Name,
		h.Method, h.OriginalSize))
	*data, err = io.ReadAll(lr)
	if err != nil {
		return true, fmt.Errorf(ERR_ARCH_EXTRACTION,
			h.Name, err)
	}
	written := uint32(len(*data))
	if written != h.OriginalSize {
		return true, fmt.Errorf(ERR_ARCH_SIZE_MISMATCH,
			h.Name, h.OriginalSize, written)
	}
	if n == h.Name { // found our file
		return true, nil
	}
	return false, nil
}
