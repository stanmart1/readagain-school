import { useState, useEffect } from 'react';
import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000/api/v1';

export default function GroupMembers({ group, onClose }) {
  const [members, setMembers] = useState([]);
  const [users, setUsers] = useState([]);
  const [selectedUsers, setSelectedUsers] = useState([]);
  const [showAddModal, setShowAddModal] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchMembers();
    fetchUsers();
  }, []);

  const fetchMembers = async () => {
    try {
      const token = localStorage.getItem('token');
      const { data } = await axios.get(`${API_URL}/admin/groups/${group.id}/members`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      setMembers(data);
    } catch (error) {
      console.error('Failed to fetch members:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchUsers = async () => {
    try {
      const token = localStorage.getItem('token');
      const { data } = await axios.get(`${API_URL}/users`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      setUsers(data.users || data);
    } catch (error) {
      console.error('Failed to fetch users:', error);
    }
  };

  const handleAddMembers = async () => {
    if (selectedUsers.length === 0) return;

    try {
      const token = localStorage.getItem('token');
      await axios.post(
        `${API_URL}/admin/groups/${group.id}/members/bulk`,
        { user_ids: selectedUsers },
        { headers: { Authorization: `Bearer ${token}` } }
      );
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
      const token = localStorage.getItem('token');
      await axios.delete(`${API_URL}/admin/groups/${group.id}/members/${userId}`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      fetchMembers();
    } catch (error) {
      alert('Failed to remove member');
    }
  };

  const availableUsers = users.filter(
    (user) => !members.some((member) => member.id === user.id)
  );

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-full max-w-3xl max-h-[80vh] overflow-y-auto">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-bold">
            {group.name} - Members ({members.length})
          </h2>
          <button
            onClick={() => setShowAddModal(true)}
            className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
          >
            Add Members
          </button>
        </div>

        {loading ? (
          <div className="text-center py-12">Loading...</div>
        ) : members.length === 0 ? (
          <div className="text-center py-12 text-gray-500">
            No members yet. Add some members to this group.
          </div>
        ) : (
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Role</th>
                <th className="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Actions</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {members.map((member) => (
                <tr key={member.id}>
                  <td className="px-4 py-3 whitespace-nowrap">
                    {member.first_name} {member.last_name}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-600">
                    {member.email}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap text-sm">
                    <span className="px-2 py-1 bg-gray-100 rounded text-xs">
                      {member.role?.name || 'N/A'}
                    </span>
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap text-right">
                    <button
                      onClick={() => handleRemoveMember(member.id)}
                      className="text-red-600 hover:text-red-900 text-sm"
                    >
                      Remove
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}

        <div className="mt-6 flex justify-end">
          <button
            onClick={onClose}
            className="px-4 py-2 border rounded hover:bg-gray-50"
          >
            Close
          </button>
        </div>

        {showAddModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-md max-h-[60vh] overflow-y-auto">
              <h3 className="text-lg font-bold mb-4">Add Members</h3>

              {availableUsers.length === 0 ? (
                <p className="text-gray-500 text-center py-4">All users are already members</p>
              ) : (
                <div className="space-y-2 mb-4">
                  {availableUsers.map((user) => (
                    <label key={user.id} className="flex items-center space-x-2 p-2 hover:bg-gray-50 rounded">
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
                        className="rounded"
                      />
                      <span className="flex-1">
                        {user.first_name} {user.last_name} ({user.email})
                      </span>
                    </label>
                  ))}
                </div>
              )}

              <div className="flex justify-end space-x-2">
                <button
                  onClick={() => {
                    setShowAddModal(false);
                    setSelectedUsers([]);
                  }}
                  className="px-4 py-2 border rounded hover:bg-gray-50"
                >
                  Cancel
                </button>
                <button
                  onClick={handleAddMembers}
                  disabled={selectedUsers.length === 0}
                  className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
                >
                  Add Selected ({selectedUsers.length})
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
