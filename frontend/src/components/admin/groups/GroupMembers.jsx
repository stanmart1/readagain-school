import { useState, useEffect } from 'react';
import { useGroups } from '../../../hooks/useGroups';
import api from '../../../lib/api';

export default function GroupMembers({ group, onClose }) {
  const { getMembers, addMembers, removeMember, loading } = useGroups();
  const [members, setMembers] = useState([]);
  const [users, setUsers] = useState([]);
  const [selectedUsers, setSelectedUsers] = useState([]);
  const [showAddModal, setShowAddModal] = useState(false);
  const [loadingMembers, setLoadingMembers] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');

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

  const filteredUsers = availableUsers.filter((user) =>
    `${user.first_name} ${user.last_name} ${user.email}`
      .toLowerCase()
      .includes(searchTerm.toLowerCase())
  );

  return (
    <div 
      className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
      onClick={onClose}
    >
      <div 
        className="bg-white rounded-2xl shadow-2xl w-full max-w-4xl max-h-[85vh] overflow-hidden flex flex-col"
        onClick={(e) => e.stopPropagation()}
      >
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
            <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase">Name</th>
                    <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase">Email</th>
                    <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase">Role</th>
                    <th className="px-6 py-3 text-right text-xs font-semibold text-gray-700 uppercase">Actions</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {members.map((member) => (
                    <tr key={member.id} className="hover:bg-gray-50 transition-colors">
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="flex items-center">
                          <div className="flex-shrink-0 h-10 w-10 bg-gradient-to-br from-blue-500 to-purple-500 rounded-full flex items-center justify-center">
                            <span className="text-white font-semibold text-sm">
                              {member.first_name?.[0]}{member.last_name?.[0]}
                            </span>
                          </div>
                          <div className="ml-4">
                            <div className="text-sm font-semibold text-gray-900">
                              {member.first_name} {member.last_name}
                            </div>
                          </div>
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                        {member.email}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                          {member.role?.name || 'N/A'}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-right">
                        <button
                          onClick={() => handleRemoveMember(member.id)}
                          className="text-red-600 hover:text-red-900 transition-colors font-medium text-sm"
                        >
                          <i className="ri-user-unfollow-line mr-1"></i>
                          Remove
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
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
          <div 
            className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
            onClick={() => {
              setShowAddModal(false);
              setSelectedUsers([]);
              setSearchTerm('');
            }}
          >
            <div 
              className="bg-white rounded-2xl shadow-2xl w-full max-w-md max-h-[70vh] overflow-hidden flex flex-col"
              onClick={(e) => e.stopPropagation()}
            >
              <div className="bg-gradient-to-r from-blue-600 to-purple-600 px-6 py-4">
                <h3 className="text-lg font-bold text-white flex items-center gap-2">
                  <i className="ri-user-add-line"></i>
                  Add Members to Group
                </h3>
              </div>

              <div className="p-4 border-b border-gray-200">
                <input
                  type="text"
                  placeholder="Search users..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>

              <div className="flex-1 overflow-y-auto p-6">
                {filteredUsers.length === 0 ? (
                  <p className="text-gray-500 text-center py-8">
                    {searchTerm ? 'No users found matching your search' : 'All users are already members'}
                  </p>
                ) : (
                  <div className="space-y-2">
                    {filteredUsers.map((user) => (
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
                  className="px-5 py-2.5 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg hover:shadow-lg transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed font-medium"
                >
                  {loading ? (
                    <span className="flex items-center gap-2">
                      <i className="ri-loader-4-line animate-spin"></i>
                      Adding...
                    </span>
                  ) : (
                    `Add ${selectedUsers.length} ${selectedUsers.length === 1 ? 'Member' : 'Members'}`
                  )}
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
