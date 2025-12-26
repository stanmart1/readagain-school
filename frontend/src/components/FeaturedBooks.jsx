import { useState } from 'react';
import { motion } from 'framer-motion';
import { getImageUrl } from '../lib/fileService';
import { Link } from 'react-router-dom';
import { useBooks } from '../hooks';
import ProgressiveImage from './ProgressiveImage';
import api from '../lib/api';

export default function FeaturedBooks() {
  const [selectedCategory, setSelectedCategory] = useState('featured');
  const [addedToCart, setAddedToCart] = useState(new Set());

  const params = {
    page: 1,
    limit: 8,
    status: 'published',
    ...(selectedCategory === 'featured' && { is_featured: true }),
    ...(selectedCategory === 'bestsellers' && { is_bestseller: true }),
    ...(selectedCategory === 'new' && { is_new_release: true })
  };

  const { books, loading } = useBooks(params);

  const handleAddToCart = async (book, e) => {
    e.preventDefault();
    try {
      await api.post('/library', { book_id: book.id });
      setAddedToCart(prev => new Set(prev).add(book.id));
      setTimeout(() => {
        setAddedToCart(prev => {
          const newSet = new Set(prev);
          newSet.delete(book.id);
          return newSet;
        });
      }, 2000);
    } catch (error) {
      console.error('Error adding to library:', error);
      alert('Failed to add to library');
    }
  };

  const categories = [
    { id: 'featured', name: 'Featured', icon: 'ri-star-fill', color: 'from-yellow-400 to-orange-500' },
    { id: 'bestsellers', name: 'Bestsellers', icon: 'ri-fire-fill', color: 'from-red-400 to-pink-500' },
    { id: 'new', name: 'New Releases', icon: 'ri-flashlight-fill', color: 'from-blue-400 to-cyan-500' }
  ];

  const getBookImage = (book) => {
    return getImageUrl(book.cover_image_url || book.cover_image, `https://picsum.photos/seed/${book.id}/400/600`);
  };

  return (
    <section className="py-20 bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Section Header */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-center mb-16"
        >
          <h2 className="text-4xl lg:text-5xl font-bold text-gray-900 mb-4">
            Explore Our <span className="text-primary-600">Collection</span>
          </h2>
          <p className="text-lg text-gray-600 max-w-2xl mx-auto">
            Handpicked books for every reader - from timeless classics to the latest releases
          </p>
        </motion.div>

        {/* Category Pills */}
        <div className="flex justify-center gap-3 mb-12 flex-wrap">
          {categories.map((category, index) => (
            <motion.button
              key={category.id}
              initial={{ opacity: 0, scale: 0.8 }}
              whileInView={{ opacity: 1, scale: 1 }}
              viewport={{ once: true }}
              transition={{ delay: index * 0.1 }}
              onClick={() => setSelectedCategory(category.id)}
              className={`relative px-6 py-3 rounded-full font-semibold transition-all overflow-hidden ${
                selectedCategory === category.id
                  ? 'text-white shadow-lg scale-105'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              {selectedCategory === category.id && (
                <motion.div
                  layoutId="activeCategory"
                  className={`absolute inset-0 bg-gradient-to-r ${category.color}`}
                  transition={{ type: "spring", bounce: 0.2, duration: 0.6 }}
                />
              )}
              <span className="relative flex items-center gap-2">
                <i className={category.icon}></i>
                {category.name}
              </span>
            </motion.button>
          ))}
        </div>

        {/* Books Grid */}
        {loading ? (
          <div className="flex justify-center py-20">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
          </div>
        ) : books.length === 0 ? (
          <div className="text-center py-20">
            <i className="ri-book-line text-6xl text-gray-300 mb-4"></i>
            <p className="text-xl text-gray-500">No books found in this category</p>
          </div>
        ) : (
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
            {books.map((book, index) => (
              <motion.div
                key={book.id}
                initial={{ opacity: 0, y: 30 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.05 }}
              >
                <Link
                  to={`/books/${book.id}`}
                  className="group block"
                >
                  <div className="relative bg-gray-50 rounded-2xl overflow-hidden shadow-sm hover:shadow-xl transition-all duration-300">
                    
                    {/* Book Cover */}
                    <div className="aspect-[3/4] overflow-hidden bg-gray-200">
                      <ProgressiveImage
                        src={getBookImage(book)}
                        alt={book.title}
                        className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
                      />
                    </div>

                    {/* Badge */}
                    {selectedCategory === 'featured' && (
                      <div className="absolute top-3 right-3 bg-yellow-400 text-yellow-900 px-3 py-1 rounded-full text-xs font-bold shadow-lg">
                        <i className="ri-star-fill mr-1"></i>
                        Featured
                      </div>
                    )}
                    {selectedCategory === 'bestsellers' && (
                      <div className="absolute top-3 right-3 bg-red-500 text-white px-3 py-1 rounded-full text-xs font-bold shadow-lg">
                        <i className="ri-fire-fill mr-1"></i>
                        Hot
                      </div>
                    )}
                    {selectedCategory === 'new' && (
                      <div className="absolute top-3 right-3 bg-blue-500 text-white px-3 py-1 rounded-full text-xs font-bold shadow-lg">
                        <i className="ri-flashlight-fill mr-1"></i>
                        New
                      </div>
                    )}

                    {/* Book Info */}
                    <div className="p-4">
                      <h3 className="font-bold text-gray-900 mb-1 line-clamp-2 group-hover:text-primary-600 transition-colors">
                        {book.title}
                      </h3>
                      <p className="text-sm text-gray-600 mb-3 line-clamp-1">
                        {book.author?.name || 'Unknown Author'}
                      </p>

                      {/* Action Button */}
                      <button
                        onClick={(e) => handleAddToCart(book, e)}
                        disabled={addedToCart.has(book.id)}
                        className={`w-full py-2 rounded-lg font-medium transition-all ${
                          addedToCart.has(book.id)
                            ? 'bg-green-500 text-white'
                            : 'bg-primary-600 text-white hover:bg-primary-700'
                        }`}
                      >
                        {addedToCart.has(book.id) ? (
                          <span className="flex items-center justify-center gap-2">
                            <i className="ri-check-line"></i>
                            Added
                          </span>
                        ) : (
                          <span className="flex items-center justify-center gap-2">
                            <i className="ri-add-line"></i>
                            Add to Library
                          </span>
                        )}
                      </button>
                    </div>
                  </div>
                </Link>
              </motion.div>
            ))}
          </div>
        )}

        {/* View All Button */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-center mt-12"
        >
          <Link
            to="/books"
            className="inline-flex items-center gap-2 px-8 py-4 bg-primary-600 text-white rounded-full font-semibold hover:bg-primary-700 transition-all shadow-lg hover:shadow-xl"
          >
            <span>View All Books</span>
            <i className="ri-arrow-right-line"></i>
          </Link>
        </motion.div>
      </div>
    </section>
  );
}
