package protracker

// Start, Length and End do not refer to the Sample but rather the
// forward loop if present. Obviously, Start < End,
// Start + Length = End, End < file size.
// Max file size = 65535*2 = 131070, needs a uint32 to store them.
// Length = 0 for no loop.
type Sample struct {
	Id uint8
	Source string
	Name string
	Start uint32
	End uint32
	Length uint32
}

func SampleFactory() Sample {
	return Sample{
		Id: 0,
		Source: "",
		Name: "",
		Start: 0,
		End: 0,
		Length: 0,
	}
}

// TODO: calculate missing loop value (End or Length). Choose Length over End.
