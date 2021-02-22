package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type (
	Request struct {
		Year       string
		ReportType string
		Periode    string
		KodeEmiten string
	}

	Response struct {
		Search      Search   `json:"Search"`
		ResultCount int64    `json:"ResultCount"`
		Results     []Result `json:"Results"`
	}

	Result struct {
		KodeEmiten   string       `json:"KodeEmiten"`
		FileModified string       `json:"File_Modified"`
		ReportPeriod string       `json:"Report_Period"`
		ReportYear   string       `json:"Report_Year"`
		NamaEmiten   string       `json:"NamaEmiten"`
		Attachments  []Attachment `json:"Attachments"`
	}

	Attachment struct {
		EmitenCode   string `json:"Emiten_Code"`
		FileID       string `json:"File_ID"`
		FileModified string `json:"File_Modified"`
		FileName     string `json:"File_Name"`
		FilePath     string `json:"File_Path"`
		FileSize     int64  `json:"File_Size"`
		FileType     string `json:"File_Type"`
		ReportPeriod string `json:"Report_Period"`
		ReportType   string `json:"Report_Type"`
		ReportYear   string `json:"Report_Year"`
		NamaEmiten   string `json:"NamaEmiten"`
	}

	Search struct {
		ReportType string `json:"ReportType"`
		KodeEmiten string `json:"KodeEmiten"`
		Year       string `json:"Year"`
		Periode    string `json:"Periode"`
		Indexfrom  int64  `json:"indexfrom"`
		Pagesize   int64  `json:"pagesize"`
	}
)

func GetLinks(req *Request) (*Response, error) {
	endpoint := "https://idx.co.id/umbraco/Surface/ListedCompany/GetFinancialReport?indexFrom=0&pageSize=10&year=%v&reportType=%v&periode=%v&kodeEmiten=%v"
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	endpoint = fmt.Sprintf(endpoint, req.Year, req.ReportType, req.Periode, req.KodeEmiten)

	request, err := http.NewRequest("GET", fmt.Sprintf("%s", endpoint), bytes.NewBuffer(nil))
	if err != nil {
		log.Printf("error when creating new request: %v", err)
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil || resp.StatusCode != 200 {
		log.Printf("error when creating new request: %v", err)
		return nil, err
	}

	var result = new(Response)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return result, nil
}
