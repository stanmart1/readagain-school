import { useState, useEffect } from 'react';
import api from '../lib/api';

export const useLibrary = () => {
  const [books, setBooks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchLibrary();
  }, []);

  const fetchLibrary = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await api.get('/library');
      // Backend returns: { data: [...books], meta: {...} }
      setBooks(response.data.data || []);
    } catch (err) {
      setError(err.message);
      console.error('Error fetching library:', err);
    } finally {
      setLoading(false);
    }
  };

  return { books, loading, error, refetch: fetchLibrary };
};
