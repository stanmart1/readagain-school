import { useState, useEffect } from 'react';
import { useGroups } from '../../../hooks/useGroups';
import api from '../../../lib/api';
import { getFileUrl } from '../../../lib/fileService';

export default function AssignBooksModal({ group, onClose }) {
  const { assignBooks, loading } = useGroups();
  const [books, setBooks] = useState([]);
  const [selectedBooks, setSelectedBooks] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [loadingBooks, setLoadingBooks] = useState(true);

  useEffect(() => {
    fetchBooks();
  }, []);

  const fetchBooks = async () => {
    try {
      const { data } = await api.get('/books');
      setBooks(data.books || data);
    } catch (error) {
      console.error('Failed to fetch books:', error);
    } finally {
      setLoadingBooks(false);
    }
  };

  const getBookCoverUrl = (coverImage) => {
    if (!coverImage) return null;
    if (coverImage.startsWith('http')) return coverImage;
    return `${UPLOAD_API_URL}/files/${coverImage}`;
  };

  const handleAssignBooks = async () => {
    if (selectedBooks.length === 0) return;
    try {
      await assignBooks(group.id, selectedBooks);
      onClose();
    } catch (error) {
      alert('Failed to assign books');
    }
  };

  const filteredBooks = books.filter((book) =>
    `${book.title} ${book.author?.business_name || ''}`
      .toLowerCase()
      .includes(searchTerm.toLowerCase())
  );

  return (
    <div 
      className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
      onClick={onClose}
    >
      <div 
        className="bg-white rounded-2xl shadow-2xl w-full max-w-3xl max-h-[80vh] overflow-hidden flex flex-col"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="bg-gradient-to-r from-green-600 to-blue-600 px-6 py-4">
          <h3 className="text-lg font-bold text-white flex items-center gap-2">
            <i className="ri-book-line"></i>
            Assign Books to {group.name}
          </h3>
          <p className="text-green-100 text-sm mt-1">
            Books will be added to all {group.member_count} members' libraries
          </p>
        </div>

        <div className="p-4 border-b border-gray-200">
          <input
            type="text"
            placeholder="Search books by title or author..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent"
          />
        </div>

        <div className="flex-1 overflow-y-auto p-6">
          {loadingBooks ? (
            <div className="flex justify-center items-center py-20">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-green-600"></div>
            </div>
          ) : filteredBooks.length === 0 ? (
            <p className="text-gray-500 text-center py-8">
              {searchTerm ? 'No books found matching your search' : 'No books available'}
            </p>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {filteredBooks.map((book) => (
                <label
                  key={book.id}
                  className="flex items-start space-x-3 p-4 hover:bg-gray-50 rounded-lg cursor-pointer transition-colors border border-gray-200 hover:border-green-300"
                >
                  <input
                    type="checkbox"
                    checked={selectedBooks.includes(book.id)}
                    onChange={(e) => {
                      if (e.target.checked) {
                        setSelectedBooks([...selectedBooks, book.id]);
                      } else {
                        setSelectedBooks(selectedBooks.filter((id) => id !== book.id));
                      }
                    }}
                    className="mt-1 rounded text-green-600 focus:ring-green-500"
                  />
                  <div className="flex-1 min-w-0">
                    <div className="flex items-start gap-3">
                      {book.cover_image && (
                        <img
                          src={getFileUrl(book.cover_image)}
                          alt={book.title}
                          className="w-12 h-16 object-cover rounded"
                        />
                      )}
                      <div className="flex-1">
                        <div className="text-sm font-semibold text-gray-900 line-clamp-2">
                          {book.title}
                        </div>
                        <div className="text-xs text-gray-500 mt-1">
                          {book.author?.business_name || 'Unknown Author'}
                        </div>
                        {book.category && (
                          <span className="inline-block mt-1 px-2 py-0.5 text-xs bg-gray-100 text-gray-700 rounded">
                            {book.category.name}
                          </span>
                        )}
                      </div>
                    </div>
                  </div>
                </label>
              ))}
            </div>
          )}
        </div>

        <div className="border-t border-gray-200 px-6 py-4 bg-gray-50 flex justify-between items-center">
          <div className="text-sm text-gray-600">
            {selectedBooks.length} {selectedBooks.length === 1 ? 'book' : 'books'} selected
          </div>
          <div className="flex gap-3">
            <button
              onClick={onClose}
              className="px-5 py-2.5 border border-gray-300 rounded-lg hover:bg-white transition-colors font-medium"
            >
              Cancel
            </button>
            <button
              onClick={handleAssignBooks}
              disabled={selectedBooks.length === 0 || loading}
              className="px-5 py-2.5 bg-gradient-to-r from-green-600 to-blue-600 text-white rounded-lg hover:shadow-lg transition-all duration-200 disabled:opacity-50 font-medium"
            >
              {loading ? 'Assigning...' : `Assign ${selectedBooks.length} ${selectedBooks.length === 1 ? 'Book' : 'Books'}`}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
