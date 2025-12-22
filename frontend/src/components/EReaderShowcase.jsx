import { useState } from 'react';
import { motion } from 'framer-motion';

export default function EReaderShowcase() {
  const [fontSize, setFontSize] = useState(16);
  const [theme, setTheme] = useState('light');

  const features = [
    {
      icon: 'ri-book-open-line',
      title: 'Advanced E-Reader',
      description: 'Customizable reading experience with adjustable fonts, themes, and bookmarks. Read comfortably on any device.',
      color: 'from-blue-500 to-cyan-500'
    },
    {
      icon: 'ri-team-line',
      title: 'Book Assignments',
      description: 'Teachers can assign books to classes and track which students have completed their reading assignments.',
      color: 'from-purple-500 to-pink-500'
    },
    {
      icon: 'ri-bar-chart-line',
      title: 'Reading Analytics',
      description: 'Track reading progress, time spent, and completion rates. Get insights into student reading habits.',
      color: 'from-green-500 to-emerald-500'
    },
    {
      icon: 'ri-library-line',
      title: 'Digital Library',
      description: 'Access thousands of books instantly. Build your school library with books available to all students.',
      color: 'from-orange-500 to-red-500'
    },
    {
      icon: 'ri-download-cloud-line',
      title: 'Offline Reading',
      description: 'Download books for offline access. Students can read anytime, anywhere without internet connection.',
      color: 'from-indigo-500 to-blue-500'
    },
    {
      icon: 'ri-device-line',
      title: 'Multi-Device Sync',
      description: 'Start reading on one device and continue on another. Progress syncs automatically across all devices.',
      color: 'from-pink-500 to-rose-500'
    }
  ];

  const sampleText = `Welcome to your school's digital library! Access thousands of books, track your reading progress, and discover new stories every day.

Teachers can assign books to classes, monitor student progress, and get detailed analytics on reading habits. Students can read on any device, download books for offline access, and customize their reading experience.

Our platform makes it easy to build a culture of reading in your school. With features designed for both educators and learners, every student can find books they love and track their reading journey.`;

  const themes = {
    light: { bg: 'bg-white', text: 'text-gray-900', border: 'border-gray-200' },
    dark: { bg: 'bg-gray-900', text: 'text-gray-100', border: 'border-gray-700' },
    sepia: { bg: 'bg-amber-50', text: 'text-amber-900', border: 'border-amber-200' }
  };

  const currentTheme = themes[theme];

  const cycleTheme = () => {
    setTheme(theme === 'light' ? 'dark' : theme === 'dark' ? 'sepia' : 'light');
  };

  return (
    <section className="py-20 bg-gradient-to-br from-blue-50 to-purple-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section Header */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-center mb-16"
        >
          <span className="inline-flex items-center px-4 py-2 rounded-full text-sm font-medium bg-purple-100 text-purple-800 mb-4">
            <i className="ri-star-line mr-2"></i>
            Platform Features
          </span>
          <h2 className="text-4xl md:text-5xl font-bold text-gray-900 mb-6">
            Everything Your School Needs
          </h2>
          <p className="text-lg md:text-xl text-gray-600 max-w-3xl mx-auto">
            Powerful features designed for schools, teachers, and students. 
            From digital library management to reading analytics and more.
          </p>
        </motion.div>
        
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 md:gap-12 items-center">
          {/* Left Column - Features */}
          <div className="space-y-4 md:space-y-6 order-2 lg:order-1">
            {features.map((benefit, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, x: -50 }}
                whileInView={{ opacity: 1, x: 0 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.1 }}
                className="flex items-start space-x-4 p-4 rounded-xl hover:bg-white/50 transition-all duration-300"
              >
                <div className={`w-12 h-12 rounded-xl flex items-center justify-center bg-gradient-to-r ${benefit.color} shadow-lg flex-shrink-0`}>
                  <i className={`${benefit.icon} text-white text-xl`}></i>
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 mb-2">
                    {benefit.title}
                  </h3>
                  <p className="text-gray-600 leading-relaxed">
                    {benefit.description}
                  </p>
                </div>
              </motion.div>
            ))}
          </div>
          
          {/* Right Column - E-Reader Simulator */}
          <motion.div
            initial={{ opacity: 0, x: 50 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
            className="relative order-1 lg:order-2"
          >
            <div className={`${currentTheme.bg} ${currentTheme.border} border-2 rounded-2xl shadow-2xl overflow-hidden w-full max-w-md mx-auto transition-all duration-300`}>
              {/* Reader Header */}
              <div className={`${currentTheme.bg} ${currentTheme.border} border-b px-3 md:px-4 py-2 md:py-3 flex items-center justify-between`}>
                <div className="flex items-center space-x-2 md:space-x-3">
                  <i className="ri-book-line text-base md:text-lg text-blue-600"></i>
                  <span className={`font-medium text-sm md:text-base ${currentTheme.text}`}>Sample Book</span>
                </div>
                <div className="flex items-center space-x-1 md:space-x-2">
                  <button
                    onClick={() => setFontSize(Math.max(12, fontSize - 2))}
                    className={`p-1 rounded text-sm md:text-base ${currentTheme.text} hover:bg-gray-100 transition-colors`}
                  >
                    <i className="ri-subtract-line"></i>
                  </button>
                  <span className={`text-xs md:text-sm ${currentTheme.text}`}>{fontSize}px</span>
                  <button
                    onClick={() => setFontSize(Math.min(24, fontSize + 2))}
                    className={`p-1 rounded text-sm md:text-base ${currentTheme.text} hover:bg-gray-100 transition-colors`}
                  >
                    <i className="ri-add-line"></i>
                  </button>
                </div>
              </div>

              {/* Reader Content */}
              <div className="p-4 md:p-6 h-[300px] md:h-[400px] overflow-y-auto">
                <div 
                  className={`${currentTheme.text} leading-relaxed transition-all duration-300 select-text`}
                  style={{ fontSize: `${fontSize}px` }}
                >
                  {sampleText}
                </div>
              </div>

              {/* Reader Footer */}
              <div className={`${currentTheme.bg} ${currentTheme.border} border-t px-3 md:px-4 py-2 md:py-3 flex items-center justify-between`}>
                <div className="flex items-center space-x-2 md:space-x-4">
                  <button
                    onClick={cycleTheme}
                    className={`p-1.5 md:p-2 rounded-lg text-sm md:text-base ${currentTheme.text} hover:bg-gray-100 transition-colors`}
                    title="Change theme"
                  >
                    <i className="ri-contrast-2-line"></i>
                  </button>
                  <button className={`p-1.5 md:p-2 rounded-lg text-sm md:text-base ${currentTheme.text} hover:bg-gray-100 transition-colors`} title="Bookmark">
                    <i className="ri-bookmark-line"></i>
                  </button>
                  <button className={`p-1.5 md:p-2 rounded-lg text-sm md:text-base ${currentTheme.text} hover:bg-gray-100 transition-colors`} title="Highlight">
                    <i className="ri-edit-line"></i>
                  </button>
                </div>
                <div className={`text-xs md:text-sm ${currentTheme.text}`}>
                  Page 1 of 156
                </div>
              </div>
            </div>

            {/* Floating Badge */}
            <motion.div
              animate={{ y: [0, -10, 0] }}
              transition={{ duration: 2, repeat: Infinity }}
              className="hidden md:flex absolute -bottom-4 -left-4 w-12 h-12 bg-gradient-to-r from-green-400 to-blue-500 rounded-full items-center justify-center shadow-lg"
            >
              <i className="ri-bookmark-line text-white text-lg"></i>
            </motion.div>
          </motion.div>
        </div>
      </div>
    </section>
  );
}
