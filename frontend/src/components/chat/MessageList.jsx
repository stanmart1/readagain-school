import { useEffect, useRef } from 'react';
import { formatDistanceToNow } from 'date-fns';

const MessageList = ({ messages, currentUserId, onReact, onDelete, onEdit }) => {
  const messagesEndRef = useRef(null);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  return (
    <div className="flex-1 overflow-y-auto p-4 space-y-4">
      {messages.map((message) => {
        const isOwn = message.user_id === currentUserId;
        
        return (
          <div
            key={message.id}
            className={`flex ${isOwn ? 'justify-end' : 'justify-start'}`}
          >
            <div className={`max-w-[70%] ${isOwn ? 'items-end' : 'items-start'}`}>
              {!isOwn && (
                <div className="text-xs text-gray-600 mb-1 px-2">
                  {message.user?.username || 'Unknown'}
                </div>
              )}
              
              <div
                className={`rounded-lg px-4 py-2 ${
                  isOwn
                    ? 'bg-gradient-to-r from-purple-600 to-pink-600 text-white'
                    : 'bg-gray-100 text-gray-900'
                }`}
              >
                {message.reply_to && (
                  <div className="text-xs opacity-75 mb-2 pb-2 border-b border-white/20">
                    Replying to: {message.reply_to.message}
                  </div>
                )}
                
                <p className="whitespace-pre-wrap break-words">{message.message}</p>
                
                {message.is_edited && (
                  <span className="text-xs opacity-75 ml-2">(edited)</span>
                )}
              </div>
              
              <div className="flex items-center gap-2 mt-1 px-2">
                <span className="text-xs text-gray-500">
                  {formatDistanceToNow(new Date(message.created_at), { addSuffix: true })}
                </span>
                
                {isOwn && (
                  <div className="flex gap-1">
                    <button
                      onClick={() => onEdit(message)}
                      className="text-xs text-gray-500 hover:text-gray-700"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => onDelete(message.id)}
                      className="text-xs text-red-500 hover:text-red-700"
                    >
                      Delete
                    </button>
                  </div>
                )}
              </div>
            </div>
          </div>
        );
      })}
      <div ref={messagesEndRef} />
    </div>
  );
};

export default MessageList;
