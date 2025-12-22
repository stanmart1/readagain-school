import { API_BASE_URL } from './api';

const UPLOAD_API_URL = import.meta.env.VITE_UPLOAD_API_URL || 'http://localhost:8001';

/**
 * Upload book cover image
 * @param {File} file - Image file to upload
 * @returns {Promise<{filename: string, path: string, url: string, size: number}>}
 */
export const uploadCover = async (file) => {
  const formData = new FormData();
  formData.append('file', file);

  const response = await fetch(`${UPLOAD_API_URL}/api/upload/cover`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to upload cover');
  }

  return response.json();
};

/**
 * Upload book file (PDF, EPUB, HTML)
 * @param {File} file - Book file to upload
 * @returns {Promise<{filename: string, path: string, url: string, size: number}>}
 */
export const uploadBook = async (file) => {
  const formData = new FormData();
  formData.append('file', file);

  const response = await fetch(`${UPLOAD_API_URL}/api/upload/book`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to upload book');
  }

  return response.json();
};

/**
 * Upload profile picture
 * @param {File} file - Image file to upload
 * @returns {Promise<{filename: string, path: string, url: string, size: number}>}
 */
export const uploadProfile = async (file) => {
  const formData = new FormData();
  formData.append('file', file);

  const response = await fetch(`${UPLOAD_API_URL}/api/upload/profile`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to upload profile picture');
  }

  return response.json();
};

/**
 * Delete file from upload API
 * @param {string} filename - Filename to delete
 * @returns {Promise<{message: string}>}
 */
export const deleteFile = async (filename) => {
  const response = await fetch(`${UPLOAD_API_URL}/api/files/${filename}`, {
    method: 'DELETE',
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to delete file');
  }

  return response.json();
};

/**
 * Get full URL for file serving from upload API
 * @param {string} filename - Filename or path (e.g., 'abc-123.jpg' or 'covers/abc-123.jpg')
 * @returns {string} Full URL to access the file
 */
export const getFileUrl = (filename) => {
  if (!filename) return null;
  
  // If already a full URL, return as is
  if (filename.startsWith('http://') || filename.startsWith('https://')) {
    return filename;
  }
  
  // Extract just the filename if path is provided
  const cleanFilename = filename.includes('/') ? filename.split('/').pop() : filename;
  
  // Return full URL to upload API
  return `${UPLOAD_API_URL}/api/files/${cleanFilename}`;
};

/**
 * Get image URL with fallback
 * @param {string} filename - Filename or path from backend
 * @param {string} fallback - Fallback image path (default: '/placeholder-book.png')
 * @returns {string} Image URL or fallback
 */
export const getImageUrl = (filename, fallback = '/placeholder-book.png') => {
  return getFileUrl(filename) || fallback;
};

