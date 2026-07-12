package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	doh "github.com/FatmanUK/fatgo/docopt_helpers"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"path/filepath"
	pt "ptgen/internal/protracker"
	"regexp"
	"strings"
	"time"
)

type ScanData struct {
	pattern   pt.Pattern
	scanner   *bufio.Scanner
	regex     *regexp.Regexp
	isHexRows bool
	rowsRead  uint
}

func panicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}

func scanLoop(logs chan string, d *ScanData) (pt.Pattern, error) {
	for d.scanner.Scan() {
		m := d.regex.FindStringSubmatch(d.scanner.Text())
		if len(m) == 0 {
			continue
		}
		d.rowsRead++
		n, err := rowNumFromRowStr(m[1], d.isHexRows)
		if err != nil {
			return d.pattern, err
		}
		m2 := [4]string{m[2], m[3], m[4], m[5]}
		d.pattern, err = d.pattern.EmplaceRow(n, m2)
		if err != nil {
			return d.pattern, err
		}
		logs <- fmt.Sprintf("%s", d.pattern[n].Save(n))
		if err != nil {
			return d.pattern, err
		}
	}
	return d.pattern, nil
}

func readPattern(file io.Reader, logs chan string,
	isHexRows bool) (pt.Pattern, error) {
	var err error
	data := ScanData{
		pattern:   pt.PatternFactory(),
		scanner:   bufio.NewScanner(file),
		regex:     regexp.MustCompile(pt.RowRegexFactory()),
		isHexRows: isHexRows,
		rowsRead:  0,
	}
	data.pattern, err = scanLoop(logs, &data)
	if err != nil {
		return data.pattern, err
	}
	logs <- fmt.Sprintf(MSG_TOTAL_ROWS_PROCESSED, data.rowsRead)
	return data.pattern, data.scanner.Err()
}

func readMetadata(proj *pt.ModProject, logs chan string,
	file *os.File) (*pt.ModProject, error) {
	var err error
	content, err := io.ReadAll(file)
	if err != nil {
		return proj, err
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return proj, err
	}
	if json.Valid(content) {
		logs <- MSG_JSON_DETECTED
		err = json.NewDecoder(file).Decode(proj)
	} else {
		var node yaml.Node
		err = yaml.Unmarshal(content, &node)
		if err == nil {
			logs <- MSG_YAML_DETECTED
			err = yaml.NewDecoder(file).Decode(proj)
		}
	}
	return proj, err
}

func defaultOrderlist(proj *pt.ModProject) []uint8 {
	var orderList []uint8
	p := len(proj.Patterns)
	for n := 0; n < p; n++ {
		orderList = append(orderList, uint8(n))
	}
	return orderList
}

func calculateLength(proj *pt.ModProject, logs chan string) {
	for k := range proj.Instruments {
		(&proj.Instruments[k]).CalculateLength()
		//logs <- proj.Instruments[k].Save()
	}
}

// The song metadata in JSON or YAML format, including order list.
func populateMetadata(proj *pt.ModProject, logs chan string,
	fileName string) error {
	var err error
	proj.OrderList = defaultOrderlist(proj)
	if fileName != "<nil>" {
		logs <- MSG_LOADING_METADATA
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		proj, err = readMetadata(proj, logs, file)
		if err != nil {
			return err
		}
		calculateLength(proj, logs)
	}
	if !proj.IsOrderListValid() {
		err = errors.New(ERR_MOD_OLST)
	}
	return err
}

func countPatterns(path string) (uint8, error) {
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

func locateFile(expected string, entries []os.DirEntry) bool {
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if entry.Name() == expected {
			return true
		}
	}
	return false
}

// check all required pattern files exist
func isSequentialPatterns(path string, count uint8) error {
	var err error
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for c := 0; uint8(c) < count; c++ {
		expected := fmt.Sprintf(MSG_PTTN_FILE, c)
		if !locateFile(expected, entries) {
			return errors.New(MSG_PTTN_NOT_SEQ)
		}
	}
	return nil
}

func prepareRegexes() (*regexp.Regexp, *regexp.Regexp,
	*regexp.Regexp) {
	reRow := regexp.MustCompile(pt.RowRegexFactory())
	reHex := regexp.MustCompile(RGX_HEX_EVIDENCE)
	rePttn := regexp.MustCompile(RGX_PTTN_FILE)
	return reRow, reHex, rePttn
}

// Row notation must be consistent across pattern files, but I don't
// see a way to enforce it.
func isHexRowNotationDetected(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	reRow, reHex, rePttn := prepareRegexes()
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !rePttn.MatchString(name) {
			continue
		}
		filePath := filepath.Join(path, name)
		detected, err := checkFileForHexRow(filePath,
			reRow, reHex)
		if err != nil {
			return false, err
		}
		if detected {
			return true, nil
		}
	}
	return false, nil
}

