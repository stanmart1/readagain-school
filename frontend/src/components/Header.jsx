import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import { useAuth } from '../hooks';

export default function Header() {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [isProfileOpen, setIsProfileOpen] = useState(false);
  const { isAuthenticated, getUser, logout } = useAuth();
  const user = getUser();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  return (
    <>
      <motion.header 
        initial={{ y: -100 }}
        animate={{ y: 0 }}
        className="bg-white shadow-sm border-b border-gray-200 fixed top-0 left-0 right-0 z-50"
      >
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            {/* Logo */}
            <Link to="/" className="flex items-center space-x-2">
              <div className="w-8 h-8 bg-gradient-to-r from-primary-500 to-primary-600 rounded-lg flex items-center justify-center">
                <i className="ri-book-line text-white text-lg"></i>
              </div>
              <span className="text-xl font-bold font-pacifico text-gray-900">
                ReadAgain
              </span>
            </Link>

            {/* Desktop Navigation */}
            <nav className="hidden md:flex items-center space-x-1">
              {['Home', 'Books', 'Blog', 'FAQ', 'About', 'Contact'].map((item) => (
                <Link
                  key={item}
                  to={item === 'Home' ? '/' : `/${item.toLowerCase()}`}
                  className="relative px-3 py-2 text-gray-600 hover:text-primary-600 hover:bg-primary-50 rounded-lg transition-all duration-200 font-medium group"
                >
                  <span>{item}</span>
                  <div className="absolute bottom-0 left-1/2 w-0 h-0.5 bg-primary-600 transition-all duration-200 group-hover:w-4/5 transform -translate-x-1/2"></div>
                </Link>
              ))}
            </nav>

            {/* Right side */}
            <div className="flex items-center space-x-4">
              
              {isAuthenticated() ? (
                <div className="relative">
                  <button 
                    onClick={() => setIsProfileOpen(!isProfileOpen)}
                    className="flex items-center gap-2 px-3 py-2 hover:bg-gray-100 rounded-lg transition-all"
                  >
                    <div className="w-8 h-8 bg-gradient-to-r from-primary-600 to-primary-700 rounded-full flex items-center justify-center text-white font-semibold">
                      {user?.first_name?.[0] || 'U'}
                    </div>
                    <span className="hidden md:block font-medium text-gray-700">{user?.first_name}</span>
                    <i className="ri-arrow-down-s-line text-gray-600"></i>
                  </button>
                  {isProfileOpen && (
                    <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg py-2 z-50">
                      <div className="px-4 py-2 border-b">
                        <p className="font-semibold text-gray-900">{user?.first_name} {user?.last_name}</p>
                        <p className="text-sm text-gray-500">{user?.email}</p>
                      </div>
                      {(user?.role?.name === 'platform_admin' || user?.role?.name === 'school_admin') && (
                        <Link to="/admin" onClick={() => setIsProfileOpen(false)} className="block px-4 py-2 text-primary-600 hover:bg-primary-50">
                          <i className="ri-admin-line mr-2"></i>Admin Dashboard
                        </Link>
                      )}
                      <Link to="/dashboard" onClick={() => setIsProfileOpen(false)} className="block px-4 py-2 text-gray-700 hover:bg-gray-100">
                        <i className="ri-dashboard-line mr-2"></i>Dashboard
                      </Link>
                      <Link to="/dashboard/library" onClick={() => setIsProfileOpen(false)} className="block px-4 py-2 text-gray-700 hover:bg-gray-100">
                        <i className="ri-book-line mr-2"></i>My Library
                      </Link>
                      <button
                        onClick={handleLogout}
                        className="w-full text-left px-4 py-2 text-red-600 hover:bg-red-50"
                      >
                        <i className="ri-logout-box-line mr-2"></i>Logout
                      </button>
                    </div>
                  )}
                </div>
              ) : (
                <>
                  <Link to="/login" className="hidden md:block px-6 py-2 text-primary-600 hover:bg-primary-50 rounded-lg transition-all font-medium">
                    Login
                  </Link>
                  <Link to="/signup" className="hidden md:block px-6 py-2 bg-gradient-to-r from-primary-600 to-primary-700 text-white rounded-lg hover:from-primary-700 hover:to-primary-800 transition-all font-medium">
                    Sign Up
                  </Link>
                </>
              )}
              
              {/* Mobile menu button */}
              <button
                onClick={() => setIsMenuOpen(!isMenuOpen)}
                className="md:hidden p-2 text-gray-700 hover:bg-gray-100 rounded-lg"
              >
                <i className={`ri-${isMenuOpen ? 'close' : 'menu'}-line text-2xl`}></i>
              </button>
            </div>
          </div>
        </div>
      </motion.header>

      {/* Mobile Sidebar Overlay */}
      {isMenuOpen && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          onClick={() => setIsMenuOpen(false)}
          className="fixed inset-0 bg-black bg-opacity-50 z-[60] md:hidden"
        />
      )}

      {/* Mobile Sidebar */}
      <motion.div
        initial={{ x: '100%' }}
        animate={{ x: isMenuOpen ? '0%' : '100%' }}
        transition={{ type: 'spring', stiffness: 300, damping: 30 }}
        className="fixed right-0 top-0 h-full w-64 bg-white z-[70] md:hidden shadow-lg overflow-y-auto overflow-x-hidden"
      >
        {/* Sidebar Header */}
        <div className="flex items-center justify-between p-4 border-b border-gray-200">
          <span className="text-lg font-bold text-gray-900">Menu</span>
          <button
            onClick={() => setIsMenuOpen(false)}
            className="p-2 text-gray-700 hover:bg-gray-100 rounded-lg"
          >
            <i className="ri-close-line text-2xl"></i>
          </button>
        </div>

        {/* Navigation Links */}
        <div className="p-4 space-y-2">
          {['Home', 'Books', 'Blog', 'FAQ', 'About', 'Contact'].map((item) => (
            <Link
              key={item}
              to={item === 'Home' ? '/' : `/${item.toLowerCase()}`}
              className="block px-4 py-3 text-gray-700 hover:bg-primary-50 hover:text-primary-600 rounded-lg transition-all font-medium"
              onClick={() => setIsMenuOpen(false)}
            >
              {item}
            </Link>
          ))}
        </div>

        {/* Auth Links */}
        {!isAuthenticated() && (
          <div className="p-4 border-t border-gray-200 space-y-3">
            <Link
              to="/login"
              className="block w-full text-center px-4 py-2 text-primary-600 border border-primary-600 rounded-lg font-medium hover:bg-primary-50 transition-all"
              onClick={() => setIsMenuOpen(false)}
            >
              Login
            </Link>
            <Link
              to="/signup"
              className="block w-full text-center px-4 py-2 bg-gradient-to-r from-primary-600 to-primary-700 text-white rounded-lg font-medium hover:from-primary-700 hover:to-primary-800 transition-all"
              onClick={() => setIsMenuOpen(false)}
            >
              Sign Up
            </Link>
          </div>
        )}
      </motion.div>
    </>
  );
}
