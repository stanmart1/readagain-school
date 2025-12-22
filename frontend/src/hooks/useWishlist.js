import { useState, useEffect } from 'react';
import api from '../lib/api';

export const useWishlist = () => {
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchWishlist();
  }, []);

  const fetchWishlist = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await api.get('/wishlist');
      // Backend returns: { data: [...items] }
      setItems(response.data.data || []);
    } catch (err) {
      setError(err.message);
      console.error('Error fetching wishlist:', err);
    } finally {
      setLoading(false);
    }
  };

  const addToWishlist = async (bookId) => {
    try {
      await api.post('/wishlist', { book_id: bookId });
      await fetchWishlist();
      return true;
      return false;
    } catch (err) {
      console.error('Error adding to wishlist:', err);
      return false;
    }
  };

  const removeFromWishlist = async (itemId) => {
    try {
      await api.delete(`/wishlist/${itemId}`);
      await fetchWishlist();
      return true;
      return false;
    } catch (err) {
      console.error('Error removing from wishlist:', err);
      return false;
    }
  };

  return { items, loading, error, addToWishlist, removeFromWishlist, refetch: fetchWishlist };
};
