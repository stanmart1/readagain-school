import { useState, useEffect } from 'react';
import DashboardLayout from '../../components/DashboardLayout';
import { useGroups } from '../../hooks/useGroups';
import Pagination from '../../components/admin/Pagination';
import GroupCard from '../../components/dashboard/groups/GroupCard';
import CreateGroupModal from '../../components/dashboard/groups/CreateGroupModal';
import ManageMembersModal from '../../components/dashboard/groups/ManageMembersModal';

export default function MyGroups() {
  const { groups, loading, pagination, fetchGroups, deleteGroup } = useGroups();
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedGroup, setSelectedGroup] = useState(null);
  const [showMembers, setShowMembers] = useState(false);
  const [showCreateForm, setShowCreateForm] = useState(false);

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
      </div>
    </DashboardLayout>
  );
}

export default function MyGroups() {
  const { groups, loading, pagination, fetchGroups, createGroup, deleteGroup } = useGroups();
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedGroup, setSelectedGroup] = useState(null);
  const [showMembers, setShowMembers] = useState(false);
  const [showCreateForm, setShowCreateForm] = useState(false);

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
                <div
                  key={group.id}
                  className="bg-white rounded-xl shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow"
                >
                  <div className="flex items-start justify-between mb-4">
                    <div className="flex items-center gap-3">
                      <div className="h-12 w-12 bg-gradient-to-br from-blue-500 to-purple-500 rounded-lg flex items-center justify-center">
                        <i className="ri-group-line text-white text-xl"></i>
                      </div>
                      <div>
                        <h3 className="font-semibold text-gray-900">{group.name}</h3>
                        <p className="text-sm text-gray-500">
                          {group.member_count} {group.member_count === 1 ? 'member' : 'members'}
                        </p>
                      </div>
                    </div>
                  </div>

                  {group.description && (
                    <p className="text-sm text-gray-600 mb-4 line-clamp-2">
                      {group.description}
                    </p>
                  )}

                  <div className="flex items-center justify-between pt-4 border-t border-gray-100">
                    <button
                      onClick={() => handleViewMembers(group)}
                      className="text-blue-600 hover:text-blue-700 text-sm font-medium flex items-center gap-1"
                    >
                      <i className="ri-team-line"></i>
                      Manage Members
                    </button>
                    <button
                      onClick={() => handleDelete(group.id)}
                      className="text-red-600 hover:text-red-700 text-sm font-medium flex items-center gap-1"
                    >
                      <i className="ri-delete-bin-line"></i>
                      Delete
                    </button>
                  </div>
                </div>
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
          <MembersModal
            group={selectedGroup}
            onClose={() => {
              setShowMembers(false);
              setSelectedGroup(null);
              fetchGroups(pagination.page, pagination.limit, searchTerm);
            }}
          />
        )}
      </div>
    </DashboardLayout>
  );
}

