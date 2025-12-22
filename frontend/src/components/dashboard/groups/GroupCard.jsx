import { useNavigate } from 'react-router-dom';

export default function GroupCard({ group, onViewMembers, onDelete, onAssignBooks }) {
  const navigate = useNavigate();

  return (
    <div 
      onClick={() => navigate(`/dashboard/groups/${group.id}/chat`)}
      className="bg-white rounded-xl shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow cursor-pointer"
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

      <div className="flex items-center justify-between pt-4 border-t border-gray-100 gap-2">
        <button
          onClick={(e) => {
            e.stopPropagation();
            onAssignBooks(group);
          }}
          className="text-green-600 hover:text-green-700 text-sm font-medium flex items-center gap-1"
        >
          <i className="ri-book-line"></i>
          Assign Books
        </button>
        <button
          onClick={(e) => {
            e.stopPropagation();
            onViewMembers(group);
          }}
          className="text-blue-600 hover:text-blue-700 text-sm font-medium flex items-center gap-1"
        >
          <i className="ri-team-line"></i>
          Members
        </button>
        <button
          onClick={(e) => {
            e.stopPropagation();
            onDelete(group.id);
          }}
          className="text-red-600 hover:text-red-700 text-sm font-medium flex items-center gap-1"
        >
          <i className="ri-delete-bin-line"></i>
          Delete
        </button>
      </div>
    </div>
  );
}
