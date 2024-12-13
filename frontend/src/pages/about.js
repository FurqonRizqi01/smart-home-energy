import React from 'react';

const About = () => {
  return (
    <div className="container mx-auto px-4 py-8 flex-grow">
      <div className="max-w-xl mx-auto">
        <h1 className="text-3xl font-bold mb-4">About Data Analysis Chatbot</h1>
        <p className="mb-4">
          This is an advanced AI-powered data analysis and chat application 
          that helps you gain insights from your uploaded files and interact 
          with an intelligent AI assistant.
        </p>
        <h2 className="text-2xl font-semibold mb-2">Features</h2>
        <ul className="list-disc pl-5">
          <li>File Upload Analysis</li>
          <li>Tapas AI Table Querying</li>
          <li>Mini Chat AI Interaction</li>
        </ul>
      </div>
    </div>
  );
};

export default About;