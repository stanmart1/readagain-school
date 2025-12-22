import { useState, useEffect } from 'react';
import { useGroups } from '../../../hooks/useGroups';

export default function GroupForm({ group, onClose }) {
  const { createGroup, updateGroup, loading } = useGroups();
  const [formData, setFormData] = useState({
    name: '',
    description: ''
  });

  useEffect(() => {
    if (group) {
      setFormData({
        name: group.name,
        description: group.description || ''
      });
    }
  }, [group]);

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      if (group) {
        await updateGroup(group.id, formData);
      } else {
        await createGroup(formData);
      }
      onClose();
    } catch (error) {
      alert('Failed to save group');
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md transform transition-all">
        <div className="bg-gradient-to-r from-blue-600 to-purple-600 px-6 py-4 rounded-t-2xl">
          <h2 className="text-xl font-bold text-white flex items-center gap-2">
            <i className={`ri-${group ? 'edit' : 'add'}-line`}></i>
            {group ? 'Edit Group' : 'Create New Group'}
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
              className="w-full border border-gray-300 rounded-lg px-4 py-2.5 focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
              placeholder="e.g., Grade 5A, Science Club"
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
              className="w-full border border-gray-300 rounded-lg px-4 py-2.5 focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all resize-none"
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
              className="px-5 py-2.5 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg hover:shadow-lg transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed font-medium"
            >
              {loading ? (
                <span className="flex items-center gap-2">
                  <i className="ri-loader-4-line animate-spin"></i>
                  Saving...
                </span>
              ) : (
                <span className="flex items-center gap-2">
                  <i className="ri-save-line"></i>
                  {group ? 'Update' : 'Create'} Group
                </span>
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
