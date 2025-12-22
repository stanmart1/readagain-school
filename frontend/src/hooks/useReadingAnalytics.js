import { useState } from 'react';
import api from '../lib/api';

export const useReadingAnalytics = () => {
  const [analyticsData, setAnalyticsData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchReadingAnalytics = async (period = 'month', classLevel = '') => {
    try {
      setLoading(true);
      setError(null);
      const params = new URLSearchParams({ period });
      if (classLevel) params.append('class_level', classLevel);
      
      const response = await api.get(`/analytics/reading?${params}`);
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
