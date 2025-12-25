import { motion } from 'framer-motion';
import { useState, useRef, useEffect } from 'react';

export default function FeaturesSection() {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [cardsPerView, setCardsPerView] = useState(3);
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
      color: 'from-primary-600 to-accent-500'
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

  const infiniteFeatures = [...features, ...features, ...features];

  const nextSlide = () => {
    setCurrentIndex((prev) => prev + 1);
  };

  const prevSlide = () => {
    setCurrentIndex((prev) => prev - 1);
  };

  useEffect(() => {
    const handleResize = () => {
      if (window.innerWidth >= 1024) {
        setCardsPerView(3);
      } else if (window.innerWidth >= 768) {
        setCardsPerView(2);
      } else {
        setCardsPerView(1);
      }
    };

    handleResize();
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  useEffect(() => {
    setCurrentIndex(features.length);
  }, [features.length]);

  useEffect(() => {
    const interval = setInterval(nextSlide, 4000);
    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    if (currentIndex < features.length) {
      const carousel = carouselRef.current?.querySelector('div');
      if (carousel) {
        carousel.style.transition = 'none';
        setCurrentIndex(features.length * 2 - 1);
        setTimeout(() => {
          carousel.style.transition = '';
        }, 50);
      }
    }
  }, [currentIndex, features.length]);

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
          {/* Navigation Arrows */}
          <button
            onClick={prevSlide}
            className="absolute left-0 top-1/2 -translate-y-1/2 z-10 bg-white/90 hover:bg-white p-3 rounded-full shadow-lg transition-all"
          >
            <i className="ri-arrow-left-s-line text-2xl text-gray-800"></i>
          </button>
          <button
            onClick={nextSlide}
            className="absolute right-0 top-1/2 -translate-y-1/2 z-10 bg-white/90 hover:bg-white p-3 rounded-full shadow-lg transition-all"
          >
            <i className="ri-arrow-right-s-line text-2xl text-gray-800"></i>
          </button>

          <div className="overflow-hidden" ref={carouselRef}>
            <div
              className="flex transition-transform duration-500 ease-out"
              style={{ transform: `translateX(-${currentIndex * (100 / cardsPerView)}%)` }}
              onTransitionEnd={() => {
                if (currentIndex >= features.length * 2) {
                  const carousel = carouselRef.current?.querySelector('div');
                  if (carousel) {
                    carousel.style.transition = 'none';
                    setCurrentIndex(features.length);
                    setTimeout(() => {
                      carousel.style.transition = '';
                    }, 50);
                  }
                }
              }}
            >
              {infiniteFeatures.map((feature, index) => (
                <div 
                  key={index} 
                  className="px-4"
                  style={{ minWidth: `${100 / cardsPerView}%` }}
                >
                  <motion.div
                    initial={{ opacity: 0, y: 30 }}
                    whileInView={{ opacity: 1, y: 0 }}
                    viewport={{ once: true }}
                    className="bg-white rounded-2xl p-8 shadow-lg h-80 flex flex-col"
                  >
                    <div className={`w-16 h-16 rounded-xl flex items-center justify-center bg-gradient-to-r ${feature.color} shadow-lg mb-6 flex-shrink-0`}>
                      <i className={`${feature.icon} text-white text-3xl`}></i>
                    </div>
                    <h3 className="text-xl font-bold text-gray-900 mb-3 flex-shrink-0">
                      {feature.title}
                    </h3>
                    <p className="text-gray-600 leading-relaxed flex-grow">
                      {feature.description}
                    </p>
                  </motion.div>
                </div>
              ))}
            </div>
          </div>

          {/* Dots Indicator */}
          <div className="flex justify-center gap-2 mt-8">
            {features.map((_, index) => (
              <button
                key={index}
                onClick={() => setCurrentIndex(features.length + index)}
                className={`w-2 h-2 rounded-full transition-all ${
                  (currentIndex % features.length) === index ? 'bg-purple-600 w-8' : 'bg-gray-300'
                }`}
              />
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
