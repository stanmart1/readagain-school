import { useState, useEffect } from 'react';
import DashboardLayout from '../../components/DashboardLayout';
import { useGroups } from '../../hooks/useGroups';
import Pagination from '../../components/admin/Pagination';
import GroupCard from '../../components/dashboard/groups/GroupCard';
import CreateGroupModal from '../../components/dashboard/groups/CreateGroupModal';
import ManageMembersModal from '../../components/dashboard/groups/ManageMembersModal';
import AssignBooksModal from '../../components/dashboard/groups/AssignBooksModal';

export default function MyGroups() {
  const { groups, loading, pagination, fetchGroups, deleteGroup } = useGroups();
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedGroup, setSelectedGroup] = useState(null);
  const [showMembers, setShowMembers] = useState(false);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [showAssignBooks, setShowAssignBooks] = useState(false);

  useEffect(() => {
    fetchGroups(1, 10, searchTerm);
  }, [fetchGroups, searchTerm]);

  const handlePageChange = (page) => {
    fetchGroups(page, pagination.limit, searchTerm);
  };

  const handleSearch = (e) => {
    setSearchTerm(e.target.value);
  };

  const handleViewMembers = (group) => {
    setSelectedGroup(group);
    setShowMembers(true);
  };

  const handleAssignBooks = (group) => {
    setSelectedGroup(group);
    setShowAssignBooks(true);
  };

  const handleCreateGroup = () => {
    setShowCreateForm(true);
  };

  const handleDelete = async (id) => {
    if (!confirm('Delete this group? Members will not be deleted.')) return;
    try {
      await deleteGroup(id);
    } catch (error) {
      alert('Failed to delete group');
    }
  };

  return (
    <DashboardLayout>
      <div className="p-6">
        <div className="flex justify-between items-center mb-6">
          <div>
            <h1 className="text-3xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
              My Groups
            </h1>
            <p className="text-gray-600 mt-1">Create and manage your groups</p>
          </div>
          <button
            onClick={handleCreateGroup}
            className="bg-gradient-to-r from-blue-600 to-purple-600 text-white px-6 py-2.5 rounded-lg hover:shadow-lg transition-all duration-200 flex items-center gap-2"
          >
            <i className="ri-add-line"></i>
            Create Group
          </button>
        </div>

        <div className="mb-4">
          <input
            type="text"
            placeholder="Search groups..."
            value={searchTerm}
            onChange={handleSearch}
            className="w-full max-w-md px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        {loading ? (
          <div className="flex justify-center items-center py-20">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          </div>
        ) : groups.length === 0 ? (
          <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-16 text-center">
            <i className="ri-group-line text-6xl text-gray-300 mb-4"></i>
            <p className="text-gray-500 text-lg font-medium">No groups found</p>
            <p className="text-gray-400 text-sm mt-1">Create your first group to get started</p>
          </div>
        ) : (
          <>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {groups.map((group) => (
                <GroupCard
                  key={group.id}
                  group={group}
                  onViewMembers={handleViewMembers}
                  onAssignBooks={handleAssignBooks}
                  onDelete={handleDelete}
                />
              ))}
            </div>

            {pagination.totalPages > 1 && (
              <div className="mt-6">
                <Pagination
                  currentPage={pagination.page}
                  totalPages={pagination.totalPages}
                  onPageChange={handlePageChange}
                />
              </div>
            )}
          </>
        )}

        {showCreateForm && (
          <CreateGroupModal
            onClose={() => {
              setShowCreateForm(false);
              fetchGroups(pagination.page, pagination.limit, searchTerm);
            }}
          />
        )}

        {showMembers && selectedGroup && (
          <ManageMembersModal
            group={selectedGroup}
            onClose={() => {
              setShowMembers(false);
              setSelectedGroup(null);
              fetchGroups(pagination.page, pagination.limit, searchTerm);
            }}
          />
        )}

        {showAssignBooks && selectedGroup && (
          <AssignBooksModal
            group={selectedGroup}
            onClose={() => {
              setShowAssignBooks(false);
              setSelectedGroup(null);
            }}
          />
        )}
      </div>
    </DashboardLayout>
  );
}
