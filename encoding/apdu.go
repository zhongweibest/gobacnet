package encoding

import (
	bactype "github.com/alexbeltran/gobacnet/types"
)

func (e *Encoder) APDU(a bactype.APDU) {
	meta := APDUMetadata(0)
	meta.setDataType(a.DataType)
	meta.setMoreFollows(a.MoreFollows)
	meta.setSegmentedMessage(a.SegmentedMessage)
	meta.setSegmentedAccepted(a.SegmentedResponseAccepted)
	e.write(meta)

	if a.DataType == bactype.ComplexAck {
		e.apduCompledAck(a)
		return
	}

	e.maxSegsMaxApdu(a.MaxSegs, a.MaxApdu)
	e.write(a.InvokeId)
	if a.SegmentedMessage {
		e.write(a.Sequence)
		e.write(a.WindowNumber)
	}

	e.write(a.Service)
}

func (e *Encoder) apduCompledAck(a bactype.APDU) {
	e.write(a.InvokeId)
	e.write(a.Service)
}

func (d *Decoder) APDU(a *bactype.APDU) error {
	var meta APDUMetadata
	d.decode(&meta)
	a.SegmentedMessage = meta.isSegmentedMessage()
	a.SegmentedResponseAccepted = meta.segmentedResponseAccepted()
	a.MoreFollows = meta.moreFollows()
	a.DataType = meta.DataType()

	if a.DataType == bactype.ComplexAck {
		d.decode(&a.InvokeId)
		d.decode(&a.Service)
		return d.Error()
	}

	a.MaxSegs, a.MaxApdu = d.maxSegsMaxApdu()

	d.decode(&a.InvokeId)
	if a.SegmentedMessage {
		d.decode(&a.Sequence)
		d.decode(&a.WindowNumber)
	}

	d.decode(&a.Service)
	if d.len() > 0 {
		a.Data = make([]byte, d.len())
		d.decode(&a.Data)
	}

	return d.Error()
}

type APDUMetadata byte

const (
	apduMaskSegmented         = 1 << 3
	apduMaskMoreFollows       = 1 << 2
	apduMaskSegmentedAccepted = 1 << 1
	// Bit 0 is reserved
)

func (meta *APDUMetadata) setInfoMask(b bool, mask byte) {
	*meta = APDUMetadata(setInfoMask(byte(*meta), b, mask))
}

// CheckMask uses mask to check bit position
func (meta *APDUMetadata) checkMask(mask byte) bool {
	return (*meta & APDUMetadata(mask)) > 0
}

func (meta *APDUMetadata) isSegmentedMessage() bool {
	return meta.checkMask(apduMaskSegmented)
}

func (meta *APDUMetadata) moreFollows() bool {
	return meta.checkMask(apduMaskMoreFollows)
}

func (meta *APDUMetadata) segmentedResponseAccepted() bool {
	return meta.checkMask(apduMaskSegmentedAccepted)
}

func (meta *APDUMetadata) setSegmentedMessage(b bool) {
	meta.setInfoMask(b, apduMaskSegmented)
}

func (meta *APDUMetadata) setMoreFollows(b bool) {
	meta.setInfoMask(b, apduMaskMoreFollows)
}

func (meta *APDUMetadata) setSegmentedAccepted(b bool) {
	meta.setInfoMask(b, apduMaskSegmentedAccepted)
}

func (meta *APDUMetadata) setDataType(t bactype.PDUType) {
	// clean the first 4 bits
	*meta = (*meta & APDUMetadata(0xF0)) | APDUMetadata(t)
}
func (meta *APDUMetadata) DataType() bactype.PDUType {
	// clean the first 4 bits
	return bactype.PDUType(0xF0) & bactype.PDUType(*meta)
}
