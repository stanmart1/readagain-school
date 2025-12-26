import { motion } from 'framer-motion';

export default function MasonryFeatures() {
  const features = [
    {
      icon: 'ri-book-open-line',
      title: 'Advanced E-Reader',
      description: 'Customizable reading experience with adjustable fonts, themes, and bookmarks. Read comfortably on any device.',
      color: 'from-blue-500 to-cyan-500',
      size: 'tall'
    },
    {
      icon: 'ri-team-line',
      title: 'Book Assignments',
      description: 'Teachers can assign books to classes and track which students have completed their reading assignments.',
      color: 'from-primary-600 to-accent-500',
      size: 'tall'
    },
    {
      icon: 'ri-bar-chart-line',
      title: 'Reading Analytics',
      description: 'Track reading progress, time spent, and completion rates. Get insights into student reading habits.',
      color: 'from-green-500 to-emerald-500',
      size: 'short'
    },
    {
      icon: 'ri-building-2-line',
      title: 'Digital Library',
      description: 'Access thousands of books instantly. Build your school library with books available to all students.',
      color: 'from-orange-500 to-red-500',
      size: 'short'
    },
    {
      icon: 'ri-download-cloud-line',
      title: 'Offline Reading',
      description: 'Download books for offline access. Students can read anytime, anywhere without internet connection.',
      color: 'from-indigo-500 to-blue-500',
      size: 'short'
    },
    {
      icon: 'ri-device-line',
      title: 'Multi-Device Sync',
      description: 'Start reading on one device and continue on another. Progress syncs automatically across all devices.',
      color: 'from-pink-500 to-rose-500',
      size: 'short'
    }
  ];

  return (
    <section className="py-20 bg-gradient-to-br from-gray-50 to-primary-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section Header */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-center mb-16"
        >
          <span className="inline-flex items-center px-4 py-2 rounded-full text-sm font-medium bg-primary-100 text-primary-800 mb-4">
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
        
        {/* Masonry Grid - Desktop, Regular Grid - Mobile */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 lg:auto-rows-[240px]">
          {features.map((feature, index) => (
            <motion.div
              key={index}
              initial={{ opacity: 0, y: 30 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: index * 0.1 }}
              whileHover={{ y: -8, scale: 1.02 }}
              className={`bg-white rounded-2xl p-6 shadow-lg hover:shadow-2xl transition-all duration-300 flex flex-col ${
                feature.size === 'tall' ? 'lg:row-span-2' : 'lg:row-span-1'
              }`}
            >
              <div className={`w-14 h-14 rounded-xl flex items-center justify-center bg-gradient-to-r ${feature.color} shadow-lg mb-4 flex-shrink-0`}>
                <i className={`${feature.icon} text-white text-2xl`}></i>
              </div>
              <h3 className="text-lg font-bold text-gray-900 mb-2 flex-shrink-0">
                {feature.title}
              </h3>
              <p className="text-gray-600 leading-relaxed text-sm overflow-hidden">
                {feature.description}
              </p>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
