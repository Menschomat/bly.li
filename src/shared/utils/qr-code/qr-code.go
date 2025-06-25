package qrcode

import (
	"bytes"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// bufWriteCloser wraps a bytes.Buffer to implement io.WriteCloser
type bufWriteCloser struct {
	*bytes.Buffer
}

func (bwc *bufWriteCloser) Close() error { return nil }

// QRImageType is a type-safe enum for supported output formats
type QRImageType int

const (
	QRPNG QRImageType = iota
	QRJPEG
)

// String returns the string name (for debugging/logging, optional)
func (t QRImageType) String() string {
	switch t {
	case QRPNG:
		return "png"
	case QRJPEG:
		return "jpeg"
	default:
		return "unknown"
	}
}

// GenerateQRCode generates a QR code as []byte in PNG or JPEG format
func GenerateQRCode(data string, imageType QRImageType) ([]byte, error) {
	qrc, err := qrcode.New(data)
	if err != nil {
		return nil, err
	}
	buf := &bufWriteCloser{bytes.NewBuffer(nil)}

	imgType := standard.PNG_FORMAT
	if imageType == QRJPEG {
		imgType = standard.JPEG_FORMAT
	}

	w := standard.NewWithWriter(
		buf,
		standard.WithBuiltinImageEncoder(imgType),
	)
	if err := qrc.Save(w); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
