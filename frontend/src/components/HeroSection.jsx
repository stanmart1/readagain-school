import { motion } from 'framer-motion';
import { Link } from 'react-router-dom';

export default function HeroSection() {
  const floatingBooks = [
    { id: 1, rotation: -8, delay: 0, y: 20 },
    { id: 2, rotation: 5, delay: 0.3, y: -30 },
    { id: 3, rotation: -3, delay: 0.6, y: 10 },
    { id: 4, rotation: 8, delay: 0.9, y: -20 }
  ];

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
            className="text-left"
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
          
          {/* Right Side - Floating Books */}
          <div className="relative hidden lg:block h-[600px]">
            {floatingBooks.map((book, index) => (
              <motion.div
                key={book.id}
                initial={{ opacity: 0, scale: 0.5, rotate: 0 }}
                animate={{ 
                  opacity: 1, 
                  scale: 1, 
                  rotate: book.rotation,
                  y: [book.y, book.y - 20, book.y]
                }}
                transition={{ 
                  opacity: { duration: 0.6, delay: book.delay },
                  scale: { duration: 0.6, delay: book.delay },
                  rotate: { duration: 0.6, delay: book.delay },
                  y: { duration: 3, repeat: Infinity, ease: "easeInOut", delay: book.delay }
                }}
                className="absolute"
                style={{
                  top: `${index * 18}%`,
                  left: `${(index % 2) * 40 + 10}%`,
                  zIndex: 5 - index
                }}
              >
                <motion.div
                  whileHover={{ scale: 1.1, rotate: 0, zIndex: 50 }}
                  className="relative group cursor-pointer"
                >
                  <div className="w-32 h-44 bg-gradient-to-br from-white to-gray-100 rounded-lg shadow-2xl overflow-hidden">
                    <img
                      src={`https://picsum.photos/seed/book${book.id}/200/300`}
                      alt={`Book ${book.id}`}
                      className="w-full h-full object-cover"
                      loading="lazy"
                    />
                  </div>
                  {/* Glow effect */}
                  <div className="absolute inset-0 bg-yellow-400/20 rounded-lg blur-xl group-hover:bg-yellow-400/40 transition-all"></div>
                </motion.div>
              </motion.div>
            ))}
          </div>

        </div>
      </div>
    </div>
  );
}
