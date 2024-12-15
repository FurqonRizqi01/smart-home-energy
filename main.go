package main

import (
	"log"
	"net/http"
	"os"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"path/filepath"
	"strings"
	"math"
	"a21hc3NpZ25tZW50/service"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// Initialize the services
var fileService = &service.FileService{}
var aiService = &service.AIService{Client: &http.Client{}}
var store = sessions.NewCookieStore([]byte("my-key"))

func getSession(r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "chat-session")
	return session
}

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the Hugging Face token from the environment variables
	token := os.Getenv("HUGGINGFACE_TOKEN")
	if token == "" {
		log.Fatal("HUGGINGFACE_TOKEN is not set in the .env file")
	}

	// Set up the router
	router := mux.NewRouter()

	router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20) // 10 MB maks
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
	
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving the file", http.StatusBadRequest)
			return
		}
		defer file.Close()
	
		// Cek ekstensi file
		filename := header.Filename
		ext := filepath.Ext(filename)
		if ext != ".csv" {
			http.Error(w, "Only CSV files are supported", http.StatusBadRequest)
			return
		}
	
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
	
		processedData, err := fileService.ProcessFile(string(fileBytes))
		if err != nil {
			http.Error(w, "Error processing file: "+err.Error(), http.StatusBadRequest)
			return
		}
	
		analysis := analyzeEnergyConsumption(processedData)
	
		query := r.FormValue("query")
		var aiResponse string
		if query != "" {
			aiContext := fmt.Sprintf("I've uploaded energy consumption data. %s", analysis)
			chatResponse, err := aiService.ChatWithAI(aiContext, query, token)
			if err == nil {
				aiResponse = chatResponse.GeneratedText
			}
		}
	
		response := map[string]string{
			"status":     "success",
			"analysis":   analysis,
			"aiResponse": aiResponse,
		}
	
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("POST")
		
	// Chat endpoint
	router.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		var chatReq struct {
			Context string `json:"context"`
			Query   string `json:"query"`
		}
		err := json.NewDecoder(r.Body).Decode(&chatReq)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
	
		token := os.Getenv("HUGGINGFACE_TOKEN")
	
		chatResponse, err := aiService.ChatWithAI(chatReq.Context, chatReq.Query, token)
		if err != nil {
			http.Error(w, "Error chatting with AI: "+err.Error(), http.StatusInternalServerError)
			return
		}
	
		response := map[string]string{
			"status": "success",
			"answer": chatResponse.GeneratedText,
		}
	
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("POST")

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Allow your React app's origin
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler(router)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}

type deviceStats struct {
    TotalConsumption     float64
    OnTime               int
    OffTime              int
    AverageConsumption   float64
}

type roomStats struct {
    TotalConsumption     float64
    DeviceCount          int
    MostUsedDevice       string
}

func analyzeEnergyConsumption(processedData map[string][]string) string {
    // Validasi kolom
    requiredColumns := []string{"Date", "Time", "Appliance", "Energy_Consumption", "Room", "Status"}
    for _, col := range requiredColumns {
        if _, exists := processedData[col]; !exists {
            return "Error: Missing required columns for analysis"
        }
    }

    analysis := struct {
        Summary struct {
            TotalConsumption     float64
            UniqueDevices        int
            HighestConsumption   string
            LowestConsumption    string
        }
        DeviceBreakdown map[string]deviceStats
        RoomBreakdown   map[string]roomStats
    }{
        DeviceBreakdown: make(map[string]deviceStats),
        RoomBreakdown:   make(map[string]roomStats),
    }

    // Proses data
    for i := 0; i < len(processedData["Appliance"]); i++ {
        appliance := processedData["Appliance"][i]
        room := processedData["Room"][i]
        status := processedData["Status"][i]
        consumption, _ := strconv.ParseFloat(processedData["Energy_Consumption"][i], 64)

        // Update device stats
        device := analysis.DeviceBreakdown[appliance]
        device.TotalConsumption += consumption
        if status == "On" {
            device.OnTime++
        } else {
            device.OffTime++
        }
        analysis.DeviceBreakdown[appliance] = device

        // Update room stats
        roomData := analysis.RoomBreakdown[room]
        roomData.TotalConsumption += consumption
        roomData.DeviceCount++
        analysis.RoomBreakdown[room] = roomData

        // Update total consumption
        analysis.Summary.TotalConsumption += consumption
    }

    // Proses statistik akhir
    var highestConsumptionDevice, lowestConsumptionDevice string
    var highestConsumption, lowestConsumption float64 = 0, math.Inf(1)

    for device, stats := range analysis.DeviceBreakdown {
        // Hitung rata-rata konsumsi
        stats.AverageConsumption = stats.TotalConsumption / float64(stats.OnTime + stats.OffTime)
        analysis.DeviceBreakdown[device] = stats

        if stats.TotalConsumption > highestConsumption {
            highestConsumption = stats.TotalConsumption
            highestConsumptionDevice = device
        }
        if stats.TotalConsumption < lowestConsumption && stats.TotalConsumption > 0 {
            lowestConsumption = stats.TotalConsumption
            lowestConsumptionDevice = device
        }
    }

    var report strings.Builder
    report.WriteString("🔋 Smart Energy Consumption Analysis 🔋\n")
    report.WriteString("=====================================\n\n")

    // Ringkasan Utama
    report.WriteString("📊 Overall Summary:\n")
    report.WriteString(fmt.Sprintf("   Total Energy Consumption: %.2f kWh\n", analysis.Summary.TotalConsumption))
    report.WriteString(fmt.Sprintf("   Highest Consumption Device: %s (%.2f kWh)\n", highestConsumptionDevice, highestConsumption))
    report.WriteString(fmt.Sprintf("   Lowest Consumption Device: %s (%.2f kWh)\n\n", lowestConsumptionDevice, lowestConsumption))

    // Breakdown Perangkat
    report.WriteString("🔌 Device Energy Breakdown:\n")
    for device, stats := range analysis.DeviceBreakdown {
        report.WriteString(fmt.Sprintf("   %s:\n", device))
        report.WriteString(fmt.Sprintf("     - Total Consumption: %.2f kWh\n", stats.TotalConsumption))
        report.WriteString(fmt.Sprintf("     - On Time: %d hours\n", stats.OnTime))
        report.WriteString(fmt.Sprintf("     - Off Time: %d hours\n", stats.OffTime))
        report.WriteString(fmt.Sprintf("     - Average Consumption: %.2f kWh\n\n", stats.AverageConsumption))
    }

    // Breakdown Ruangan
    report.WriteString("🏠 Room Energy Distribution:\n")
    for room, stats := range analysis.RoomBreakdown {
        report.WriteString(fmt.Sprintf("   %s:\n", room))
        report.WriteString(fmt.Sprintf("     - Total Consumption: %.2f kWh\n", stats.TotalConsumption))
        report.WriteString(fmt.Sprintf("     - Devices Used: %d\n\n", stats.DeviceCount))
    }

    // Rekomendasi Efisiensi Energi
    report.WriteString("💡 Energy Efficiency Recommendations:\n")
    if highestConsumptionDevice != "" {
        report.WriteString(fmt.Sprintf("   - Consider optimizing %s usage\n", highestConsumptionDevice))
    }
    report.WriteString("   - Turn off devices when not in use\n")
    report.WriteString("   - Use energy-efficient appliances\n")

    return report.String()
}