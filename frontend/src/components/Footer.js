import React from 'react';
import { FaGithub, FaEnvelope, FaArrowUp } from 'react-icons/fa';

function Footer() {
  return (
    <footer className="bg-gray-800 text-white py-6">
      <div className="container mx-auto px-4">
        <div className="flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">
          {/* Left Section */}
          <div className="text-left">
            © {new Date().getFullYear()} Muhammad Furqon Rizqi - All rights reserved.
          </div>
          
          {/* Right Section */}
          <div className="flex space-x-4">
            <a 
              href="https://github.com/FurqonRizqi01" 
              target="_blank" 
              rel="noopener noreferrer"
              className="hover:text-purple-400 transition duration-300"
              aria-label="GitHub Repository"
            >
              <FaGithub size={24} />
            </a>
            <a 
              href="mailto:contact@example.com"
              className="hover:text-purple-400 transition duration-300"
              aria-label="Contact Email"
            >
              <FaEnvelope size={24} />
            </a>
            <button 
              onClick={() => window.scrollTo(0, 0)}
              className="hover:text-purple-400 transition duration-300"
              aria-label="Back to Top"
            >
              <FaArrowUp size={24} />
            </button>
          </div>
        </div>
      </div>
    </footer>
  );
}

export default Footer;