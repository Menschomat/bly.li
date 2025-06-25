package qrcode

import (
	"bytes"
	"testing"
)

func TestGenerateQRCode_PNG_Success(t *testing.T) {
	png, err := GenerateQRCode("https://example.com", QRPNG)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(png) < 8 {
		t.Fatalf("expected PNG bytes, got only %d bytes", len(png))
	}
	// PNG signature: 89 50 4E 47 0D 0A 1A 0A
	wantHeader := []byte("\x89PNG\r\n\x1a\n")
	if !bytes.HasPrefix(png, wantHeader) {
		t.Errorf("expected PNG header %x, got %x", wantHeader, png[:8])
	}
}

func TestGenerateQRCode_JPEG_Success(t *testing.T) {
	jpg, err := GenerateQRCode("https://example.com", QRJPEG)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(jpg) < 3 {
		t.Fatalf("expected JPEG bytes, got only %d bytes", len(jpg))
	}
	// JPEG signature: FF D8 FF
	wantHeader := []byte("\xff\xd8\xff")
	if !bytes.HasPrefix(jpg, wantHeader) {
		t.Errorf("expected JPEG header %x, got %x", wantHeader, jpg[:3])
	}
}

func TestGenerateQRCode_InvalidData(t *testing.T) {
	_, err := GenerateQRCode("", QRPNG)
	if err != nil {
		t.Fatalf("expected no error with empty string, got %v", err)
	}
}

func TestQRImageType_String(t *testing.T) {
	if QRPNG.String() != "png" {
		t.Errorf("expected QRPNG.String() == \"png\", got %q", QRPNG.String())
	}
	if QRJPEG.String() != "jpeg" {
		t.Errorf("expected QRJPEG.String() == \"jpeg\", got %q", QRJPEG.String())
	}
}
