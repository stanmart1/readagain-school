export default function GroupsList({ groups, onEdit, onDelete, onViewMembers }) {
  return (
    <div className="bg-white rounded-lg shadow overflow-hidden">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Description</th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Members</th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Created By</th>
            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Actions</th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {groups.length === 0 ? (
            <tr>
              <td colSpan="5" className="px-6 py-12 text-center text-gray-500">
                No groups found. Create your first group!
              </td>
            </tr>
          ) : (
            groups.map((group) => (
              <tr key={group.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap font-medium">{group.name}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{group.description || '-'}</td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className="px-2 py-1 text-xs bg-blue-100 text-blue-800 rounded">
                    {group.member_count} members
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                  {group.creator?.first_name} {group.creator?.last_name}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-right text-sm space-x-2">
                  <button
                    onClick={() => onViewMembers(group)}
                    className="text-blue-600 hover:text-blue-900"
                  >
                    Members
                  </button>
                  <button
                    onClick={() => onEdit(group)}
                    className="text-indigo-600 hover:text-indigo-900"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => onDelete(group.id)}
                    className="text-red-600 hover:text-red-900"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
