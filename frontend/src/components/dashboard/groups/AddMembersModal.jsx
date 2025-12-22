import { useState, useEffect } from 'react';
import { useGroups } from '../../../hooks/useGroups';
import api from '../../../lib/api';

export default function AddMembersModal({ group, existingMembers, onClose }) {
  const { addMembers, loading } = useGroups();
  const [users, setUsers] = useState([]);
  const [selectedUsers, setSelectedUsers] = useState([]);

  useEffect(() => {
    fetchUsers();
  }, []);

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
      onClose();
    } catch (error) {
      alert('Failed to add members');
    }
  };

  const availableUsers = users.filter(
    (user) => !existingMembers.some((member) => member.id === user.id)
  );

  return (
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
            onClick={onClose}
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
  );
}
