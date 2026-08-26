package protracker

import (
	"path/filepath"
)

// TODO: replace, look up archiveMap dynamically
var archiveMap = map[string]Archive{
	"st01": {
		Url: "https://aminet.net/mods/inst/st-01.lha",
		Checksum: SHA256SumFactory(
			[]byte("8bd8c62d542de794a5f843b5351de1c6"),
			[]byte("f1bed8f7d7643a24a04d0ad8ce962553"),
		),
		File: "st-01.lha",
	},
	"st02": {
		Url: "https://aminet.net/mods/inst/st-02.lha",
		Checksum: SHA256SumFactory(
			[]byte("3ffbbf30e652a65aa47d08afd976990f"),
			[]byte("67e2f9de5cbce171706d2c813d8dbce9"),
		),
		File: "st-02.lha",
	},
}

type SHA256Sum struct {
	BytesHex []byte
}

func SHA256SumFactory(blob1 []byte, blob2 []byte) SHA256Sum {
	return SHA256Sum{BytesHex: append(blob1, blob2...)}
}

type Archive struct {
	Url       string
	Checksum  SHA256Sum
	File      string
	logs chan string
}

func ArchiveFactory(i Instrument, l chan string, p string) Archive {
	// TODO: replace, look up archiveMap dynamically
	return Archive{
		Url: archiveMap[i.Source].Url,
		Checksum: archiveMap[i.Source].Checksum,
		File: filepath.Join(p, archiveMap[i.Source].File),
		logs: l,
	}
}

func (a Archive) Download() error {
	return nil
}
