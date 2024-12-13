import React from 'react';

function Footer() {
  return (
    <footer className="bg-gray-800 text-white py-6">
      <div className="container mx-auto px-4">
        <div className="flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">
          {/* Bagian Kiri */}
          <div className="text-left">
            © {new Date().getFullYear()} Muhammad Furqon Rizqi - All rights reserved.
          </div>
          
          {/* Bagian Kanan */}
          <div className="flex space-x-4">
            <a 
              href="https://github.com/yourusername/project" 
              target="_blank" 
              rel="noopener noreferrer"
              className="hover:text-purple-400 transition duration-300"
            >
              GitHub
            </a>
            <a 
              href="mailto:contact@example.com"
              className="hover:text-purple-400 transition duration-300"
            >
              Contact
            </a>
            <button 
              onClick={() => window.scrollTo(0, 0)}
              className="hover:text-purple-400 transition duration-300"
            >
              Back to Top
            </button>
          </div>
        </div>
      </div>
    </footer>
  );
}

export default Footer;