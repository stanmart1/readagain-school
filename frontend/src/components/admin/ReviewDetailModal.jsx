import { getImageUrl } from '../../lib/fileService';

const ReviewDetailModal = ({ review, onClose }) => {
  if (!review) return null;

  const getBookImage = (review) => {
    return getImageUrl(review.book?.cover_image);
  };

  const getBookTitle = (review) => {
    return review.book?.title || 'Unknown Book';
  };

  const getBookAuthor = (review) => {
    if (review.book?.author?.business_name) return review.book.author.business_name;
    if (review.book?.author?.user) {
      const { first_name, last_name } = review.book.author.user;
      return `${first_name || ''} ${last_name || ''}`.trim();
    }
    return null;
  };

  const getUserName = (review) => {
    if (review.user) {
      return `${review.user.first_name || ''} ${review.user.last_name || ''}`.trim();
    }
    return 'Unknown User';
  };

  const getUserEmail = (review) => {
    return review.user?.email || '';
  };

  const getClassLevel = (review) => {
    return review.user?.class_level || null;
  };

  const getSchoolName = (review) => {
    return review.user?.school_name || null;
  };

  const getRatingStars = (rating) => {
    return (
      <div className="flex items-center">
        {[1, 2, 3, 4, 5].map((star) => (
          <i
            key={star}
            className={`ri-star-fill text-sm ${
              star <= rating ? 'text-yellow-400' : 'text-gray-300'
            }`}
          />
        ))}
        <span className="ml-1 text-sm text-gray-600">({rating})</span>
      </div>
    );
  };

  const getStatusBadge = (status, isFeatured) => {
    if (isFeatured) {
      return (
        <span className="px-2 py-1 text-xs font-semibold rounded-full bg-purple-100 text-purple-800 flex items-center">
          <i className="ri-star-fill mr-1"></i>Featured
        </span>
      );
    }
    
    const badges = {
      pending: 'bg-yellow-100 text-yellow-800',
      approved: 'bg-green-100 text-green-800',
      rejected: 'bg-red-100 text-red-800'
    };
    
    return (
      <span className={`px-2 py-1 text-xs font-semibold rounded-full ${badges[status]}`}>
        {status.charAt(0).toUpperCase() + status.slice(1)}
      </span>
    );
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto p-6">
        <div className="flex items-center justify-between mb-6">
          <h3 className="text-xl font-bold text-gray-900">Review Details</h3>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <i className="ri-close-line text-2xl"></i>
          </button>
        </div>
        
        <div className="space-y-4">
          <div className="flex items-start gap-4">
            <img
              src={getBookImage(review)}
              alt={getBookTitle(review)}
              className="w-24 h-32 object-cover rounded"
            />
            <div className="flex-1">
              <h4 className="text-lg font-semibold text-gray-900">{getBookTitle(review)}</h4>
              {getBookAuthor(review) && (
                <p className="text-gray-600">by {getBookAuthor(review)}</p>
              )}
              <div className="mt-2">
                {getRatingStars(review.rating)}
              </div>
              <div className="mt-2">
                {getStatusBadge(review.status, review.is_featured)}
              </div>
            </div>
          </div>

          {review.comment && (
            <div>
              <h5 className="font-semibold text-gray-900 mb-1">Review</h5>
              <p className="text-gray-700 whitespace-pre-wrap">{review.comment}</p>
            </div>
          )}

          <div className="border-t pt-4">
            <h5 className="font-semibold text-gray-900 mb-2">Reviewer Information</h5>
            <p className="text-gray-700">
              {getUserName(review)}
            </p>
            <p className="text-gray-600 text-sm">{getUserEmail(review)}</p>
            {getClassLevel(review) && (
              <p className="text-gray-600 text-sm">Class: {getClassLevel(review)}</p>
            )}
            {getSchoolName(review) && (
              <p className="text-gray-600 text-sm">School: {getSchoolName(review)}</p>
            )}
            <p className="text-gray-500 text-sm mt-1">
              Reviewed on {new Date(review.created_at).toLocaleDateString()}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ReviewDetailModal;
