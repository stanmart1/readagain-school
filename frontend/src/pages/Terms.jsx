import { motion } from 'framer-motion';
import Header from '../components/Header';
import Footer from '../components/Footer';

export default function Terms() {
  const sections = [
    {
      title: "1. Acceptance of Terms",
      content: "By accessing and using ReadAgain, you accept and agree to be bound by these Terms of Service. If you do not agree to these terms, please do not use our services."
    },
    {
      title: "2. User Accounts",
      content: "To access our platform, you must create an account. You agree to provide accurate information, maintain account security, notify us of unauthorized access, and be responsible for all activities under your account."
    },
    {
      title: "3. Use of Service",
      content: "You may use ReadAgain for lawful purposes only. You agree not to misuse the service, attempt unauthorized access, distribute malware, or violate any applicable laws."
    },
    {
      title: "4. Content and Copyright",
      content: "All books and content on ReadAgain are protected by copyright. You may not copy, distribute, or modify any content without permission. Your use is limited to personal, educational purposes."
    },
    {
      title: "5. Privacy",
      content: "Your privacy is important to us. Please review our Privacy Policy to understand how we collect, use, and protect your information."
    },
    {
      title: "6. Termination",
      content: "We reserve the right to suspend or terminate your account if you violate these terms. You may also terminate your account at any time by contacting us."
    },
    {
      title: "7. Changes to Terms",
      content: "We may update these terms from time to time. Continued use of the service after changes constitutes acceptance of the new terms."
    },
    {
      title: "8. Contact",
      content: "If you have questions about these terms, please contact us through our support channels."
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
              Terms of Service
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
            <h3 className="text-2xl font-bold text-gray-900 mb-3">Have Questions?</h3>
            <p className="text-gray-600 mb-6">
              If you need clarification on any of these terms, we're here to help.
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
