import { useState, useRef, useEffect } from 'react';
import { motion } from 'framer-motion';
import { getImageUrl } from '../lib/fileService';
import { Link } from 'react-router-dom';
import { useBooks } from '../hooks';
import ProgressiveImage from './ProgressiveImage';
import api from '../lib/api';

export default function FeaturedBooks() {
  const [addedToCart, setAddedToCart] = useState(new Set());
  const desktopCarouselRef = useRef(null);
  const mobileCarouselRef = useRef(null);
  const [isAutoScrolling, setIsAutoScrolling] = useState(true);
  const [isUserScrolling, setIsUserScrolling] = useState(false);
  const scrollTimeoutRef = useRef(null);

  const params = {
    page: 1,
    limit: 8,
    status: 'published'
  };

  const { books, loading } = useBooks(params);
  
  const infiniteBooks = books.length > 2 ? [...books, ...books, ...books] : books;

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

  const getBookImage = (book) => {
    return getImageUrl(book.cover_image_url || book.cover_image, `https://picsum.photos/seed/${book.id}/400/600`);
  };

  const scroll = (direction) => {
    const carousel = window.innerWidth >= 640 ? desktopCarouselRef.current : mobileCarouselRef.current;
    if (carousel) {
      const scrollAmount = carousel.offsetWidth * 0.8;
      carousel.scrollBy({
        left: direction === 'left' ? -scrollAmount : scrollAmount,
        behavior: 'smooth'
      });
    }
  };

  useEffect(() => {
    if (books.length > 0) {
      const cardWidth = 296;
      if (desktopCarouselRef.current) {
        desktopCarouselRef.current.scrollLeft = cardWidth * books.length;
      }
      if (mobileCarouselRef.current) {
        mobileCarouselRef.current.scrollLeft = cardWidth * books.length;
      }
    }
  }, [books.length]);

  useEffect(() => {
    if (!isAutoScrolling || books.length === 0) return;

    const interval = setInterval(() => {
      const carousel = window.innerWidth >= 640 ? desktopCarouselRef.current : mobileCarouselRef.current;
      if (carousel) {
        const cardWidth = 296;
        carousel.scrollBy({ left: cardWidth, behavior: 'smooth' });
      }
    }, 3000);

    return () => clearInterval(interval);
  }, [isAutoScrolling, books.length]);

  useEffect(() => {
    const handleScroll = () => {
      const carousel = window.innerWidth >= 640 ? desktopCarouselRef.current : mobileCarouselRef.current;
      if (!carousel || books.length === 0) return;

      const cardWidth = 296;
      const scrollLeft = carousel.scrollLeft;
      const maxScroll = cardWidth * books.length * 2;

      if (scrollLeft <= cardWidth) {
        carousel.scrollLeft = cardWidth * books.length + scrollLeft;
      } else if (scrollLeft >= maxScroll) {
        carousel.scrollLeft = cardWidth * books.length + (scrollLeft - maxScroll);
      }
    };

    const carousel = window.innerWidth >= 640 ? desktopCarouselRef.current : mobileCarouselRef.current;
    if (carousel) {
      carousel.addEventListener('scroll', handleScroll);
      return () => carousel.removeEventListener('scroll', handleScroll);
    }
  }, [books.length]);

  const handleUserScroll = () => {
    setIsAutoScrolling(false);
    setIsUserScrolling(true);
    
    if (scrollTimeoutRef.current) {
      clearTimeout(scrollTimeoutRef.current);
    }
    
    scrollTimeoutRef.current = setTimeout(() => {
      setIsUserScrolling(false);
      setIsAutoScrolling(true);
    }, 5000);
  };

  return (
    <section className="py-20 bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
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
          <>
            <div className="hidden sm:block relative">
              <button
                onClick={() => scroll('left')}
                className="absolute left-0 top-1/2 -translate-y-1/2 z-10 bg-white/90 backdrop-blur-sm p-3 rounded-full shadow-lg hover:bg-white transition-all"
              >
                <i className="ri-arrow-left-s-line text-2xl text-gray-700"></i>
              </button>
              
              <div
                ref={desktopCarouselRef}
                onScroll={handleUserScroll}
                className="flex gap-6 overflow-x-auto scrollbar-hide scroll-smooth px-2"
                style={{ scrollbarWidth: 'none', msOverflowStyle: 'none' }}
              >
                {infiniteBooks.map((book, index) => (
                  <motion.div
                    key={`${book.id}-${index}`}
                    initial={{ opacity: 0, scale: 0.9 }}
                    whileInView={{ opacity: 1, scale: 1 }}
                    viewport={{ once: true }}
                    transition={{ delay: (index % books.length) * 0.05 }}
                    className="flex-shrink-0 w-[280px]"
                  >
                    <Link to={`/books/${book.id}`} className="group block">
                      <div className="relative bg-gray-50 rounded-2xl overflow-hidden shadow-sm hover:shadow-xl transition-all duration-300">
                        
                        <div className="aspect-[2/3] overflow-hidden bg-gray-200">
                          <ProgressiveImage
                            src={getBookImage(book)}
                            alt={book.title}
                            className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
                          />
                        </div>

                        <div className="p-4">
                          <h3 className="font-bold text-gray-900 mb-1 line-clamp-2 group-hover:text-primary-600 transition-colors">
                            {book.title}
                          </h3>
                          <p className="text-sm text-gray-600 mb-3 line-clamp-1">
                            {book.author?.name || 'Unknown Author'}
                          </p>

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

              <button
                onClick={() => scroll('right')}
                className="absolute right-0 top-1/2 -translate-y-1/2 z-10 bg-white/90 backdrop-blur-sm p-3 rounded-full shadow-lg hover:bg-white transition-all"
              >
                <i className="ri-arrow-right-s-line text-2xl text-gray-700"></i>
              </button>
            </div>

            <div className="sm:hidden">
              <div
                ref={mobileCarouselRef}
                onScroll={handleUserScroll}
                className="flex gap-4 overflow-x-auto scrollbar-hide scroll-smooth px-2"
                style={{ scrollbarWidth: 'none', msOverflowStyle: 'none' }}
              >
                {infiniteBooks.map((book, index) => (
                  <div key={`${book.id}-${index}`} className="flex-shrink-0 w-[160px]">
                    <Link to={`/books/${book.id}`} className="group block">
                      <div className="relative bg-gray-50 rounded-2xl overflow-hidden shadow-sm">
                        <div className="aspect-[2/3] overflow-hidden bg-gray-200">
                          <ProgressiveImage
                            src={getBookImage(book)}
                            alt={book.title}
                            className="w-full h-full object-cover"
                          />
                        </div>
                        <div className="p-3">
                          <h3 className="font-bold text-sm text-gray-900 line-clamp-2">
                            {book.title}
                          </h3>
                        </div>
                      </div>
                    </Link>
                  </div>
                ))}
              </div>
            </div>
          </>
        )}

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
