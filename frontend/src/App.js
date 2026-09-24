import React, { useState, useCallback, useRef } from "react";
import axios from "axios";
import { motion } from "framer-motion";
import ReactMarkdown from "react-markdown";
import remarkBreaks from "remark-breaks";
import remarkGfm from "remark-gfm";
import Header from "./components/Header";
import Footer from "./components/Footer";

const markdownComponents = {
  h1: ({ children }) => <h1 className="mt-6 mb-3 text-2xl font-bold text-gray-900 dark:text-white">{children}</h1>,
  h2: ({ children }) => <h2 className="mt-6 mb-3 text-xl font-bold text-gray-900 dark:text-white">{children}</h2>,
  h3: ({ children }) => <h3 className="mt-5 mb-2 text-lg font-semibold text-gray-900 dark:text-white">{children}</h3>,
  p: ({ children }) => <p className="mb-4 break-words leading-7">{children}</p>,
  strong: ({ children }) => <strong className="font-bold text-gray-900 dark:text-white">{children}</strong>,
  ul: ({ children }) => <ul className="mb-4 list-disc space-y-1 pl-6">{children}</ul>,
  ol: ({ children }) => <ol className="mb-4 list-decimal space-y-1 pl-6">{children}</ol>,
  li: ({ children }) => <li className="break-words pl-1">{children}</li>,
  blockquote: ({ children }) => (
    <blockquote className="my-4 border-l-4 border-yellow-500 bg-gray-100 px-4 py-2 italic dark:bg-gray-700">{children}</blockquote>
  ),
  hr: () => <hr className="my-6 border-gray-300 dark:border-gray-600" />,
  table: ({ children }) => (
    <div className="my-5 max-w-full overflow-x-auto rounded-lg border border-gray-300 dark:border-gray-600">
      <table className="w-full min-w-[640px] border-collapse text-left text-sm">{children}</table>
    </div>
  ),
  thead: ({ children }) => <thead className="bg-gray-200 dark:bg-gray-700">{children}</thead>,
  th: ({ children }) => <th className="border-b border-r border-gray-300 px-4 py-3 font-semibold last:border-r-0 dark:border-gray-600">{children}</th>,
  td: ({ children }) => <td className="border-b border-r border-gray-200 px-4 py-3 align-top last:border-r-0 dark:border-gray-700">{children}</td>,
  pre: ({ children }) => <pre className="my-4 max-w-full overflow-x-auto rounded-lg bg-gray-950 p-4 text-sm text-gray-100">{children}</pre>,
  code: ({ children }) => <code className="rounded bg-gray-200 px-1 py-0.5 text-sm dark:bg-gray-700">{children}</code>,
  a: ({ children, href }) => <a className="text-blue-600 underline hover:text-blue-500 dark:text-blue-400" href={href} target="_blank" rel="noreferrer">{children}</a>,
};

function App() {
  const fileInputRef = useRef(null);
  const [file, setFile] = useState(null);
  const [tapasQuery, setTapasQuery] = useState("");
  const [miniChatQuery, setMiniChatQuery] = useState("");
  const [response, setResponse] = useState("");
  const [aiContext, setAiContext] = useState("");
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
            e.target.value = "";
            return;
        }
        
        setFile(selectedFile);
        setAiContext("");
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
        timeout: 60000
      });
      
      const fullResponse =
        `## 📊 File Analysis\n\n${res.data.analysis || 'No analysis available'}\n\n---\n\n` +
        `## 🤖 AI Insights\n\n${res.data.aiResponse || res.data.answer || 'No additional insights'}`;
      
      setAiContext(res.data.context || "");
      setResponse(fullResponse);
    } catch (error) {
      console.error('Error uploading file:', error);
      
      if (error.response) {
        // Server responded with an error
        const serverMessage = typeof error.response.data === "string"
          ? error.response.data
          : error.response.data?.message;
        setErrorMessage(`Upload Error: ${serverMessage || 'Server error'}`);
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
        context: aiContext,
        query: trimmedQuery 
      }, {
        timeout: 60000
      });
      
      const fullResponse = `## 🤖 Mini Chat Response\n\n${res.data.answer || 'No response received'}`;
      setResponse(fullResponse);
    } catch (error) {
      console.error("Error querying chat:", error);
      
      // Detailed error handling
      if (error.response) {
        const serverMessage = typeof error.response.data === "string"
          ? error.response.data
          : error.response.data?.message;
        setErrorMessage(`Chat Error: ${serverMessage || 'Server error'}`);
      } else if (error.request) {
        setErrorMessage("No response from server. Please check your connection.");
      } else {
        setErrorMessage("Error setting up the chat request.");
      }
    } finally {
      setIsLoading(false);
    }
  }, [aiContext, miniChatQuery]);

  // Reset fungsi
  const resetForm = () => {
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
    setFile(null);
    setTapasQuery("");
    setMiniChatQuery("");
    setResponse("");
    setAiContext("");
    setErrorMessage("");
  };

  return (
    <div className="flex flex-col min-h-screen bg-gray-50 dark:bg-gray-900">
      <Header />
      <main className="flex-grow container mx-auto px-4 py-8">
      <motion.div
          initial={{ opacity: 0, y: -50 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1 }}
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
                  <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
                    <input
                      ref={fileInputRef}
                      type="file"
                      onChange={handleFileChange}
                      className="block min-w-0 flex-1 text-sm text-gray-500 
                        file:mr-4 file:py-2 file:px-4 
                        file:rounded-full file:border-0 
                        file:text-sm file:font-semibold 
                        file:bg-yellow-50 file:text-yellow-700 
                        hover:file:bg-purple-100"
                    />
                    {file && (
                      <button
                        type="button"
                        onClick={resetForm}
                        className="shrink-0 rounded-md border border-red-500 px-4 py-2 text-sm font-medium text-red-500 transition-colors hover:bg-red-500 hover:text-white"
                      >
                        Clear
                      </button>
                    )}
                  </div>
                </div>

                <div className="space-y-2">
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    AI Question About CSV
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
                    Follow-up Chat {aiContext ? '(uses uploaded CSV)' : '(general)'}
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
              <div className="min-w-0 bg-white dark:bg-gray-800 shadow rounded-lg">
                <div className="min-w-0 max-w-full px-4 py-5 text-gray-700 dark:text-gray-300 sm:p-6">
                  {response ? (
                    <ReactMarkdown
                      remarkPlugins={[remarkGfm, remarkBreaks]}
                      components={markdownComponents}
                    >
                      {response}
                    </ReactMarkdown>
                  ) : (
                    <p>Your response will appear here.</p>
                  )}
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
