export default function GroupsList({ groups, onEdit, onDelete, onViewMembers }) {
  return (
    <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gradient-to-r from-gray-50 to-gray-100">
            <tr>
              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Group Name
              </th>
              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Description
              </th>
              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Members
              </th>
              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Created By
              </th>
              <th className="px-6 py-4 text-right text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {groups.length === 0 ? (
              <tr>
                <td colSpan="5" className="px-6 py-16 text-center">
                  <div className="flex flex-col items-center justify-center">
                    <i className="ri-group-line text-6xl text-gray-300 mb-4"></i>
                    <p className="text-gray-500 text-lg font-medium">No groups found</p>
                    <p className="text-gray-400 text-sm mt-1">Create your first group to get started!</p>
                  </div>
                </td>
              </tr>
            ) : (
              groups.map((group) => (
                <tr key={group.id} className="hover:bg-gray-50 transition-colors">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      <div className="flex-shrink-0 h-10 w-10 bg-gradient-to-br from-blue-500 to-purple-500 rounded-lg flex items-center justify-center">
                        <i className="ri-group-line text-white text-lg"></i>
                      </div>
                      <div className="ml-4">
                        <div className="text-sm font-semibold text-gray-900">{group.name}</div>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4">
                    <div className="text-sm text-gray-600 max-w-xs truncate">
                      {group.description || <span className="text-gray-400 italic">No description</span>}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-gradient-to-r from-blue-100 to-purple-100 text-blue-800">
                      <i className="ri-user-line mr-1"></i>
                      {group.member_count} {group.member_count === 1 ? 'member' : 'members'}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                    {group.creator?.first_name} {group.creator?.last_name}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-3">
                    <button
                      onClick={() => onViewMembers(group)}
                      className="text-blue-600 hover:text-blue-900 transition-colors"
                    >
                      <i className="ri-team-line mr-1"></i>
                      Members
                    </button>
                    <button
                      onClick={() => onEdit(group)}
                      className="text-indigo-600 hover:text-indigo-900 transition-colors"
                    >
                      <i className="ri-edit-line mr-1"></i>
                      Edit
                    </button>
                    <button
                      onClick={() => onDelete(group.id)}
                      className="text-red-600 hover:text-red-900 transition-colors"
                    >
                      <i className="ri-delete-bin-line mr-1"></i>
                      Delete
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
