import { useState, useEffect } from 'react';
import AdminLayout from '../../components/AdminLayout';
import GroupsList from '../../components/admin/groups/GroupsList';
import GroupForm from '../../components/admin/groups/GroupForm';
import GroupMembers from '../../components/admin/groups/GroupMembers';
import { useGroups } from '../../hooks/useGroups';

export default function Groups() {
  const { groups, loading, fetchGroups, deleteGroup } = useGroups();
  const [selectedGroup, setSelectedGroup] = useState(null);
  const [showForm, setShowForm] = useState(false);
  const [showMembers, setShowMembers] = useState(false);

  useEffect(() => {
    fetchGroups();
  }, [fetchGroups]);

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
      await deleteGroup(id);
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
          <div>
            <h1 className="text-3xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
              Groups Management
            </h1>
            <p className="text-gray-600 mt-1">Create and manage student groups</p>
          </div>
          <button
            onClick={handleCreate}
            className="bg-gradient-to-r from-blue-600 to-purple-600 text-white px-6 py-2.5 rounded-lg hover:shadow-lg transition-all duration-200 flex items-center gap-2"
          >
            <i className="ri-add-line"></i>
            Create Group
          </button>
        </div>

        {loading ? (
          <div className="flex justify-center items-center py-20">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          </div>
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
