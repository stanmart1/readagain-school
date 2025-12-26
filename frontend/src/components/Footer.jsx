import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';

export default function Footer() {
  return (
    <footer className="bg-gray-900 text-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {/* Brand Section */}
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            className="lg:col-span-1"
          >
            <div className="flex items-center space-x-2 mb-2">
              <div className="w-8 h-8 bg-gradient-to-r from-primary-500 to-primary-600 rounded-lg flex items-center justify-center">
                <i className="ri-book-line text-white text-lg"></i>
              </div>
              <span className="text-xl font-bold font-pacifico">ReadAgain</span>
            </div>
            <p className="text-gray-400 mb-3 leading-relaxed text-sm">
              Empowering schools and students with digital reading resources.
            </p>
            <div className="flex space-x-3">
              {['facebook', 'twitter', 'instagram', 'linkedin'].map((social) => (
                <motion.a
                  key={social}
                  whileHover={{ scale: 1.1, y: -2 }}
                  href="#"
                  className="w-8 h-8 bg-gray-800 rounded-full flex items-center justify-center hover:bg-primary-600 transition-colors"
                >
                  <i className={`ri-${social}-fill text-sm`}></i>
                </motion.a>
              ))}
            </div>
          </motion.div>
          
          <div className="lg:col-span-2"></div>
          
          {/* Links */}
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.2 }}
            className="lg:col-span-1"
          >
            <h3 className="text-base font-semibold mb-2">Legal & Info</h3>
            <ul className="space-y-1 text-sm">
              <li>
                <Link to="/terms" className="text-gray-400 hover:text-white transition-colors">
                  Terms of Service
                </Link>
              </li>
              <li>
                <Link to="/privacy" className="text-gray-400 hover:text-white transition-colors">
                  Privacy Policy
                </Link>
              </li>
              <li>
                <Link to="/about" className="text-gray-400 hover:text-white transition-colors">
                  About Us
                </Link>
              </li>
              <li>
                <Link to="/contact" className="text-gray-400 hover:text-white transition-colors">
                  Contact Us
                </Link>
              </li>
            </ul>
          </motion.div>
        </div>
        
        {/* Bottom Bar */}
        <div className="border-t border-gray-800 mt-4 pt-3 flex flex-col md:flex-row justify-between items-center gap-2">
          <p className="text-gray-400 text-xs">
            © 2025 ReadAgain. All rights reserved.
          </p>
          <div className="flex gap-4 text-xs">
            <Link to="/privacy" className="text-gray-400 hover:text-white transition-colors">
              Privacy Policy
            </Link>
            <Link to="/terms" className="text-gray-400 hover:text-white transition-colors">
              Terms of Service
            </Link>
          </div>
        </div>
        
        {/* Created with Love */}
        <div className="border-t border-gray-800 mt-2 pt-2 flex justify-center">
          <p className="text-gray-400 text-xs flex items-center space-x-1">
            <span>Created with</span>
            <motion.i
              animate={{ scale: [1, 1.2, 1] }}
              transition={{ duration: 1, repeat: Infinity }}
              className="ri-heart-fill text-red-500"
            />
            <span>by</span>
            <a
              href="tel:+2347062750516"
              className="text-primary-400 hover:text-primary-300 transition-colors font-medium"
            >
              ScaleITPro
            </a>
          </p>
        </div>
      </div>
    </footer>
  );
}
