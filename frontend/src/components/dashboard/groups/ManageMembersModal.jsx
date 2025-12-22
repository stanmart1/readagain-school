import { useState, useEffect } from 'react';
import { useGroups } from '../../../hooks/useGroups';
import api from '../../../lib/api';
import AddMembersModal from './AddMembersModal';

export default function ManageMembersModal({ group, onClose }) {
  const { getMembers, removeMember } = useGroups();
  const [members, setMembers] = useState([]);
  const [showAddModal, setShowAddModal] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchMembers();
  }, []);

  const fetchMembers = async () => {
    try {
      const data = await getMembers(group.id);
      setMembers(data);
    } catch (error) {
      console.error('Failed to fetch members:', error);
    } finally {
      setLoading(false);
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
          {loading ? (
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
          <AddMembersModal
            group={group}
            existingMembers={members}
            onClose={() => {
              setShowAddModal(false);
              fetchMembers();
            }}
          />
        )}
      </div>
    </div>
  );
}