function CreateGroupModal({ onClose }) {
  const { createGroup, loading } = useGroups();
  const [formData, setFormData] = useState({
    name: '',
    description: ''
  });

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await createGroup(formData);
      onClose();
    } catch (error) {
      alert('Failed to create group');
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md">
        <div className="bg-gradient-to-r from-blue-600 to-purple-600 px-6 py-4 rounded-t-2xl">
          <h2 className="text-xl font-bold text-white flex items-center gap-2">
            <i className="ri-add-line"></i>
            Create New Group
          </h2>
        </div>

        <form onSubmit={handleSubmit} className="p-6">
          <div className="mb-5">
            <label className="block text-sm font-semibold text-gray-700 mb-2">
              Group Name <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="w-full border border-gray-300 rounded-lg px-4 py-2.5 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="e.g., Study Group, Project Team"
              required
            />
          </div>

          <div className="mb-6">
            <label className="block text-sm font-semibold text-gray-700 mb-2">
              Description
            </label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              className="w-full border border-gray-300 rounded-lg px-4 py-2.5 focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
              rows="4"
              placeholder="Brief description of the group..."
            />
          </div>

          <div className="flex justify-end gap-3">
            <button
              type="button"
              onClick={onClose}
              className="px-5 py-2.5 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors font-medium"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading}
              className="px-5 py-2.5 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg hover:shadow-lg transition-all duration-200 disabled:opacity-50 font-medium"
            >
              {loading ? 'Creating...' : 'Create Group'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

function MembersModal({ group, onClose }) {
  const { getMembers, addMembers, removeMember, loading } = useGroups();
  const [members, setMembers] = useState([]);
  const [users, setUsers] = useState([]);
  const [selectedUsers, setSelectedUsers] = useState([]);
  const [showAddModal, setShowAddModal] = useState(false);
  const [loadingMembers, setLoadingMembers] = useState(true);

  useEffect(() => {
    fetchMembers();
    fetchUsers();
  }, []);

  const fetchMembers = async () => {
    try {
      const data = await getMembers(group.id);
      setMembers(data);
    } catch (error) {
      console.error('Failed to fetch members:', error);
    } finally {
      setLoadingMembers(false);
    }
  };

  const fetchUsers = async () => {
    try {
      const { data } = await api.get('/users');
      setUsers(data.users || data);
    } catch (error) {
      console.error('Failed to fetch users:', error);
    }
  };

  const handleAddMembers = async () => {
    if (selectedUsers.length === 0) return;
    try {
      await addMembers(group.id, selectedUsers);
      setShowAddModal(false);
      setSelectedUsers([]);
      fetchMembers();
    } catch (error) {
      alert('Failed to add members');
    }
  };

  const handleRemoveMember = async (userId) => {
    if (!confirm('Remove this member from the group?')) return;
    try {
      await removeMember(group.id, userId);
      fetchMembers();
    } catch (error) {
      alert('Failed to remove member');
    }
  };

  const availableUsers = users.filter(
    (user) => !members.some((member) => member.id === user.id)
  );

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-4xl max-h-[85vh] overflow-hidden flex flex-col">
        <div className="bg-gradient-to-r from-blue-600 to-purple-600 px-6 py-4">
          <div className="flex justify-between items-center">
            <div>
              <h2 className="text-xl font-bold text-white flex items-center gap-2">
                <i className="ri-team-line"></i>
                {group.name}
              </h2>
              <p className="text-blue-100 text-sm mt-1">{members.length} members</p>
            </div>
            <button
              onClick={() => setShowAddModal(true)}
              className="bg-white text-blue-600 px-4 py-2 rounded-lg hover:bg-blue-50 transition-colors font-medium flex items-center gap-2"
            >
              <i className="ri-user-add-line"></i>
              Add Members
            </button>
          </div>
        </div>

        <div className="flex-1 overflow-y-auto p-6">
          {loadingMembers ? (
            <div className="flex justify-center items-center py-20">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
            </div>
          ) : members.length === 0 ? (
            <div className="text-center py-16">
              <i className="ri-user-line text-6xl text-gray-300 mb-4"></i>
              <p className="text-gray-500 text-lg font-medium">No members yet</p>
              <p className="text-gray-400 text-sm mt-1">Add some members to this group</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {members.map((member) => (
                <div
                  key={member.id}
                  className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
                >
                  <div className="flex items-center gap-3">
                    <div className="flex-shrink-0 h-12 w-12 bg-gradient-to-br from-blue-500 to-purple-500 rounded-full flex items-center justify-center">
                      <span className="text-white font-semibold text-sm">
                        {member.first_name?.[0]}{member.last_name?.[0]}
                      </span>
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="text-sm font-semibold text-gray-900 truncate">
                        {member.first_name} {member.last_name}
                      </div>
                      <div className="text-xs text-gray-500 truncate">{member.email}</div>
                    </div>
                  </div>
                  <button
                    onClick={() => handleRemoveMember(member.id)}
                    className="text-red-600 hover:text-red-700 p-2 rounded-lg hover:bg-red-50 transition-colors"
                  >
                    <i className="ri-user-unfollow-line text-lg"></i>
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="border-t border-gray-200 px-6 py-4 bg-gray-50">
          <button
            onClick={onClose}
            className="px-5 py-2.5 border border-gray-300 rounded-lg hover:bg-white transition-colors font-medium"
          >
            Close
          </button>
        </div>

        {showAddModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md max-h-[70vh] overflow-hidden flex flex-col">
              <div className="bg-gradient-to-r from-blue-600 to-purple-600 px-6 py-4">
                <h3 className="text-lg font-bold text-white flex items-center gap-2">
                  <i className="ri-user-add-line"></i>
                  Add Members to Group
                </h3>
              </div>

              <div className="flex-1 overflow-y-auto p-6">
                {availableUsers.length === 0 ? (
                  <p className="text-gray-500 text-center py-8">All users are already members</p>
                ) : (
                  <div className="space-y-2">
                    {availableUsers.map((user) => (
                      <label
                        key={user.id}
                        className="flex items-center space-x-3 p-3 hover:bg-gray-50 rounded-lg cursor-pointer transition-colors border border-transparent hover:border-gray-200"
                      >
                        <input
                          type="checkbox"
                          checked={selectedUsers.includes(user.id)}
                          onChange={(e) => {
                            if (e.target.checked) {
                              setSelectedUsers([...selectedUsers, user.id]);
                            } else {
                              setSelectedUsers(selectedUsers.filter((id) => id !== user.id));
                            }
                          }}
                          className="rounded text-blue-600 focus:ring-blue-500"
                        />
                        <div className="flex-1">
                          <div className="text-sm font-medium text-gray-900">
                            {user.first_name} {user.last_name}
                          </div>
                          <div className="text-xs text-gray-500">{user.email}</div>
                        </div>
                      </label>
                    ))}
                  </div>
                )}
              </div>

              <div className="border-t border-gray-200 px-6 py-4 bg-gray-50 flex justify-end gap-3">
                <button
                  onClick={() => {
                    setShowAddModal(false);
                    setSelectedUsers([]);
                  }}
                  className="px-5 py-2.5 border border-gray-300 rounded-lg hover:bg-white transition-colors font-medium"
                >
                  Cancel
                </button>
                <button
                  onClick={handleAddMembers}
                  disabled={selectedUsers.length === 0 || loading}
                  className="px-5 py-2.5 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg hover:shadow-lg transition-all duration-200 disabled:opacity-50 font-medium"
                >
                  {loading ? 'Adding...' : `Add ${selectedUsers.length} ${selectedUsers.length === 1 ? 'Member' : 'Members'}`}
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
