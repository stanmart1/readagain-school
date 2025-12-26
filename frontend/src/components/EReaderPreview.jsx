import { motion } from 'framer-motion';
import { useState } from 'react';

export default function EReaderPreview() {
  const [fontSize, setFontSize] = useState(16);
  const [theme, setTheme] = useState('light');
  const [page, setPage] = useState(1);
  const totalPages = 245;
  const progress = Math.round((page / totalPages) * 100);

  const themes = {
    light: { bg: 'bg-white', text: 'text-gray-700', heading: 'text-gray-900' },
    sepia: { bg: 'bg-amber-50', text: 'text-amber-900', heading: 'text-amber-950' },
    dark: { bg: 'bg-gray-900', text: 'text-gray-300', heading: 'text-white' }
  };

  const currentTheme = themes[theme];

  const nextPage = () => {
    if (page < totalPages) setPage(page + 1);
  };

  const prevPage = () => {
    if (page > 1) setPage(page - 1);
  };

  const cycleFontSize = () => {
    setFontSize(prev => prev >= 20 ? 14 : prev + 2);
  };

  const cycleTheme = () => {
    const themeOrder = ['light', 'sepia', 'dark'];
    const currentIndex = themeOrder.indexOf(theme);
    setTheme(themeOrder[(currentIndex + 1) % themeOrder.length]);
  };

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.8, delay: 0.3 }}
      className="relative"
    >
      {/* E-Reader Device */}
      <div className="relative w-full max-w-md mx-auto">
        {/* Device Frame */}
        <div className="bg-gray-800 rounded-3xl p-4 shadow-2xl">
          {/* Screen */}
          <div className={`${currentTheme.bg} rounded-2xl overflow-hidden shadow-inner transition-colors duration-300`}>
            {/* Reader Header */}
            <div className="bg-primary-600 text-white px-4 py-3 flex items-center justify-between">
              <div className="flex items-center gap-2">
                <i className="ri-book-open-line text-xl"></i>
                <span className="font-semibold text-sm">ReadAgain Reader</span>
              </div>
              <div className="flex items-center gap-3">
                <button className="hover:scale-110 transition-transform">
                  <i className="ri-bookmark-line text-lg"></i>
                </button>
                <button className="hover:scale-110 transition-transform">
                  <i className="ri-settings-3-line text-lg"></i>
                </button>
              </div>
            </div>

            {/* Book Content */}
            <div className="p-6 h-96 overflow-hidden relative">
              <h2 className={`text-2xl font-bold ${currentTheme.heading} mb-4 transition-colors duration-300`} style={{ fontSize: `${fontSize + 8}px` }}>
                Chapter 1: The Beginning
              </h2>
              <div className={`space-y-4 ${currentTheme.text} leading-relaxed transition-colors duration-300`} style={{ fontSize: `${fontSize}px` }}>
                <p className="text-justify">
                  In the heart of the ancient library, where dust motes danced in shafts of golden sunlight, 
                  young Emma discovered a book that would change her life forever. The leather-bound tome 
                  seemed to whisper secrets of forgotten worlds and untold adventures.
                </p>
                <p className="text-justify">
                  As she carefully opened the first page, the scent of aged paper filled her senses, 
                  and the words began to glow with an ethereal light. This was no ordinary book—it was 
                  a gateway to knowledge, imagination, and endless possibilities.
                </p>
                <p className="text-justify opacity-50">
                  The library keeper had warned her about the power of books, but Emma never imagined...
                </p>
              </div>

              {/* Navigation Arrows */}
              <div className="absolute bottom-20 left-0 right-0 flex justify-between px-6">
                <button
                  onClick={prevPage}
                  disabled={page === 1}
                  className={`w-10 h-10 rounded-full flex items-center justify-center transition-all ${
                    page === 1 
                      ? 'bg-gray-200 text-gray-400 cursor-not-allowed' 
                      : 'bg-primary-600 text-white hover:bg-primary-700 hover:scale-110'
                  }`}
                >
                  <i className="ri-arrow-left-s-line text-xl"></i>
                </button>
                <button
                  onClick={nextPage}
                  disabled={page === totalPages}
                  className={`w-10 h-10 rounded-full flex items-center justify-center transition-all ${
                    page === totalPages 
                      ? 'bg-gray-200 text-gray-400 cursor-not-allowed' 
                      : 'bg-primary-600 text-white hover:bg-primary-700 hover:scale-110'
                  }`}
                >
                  <i className="ri-arrow-right-s-line text-xl"></i>
                </button>
              </div>

              {/* Progress Bar */}
              <div className="absolute bottom-6 left-6 right-6">
                <div className={`flex items-center justify-between text-xs mb-2 ${currentTheme.text} transition-colors duration-300`}>
                  <span>Page {page} of {totalPages}</span>
                  <span>{progress}% Complete</span>
                </div>
                <div className="h-1.5 bg-gray-200 rounded-full overflow-hidden border border-gray-300">
                  <motion.div
                    animate={{ width: `${progress}%` }}
                    transition={{ duration: 0.3 }}
                    className="h-full bg-gradient-to-r from-primary-500 to-primary-600"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Floating Action Buttons */}
        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.8 }}
          className="absolute -right-4 top-1/4 flex flex-col gap-3"
        >
          <button
            onClick={cycleFontSize}
            className="w-12 h-12 bg-white rounded-full shadow-lg flex items-center justify-center text-primary-600 hover:bg-primary-50 transition-all hover:scale-110 active:scale-95"
            title={`Font size: ${fontSize}px`}
          >
            <i className="ri-font-size text-xl"></i>
          </button>
          <button
            onClick={cycleTheme}
            className="w-12 h-12 bg-white rounded-full shadow-lg flex items-center justify-center text-primary-600 hover:bg-primary-50 transition-all hover:scale-110 active:scale-95"
            title={`Theme: ${theme}`}
          >
            <i className="ri-palette-line text-xl"></i>
          </button>
        </motion.div>

        {/* Glow Effect */}
        <div className="absolute inset-0 bg-primary-400/20 rounded-3xl blur-3xl -z-10"></div>
      </div>
    </motion.div>
  );
}
