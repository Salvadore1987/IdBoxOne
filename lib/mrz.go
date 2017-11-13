package lib

import (
	"strings"
)

type MRZ struct {
	line string
	mrz_info string
	checksum bool
}

func NewMRZ(mrz string) *MRZ {
	mrzIns := new(MRZ)
	lines := strings.Split(mrz, "\r")
	mrzIns.line = lines[1]
	mrzIns.checksum = false
	mrzIns.extractMRZInfo()
	return mrzIns
}

func (t *MRZ) extractMRZInfo()  {
	t.mrz_info = t.line[0:10] + t.line[13:20] + t.line[21:28]
	if  CalculateChecksum(t.line[0:10]) ||
		CalculateChecksum(t.line[13:20]) ||
		CalculateChecksum(t.line[21:28]) {
		t.checksum = true
	} else {
		t.checksum = false
	}
}

func (t *MRZ) GetMRZInfo() string {
	return t.mrz_info
}

func (t* MRZ) GetChecksum() bool {
	return t.checksum;
}
