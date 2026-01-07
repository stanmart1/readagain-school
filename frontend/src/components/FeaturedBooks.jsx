import { useState, useRef, useEffect } from 'react';
import { motion } from 'framer-motion';
import { useBooks } from '../hooks';
import BookCard from './BookCard';

export default function FeaturedBooks() {
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
                    <BookCard book={book} />
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
                  <div key={`${book.id}-${index}`} className="flex-shrink-0 w-[200px]">
                    <BookCard book={book} />
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
