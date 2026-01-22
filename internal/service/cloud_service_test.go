package service

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

// MockRoundTripper allows mocking the HTTP response
type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) *http.Response
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req), nil
}

func TestGetMachineID(t *testing.T) {
	svc := NewCloudService()
	id, err := svc.GetMachineID()
	if err != nil {
		t.Fatalf("GetMachineID returned error: %v", err)
	}

	// Verify it is valid Hex
	decoded, err := hex.DecodeString(id)
	if err != nil {
		t.Errorf("GetMachineID returned invalid Hex: %v", err)
	}

	// SHA256 is 32 bytes
	if len(decoded) != 32 {
		t.Errorf("GetMachineID decoded length = %d, expected 32", len(decoded))
	}

	// Hex encoded length for 32 bytes should be 64 chars
	if len(id) != 64 {
		t.Errorf("GetMachineID string length = %d, expected 64", len(id))
	}
}

func TestActivateLicense_RequestEncoding(t *testing.T) {
	// Setup Mock
	capturedBody := []byte{}
	
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) *http.Response {
			// Read body
			capturedBody, _ = io.ReadAll(req.Body)
			req.Body.Close()

			// Return Success Response
			respBody := `{"token": "valid.token.sig", "message": "Success"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(respBody)),
				Header:     make(http.Header),
			}
		},
	}

	svc := NewCloudService()
	svc.Client.Transport = mockTransport

	// Execute
	rawKey := " KSH-AAAA-BBBB-CCCC " // With spaces to test trimming
	_, _ = svc.ActivateLicense(rawKey) 

	// Verify Body
	var reqMap map[string]string
	if err := json.Unmarshal(capturedBody, &reqMap); err != nil {
		t.Fatalf("Failed to unmarshal request body: %v", err)
	}

	// Check LicenseKey (Should be Trimmed Only)
	expectedKey := "KSH-AAAA-BBBB-CCCC"
	if reqMap["license_key"] != expectedKey {
		t.Errorf("Request LicenseKey = %s, expected %s", reqMap["license_key"], expectedKey)
	}

	// Check MachineID (Should be Hex)
	if _, err := hex.DecodeString(reqMap["machine_id"]); err != nil {
		t.Errorf("Request MachineID is not valid Hex: %s", reqMap["machine_id"])
	}
}
