import { motion } from 'framer-motion';
import Header from '../components/Header';
import Footer from '../components/Footer';

export default function Privacy() {
  const sections = [
    {
      title: "1. Information We Collect",
      content: "We collect information you provide when creating an account (name, email, school information), usage data (reading progress, books accessed), and technical data (device information, IP address) to provide and improve our services."
    },
    {
      title: "2. How We Use Your Information",
      content: "We use your information to provide our services, track reading progress, personalize your experience, communicate with you, improve our platform, and ensure security."
    },
    {
      title: "3. Information Sharing",
      content: "We do not sell your personal information. We may share data with service providers who help us operate the platform, with schools for educational purposes, and when required by law."
    },
    {
      title: "4. Data Security",
      content: "We implement industry-standard security measures to protect your information, including encryption, secure servers, and regular security audits. However, no system is completely secure."
    },
    {
      title: "5. Your Rights",
      content: "You have the right to access your data, request corrections, delete your account, opt-out of communications, and export your data. Contact us to exercise these rights."
    },
    {
      title: "6. Cookies",
      content: "We use cookies and similar technologies to maintain your session, remember preferences, and analyze usage. You can control cookies through your browser settings."
    },
    {
      title: "7. Children's Privacy",
      content: "Our service is designed for educational use. We comply with applicable privacy laws regarding children's data and require parental or school consent for users under 13."
    },
    {
      title: "8. Changes to Policy",
      content: "We may update this privacy policy periodically. We will notify you of significant changes via email or platform notification."
    }
  ];

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white">
      <Header />
      
      {/* Hero Section */}
      <section className="bg-gradient-to-r from-primary-600 to-primary-700 text-white py-16 mt-20">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="text-center"
          >
            <div className="inline-block bg-white/20 backdrop-blur-sm px-4 py-2 rounded-full mb-4">
              <span className="text-sm font-semibold">Legal</span>
            </div>
            <h1 className="text-4xl md:text-5xl font-bold mb-4">
              Privacy Policy
            </h1>
            <p className="text-lg text-primary-100">
              Last updated: January 2025
            </p>
          </motion.div>
        </div>
      </section>

      {/* Content */}
      <section className="py-16">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="space-y-6"
          >
            {sections.map((section, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.1 * index }}
                className="bg-white rounded-xl shadow-md hover:shadow-lg transition-shadow p-6 border border-gray-100"
              >
                <h2 className="text-xl font-bold text-gray-900 mb-3 flex items-center gap-3">
                  <span className="w-8 h-8 bg-primary-100 text-primary-600 rounded-full flex items-center justify-center text-sm font-bold">
                    {index + 1}
                  </span>
                  {section.title.replace(/^\d+\.\s*/, '')}
                </h2>
                <p className="text-gray-600 leading-relaxed">
                  {section.content}
                </p>
              </motion.div>
            ))}
          </motion.div>

          {/* Contact CTA */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.8 }}
            className="mt-12 bg-gradient-to-r from-primary-50 to-primary-100 rounded-xl p-8 text-center border border-primary-200"
          >
            <h3 className="text-2xl font-bold text-gray-900 mb-3">Privacy Questions?</h3>
            <p className="text-gray-600 mb-6">
              We're committed to protecting your privacy. Contact us if you have any concerns.
            </p>
            <a
              href="/contact"
              className="inline-flex items-center gap-2 px-6 py-3 bg-primary-600 text-white rounded-full font-semibold hover:bg-primary-700 transition-all shadow-lg hover:shadow-xl"
            >
              <i className="ri-mail-line"></i>
              Contact Us
            </a>
          </motion.div>
        </div>
      </section>

      <Footer />
    </div>
  );
}
