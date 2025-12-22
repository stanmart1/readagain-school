import { useState } from 'react';
import api from '../lib/api';

export const useAuditLogs = () => {
  const [auditLogs, setAuditLogs] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [totalLogs, setTotalLogs] = useState(0);

  const fetchAuditLogs = async (page = 1, filters = {}) => {
    try {
      setLoading(true);
      setError(null);

      const params = {
        skip: (page - 1) * 20,
        limit: 20,
        ...(filters.userId && { user_id: filters.userId }),
        ...(filters.action && { action: filters.action }),
        ...(filters.resource && { resource: filters.resource })
      };

      const response = await api.get('/admin/audit', { params });
      // Backend returns: { data: [...logs], meta: {...} }
      const data = response.data;

      setAuditLogs(data.data || []);
      setTotalLogs(data.meta?.total || 0);
      
      return { success: true, data };
    } catch (err) {
      setError(err.message);
      console.error('Error fetching audit logs:', err);
      return { success: false, error: err.message };
    } finally {
      setLoading(false);
    }
  };

  const exportAuditLogs = async (filters = {}) => {
    try {
      const params = {
        ...(filters.userId && { user_id: filters.userId }),
        ...(filters.action && { action: filters.action }),
        ...(filters.resource && { resource: filters.resource })
      };

      const response = await api.get('/admin/audit', { params });
      return { success: true, data: response.data.data || [] };
    } catch (err) {
      console.error('Error exporting audit logs:', err);
      return { success: false, error: err.message };
    }
  };

  return {
    auditLogs,
    loading,
    error,
    totalLogs,
    fetchAuditLogs,
    exportAuditLogs
  };
};