// TODO: row number '5x' or '6x' is an instant false.
func checkFileForHexRow(filePath string, reRow *regexp.Regexp,
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

func loadPattern(logs chan string, fileName string,
	isHexRows bool) (pt.Pattern, error) {
	var err error
	var pattern pt.Pattern
	file, err := os.Open(fileName)
	if err != nil {
		return pt.PatternFactory(), err
	}
	defer file.Close()
	var p any
	p, err = readPattern(file, logs, isHexRows)
	pattern = p.(pt.Pattern)
	return pattern, err
}

func validatePatterns(path string) (uint8, error) {
	numPttns, err := countPatterns(path)
	if err != nil {
		return numPttns, err
	}
	err = pt.IsTooManyPatterns(numPttns)
	if err != nil {
		return numPttns, err
	}
	err = isSequentialPatterns(path, numPttns)
	if err != nil {
		return numPttns, err
	}
	return numPttns, nil
}

// No monolithic pattern files. Too annoying.
// Requires the patterns formatted in plain text. Screw JSON.
func populatePatterns(proj *pt.ModProject, logs chan string,
	path string) error {
	logs <- MSG_LOADING_PTTNS
	numPatterns, err := validatePatterns(path)
	if err != nil {
		return err
	}
	isHex, err := isHexRowNotationDetected(path)
	if err != nil {
		return err
	}
	patterns := make([]pt.Pattern, numPatterns)
	for i := range patterns {
		fileName := fmt.Sprintf(MSG_PTTN_FILE, i)
		fullPath := filepath.Join(path, fileName)
		logs <- fmt.Sprintf(MSG_LOADING_FILE, fileName)
		patterns[i], err = loadPattern(logs, fullPath, isHex)
		if err != nil {
			return err
		}
	}
	proj.Patterns = patterns
	return err
}

// Find the user cache dir. Create source mappings
// ("ref"=>("https://url","local/file")). Decide on mappings depending
// which sample refs are selected.
func decideSources(samples []pt.Instrument) (map[string]Download, error) {
	st01 := []byte(SHA256_ST01)
	st02 := []byte(SHA256_ST02)
	sources := []FileSource{
		{"st01", URL_ST01, SHA256SumFactory(st01)},
		{"st02", URL_ST02, SHA256SumFactory(st02)},
	}
	downloads := map[string]Download{}
	chosen := map[string]Download{}
	cache, err := os.UserCacheDir()
	if err != nil {
		return chosen, err
	}
	for _, m := range sources {
		downloads[m.Ref] = DownloadFactory(m.URL,
			m.ExpectedSum, cache)
	}
	for _, s := range samples {
		chosen[s.Source] = downloads[s.Source]
	}
	return chosen, nil
}

// Has delay to prevent abuse/damage to Aminet. Also checks checksums.
func download(downloads map[string]Download, logs chan string) error {
	for _, v := range downloads {
		exists, err := checkFileExistsWithMkdir(v.Local)
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		logs <- fmt.Sprintf(MSG_DOWNLOAD, v.Remote, v.Local)
		err = v.BlobFromURI(logs)
		if err != nil {
			return err
		}
		logs <- fmt.Sprintf(MSG_SAFETY_PAUSE, SAFETY_PAUSE_S)
		time.Sleep(SAFETY_PAUSE_S * time.Second)
	}
	return nil
}

func outputEverything(proj *pt.ModProject, logs chan string) error {
	var buf bytes.Buffer
	var err error
	info := proj.ModInfoFactory()
	output := mustPrepTemplate("output", metadataString, info)
	logs <- string(output)
	if len(proj.Title) > 20 { // title less than 21 bytes
		return errors.New(ERR_MOD_TITLE_LONG)
	}
	downloads, err := decideSources(proj.Instruments)
	if err != nil {
		return err
	}
	err = download(downloads, logs)
	if err != nil {
		return err
	}
	err = pt.WriteMod(&buf, proj)
	if err != nil {
		return err
	}
	if buf.Len() > MOD_TOO_BIG_BYTES {
		logs <- fmt.Sprintf(MSG_WRN_TOO_BIG_KB,
			buf.Len()/1024, MOD_TOO_BIG_KB)
	}
	err = binary.Write(os.Stdout, binary.BigEndian, buf.Bytes())
	return err
}

func threadGenerate(logs chan string, metaFile any, patternPath any) {
	defer close(logs)
	var err error
	proj := pt.ModProjectFactory()
	pttns := stringFromAny(patternPath)
	err = populatePatterns(&proj, logs, pttns)
	panicIfNotNil(err)
	meta := stringFromAny(metaFile)
	err = populateMetadata(&proj, logs, meta)
	panicIfNotNil(err)
	err = outputEverything(&proj, logs)
	panicIfNotNil(err)
}

func main() {
	dotv := DocOptVarsFactory()
	ds := string(mustPrepTemplate("docopt", docoptString, dotv))
	appVer := dotv.Name + " " + dotv.Version
	args, err := doh.NoExitParser.ParseArgs(ds, nil, appVer)
	if args == nil { // All done, exit here.
		return
	}
	panicIfNotNil(err)
	logs := make(chan string)
	go threadGenerate(logs, args["-m"], args["-p"])
	for msg := range logs {
		msgs := strings.Split(msg, "\\n")
		for _, m := range removeEmptyStrings(msgs) {
			log.Println(m)
		}
	}
}
