package main

import (
	"log"
	"net/http"
	"os"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
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
		err := r.ParseMultipartForm(10 << 20)
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

		filename := "uploaded_" + header.Filename
		err = fileService.Repo.SaveFile(filename, fileBytes)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		// Analisis konsumsi energi
		highestConsumption := make(map[string]string)
		lowestConsumption := make(map[string]string)
		maxConsumption := make(map[string]float64)
		minConsumption := make(map[string]float64)

		for _, energySource := range []string{"Electric", "Solar", "Battery"} {
			maxConsumption[energySource] = 0.0
			minConsumption[energySource] = math.Inf(1)
		}

		for i, appliance := range processedData["Appliance"] {
			consumption, _ := strconv.ParseFloat(processedData["Energy_Consumption"][i], 64)
			energySource := processedData["Energy_Source"][i]

			if consumption > maxConsumption[energySource] {
				maxConsumption[energySource] = consumption
				highestConsumption[energySource] = appliance
			}
			if consumption < minConsumption[energySource] {
				minConsumption[energySource] = consumption
				lowestConsumption[energySource] = appliance
			}
		}

		// Tambahkan analisis untuk setiap sumber energi
		analysis := "Energy Consumption Analysis:\n"
		for _, energySource := range []string{"Electric", "Solar", "Battery"} {
			if maxConsumption[energySource] > 0 {
				analysis += fmt.Sprintf(
					"%s: Most consumption by %s (%.2f), Least consumption by %s (%.2f)\n",
					energySource,
					highestConsumption[energySource],
					maxConsumption[energySource],
					lowestConsumption[energySource],
					minConsumption[energySource],
				)
			}
		}

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
