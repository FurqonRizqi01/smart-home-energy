import React, { useState, useCallback } from "react";
import axios from "axios";
import { motion } from "framer-motion";
import Header from "./components/Header";
import Footer from "./components/Footer";

function App() {
  const [file, setFile] = useState(null);
  const [tapasQuery, setTapasQuery] = useState("");
  const [miniChatQuery, setMiniChatQuery] = useState("");
  const [response, setResponse] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");

  // Fungsi untuk menangani perubahan file
  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    const allowedTypes = ['text/csv', 'application/vnd.ms-excel'];
    
    if (selectedFile) {
        if (!allowedTypes.includes(selectedFile.type)) {
            setErrorMessage("Only CSV files are allowed!");
            setFile(null);
            return;
        }
        
        setFile(selectedFile);
        setErrorMessage("");
        setResponse("");
    }
};

  // Fungsi upload
  const handleUpload = useCallback(async () => {
    if (!file) {
      setErrorMessage("Please select a file first.");
      return;
    }

    setIsLoading(true);
    setErrorMessage("");
    setResponse("");

    const formData = new FormData();
    formData.append("file", file);
    
    //query Tapas
    if (tapasQuery) {
      formData.append("query", tapasQuery);
    }

    try {
      const res = await axios.post('http://localhost:8080/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        timeout: 30000
      });
      
      const fullResponse = 
        `📊 File Analysis:\n${res.data.analysis || 'No analysis available'}\n\n` +
        `🤖 AI Insights: ${res.data.aiResponse || res.data.answer || 'No additional insights'}`;
      
      setResponse(fullResponse);
    } catch (error) {
      console.error('Error uploading file:', error);
      
      if (error.response) {
        // Server responded with an error
        setErrorMessage(`Upload Error: ${error.response.data.message || 'Server error'}`);
      } else if (error.request) {
        // Request made but no response received
        setErrorMessage("No response from server. Please check your connection.");
      } else {
        // Something happened in setting up the request
        setErrorMessage("Error setting up the upload request.");
      }
    } finally {
      setIsLoading(false);
    }
  }, [file, tapasQuery]);

  // Fungsi chat
  const handleMiniChat = useCallback(async () => {
    const trimmedQuery = miniChatQuery.trim();
    if (!trimmedQuery) {
      setErrorMessage("Please enter a query.");
      return;
    }

    setIsLoading(true);
    setErrorMessage("");
    
    try {
      const res = await axios.post("http://localhost:8080/chat", { 
        context: "",
        query: trimmedQuery 
      }, {
        timeout: 30000
      });
      
      const fullResponse = `🤖 Mini Chat Response:\n${res.data.answer || 'No response received'}`;
      setResponse(fullResponse);
    } catch (error) {
      console.error("Error querying chat:", error);
      
      // Detailed error handling
      if (error.response) {
        setErrorMessage(`Chat Error: ${error.response.data.message || 'Server error'}`);
      } else if (error.request) {
        setErrorMessage("No response from server. Please check your connection.");
      } else {
        setErrorMessage("Error setting up the chat request.");
      }
    } finally {
      setIsLoading(false);
    }
  }, [miniChatQuery]);

  // Reset fungsi
  const resetForm = () => {
    setFile(null);
    setTapasQuery("");
    setMiniChatQuery("");
    setResponse("");
    setErrorMessage("");
  };

  return (
    <div className="flex flex-col min-h-screen bg-gray-50 dark:bg-gray-900">
      <Header />
      <main className="flex-grow container mx-auto px-4 py-8">
      <motion.div
          initial={{ opacity: 0, y: -50 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          className="text-center mb-8"
        >
          <h1 className="text-4xl md:text-5xl font-bold text-yellow-600 mb-2 animate-pulse">
            Welcome to SmartEnergy Hub!
          </h1>
          <p className="text-xl md:text-2xl text-gray-600 dark:text-gray-300">
            Mengelola Energy, Mengelola Hidup!
          </p>
        </motion.div>
        <div className="max-w-4xl mx-auto">
          <div className="bg-white dark:bg-gray-800 shadow-xl rounded-lg overflow-hidden">
            <div className="p-6 md:p-8">
              <h2 className="text-2xl font-semibold mb-6 text-gray-800 dark:text-white">
                AI Data Analysis Chatbot
              </h2>

              {errorMessage && (
                <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
                  {errorMessage}
                </div>
              )}

              <div className="space-y-6">
                <div className="space-y-2">
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    Upload File
                  </label>
                  <div className="flex items-center space-x-2">
                    <input
                      type="file"
                      onChange={handleFileChange}
                      className="block w-full text-sm text-gray-500 
                        file:mr-4 file:py-2 file:px-4 
                        file:rounded-full file:border-0 
                        file:text-sm file:font-semibold 
                        file:bg-yellow-50 file:text-yellow-700 
                        hover:file:bg-purple-100"
                    />
                    {file && (
                      <button 
                        onClick={resetForm}
                        className="text-red-500 hover:text-red-700"
                      >
                        Clear
                      </button>
                    )}
                  </div>
                </div>

                <div className="space-y-2">
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    Tapas AI Table Query
                  </label>
                  <div className="flex items-center space-x-2">
                    <input
                      type="text"
                      value={tapasQuery}
                      onChange={(e) => setTapasQuery(e.target.value)}
                      placeholder="Which appliance uses most energy?"
                      className="shadow-sm focus:ring-purple-500 
                        focus:border-purple-500 block w-full 
                        sm:text-sm border-gray-300 rounded-md 
                        dark:bg-gray-700 dark:border-gray-600 
                        dark:text-white"
                    />
                    <button
                      onClick={handleUpload}
                      disabled={!file || isLoading}
                      className="inline-flex items-center px-4 py-2 
                        border border-transparent text-sm font-medium 
                        rounded-md shadow-sm text-white bg-yellow-600 
                        hover:bg-yellow-700 focus:outline-none 
                        focus:ring-2 focus:ring-offset-2 
                        focus:ring-yellow-500 disabled:opacity-50"
                    >
                      {isLoading ? 'Uploading...' : 'Upload and Analyze'}
                    </button>
                  </div>
                 </div>

                <div className="space-y-2">
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    Mini Chat AI General Query
                  </label>
                  <div className="flex items-center space-x-2">
                    <input
                      type="text"
                      value={miniChatQuery}
                      onChange={(e) => setMiniChatQuery(e.target.value)}
                      placeholder="Ask a general question"
                      className="shadow-sm focus:ring-purple-500 
                        focus:border-purple-500 block w-full 
                        sm:text-sm border-gray-300 rounded-md 
                        dark:bg-gray-700 dark:border-gray-600 
                        dark:text-white"
                    />
                    <button
                      onClick={handleMiniChat}
                      disabled={!miniChatQuery.trim() || isLoading}
                      className="inline-flex items-center px-4 py-2 
                        border border-transparent text-sm font-medium 
                        rounded-md shadow-sm text-white bg-green-600 
                        hover:bg-green-700 focus:outline-none 
                        focus:ring-2 focus:ring-offset-2 
                        focus:ring-green-500 disabled:opacity-50"
                    >
                      {isLoading ? 'Processing...' : 'Chat with Mini AI'}
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <div className="bg-gray-50 dark:bg-gray-700 px-6 py-4 md:px-8 md:py-6">
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">
                Response
              </h3>
              <div className="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
                <div className="px-4 py-5 sm:p-6">
                  <p className="text-gray-700 dark:text-gray-300 whitespace-pre-wrap">
                    {response || 'Your response will appear here.'}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
      <Footer />
    </div>
  );
}

export default App;