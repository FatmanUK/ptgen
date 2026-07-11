package protracker

type Instrument struct {
	ID          uint8  `json:"id"` // The target slot (1-31)
	Name        string `json:"name"`
	Volume      uint8  `json:"volume"`       // 0-64
	Finetune    uint8  `json:"finetune"`     // 0-15 (0=0, 1=1... 8=-8... 15=-1)
	RepeatStart uint32 `json:"repeat_start"` // Offset in BYTES
	RepeatLen   uint32 `json:"repeat_len"`   // Length in BYTES (0 means no loop)
	Data        []byte `json:"data"`         // Raw 8-bit SIGNED PCM
}
