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
    <section className="py-16 bg-gradient-to-b from-white to-gray-50 relative overflow-hidden">
      {/* Decorative background elements */}
      <div className="absolute top-20 right-0 w-72 h-72 bg-primary-100 rounded-full blur-3xl opacity-30"></div>
      <div className="absolute bottom-20 left-0 w-96 h-96 bg-yellow-100 rounded-full blur-3xl opacity-20"></div>
      
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        
        {/* Section Header */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-center mb-8"
        >
          <span className="inline-block px-4 py-2 bg-primary-100 text-primary-700 rounded-full text-sm font-semibold mb-4">
            Who We Are
          </span>
          <h2 className="text-5xl lg:text-6xl font-bold text-gray-900 mb-4 leading-tight">
            {aboutContent.hero?.title || 'About ReadAgain'}
          </h2>
          <div className="text-2xl text-gray-600 font-light max-w-3xl mx-auto leading-relaxed"
               {...createHTMLProps(aboutContent.hero?.subtitle || 'Empowering The Mind Through Reading')}
          />
        </motion.div>

        {/* Mission Text */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6, delay: 0.2 }}
          className="text-center mb-12"
        >
          <div className="text-xl text-gray-600 leading-relaxed max-w-3xl mx-auto mb-8"
               {...createHTMLProps(aboutContent.mission?.description || defaultContent.mission.description)}
          />
          
          <Link
            to="/about"
            className="inline-flex items-center gap-2 px-8 py-4 bg-primary-600 text-white rounded-full font-semibold hover:bg-primary-700 transition-all shadow-lg hover:shadow-xl text-lg"
          >
            <span>Learn More About Us</span>
            <i className="ri-arrow-right-line"></i>
          </Link>
        </motion.div>
      </div>
    </section>
  );
}
