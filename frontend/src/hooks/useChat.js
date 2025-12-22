import { useState, useEffect, useCallback } from 'react';
import api from '../lib/api';
import websocketService from '../services/websocketService';

export const useChat = (roomId) => {
  const [messages, setMessages] = useState([]);
  const [members, setMembers] = useState([]);
  const [room, setRoom] = useState(null);
  const [loading, setLoading] = useState(false);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    if (!roomId) return;

    const token = localStorage.getItem('token');
    
    // Connect WebSocket
    websocketService.connect(roomId, token);

    // WebSocket event listeners
    websocketService.on('connected', () => setConnected(true));
    websocketService.on('disconnected', () => setConnected(false));
    
    websocketService.on('message', (data) => {
      setMessages(prev => [...prev, data.data]);
    });

    websocketService.on('typing', (data) => {
      console.log(`${data.username} is typing...`);
    });

    websocketService.on('join', (data) => {
      console.log(`${data.username} joined`);
    });

    websocketService.on('leave', (data) => {
      console.log(`${data.username} left`);
    });

    // Fetch initial data
    fetchRoom();
    fetchMessages();
    fetchMembers();

    return () => {
      websocketService.disconnect();
    };
  }, [roomId]);

  const fetchRoom = async () => {
    try {
      const response = await api.get(`/chat/rooms/${roomId}`);
      setRoom(response.data);
    } catch (error) {
      console.error('Failed to fetch room:', error);
    }
  };

  const fetchMessages = async (page = 1) => {
    setLoading(true);
    try {
      const response = await api.get(`/chat/rooms/${roomId}/messages`, {
        params: { page, limit: 50 }
      });
      setMessages(response.data.messages);
    } catch (error) {
      console.error('Failed to fetch messages:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchMembers = async () => {
    try {
      const response = await api.get(`/chat/rooms/${roomId}/members`);
      setMembers(response.data);
    } catch (error) {
      console.error('Failed to fetch members:', error);
    }
  };

  const sendMessage = useCallback(async (message, messageType = 'text', fileUrl = '', fileName = '', replyToId = null) => {
    try {
      await api.post(`/chat/rooms/${roomId}/messages`, {
        message,
        message_type: messageType,
        file_url: fileUrl,
        file_name: fileName,
        reply_to_id: replyToId
      });
    } catch (error) {
      console.error('Failed to send message:', error);
      throw error;
    }
  }, [roomId]);

  const sendTyping = useCallback(() => {
    websocketService.send('typing', { room_id: roomId });
  }, [roomId]);

  const addReaction = async (messageId, emoji) => {
    try {
      await api.post(`/chat/messages/${messageId}/reactions`, { emoji });
    } catch (error) {
      console.error('Failed to add reaction:', error);
    }
  };

  const removeReaction = async (messageId, emoji) => {
    try {
      await api.delete(`/chat/messages/${messageId}/reactions`, {
        params: { emoji }
      });
    } catch (error) {
      console.error('Failed to remove reaction:', error);
    }
  };

  const updateMessage = async (messageId, newMessage) => {
    try {
      await api.put(`/chat/messages/${messageId}`, { message: newMessage });
      fetchMessages();
    } catch (error) {
      console.error('Failed to update message:', error);
    }
  };

  const deleteMessage = async (messageId) => {
    try {
      await api.delete(`/chat/messages/${messageId}`);
      fetchMessages();
    } catch (error) {
      console.error('Failed to delete message:', error);
    }
  };

  return {
    room,
    messages,
    members,
    loading,
    connected,
    sendMessage,
    sendTyping,
    addReaction,
    removeReaction,
    updateMessage,
    deleteMessage,
    fetchMessages
  };
};

export const useChatRooms = () => {
  const [rooms, setRooms] = useState([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);

  const fetchRooms = async (page = 1, limit = 20) => {
    setLoading(true);
    try {
      const response = await api.get('/chat/rooms', {
        params: { page, limit }
      });
      setRooms(response.data.rooms);
      setTotal(response.data.total);
    } catch (error) {
      console.error('Failed to fetch rooms:', error);
    } finally {
      setLoading(false);
    }
  };

  const createRoom = async (roomData) => {
    try {
      const response = await api.post('/chat/rooms', roomData);
      fetchRooms();
      return response.data;
    } catch (error) {
      console.error('Failed to create room:', error);
      throw error;
    }
  };

  const deleteRoom = async (roomId) => {
    try {
      await api.delete(`/chat/rooms/${roomId}`);
      fetchRooms();
    } catch (error) {
      console.error('Failed to delete room:', error);
    }
  };

  const addMembers = async (roomId, userIds) => {
    try {
      await api.post(`/chat/rooms/${roomId}/members`, { user_ids: userIds });
    } catch (error) {
      console.error('Failed to add members:', error);
      throw error;
    }
  };

  const removeMember = async (roomId, memberId) => {
    try {
      await api.delete(`/chat/rooms/${roomId}/members/${memberId}`);
    } catch (error) {
      console.error('Failed to remove member:', error);
    }
  };

  const getUnreadCount = async () => {
    try {
      const response = await api.get('/chat/unread');
      return response.data.unread_count;
    } catch (error) {
      console.error('Failed to get unread count:', error);
      return 0;
    }
  };

  return {
    rooms,
    total,
    loading,
    fetchRooms,
    createRoom,
    deleteRoom,
    addMembers,
    removeMember,
    getUnreadCount
  };
};
