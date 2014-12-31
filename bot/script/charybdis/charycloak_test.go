package charybdis

import (
	"testing"
)

func TestHostCloak(t *testing.T) {
	str := "cyka.yolo-swag.com"

	output := CloakHost(str)

	if output != "dhou.yolo-swag.com" {
		t.Fatalf("Expected dhou.yolo-swag.com, got %s", output)
	}
}

func TestIPCloak(t *testing.T) {
	str := "8.8.8.8"

	output := CloakIP(str)

	if output != "8.8.l.u" {
		t.Fatalf("Expected 8.8.l.u, got %s", output)
	}
}
