import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useChat } from '../../hooks/useChat';
import { ArrowLeft, Wifi, WifiOff } from 'lucide-react';
import MessageList from '../../components/chat/MessageList';
import MessageInput from '../../components/chat/MessageInput';
import MemberList from '../../components/chat/MemberList';

const GroupChat = () => {
  const { groupId } = useParams();
  const navigate = useNavigate();
  const [roomId, setRoomId] = useState(null);
  const [onlineUsers, setOnlineUsers] = useState([]);
  
  const {
    room,
    messages,
    members,
    loading,
    connected,
    sendMessage,
    sendTyping,
    deleteMessage,
    updateMessage
  } = useChat(roomId);

  const currentUserId = parseInt(localStorage.getItem('userId'));

  useEffect(() => {
    // Fetch or create room for this group
    const initRoom = async () => {
      try {
        // First, try to find existing room
        const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/v1/chat/rooms?group_id=${groupId}`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });
        const data = await response.json();
        
        if (data.rooms && data.rooms.length > 0) {
          setRoomId(data.rooms[0].id);
        } else {
          // Create room for this group
          const createResponse = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/v1/chat/rooms`, {
            method: 'POST',
            headers: {
              'Authorization': `Bearer ${localStorage.getItem('token')}`,
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              type: 'group',
              name: `Group Chat`,
              description: 'Group discussion',
              group_id: parseInt(groupId)
            })
          });
          const newRoom = await createResponse.json();
          setRoomId(newRoom.id);
        }
      } catch (error) {
        console.error('Failed to initialize room:', error);
      }
    };

    if (groupId) {
      initRoom();
    }
  }, [groupId]);

  const handleSendMessage = async (message) => {
    try {
      await sendMessage(message);
    } catch (error) {
      console.error('Failed to send message:', error);
    }
  };

  const handleDeleteMessage = async (messageId) => {
    if (window.confirm('Delete this message?')) {
      await deleteMessage(messageId);
    }
  };

  const handleEditMessage = async (message) => {
    const newMessage = prompt('Edit message:', message.message);
    if (newMessage && newMessage !== message.message) {
      await updateMessage(message.id, newMessage);
    }
  };

  if (!roomId) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading chat...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-screen bg-white">
      {/* Header */}
      <div className="border-b p-4 flex items-center justify-between bg-white">
        <div className="flex items-center gap-4">
          <button
            onClick={() => navigate('/dashboard/groups')}
            className="p-2 hover:bg-gray-100 rounded-lg"
          >
            <ArrowLeft size={20} />
          </button>
          
          <div>
            <h1 className="text-xl font-bold text-gray-900">
              {room?.name || 'Group Chat'}
            </h1>
            <p className="text-sm text-gray-500">
              {members.length} members
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          {connected ? (
            <div className="flex items-center gap-2 text-green-600">
              <Wifi size={18} />
              <span className="text-sm">Connected</span>
            </div>
          ) : (
            <div className="flex items-center gap-2 text-red-600">
              <WifiOff size={18} />
              <span className="text-sm">Disconnected</span>
            </div>
          )}
        </div>
      </div>

      {/* Main Content */}
      <div className="flex flex-1 overflow-hidden">
        {/* Messages */}
        <div className="flex-1 flex flex-col">
          {loading ? (
            <div className="flex-1 flex items-center justify-center">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600"></div>
            </div>
          ) : (
            <>
              <MessageList
                messages={messages}
                currentUserId={currentUserId}
                onDelete={handleDeleteMessage}
                onEdit={handleEditMessage}
              />
              
              <MessageInput
                onSend={handleSendMessage}
                onTyping={sendTyping}
                disabled={!connected}
              />
            </>
          )}
        </div>

        {/* Members Sidebar */}
        <MemberList members={members} onlineUsers={onlineUsers} />
      </div>
    </div>
  );
};

export default GroupChat;
