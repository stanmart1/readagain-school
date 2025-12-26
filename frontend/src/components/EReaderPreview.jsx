import { motion, AnimatePresence } from 'framer-motion';
import { useState } from 'react';

export default function EReaderPreview() {
  const [fontSize, setFontSize] = useState(16);
  const [fontFamily, setFontFamily] = useState('serif');
  const [theme, setTheme] = useState('light');
  const [page, setPage] = useState(1);
  const [showSettings, setShowSettings] = useState(false);
  const totalPages = 2;
  const progress = Math.round((page / totalPages) * 100);

  const themes = {
    light: { bg: 'bg-white', text: 'text-gray-700', heading: 'text-gray-900' },
    sepia: { bg: 'bg-amber-50', text: 'text-amber-900', heading: 'text-amber-950' },
    dark: { bg: 'bg-gray-900', text: 'text-gray-300', heading: 'text-white' }
  };

  const fonts = {
    serif: 'font-serif',
    sans: 'font-sans',
    mono: 'font-mono'
  };

  const currentTheme = themes[theme];
  const currentFont = fonts[fontFamily];

  const bookContent = [
    {
      title: "Chapter 1: The Beginning",
      paragraphs: [
        "In the heart of the ancient library, where dust motes danced in shafts of golden sunlight, young Emma discovered a book that would change her life forever. The leather-bound tome seemed to whisper secrets of forgotten worlds and untold adventures.",
        "As she carefully opened the first page, the scent of aged paper filled her senses, and the words began to glow with an ethereal light. This was no ordinary book—it was a gateway to knowledge, imagination, and endless possibilities.",
        "The library keeper had warned her about the power of books, but Emma never imagined..."
      ]
    },
    {
      title: "Chapter 2: The Discovery",
      paragraphs: [
        "Emma's fingers trembled as she turned the page. The ancient text revealed stories of brave heroes, magical creatures, and distant lands waiting to be explored.",
        "Each word seemed to come alive, painting vivid pictures in her mind. She could almost hear the rustling of leaves in enchanted forests and feel the warmth of distant suns.",
        "This was the moment she had been waiting for—the beginning of her greatest adventure..."
      ]
    }
  ];

  const currentContent = bookContent[page - 1] || bookContent[0];

  const nextPage = () => {
    if (page < totalPages) setPage(page + 1);
  };

  const prevPage = () => {
    if (page > 1) setPage(page - 1);
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
                <span className="font-semibold text-sm">E-Reader Preview</span>
              </div>
              <div className="flex items-center gap-3">
                <button 
                  onClick={() => setShowSettings(!showSettings)}
                  className="hover:scale-110 transition-transform"
                >
                  <i className="ri-settings-3-line text-lg"></i>
                </button>
              </div>
            </div>

            {/* Book Content */}
            <div className="p-6 pb-24 h-96 overflow-hidden relative">
              <AnimatePresence mode="wait">
                <motion.div
                  key={page}
                  initial={{ opacity: 0, x: 20 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -20 }}
                  transition={{ duration: 0.3 }}
                >
                  <h2 
                    className={`font-bold ${currentTheme.heading} mb-4 transition-colors duration-300 ${currentFont}`} 
                    style={{ fontSize: `${fontSize + 8}px` }}
                  >
                    {currentContent.title}
                  </h2>
                  <div className={`space-y-4 ${currentTheme.text} leading-relaxed transition-colors duration-300 ${currentFont}`} style={{ fontSize: `${fontSize}px` }}>
                    {currentContent.paragraphs.map((para, idx) => (
                      <p key={idx} className={`text-justify ${idx === currentContent.paragraphs.length - 1 ? 'opacity-50' : ''}`}>
                        {para}
                      </p>
                    ))}
                  </div>
                </motion.div>
              </AnimatePresence>

              {/* Navigation Arrows */}
              <div className="absolute top-1/2 -translate-y-1/2 left-0 right-0 flex justify-between px-6">
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
                  disabled={page >= bookContent.length}
                  className={`w-10 h-10 rounded-full flex items-center justify-center transition-all ${
                    page >= bookContent.length
                      ? 'bg-gray-200 text-gray-400 cursor-not-allowed' 
                      : 'bg-primary-600 text-white hover:bg-primary-700 hover:scale-110'
                  }`}
                >
                  <i className="ri-arrow-right-s-line text-xl"></i>
                </button>
              </div>

              {/* Progress Bar - with theme-aware background */}
              <div className={`absolute bottom-4 left-4 right-4 ${currentTheme.bg} backdrop-blur-sm rounded-lg p-3 shadow-md border transition-colors duration-300 ${
                theme === 'dark' ? 'border-gray-700' : 'border-gray-200'
              }`}>
                <div className={`flex items-center justify-between text-xs mb-2 font-medium ${currentTheme.text} transition-colors duration-300`}>
                  <span>Page {page} of {totalPages}</span>
                  <span>{progress}% Complete</span>
                </div>
                <div className={`h-2 rounded-full overflow-hidden border transition-colors duration-300 ${
                  theme === 'dark' ? 'bg-gray-700 border-gray-600' : 'bg-gray-200 border-gray-300'
                }`}>
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

        {/* Settings Modal */}
        <AnimatePresence>
          {showSettings && (
            <>
              {/* Backdrop */}
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                onClick={() => setShowSettings(false)}
                className="fixed inset-0 z-40"
              />
              
              {/* Modal */}
              <motion.div
                initial={{ opacity: 0, scale: 0.9, y: -20 }}
                animate={{ opacity: 1, scale: 1, y: 0 }}
                exit={{ opacity: 0, scale: 0.9, y: -20 }}
                className="absolute top-16 right-4 bg-white rounded-2xl shadow-2xl p-4 w-64 z-50 border border-gray-200"
              >
                <div className="space-y-4">
                <div>
                  <label className="text-sm font-semibold text-gray-700 mb-2 block">Font Size</label>
                  <div className="flex gap-2">
                    {[14, 16, 18, 20].map(size => (
                      <button
                        key={size}
                        onClick={() => setFontSize(size)}
                        className={`flex-1 py-2 rounded-lg text-sm font-medium transition-all ${
                          fontSize === size
                            ? 'bg-primary-600 text-white'
                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                        }`}
                      >
                        {size}px
                      </button>
                    ))}
                  </div>
                </div>

                <div>
                  <label className="text-sm font-semibold text-gray-700 mb-2 block">Font Family</label>
                  <div className="flex gap-2">
                    {[
                      { key: 'serif', label: 'Serif' },
                      { key: 'sans', label: 'Sans' },
                      { key: 'mono', label: 'Mono' }
                    ].map(font => (
                      <button
                        key={font.key}
                        onClick={() => setFontFamily(font.key)}
                        className={`flex-1 py-2 rounded-lg text-sm font-medium transition-all ${
                          fontFamily === font.key
                            ? 'bg-primary-600 text-white'
                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                        }`}
                      >
                        {font.label}
                      </button>
                    ))}
                  </div>
                </div>

                <div>
                  <label className="text-sm font-semibold text-gray-700 mb-2 block">Theme</label>
                  <div className="flex gap-2">
                    {[
                      { key: 'light', label: 'Light', icon: 'ri-sun-line' },
                      { key: 'sepia', label: 'Sepia', icon: 'ri-contrast-2-line' },
                      { key: 'dark', label: 'Dark', icon: 'ri-moon-line' }
                    ].map(t => (
                      <button
                        key={t.key}
                        onClick={() => setTheme(t.key)}
                        className={`flex-1 py-2 rounded-lg text-sm font-medium transition-all flex items-center justify-center gap-1 ${
                          theme === t.key
                            ? 'bg-primary-600 text-white'
                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                        }`}
                      >
                        <i className={t.icon}></i>
                        {t.label}
                      </button>
                    ))}
                  </div>
                </div>
              </div>
            </>
          )}
        </AnimatePresence>

        {/* Glow Effect */}
        <div className="absolute inset-0 bg-primary-400/20 rounded-3xl blur-3xl -z-10"></div>
      </div>
    </motion.div>
  );
}
