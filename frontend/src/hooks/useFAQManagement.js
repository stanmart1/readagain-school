import { useState, useEffect, useCallback } from 'react';
import api from '../lib/api';

export const useFAQManagement = () => {
  const [faqs, setFaqs] = useState([]);
  const [categories, setCategories] = useState([]);
  const [stats, setStats] = useState({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchData = useCallback(async () => {
    setLoading(true);
    try {
      const faqResponse = await api.get('/admin/faqs');
      // Backend returns: { data: [...faqs] }
      setFaqs(faqResponse.data.data || []);

      const categoryResponse = await api.get('/faqs/categories');
      // Backend returns: { data: [...categories] }
      setCategories(categoryResponse.data.data || []);

      setStats({
        total_faqs: faqResponse.data?.total || 0,
        total_categories: categoryResponse.data?.length || 0,
        total_views: 0,
        recent_faqs: faqResponse.data?.faqs?.slice(0, 5) || []
      });
    } catch (err) {
      setError(err.message);
      console.error('Error fetching FAQ data:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  const createFAQ = useCallback(async (data) => {
    try {
      await api.post('/admin/faqs', data);
      await fetchData();
      return { success: true };
    } catch (err) {
      console.error('Error creating FAQ:', err);
      return { success: false, error: err.message };
    }
  }, [fetchData]);

  const updateFAQ = useCallback(async (data) => {
    try {
      await api.put(`/admin/faqs/${data.id}`, data);
      await fetchData();
      return { success: true };
    } catch (err) {
      console.error('Error updating FAQ:', err);
      return { success: false, error: err.message };
    }
  }, [fetchData]);

  const deleteFAQ = useCallback(async (id) => {
    try {
      await api.delete(`/admin/faqs/${id}`);
      await fetchData();
      return { success: true };
    } catch (err) {
      console.error('Error deleting FAQ:', err);
      return { success: false, error: err.message };
    }
  }, [fetchData]);

  const bulkDeleteFAQs = useCallback(async (ids) => {
    try {
      // Bulk operations not implemented in backend
      return { success: false, error: 'Bulk delete not available' };
      await fetchData();
      return { success: true };
    } catch (err) {
      console.error('Error bulk deleting FAQs:', err);
      return { success: false, error: err.message };
    }
  }, [fetchData]);

  const bulkUpdateFAQs = useCallback(async (ids, updates) => {
    try {
      // Bulk operations not implemented in backend
      return { success: false, error: 'Bulk update not available' };
      await fetchData();
      return { success: true };
    } catch (err) {
      console.error('Error bulk updating FAQs:', err);
      return { success: false, error: err.message };
    }
  }, [fetchData]);

  const createCategory = useCallback(async (data) => {
    try {
      // FAQ categories not separate - use main categories
      return { success: false, error: 'Use main categories endpoint' };
      await fetchData();
      return { success: true };
    } catch (err) {
      console.error('Error creating category:', err);
      return { success: false, error: err.message };
    }
  }, [fetchData]);

  const updateCategory = useCallback(async (data) => {
    try {
      // FAQ categories not separate - use main categories
      return { success: false, error: 'Use main categories endpoint' };
      await fetchData();
      return { success: true };
    } catch (err) {
      console.error('Error updating category:', err);
      return { success: false, error: err.message };
    }
  }, [fetchData]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return {
    faqs,
    categories,
    stats,
    loading,
    error,
    refetch: fetchData,
    createFAQ,
    updateFAQ,
    deleteFAQ,
    bulkDeleteFAQs,
    bulkUpdateFAQs,
    createCategory,
    updateCategory
  };
};
