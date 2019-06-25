package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	"github.com/twinj/uuid"
)

type (
	// SERVICE describe Service methods
	SERVICE interface {
		CreateStubMapping(Request, *http.Response) error
		CreateBodyFile(string, []byte) error
		CreateMapping(StubMapping) error
		GetResponseBody(*http.Response) ([]byte, error)
		GetRequestBody(*http.Request) ([]byte, error)
		GetUUID() string
		SanitizeTimestamp([]byte) []byte
	}

	// Service implement Service interface
	Service struct{}

	// StubMapping describe a wiremock stub mapping
	StubMapping struct {
		ID         string              `json:"id"`
		Persistent bool                `json:"persistent"`
		Request    StubMappingRequest  `json:"request"`
		Response   StubMappingResponse `json:"response"`
		UUID       string              `json:"uuid"`
	}

	// StubMappingRequest describe a wiremock stub mapping request object
	StubMappingRequest struct {
		Method       string                          `json:"method"`
		URL          string                          `json:"url"`
		BodyPatterns []StubMappingRequestBodyPattern `json:"bodyPatterns,omitempty"`
	}

	// StubMappingRequestBodyPattern describe the bodyPattern of a PUT/POST method request object
	StubMappingRequestBodyPattern struct {
		EqualToJSON         string `json:"equalToJson"`
		IgnoreArrayOrder    bool   `json:"ignoreArrayOrder"`
		IgnoreExtraElements bool   `json:"ignoreExtraElements"`
	}

	// StubMappingResponse describe a wiremock stub mapping response object
	StubMappingResponse struct {
		Header       map[string][]string `json:"headers"`
		BodyFileName string              `json:"bodyFileName"`
		Status       int                 `json:"status"`
	}
)

// CreateStubMapping Create a new stub mapping in wiremock
func (s *Service) CreateStubMapping(req Request, res *http.Response) error {
	resBuffer, _ := service.GetResponseBody(res)

	uuid := s.GetUUID()

	err := s.CreateBodyFile(uuid, resBuffer)

	if err != nil {
		return err
	}

	stubResponse := StubMappingResponse{
		Header:       res.Header,
		BodyFileName: uuid,
		Status:       res.StatusCode,
	}

	stubRequest := StubMappingRequest{
		Method: req.method,
		URL:    req.url.String(),
	}

	if stubRequest.Method != "GET" && stubRequest.Method != "DELETE" {
		json := s.SanitizeTimestamp(req.body)

		stubRequestBodyPattern := StubMappingRequestBodyPattern{
			EqualToJSON:         string(json),
			IgnoreArrayOrder:    true,
			IgnoreExtraElements: true,
		}

		stubRequest.BodyPatterns = []StubMappingRequestBodyPattern{stubRequestBodyPattern}
	}

	stubMapping := StubMapping{
		ID:         uuid,
		Persistent: true,
		Request:    stubRequest,
		Response:   stubResponse,
		UUID:       uuid,
	}

	err = s.CreateMapping(stubMapping)

	if err != nil {
		return err
	}

	return nil
}

// CreateBodyFile create a file for the body in the wiremock __files directory
func (s *Service) CreateBodyFile(uuid string, buffer []byte) error {
	data := ioutil.NopCloser(bytes.NewReader(buffer))

	req, err := http.NewRequest("PUT", config.APIWiremock+"/__admin/files/"+uuid, data)

	if err != nil {
		return errors.New("Failed to create new http request")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return errors.New("Failed to create stub mapping body file")
	}
	defer res.Body.Close()

	return nil
}

// CreateMapping create a stub mapping in the wiremock mappings directory
func (s *Service) CreateMapping(mapping StubMapping) error {
	buffer, _ := json.Marshal(mapping)
	data := ioutil.NopCloser(bytes.NewReader(buffer))

	req, err := http.NewRequest("POST", config.APIWiremock+"/__admin/mappings", data)

	if err != nil {
		return errors.New("Failed to create new http request")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return errors.New("Failed to create stub mapping")
	}
	defer res.Body.Close()

	return nil
}

// GetUUID return a V4 UUID string
func (s *Service) GetUUID() string {
	u := uuid.NewV4()

	return u.String()
}

// GetResponseBody return the byte content of a Response body
func (s *Service) GetResponseBody(res *http.Response) ([]byte, error) {
	defer res.Body.Close()
	buffer, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	body := ioutil.NopCloser(bytes.NewReader(buffer))
	res.Body = body
	res.ContentLength = int64(len(buffer))
	res.Header.Set("Content-Length", strconv.Itoa(len(buffer)))

	return buffer, nil
}

// GetRequestBody return the byte content of a Request object body
func (s *Service) GetRequestBody(req *http.Request) ([]byte, error) {
	defer req.Body.Close()
	buffer, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return nil, err
	}

	body := ioutil.NopCloser(bytes.NewReader(buffer))
	req.Body = body
	req.ContentLength = int64(len(buffer))
	req.Header.Set("Content-Length", strconv.Itoa(len(buffer)))

	return buffer, nil
}

// SanitizeTimestamp remove modifiedOn field from request body buffer
func (s *Service) SanitizeTimestamp(body []byte) []byte {
	var parsedBody map[string]interface{}
	err := json.Unmarshal(body, &parsedBody)

	if err != nil {
		return nil
	}

	// The magic sauce!
	// Recursively loop through all field for maps and maps in slice,
	// and delete `modifiedOn` field
	for key, value := range parsedBody {

		// If it's a map, we call SanitizeTimestamp from itself to remove modifiedOn field
		if reflect.ValueOf(value).Kind() == reflect.Map {
			var parsedMap map[string]interface{}
			mapBody, _ := json.Marshal(value)
			sanitizedMapBody := s.SanitizeTimestamp(mapBody)
			json.Unmarshal(sanitizedMapBody, &parsedMap)

			parsedBody[key] = parsedMap
		}

		// If it's a slice, we loop thru it and recursively calling SanitizeTimestamp
		if reflect.ValueOf(value).Kind() == reflect.Slice {
			var newSlice []interface{}
			for _, sliceItem := range value.([]interface{}) {
				sliceItemType := reflect.ValueOf(value).Kind()

				// Continue calling SanitizeTimestamp recursively if sliceItem is a map
				if sliceItemType != reflect.Map {
					var sliceParsedMap map[string]interface{}
					sliceMapBody, _ := json.Marshal(sliceItem)
					sliceSanitizedMapBody := s.SanitizeTimestamp(sliceMapBody)
					json.Unmarshal(sliceSanitizedMapBody, &sliceParsedMap)

					newSlice = append(newSlice, sliceParsedMap)
				} else {
					newSlice = append(newSlice, sliceItem)
				}
			}
			parsedBody[key] = newSlice
		}

		// The final touch! Remove the `modifiedOn` field
		if key == "modifiedOn" {
			delete(parsedBody, key)
		}
	}

	var newBody, _ = json.Marshal(parsedBody)

	return newBody
}
