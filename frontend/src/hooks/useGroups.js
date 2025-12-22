import { useState, useCallback } from 'react';
import api from '../lib/api';

export const useGroups = () => {
  const [groups, setGroups] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchGroups = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const { data } = await api.get('/groups');
      setGroups(data);
      return data;
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to fetch groups');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const getGroup = useCallback(async (id) => {
    setLoading(true);
    setError(null);
    try {
      const { data } = await api.get(`/groups/${id}`);
      return data;
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to fetch group');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const createGroup = useCallback(async (groupData) => {
    setLoading(true);
    setError(null);
    try {
      const { data } = await api.post('/groups', groupData);
      setGroups((prev) => [data, ...prev]);
      return data;
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to create group');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const updateGroup = useCallback(async (id, groupData) => {
    setLoading(true);
    setError(null);
    try {
      await api.put(`/groups/${id}`, groupData);
      setGroups((prev) =>
        prev.map((g) => (g.id === id ? { ...g, ...groupData } : g))
      );
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to update group');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const deleteGroup = useCallback(async (id) => {
    setLoading(true);
    setError(null);
    try {
      await api.delete(`/groups/${id}`);
      setGroups((prev) => prev.filter((g) => g.id !== id));
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to delete group');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const getMembers = useCallback(async (groupId) => {
    setLoading(true);
    setError(null);
    try {
      const { data } = await api.get(`/groups/${groupId}/members`);
      return data;
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to fetch members');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const addMember = useCallback(async (groupId, userId) => {
    setLoading(true);
    setError(null);
    try {
      await api.post(`/groups/${groupId}/members`, { user_id: userId });
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to add member');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const addMembers = useCallback(async (groupId, userIds) => {
    setLoading(true);
    setError(null);
    try {
      await api.post(`/groups/${groupId}/members/bulk`, { user_ids: userIds });
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to add members');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const removeMember = useCallback(async (groupId, userId) => {
    setLoading(true);
    setError(null);
    try {
      await api.delete(`/groups/${groupId}/members/${userId}`);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to remove member');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    groups,
    loading,
    error,
    fetchGroups,
    getGroup,
    createGroup,
    updateGroup,
    deleteGroup,
    getMembers,
    addMember,
    addMembers,
    removeMember,
  };
};
