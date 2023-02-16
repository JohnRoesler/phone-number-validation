package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/nyaruka/phonenumbers"
)

type phoneNumberResponse struct {
	PhoneNumber      string            `json:"phoneNumber,omitempty"`
	CountryCode      string            `json:"countryCode,omitempty"`
	AreaCode         string            `json:"areaCode,omitempty"`
	LocalPhoneNumber string            `json:"localPhoneNumber,omitempty"`
	Error            map[string]string `json:"error,omitempty"`
}

func main() {
	log.Println("Starting server...")

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/phone-numbers", phoneNumberHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
	log.Println("Server stopped")
}

func phoneNumberHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	incomingPhone := r.URL.Query().Get("phoneNumber")
	countryCode := r.URL.Query().Get("countryCode")

	resp := validatePhoneNumber(incomingPhone, countryCode)

	b, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	if resp.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}
}

func validatePhoneNumber(incomingPhone, countryCode string) phoneNumberResponse {
	resp := phoneNumberResponse{
		PhoneNumber: incomingPhone,
	}
	if incomingPhone == "" {
		resp.Error = map[string]string{
			"phoneNumber": "required value is missing",
		}
		return resp
	}

	parsedPhoneNumber, err := phonenumbers.Parse(incomingPhone, "")
	if err != nil {
		switch {
		case errors.Is(err, phonenumbers.ErrInvalidCountryCode):
			if countryCode == "" {
				resp.Error = map[string]string{
					"countryCode": "required value is missing",
				}
				return resp
			}

			valid := phonenumbers.GetCountryCodeForRegion(countryCode)
			if valid == 0 {
				resp.Error = map[string]string{
					"countryCode": "invalid value",
				}
				return resp
			}

			parsedPhoneNumber, err = phonenumbers.Parse(incomingPhone, countryCode)
			if err != nil {
				resp.Error = map[string]string{
					"phoneNumber": "invalid value",
				}
				return resp
			}

		default:
			resp.Error = map[string]string{
				"phoneNumber": "invalid value",
			}
			return resp
		}
	}

	resp.PhoneNumber = phonenumbers.Format(parsedPhoneNumber, phonenumbers.E164)
	resp.CountryCode = phonenumbers.GetRegionCodeForNumber(parsedPhoneNumber)

	areaCodeLength := phonenumbers.GetLengthOfGeographicalAreaCode(parsedPhoneNumber)
	nationalSignificantNumber := phonenumbers.GetNationalSignificantNumber(parsedPhoneNumber)

	resp.AreaCode = nationalSignificantNumber[:areaCodeLength]
	resp.LocalPhoneNumber = nationalSignificantNumber[areaCodeLength:]

	return resp
}
