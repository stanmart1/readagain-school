import { motion } from 'framer-motion';
import { Link } from 'react-router-dom';
import EReaderPreview from './EReaderPreview';

export default function HeroSection() {
  return (
    <div className="relative min-h-screen flex items-center overflow-hidden pt-20">
      {/* Gradient Background */}
      <div className="absolute inset-0 bg-gradient-to-br from-primary-900 via-primary-800 to-primary-700"></div>
      
      {/* Decorative circles */}
      <div className="absolute top-20 right-10 w-72 h-72 bg-primary-600/20 rounded-full blur-3xl"></div>
      <div className="absolute bottom-20 left-10 w-96 h-96 bg-primary-500/10 rounded-full blur-3xl"></div>
      
      {/* Main Content */}
      <div className="relative z-10 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 w-full py-20">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          
          {/* Left Side - Text Content */}
          <motion.div 
            initial={{ opacity: 0, x: -50 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.8 }}
            className="text-left -mt-32"
          >
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
              className="inline-block mb-4"
            >
              <span className="bg-yellow-400 text-primary-900 px-4 py-2 rounded-full text-sm font-bold">
                📚 School Library Platform
              </span>
            </motion.div>

            <motion.h1 
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.8, delay: 0.3 }}
              className="text-4xl md:text-5xl lg:text-6xl font-bold mb-6 leading-tight text-white"
            >
              Build Your School's
              <span className="block text-yellow-400 mt-2">Digital Library</span>
            </motion.h1>
            
            <motion.p 
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.8, delay: 0.4 }}
              className="text-lg lg:text-xl text-primary-100 mb-8 leading-relaxed max-w-xl"
            >
              Empower students with access to thousands of books. Track reading progress, assign books to classes, and foster a love for reading.
            </motion.p>
            
            <motion.div 
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.8, delay: 0.5 }}
              className="flex flex-col sm:flex-row gap-4"
            >
              <Link 
                to="/signup" 
                className="bg-gradient-to-r from-yellow-400 to-yellow-500 text-primary-900 px-8 py-4 rounded-lg font-bold hover:from-yellow-500 hover:to-yellow-600 transition-all duration-300 transform hover:scale-105 shadow-lg hover:shadow-xl flex items-center justify-center gap-2"
              >
                <span>Get Started Free</span>
                <i className="ri-arrow-right-line"></i>
              </Link>
              <Link 
                to="/books" 
                className="bg-white/10 backdrop-blur-sm text-white border-2 border-white/30 px-8 py-4 rounded-lg font-bold hover:bg-white/20 transition-all duration-300 flex items-center justify-center gap-2"
              >
                <i className="ri-book-line"></i>
                <span>Browse Books</span>
              </Link>
            </motion.div>
          </motion.div>
          
          {/* Right Side - E-Reader Preview */}
          <div className="relative hidden lg:block">
            <EReaderPreview />
          </div>

        </div>
      </div>
    </div>
  );
}
