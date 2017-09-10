package helpers

import "testing"

func TestConvertingSizeInBytesToHumanReadable(t *testing.T) {
	r := SizeInBytesToHumanReadable(0)
	if r != "0B" {
		t.Errorf("error asserting that SizeInBytesToHumanReadable(0) = '0B', got result %s", r)
	}
	r = SizeInBytesToHumanReadable(100)
	if r != "100B" {
		t.Errorf("error asserting that SizeInBytesToHumanReadable(100) = '100B', got result %s", r)
	}
	r = SizeInBytesToHumanReadable(10240)
	if r != "10K" {
		t.Errorf("error asserting that SizeInBytesToHumanReadable(10240) = '10K', got result %s", r)
	}
	r = SizeInBytesToHumanReadable(1024 * 1024)
	if r != "1MB" {
		t.Errorf("error asserting that SizeInBytesToHumanReadable(1024*1024) = '1MB', got result %s", r)
	}

	r = SizeInBytesToHumanReadable(1024 * 1024 * 1024)
	if r != "1GB" {
		t.Errorf("error asserting that SizeInBytesToHumanReadable(1024*1024*1024) = '1GB', got result %s", r)
	}
}
