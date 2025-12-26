import { motion } from 'framer-motion';

export default function EReaderPreview() {
  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.8, delay: 0.3 }}
      className="relative"
    >
      {/* E-Reader Device */}
      <div className="relative w-full max-w-md mx-auto">
        {/* Device Frame */}
        <div className="bg-gray-800 rounded-3xl p-4 shadow-2xl">
          {/* Screen */}
          <div className="bg-white rounded-2xl overflow-hidden shadow-inner">
            {/* Reader Header */}
            <div className="bg-primary-600 text-white px-4 py-3 flex items-center justify-between">
              <div className="flex items-center gap-2">
                <i className="ri-book-open-line text-xl"></i>
                <span className="font-semibold text-sm">ReadAgain Reader</span>
              </div>
              <div className="flex items-center gap-3">
                <i className="ri-bookmark-line text-lg"></i>
                <i className="ri-settings-3-line text-lg"></i>
              </div>
            </div>

            {/* Book Content */}
            <div className="p-6 h-96 overflow-hidden">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">
                Chapter 1: The Beginning
              </h2>
              <div className="space-y-4 text-gray-700 leading-relaxed">
                <p className="text-justify">
                  In the heart of the ancient library, where dust motes danced in shafts of golden sunlight, 
                  young Emma discovered a book that would change her life forever. The leather-bound tome 
                  seemed to whisper secrets of forgotten worlds and untold adventures.
                </p>
                <p className="text-justify">
                  As she carefully opened the first page, the scent of aged paper filled her senses, 
                  and the words began to glow with an ethereal light. This was no ordinary book—it was 
                  a gateway to knowledge, imagination, and endless possibilities.
                </p>
                <p className="text-justify opacity-50">
                  The library keeper had warned her about the power of books, but Emma never imagined...
                </p>
              </div>

              {/* Progress Bar */}
              <div className="absolute bottom-6 left-6 right-6">
                <div className="flex items-center justify-between text-xs text-gray-500 mb-2">
                  <span>Page 1 of 245</span>
                  <span>12% Complete</span>
                </div>
                <div className="h-1.5 bg-gray-200 rounded-full overflow-hidden">
                  <motion.div
                    initial={{ width: 0 }}
                    animate={{ width: '12%' }}
                    transition={{ duration: 1, delay: 0.5 }}
                    className="h-full bg-gradient-to-r from-primary-500 to-primary-600"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Floating Action Buttons */}
        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.8 }}
          className="absolute -right-4 top-1/4 flex flex-col gap-3"
        >
          <div className="w-12 h-12 bg-white rounded-full shadow-lg flex items-center justify-center text-primary-600 hover:bg-primary-50 transition-colors cursor-pointer">
            <i className="ri-font-size text-xl"></i>
          </div>
          <div className="w-12 h-12 bg-white rounded-full shadow-lg flex items-center justify-center text-primary-600 hover:bg-primary-50 transition-colors cursor-pointer">
            <i className="ri-palette-line text-xl"></i>
          </div>
        </motion.div>

        {/* Glow Effect */}
        <div className="absolute inset-0 bg-primary-400/20 rounded-3xl blur-3xl -z-10"></div>
      </div>
    </motion.div>
  );
}
