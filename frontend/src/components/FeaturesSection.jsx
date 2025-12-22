import { motion } from 'framer-motion';
import { useState, useRef, useEffect } from 'react';

export default function FeaturesSection() {
  const [currentIndex, setCurrentIndex] = useState(0);
  const carouselRef = useRef(null);

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
      icon: 'ri-building-2-line',
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

  const nextSlide = () => {
    setCurrentIndex((prev) => (prev + 1) % features.length);
  };

  const prevSlide = () => {
    setCurrentIndex((prev) => (prev - 1 + features.length) % features.length);
  };

  useEffect(() => {
    const interval = setInterval(nextSlide, 5000);
    return () => clearInterval(interval);
  }, []);

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
        
        {/* Carousel Container */}
        <div className="relative">
          <div className="overflow-hidden" ref={carouselRef}>
            <motion.div
              className="flex transition-transform duration-500 ease-out"
              style={{ transform: `translateX(-${currentIndex * 100}%)` }}
            >
              {features.map((feature, index) => (
                <div key={index} className="w-full flex-shrink-0 px-4">
                  <motion.div
                    initial={{ opacity: 0, y: 30 }}
                    whileInView={{ opacity: 1, y: 0 }}
                    viewport={{ once: true }}
                    className="bg-white rounded-2xl p-8 shadow-lg max-w-md mx-auto"
                  >
                    <div className={`w-16 h-16 rounded-xl flex items-center justify-center bg-gradient-to-r ${feature.color} shadow-lg mb-6`}>
                      <i className={`${feature.icon} text-white text-3xl`}></i>
                    </div>
                    <h3 className="text-xl font-bold text-gray-900 mb-3">
                      {feature.title}
                    </h3>
                    <p className="text-gray-600 leading-relaxed">
                      {feature.description}
                    </p>
                  </motion.div>
                </div>
              ))}
            </motion.div>
          </div>

          {/* Navigation Arrows */}
          <button
            onClick={prevSlide}
            className="absolute left-0 top-1/2 -translate-y-1/2 -translate-x-4 bg-white rounded-full p-3 shadow-lg hover:bg-gray-50 transition-colors z-10"
          >
            <i className="ri-arrow-left-s-line text-2xl text-gray-800"></i>
          </button>
          <button
            onClick={nextSlide}
            className="absolute right-0 top-1/2 -translate-y-1/2 translate-x-4 bg-white rounded-full p-3 shadow-lg hover:bg-gray-50 transition-colors z-10"
          >
            <i className="ri-arrow-right-s-line text-2xl text-gray-800"></i>
          </button>

          {/* Dots Indicator */}
          <div className="flex justify-center gap-2 mt-8">
            {features.map((_, index) => (
              <button
                key={index}
                onClick={() => setCurrentIndex(index)}
                className={`w-2 h-2 rounded-full transition-all ${
                  index === currentIndex ? 'bg-purple-600 w-8' : 'bg-gray-300'
                }`}
              />
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
