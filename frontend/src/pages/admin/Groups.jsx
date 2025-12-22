import { useState, useEffect } from 'react';
import axios from 'axios';
import AdminLayout from '../../components/AdminLayout';
import GroupsList from '../../components/admin/groups/GroupsList';
import GroupForm from '../../components/admin/groups/GroupForm';
import GroupMembers from '../../components/admin/groups/GroupMembers';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000/api/v1';

export default function Groups() {
  const [groups, setGroups] = useState([]);
  const [selectedGroup, setSelectedGroup] = useState(null);
  const [showForm, setShowForm] = useState(false);
  const [showMembers, setShowMembers] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchGroups();
  }, []);

  const fetchGroups = async () => {
    try {
      const token = localStorage.getItem('token');
      const { data } = await axios.get(`${API_URL}/groups`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      setGroups(data);
    } catch (error) {
      console.error('Failed to fetch groups:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setSelectedGroup(null);
    setShowForm(true);
  };

  const handleEdit = (group) => {
    setSelectedGroup(group);
    setShowForm(true);
  };

  const handleDelete = async (id) => {
    if (!confirm('Delete this group? Members will not be deleted.')) return;
    
    try {
      const token = localStorage.getItem('token');
      await axios.delete(`${API_URL}/groups/${id}`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      fetchGroups();
    } catch (error) {
      alert('Failed to delete group');
    }
  };

  const handleViewMembers = (group) => {
    setSelectedGroup(group);
    setShowMembers(true);
  };

  const handleFormClose = () => {
    setShowForm(false);
    setSelectedGroup(null);
    fetchGroups();
  };

  const handleMembersClose = () => {
    setShowMembers(false);
    setSelectedGroup(null);
    fetchGroups();
  };

  return (
    <AdminLayout>
      <div className="p-6">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-2xl font-bold">Groups Management</h1>
          <button
            onClick={handleCreate}
            className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
          >
            Create Group
          </button>
        </div>

        {loading ? (
          <div className="text-center py-12">Loading...</div>
        ) : (
          <GroupsList
            groups={groups}
            onEdit={handleEdit}
            onDelete={handleDelete}
            onViewMembers={handleViewMembers}
          />
        )}

        {showForm && (
          <GroupForm
            group={selectedGroup}
            onClose={handleFormClose}
          />
        )}

        {showMembers && selectedGroup && (
          <GroupMembers
            group={selectedGroup}
            onClose={handleMembersClose}
          />
        )}
      </div>
    </AdminLayout>
  );
}
