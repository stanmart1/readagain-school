import { useState } from 'react';
import { Link } from 'react-router-dom';
import { getImageUrl } from '../lib/fileService';
import api from '../lib/api';

export default function BookCard({ book }) {
  const [isHovered, setIsHovered] = useState(false);
  const [addingToCart, setAddingToCart] = useState(false);

  const getAuthorName = () => {
    if (book.author_name) return book.author_name;
    if (book.author?.business_name) return book.author.business_name;
    if (book.author?.user) {
      const { first_name, last_name } = book.author.user;
      return `${first_name || ''} ${last_name || ''}`.trim() || 'Unknown Author';
    }
    return 'Unknown Author';
  };

  const getCategoryName = () => {
    if (book.category_name) return book.category_name;
    if (book.category?.name) return book.category.name;
    return null;
  };
  
  const displayAuthor = getAuthorName();
  const displayCategory = getCategoryName();
  const displayCover = getImageUrl(book.cover_image_url || book.cover_image);
  const displayOriginalPrice = book.original_price || book.originalPrice;
  const displayRating = book.rating || 0;
  const displayReviewCount = book.review_count || 0;

  const handleAddToCart = async (e) => {
    e.preventDefault();
    try {
      setAddingToCart(true);
      await api.post('/library', { book_id: book.id });
      alert('Book added to your library!');
    } catch (error) {
      console.error('Error adding to library:', error);
      alert('Failed to add to library');
    } finally {
      setAddingToCart(false);
    }
  };

  const renderStars = (rating) => {
    const stars = [];
    for (let i = 1; i <= 5; i++) {
      stars.push(
        <i
          key={i}
          className={`ri-star-${i <= rating ? 'fill' : 'line'} text-yellow-400`}
        ></i>
      );
    }
    return stars;
  };

  return (
    <div 
      className="group cursor-pointer h-full"
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <div className="bg-white rounded-2xl shadow-lg overflow-hidden hover:shadow-2xl transition-all duration-300 transform hover:-translate-y-2 h-full flex flex-col">
        <div className="relative flex-shrink-0">
          <img 
            src={displayCover || ''} 
            alt={book.title}
            className="w-full h-72 object-cover"
          />
          
          {/* Book Type Badge */}
          <div className="absolute top-4 right-4 bg-green-800 text-white rounded-full px-3 py-1">
            <span className="text-sm font-medium">
              {book.format === 'physical' ? 'Physical' : book.format === 'both' ? 'Hybrid' : 'Ebook'}
            </span>
          </div>

          {/* Discount Badge */}
          {displayOriginalPrice && displayOriginalPrice > book.price && (
            <div className="absolute top-4 left-4 bg-blue-800 text-white text-xs px-2 py-1 rounded-full font-medium">
              {Math.round(((displayOriginalPrice - book.price) / displayOriginalPrice) * 100)}% OFF
            </div>
          )}

          {/* Category Badge */}
          {displayCategory && (
            <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/60 to-transparent p-4">
              <span className="bg-blue-600 text-white px-2 py-1 rounded-full text-xs font-medium">
                {displayCategory}
              </span>
            </div>
          )}

          {/* Hover Actions */}
          <div className="absolute inset-0 bg-black/50 flex items-center justify-center space-x-2 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
            <Link 
              to={`/books/${book.id}`}
              className="p-3 bg-white rounded-full hover:bg-gray-100 transition-colors"
              title="View Details"
            >
              <i className="ri-eye-line text-gray-600 text-lg"></i>
            </Link>
            <button 
              onClick={handleAddToCart}
              disabled={addingToCart}
              className="p-3 bg-blue-600 text-white rounded-full hover:bg-blue-700 transition-colors disabled:opacity-50"
              title="Add to Library"
            >
              <i className={`${addingToCart ? 'ri-loader-4-line animate-spin' : 'ri-shopping-cart-line'} text-lg`}></i>
            </button>
          </div>
        </div>
        
        <div className="p-6 flex-1 flex flex-col">
          <Link to={`/books/${book.id}`}>
            <h3 className="text-lg font-semibold text-gray-900 mb-2 group-hover:text-blue-600 transition-colors line-clamp-2 h-14">
              {book.title}
            </h3>
          </Link>
          <p className="text-gray-600 mb-4 h-6">
            by {displayAuthor}
          </p>
          
          {/* Rating - Fixed Height */}
          <div className="flex items-center space-x-1 mb-4 h-5">
            {displayRating > 0 ? (
              <>
                <div className="flex space-x-1">
                  {renderStars(displayRating)}
                </div>
                <span className="text-sm text-gray-500">({displayReviewCount})</span>
              </>
            ) : (
              <span className="text-sm text-gray-400">No ratings yet</span>
            )}
          </div>

          {/* Price and Actions - Fixed at Bottom */}
          <div className="mt-auto">
            <div className="flex items-center justify-between mb-4 h-8">
              <div className="flex items-center space-x-2">
                <span className="text-2xl font-bold text-blue-600">₦{(book.price || 0).toLocaleString()}</span>
                {displayOriginalPrice && displayOriginalPrice > book.price && (
                  <span className="text-sm text-gray-500 line-through">₦{displayOriginalPrice.toLocaleString()}</span>
                )}
              </div>
            </div>
            
            <button 
              onClick={handleAddToCart}
              disabled={addingToCart}
              className="w-full bg-blue-600 text-white px-4 py-2 rounded-full hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed font-medium"
            >
              {addingToCart ? 'Adding...' : 'Add to Library'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
