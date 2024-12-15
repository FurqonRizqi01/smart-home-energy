package service

import (
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"encoding/csv"
	"fmt"
	"strings"
)

type FileService struct {
	Repo *repository.FileRepository
}

func (s *FileService) ProcessFile(fileContent string) (map[string][]string, error) {
    if fileContent == "" {
        return nil, fmt.Errorf("file is empty")
    }

    reader := csv.NewReader(strings.NewReader(fileContent))
    records, err := reader.ReadAll()
    if err != nil {
        return nil, fmt.Errorf("error parsing CSV: %v", err)
    }

    if len(records) < 1 {
        return nil, fmt.Errorf("CSV file does not have a header row")
    }

    headers := records[0]
    
    columnMapping := map[string]string{
        "Energy_Consumption (kWh)": "Energy_Consumption",
        "Energy Consumption": "Energy_Consumption",
    }

    result := make(map[string][]string)
    
    // Hanya tambahkan header yang ada di CSV
    for _, header := range headers {
        standardHeader := header
        if mappedHeader, exists := columnMapping[header]; exists {
            standardHeader = mappedHeader
        }
        result[standardHeader] = []string{}
    }

    // Isi data
    for i := 1; i < len(records); i++ {
        if len(records[i]) > len(headers) {
            return nil, fmt.Errorf("invalid CSV data: too many columns")
        }

        for j, header := range headers {
            standardHeader := header
            if mappedHeader, exists := columnMapping[header]; exists {
                standardHeader = mappedHeader
            }

            if j < len(records[i]) {
                result[standardHeader] = append(result[standardHeader], records[i][j])
            } else {
                result[standardHeader] = append(result[standardHeader], "")
            }
        }
    }

    return result, nil
}