import { motion } from 'framer-motion';
import { Link } from 'react-router-dom';
import { useAbout } from '../hooks';
import { getFileUrl } from '../lib/fileService';
import ProgressiveImage from './ProgressiveImage';
import { createHTMLProps } from '../utils/htmlUtils';

export default function AboutSection() {
  const { content, loading } = useAbout();

  const defaultContent = {
    hero: {
      title: 'About ReadAgain',
      subtitle: 'Empowering The Mind Through Reading'
    },
    mission: {
      description: "Our mission is to empower schools and students with access to quality digital books. We're building a comprehensive reading platform that helps teachers track progress, students discover new books, and schools foster a culture of reading."
    },
    values: [
      {
        icon: 'ri-book-open-line',
        title: 'Accessibility',
        description: '<p>Making reading accessible to everyone</p>'
      },
      {
        icon: 'ri-lightbulb-line',
        title: 'Innovation',
        description: '<p>Cutting-edge technology for modern learning</p>'
      },
      {
        icon: 'ri-team-line',
        title: 'Community',
        description: '<p>Building a community of passionate readers</p>'
      },
      {
        icon: 'ri-line-chart-line',
        title: 'Progress',
        description: '<p>Track and celebrate reading achievements</p>'
      }
    ]
  };

  const aboutContent = content || defaultContent;
  const values = aboutContent.values || defaultContent.values;

  if (loading) {
    return (
      <section className="py-20 bg-gradient-to-b from-white to-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
          </div>
        </div>
      </section>
    );
  }

  return (
    <section className="py-20 bg-gradient-to-b from-white to-gray-50 relative overflow-hidden">
      {/* Decorative background elements */}
      <div className="absolute top-20 right-0 w-72 h-72 bg-primary-100 rounded-full blur-3xl opacity-30"></div>
      <div className="absolute bottom-20 left-0 w-96 h-96 bg-yellow-100 rounded-full blur-3xl opacity-20"></div>
      
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        
        {/* Section Header */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-center mb-16"
        >
          <span className="inline-block px-4 py-2 bg-primary-100 text-primary-700 rounded-full text-sm font-semibold mb-4">
            Who We Are
          </span>
          <h2 className="text-4xl lg:text-5xl font-bold text-gray-900 mb-4">
            {aboutContent.hero?.title || 'About ReadAgain'}
          </h2>
          <div className="text-xl text-gray-600 max-w-3xl mx-auto"
               {...createHTMLProps(aboutContent.hero?.subtitle || 'Empowering The Mind Through Reading')}
          />
        </motion.div>

        {/* Main Content Grid */}
        <div className="grid lg:grid-cols-2 gap-12 items-center mb-16">
          
          {/* Image */}
          <motion.div
            initial={{ opacity: 0, x: -50 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.6 }}
            className="relative"
          >
            <div className="relative rounded-3xl overflow-hidden shadow-2xl">
              <ProgressiveImage
                src={getFileUrl(aboutContent.hero?.image_url) || "https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=800"}
                alt="ReadAgain - Empowering minds through reading"
                className="w-full h-[500px] object-cover"
              />
              {/* Overlay gradient */}
              <div className="absolute inset-0 bg-gradient-to-t from-black/30 to-transparent"></div>
            </div>
          </motion.div>

          {/* Mission Text */}
          <motion.div
            initial={{ opacity: 0, x: 50 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.6 }}
          >
            <h3 className="text-2xl font-bold text-gray-900 mb-4">Our Mission</h3>
            <div className="text-lg text-gray-600 leading-relaxed mb-6"
                 {...createHTMLProps(aboutContent.mission?.description || defaultContent.mission.description)}
            />
            
            <Link
              to="/about"
              className="inline-flex items-center gap-2 px-6 py-3 bg-primary-600 text-white rounded-full font-semibold hover:bg-primary-700 transition-all shadow-lg hover:shadow-xl"
            >
              <span>Learn More About Us</span>
              <i className="ri-arrow-right-line"></i>
            </Link>
          </motion.div>
        </div>

        {/* Values Grid */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="grid grid-cols-2 lg:grid-cols-4 gap-6"
        >
          {values.map((value, index) => (
            <motion.div
              key={index}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: index * 0.1 }}
              className="bg-white rounded-2xl p-6 shadow-md hover:shadow-xl transition-all duration-300 group"
            >
              <div className="w-14 h-14 bg-gradient-to-br from-primary-100 to-primary-200 rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
                <i className={`${value.icon} text-2xl text-primary-600`}></i>
              </div>
              <h4 className="font-bold text-gray-900 mb-2">{value.title}</h4>
              <div className="text-sm text-gray-600" {...createHTMLProps(value.description)} />
            </motion.div>
          ))}
        </motion.div>
      </div>
    </section>
  );
}
