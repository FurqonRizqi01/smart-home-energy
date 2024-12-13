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
        return nil, fmt.Errorf("CSV file is empty")
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
    
    result := make(map[string][]string)
    for _, header := range headers {
        result[header] = []string{}
    }

    for i := 1; i < len(records); i++ {
        if len(records[i]) != len(headers) {
            return nil, fmt.Errorf("invalid CSV data: inconsistent number of columns")
        }

        for j, header := range headers {
            result[header] = append(result[header], records[i][j])
        }
    }

    return result, nil
}
