import { Users, UserCheck } from 'lucide-react';

const MemberList = ({ members, onlineUsers = [] }) => {
  return (
    <div className="w-64 border-l bg-gray-50 p-4 overflow-y-auto">
      <h3 className="font-semibold text-gray-900 mb-4 flex items-center gap-2">
        <Users size={18} />
        Members ({members.length})
      </h3>
      
      <div className="space-y-2">
        {members.map((member) => {
          const isOnline = onlineUsers.includes(member.user_id);
          
          return (
            <div
              key={member.id}
              className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-100"
            >
              <div className="relative">
                <div className="w-10 h-10 rounded-full bg-gradient-to-r from-purple-600 to-pink-600 flex items-center justify-center text-white font-semibold">
                  {member.user?.username?.[0]?.toUpperCase() || 'U'}
                </div>
                {isOnline && (
                  <div className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 rounded-full border-2 border-white"></div>
                )}
              </div>
              
              <div className="flex-1 min-w-0">
                <div className="font-medium text-gray-900 truncate">
                  {member.user?.username || 'Unknown'}
                </div>
                <div className="text-xs text-gray-500 capitalize">
                  {member.role}
                </div>
              </div>
              
              {isOnline && (
                <UserCheck size={16} className="text-green-500" />
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default MemberList;
