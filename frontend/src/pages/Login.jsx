import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [rememberMe, setRememberMe] = useState(false);
  const { login, loading, error } = useAuth();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const redirectPath = searchParams.get('redirect') || localStorage.getItem('redirectAfterLogin');
  const sessionExpired = searchParams.get('reason') === 'session_expired';

  useEffect(() => {
    return () => {
      localStorage.removeItem('redirectAfterLogin');
    };
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const result = await login(email, password);
    if (result) {
      const targetPath = redirectPath || '/dashboard';
      localStorage.removeItem('redirectAfterLogin');
      navigate(targetPath);
    }
  };

  return (
    <div className="min-h-screen flex">
      {/* Left Side - Branding & Benefits */}
      <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-primary-800 to-primary-900 p-12 flex-col justify-center relative overflow-hidden">
        {/* Background Pattern */}
        <div className="absolute inset-0 opacity-10">
          <div className="absolute top-20 left-20 w-64 h-64 bg-white rounded-full blur-3xl"></div>
          <div className="absolute bottom-20 right-20 w-96 h-96 bg-white rounded-full blur-3xl"></div>
        </div>

        <div className="relative z-10">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
          >
            <h1 className="text-5xl font-bold text-white mb-6">
              Welcome Back to<br />ReadAgain
            </h1>
            <p className="text-xl text-primary-100 mb-12">
              Continue your reading journey and explore thousands of books.
            </p>

            {/* Benefits */}
            <div className="space-y-6">
              {[
                { icon: 'ri-book-open-line', text: 'Access your entire digital library' },
                { icon: 'ri-bookmark-line', text: 'Pick up right where you left off' },
                { icon: 'ri-trophy-line', text: 'Track your reading progress' },
                { icon: 'ri-smartphone-line', text: 'Read anywhere, anytime' }
              ].map((benefit, index) => (
                <motion.div
                  key={index}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: 0.3 + index * 0.1 }}
                  className="flex items-center gap-4 text-white"
                >
                  <div className="w-12 h-12 bg-white/20 backdrop-blur-sm rounded-xl flex items-center justify-center flex-shrink-0">
                    <i className={`${benefit.icon} text-2xl`}></i>
                  </div>
                  <span className="text-lg">{benefit.text}</span>
                </motion.div>
              ))}
            </div>
          </motion.div>
        </div>
      </div>

      {/* Right Side - Login Form */}
      <div className="w-full lg:w-1/2 flex items-center justify-center p-8 bg-gray-50">
        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.4 }}
          className="w-full max-w-md"
        >
          {/* Mobile Back Button */}
          <div className="lg:hidden mb-6">
            <Link to="/" className="text-primary-600 hover:text-primary-700 inline-flex items-center">
              <i className="ri-arrow-left-line mr-2"></i>
              Back to Home
            </Link>
          </div>

          {/* Form Card */}
          <div className="bg-white rounded-2xl shadow-xl p-8 border border-gray-100">
            <div className="mb-8">
              <Link to="/" className="inline-flex items-center text-gray-600 hover:text-primary-600 transition-colors mb-6">
                <i className="ri-arrow-left-line mr-2"></i>
                Back to Home
              </Link>
              <h2 className="text-3xl font-bold text-gray-900 mb-2">Login</h2>
              <p className="text-gray-600">Enter your credentials to access your account</p>
            </div>

            {sessionExpired && (
              <div className="mb-6 p-4 bg-yellow-50 border border-yellow-200 text-yellow-800 rounded-lg flex items-start gap-3">
                <i className="ri-time-line text-xl flex-shrink-0 mt-0.5"></i>
                <span className="text-sm">Your session has expired. Please login again.</span>
              </div>
            )}

            {error && (
              <div className="mb-6 p-4 bg-red-50 border border-red-200 text-red-700 rounded-lg flex items-start gap-3">
                <i className="ri-error-warning-line text-xl flex-shrink-0 mt-0.5"></i>
                <span className="text-sm">{error}</span>
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-5">
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Email Address
                </label>
                <div className="relative">
                  <i className="ri-mail-line absolute left-4 top-1/2 -translate-y-1/2 text-gray-400"></i>
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    className="w-full pl-12 pr-4 py-3 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all"
                    placeholder="your@email.com"
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Password
                </label>
                <div className="relative">
                  <i className="ri-lock-line absolute left-4 top-1/2 -translate-y-1/2 text-gray-400"></i>
                  <input
                    type={showPassword ? 'text' : 'password'}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    className="w-full pl-12 pr-12 py-3 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all"
                    placeholder="••••••••"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors"
                  >
                    <i className={showPassword ? 'ri-eye-off-line text-xl' : 'ri-eye-line text-xl'}></i>
                  </button>
                </div>
              </div>

              <div className="flex items-center justify-between">
                <label className="flex items-center cursor-pointer group">
                  <div className="relative">
                    <input
                      type="checkbox"
                      checked={rememberMe}
                      onChange={(e) => setRememberMe(e.target.checked)}
                      className="sr-only"
                    />
                    <div className={`w-5 h-5 border-2 rounded transition-all ${
                      rememberMe 
                        ? 'bg-primary-600 border-primary-600' 
                        : 'border-gray-300 group-hover:border-primary-400'
                    }`}>
                      {rememberMe && (
                        <i className="ri-check-line text-white text-sm absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2"></i>
                      )}
                    </div>
                  </div>
                  <span className="ml-2 text-sm text-gray-700">Remember me</span>
                </label>
                <Link to="/forgot-password" className="text-sm text-primary-600 hover:text-primary-700 font-medium">
                  Forgot password?
                </Link>
              </div>

              <button
                type="submit"
                disabled={loading}
                className="w-full bg-gradient-to-r from-primary-600 to-primary-700 text-white py-3.5 rounded-xl font-semibold hover:from-primary-700 hover:to-primary-800 transition-all disabled:opacity-50 shadow-lg hover:shadow-xl flex items-center justify-center gap-2"
              >
                {loading ? (
                  <>
                    <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                    Logging in...
                  </>
                ) : (
                  <>
                    Login
                    <i className="ri-arrow-right-line"></i>
                  </>
                )}
              </button>
            </form>

            <div className="mt-8 text-center">
              <p className="text-gray-600">
                Don't have an account?{' '}
                <Link to="/signup" className="text-primary-600 hover:text-primary-700 font-semibold">
                  Sign up for free
                </Link>
              </p>
            </div>
          </div>
        </motion.div>
      </div>
    </div>
  );
}
