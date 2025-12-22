import { useState, useEffect } from 'react';
import api from '../lib/api';

export const useTestimonials = (limit = 10, featuredOnly = true) => {
  const [testimonials, setTestimonials] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchTestimonials();
  }, [limit, featuredOnly]);

  const fetchTestimonials = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await api.get('/testimonials');
      // Backend returns: { data: [...testimonials] }
      setTestimonials(response.data.data || []);
    } catch (err) {
      setError(err.message);
      console.error('Error fetching testimonials:', err);
    } finally {
      setLoading(false);
    }
  };

  return { testimonials, loading, error, refetch: fetchTestimonials };
};
