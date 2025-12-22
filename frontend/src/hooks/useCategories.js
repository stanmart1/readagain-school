import { useState, useEffect } from 'react';
import api from '../lib/api';

export const useCategories = () => {
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [initialized, setInitialized] = useState(false);

  const fetchCategories = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await api.get('/categories');
      // Backend returns: { data: [...categories] }
      const categoriesData = response.data.data || [];
      setCategories(categoriesData);
      return { success: true, data: categoriesData };
    } catch (err) {
      console.error('Error fetching categories:', err);
      setError(err.message);
      return { success: false, error: err.message };
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (!initialized) {
      fetchCategories();
      setInitialized(true);
    }
  }, []);

  return {
    categories,
    loading,
    error,
    fetchCategories
  };
};
