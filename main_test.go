package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_validatePhoneNumber(t *testing.T) {
	type args struct {
		incomingPhone string
		countryCode   string
	}
	tests := []struct {
		name string
		args args
		want phoneNumberResponse
	}{
		{
			"valid US phone number",
			args{"+1 650     253   0000", ""},
			phoneNumberResponse{
				PhoneNumber:      "+16502530000",
				CountryCode:      "US",
				AreaCode:         "650",
				LocalPhoneNumber: "2530000",
			},
		},
		{
			"valid MX phone number",
			args{"+525558910066", ""},
			phoneNumberResponse{
				PhoneNumber:      "+525558910066",
				CountryCode:      "MX",
				AreaCode:         "55",
				LocalPhoneNumber: "58910066",
			},
		},
		{
			"missing country code",
			args{"650-253-0000", ""},
			phoneNumberResponse{
				PhoneNumber: "650-253-0000",
				Error: map[string]string{
					"countryCode": "required value is missing",
				},
			},
		},
		{
			"invalid country code",
			args{"650-253-0000", "ESP"},
			phoneNumberResponse{
				PhoneNumber: "650-253-0000",
				Error: map[string]string{
					"countryCode": "invalid value",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validatePhoneNumber(tt.args.incomingPhone, tt.args.countryCode); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validatePhoneNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_phoneNumberHandler_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/phone-numbers?phoneNumber=%2B12125690123", nil)
	w := httptest.NewRecorder()
	phoneNumberHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	expected := `{"phoneNumber":"+12125690123","countryCode":"US","areaCode":"212","localPhoneNumber":"5690123"}`

	if string(data) != expected {
		t.Errorf("expected %v got %v", expected, string(data))
	}
}

func Test_phoneNumberHandler_Error(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/phone-numbers?phoneNumber=631%20311%208150", nil)
	w := httptest.NewRecorder()
	phoneNumberHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	expected := `{"phoneNumber":"631 311 8150","error":{"countryCode":"required value is missing"}}`

	if string(data) != expected {
		t.Errorf("expected %v got %v", expected, string(data))
	}
}
