import { useState } from 'react';
import api from '../lib/api';

export const useReadingAnalytics = () => {
  const [analyticsData, setAnalyticsData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchReadingAnalytics = async (period = 'month') => {
    try {
      setLoading(true);
      setError(null);
      const response = await api.get(`/analytics/reading?period=${period}`);
      setAnalyticsData(response.data);
      return { success: true, data: response.data };
    } catch (err) {
      console.error('Error fetching reading analytics:', err);
      setError('Error fetching reading analytics');
      return { success: false, error: err.message };
    } finally {
      setLoading(false);
    }
  };

  return {
    analyticsData,
    loading,
    error,
    fetchReadingAnalytics
  };
};
