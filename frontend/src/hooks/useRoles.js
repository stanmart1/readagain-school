import { useState } from 'react';
import api from '../lib/api';

export const useRoles = () => {
  const [roles, setRoles] = useState([]);
  const [loading, setLoading] = useState(false);

  const fetchRoles = async () => {
    setLoading(true);
    try {
      const response = await api.get('/roles');
      // Backend returns: { roles: [...] }
      setRoles(response.data.roles || []);
    } catch (err) {
      console.error('Error fetching roles:', err);
      setRoles([]);
    } finally {
      setLoading(false);
    }
  };

  const createRole = async (roleData) => {
    try {
      const response = await api.post('/roles', roleData);
      return { success: true, data: response.data };
    } catch (err) {
      console.error('Error creating role:', err);
      return { success: false, error: err.message };
    }
  };

  const updateRole = async (roleId, roleData) => {
    try {
      const response = await api.put(`/roles/${roleId}`, roleData);
      return { success: true, data: response.data };
    } catch (err) {
      console.error('Error updating role:', err);
      return { success: false, error: err.message };
    }
  };

  const deleteRole = async (roleId) => {
    try {
      const response = await api.delete(`/roles/${roleId}`);
      return { success: true, data: response.data };
    } catch (err) {
      console.error('Error deleting role:', err);
      return { success: false, error: err.message };
    }
  };

  return { roles, loading, fetchRoles, createRole, updateRole, deleteRole };
};
