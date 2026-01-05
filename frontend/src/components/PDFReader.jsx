import { useState, useEffect, useRef } from 'react';
import { motion } from 'framer-motion';
import * as pdfjsLib from 'pdfjs-dist';
import api from '../lib/api';
import { queueProgressUpdate, syncQueuedUpdates, isOnline } from '../lib/offlineSync';
import { getBookFileUrl } from '../lib/fileService';

pdfjsLib.GlobalWorkerOptions.workerSrc = `//cdnjs.cloudflare.com/ajax/libs/pdf.js/${pdfjsLib.version}/pdf.worker.min.js`;

export default function PDFReader({ bookId, onClose }) {
  const [pdfDoc, setPdfDoc] = useState(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const [scale, setScale] = useState(1.5);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [bookInfo, setBookInfo] = useState(null);
  const [showSettings, setShowSettings] = useState(false);
  const [showAnnotations, setShowAnnotations] = useState(false);
  const [highlights, setHighlights] = useState([]);
  const [notes, setNotes] = useState([]);
  const [showNoteInput, setShowNoteInput] = useState(false);
  const [noteContent, setNoteContent] = useState('');
  const [activeTab, setActiveTab] = useState('notes');
  const [editingNote, setEditingNote] = useState(null);
  const [toast, setToast] = useState(null);

  const canvasRef = useRef(null);
  const saveTimeoutRef = useRef(null);
  const sessionId = useRef(null);
  const sessionStartTime = useRef(null);
  const sessionDuration = useRef(0);
  const sessionInterval = useRef(null);

  useEffect(() => {
    loadPDF();
    loadAnnotations();
    startSession();

    return () => {
      if (sessionInterval.current) clearInterval(sessionInterval.current);
      if (sessionId.current && sessionDuration.current > 0) {
        api.post('/reading/sessions/end', {
          session_id: sessionId.current,
          duration: sessionDuration.current
        }).catch(console.error);
      }
    };
  }, [bookId]);

  useEffect(() => {
    if (pdfDoc && currentPage) {
      renderPage(currentPage);
      saveProgress();
    }
  }, [currentPage, scale, pdfDoc]);

  const startSession = async () => {
    try {
      const response = await api.post('/reading/sessions/start', {
        book_id: parseInt(bookId)
      });
      sessionId.current = response.data.session.id;
      sessionStartTime.current = Date.now();
      
      sessionInterval.current = setInterval(() => {
        sessionDuration.current = Math.floor((Date.now() - sessionStartTime.current) / 1000);
      }, 1000);
    } catch (err) {
      console.error('Failed to start session:', err);
    }
  };

  const loadPDF = async () => {
    try {
      setLoading(true);
      const response = await api.get('/library');
      const libraryItem = response.data.library.find(item => item.book_id === parseInt(bookId));

      if (!libraryItem) {
        throw new Error('Book not found in your library');
      }

      const bookData = libraryItem.book || libraryItem;
      setBookInfo(bookData);
      setCurrentPage(libraryItem.current_page || 1);

      const fileUrl = getBookFileUrl(bookData.file_path);
      const loadingTask = pdfjsLib.getDocument(fileUrl);
      const pdf = await loadingTask.promise;
      
      setPdfDoc(pdf);
      setTotalPages(pdf.numPages);
      setLoading(false);
    } catch (err) {
      console.error('Error loading PDF:', err);
      setError(err.message || 'Failed to load PDF');
      setLoading(false);
    }
  };

  const renderPage = async (pageNum) => {
    if (!pdfDoc || !canvasRef.current) return;

    try {
      const page = await pdfDoc.getPage(pageNum);
      const viewport = page.getViewport({ scale });
      const canvas = canvasRef.current;
      const context = canvas.getContext('2d');

      canvas.height = viewport.height;
      canvas.width = viewport.width;

      await page.render({
        canvasContext: context,
        viewport: viewport
      }).promise;
    } catch (err) {
      console.error('Error rendering page:', err);
    }
  };

  const saveProgress = () => {
    if (saveTimeoutRef.current) clearTimeout(saveTimeoutRef.current);

    saveTimeoutRef.current = setTimeout(async () => {
      const progress = (currentPage / totalPages) * 100;
      
      try {
        if (isOnline()) {
          await api.put(`/library/${bookId}/progress`, {
            current_page: currentPage,
            total_pages: totalPages,
            progress: progress
          });
        } else {
          queueProgressUpdate(bookId, currentPage, totalPages, progress);
        }
      } catch (err) {
        console.error('Failed to save progress:', err);
      }
    }, 2000);
  };

  const loadAnnotations = async () => {
    try {
      const [notesRes, highlightsRes] = await Promise.all([
        api.get(`/library/${bookId}/notes`),
        api.get(`/library/${bookId}/highlights`)
      ]);
      setNotes(notesRes.data.notes || []);
      setHighlights(highlightsRes.data.highlights || []);
    } catch (err) {
      console.error('Failed to load annotations:', err);
    }
  };

  const createNote = async () => {
    if (!noteContent.trim()) return;

    try {
      const response = await api.post(`/library/${bookId}/notes`, {
        page: currentPage,
        content: noteContent
      });
      setNotes([...notes, response.data.note]);
      setNoteContent('');
      setShowNoteInput(false);
      showToast('Note added successfully');
    } catch (err) {
      showToast('Failed to add note', 'error');
    }
  };

  const updateNote = async (noteId, content) => {
    try {
      const response = await api.put(`/library/notes/${noteId}`, { content });
      setNotes(notes.map(n => n.id === noteId ? response.data.note : n));
      setEditingNote(null);
      showToast('Note updated successfully');
    } catch (err) {
      showToast('Failed to update note', 'error');
    }
  };

  const deleteNote = async (noteId) => {
    try {
      await api.delete(`/library/notes/${noteId}`);
      setNotes(notes.filter(n => n.id !== noteId));
      showToast('Note deleted successfully');
    } catch (err) {
      showToast('Failed to delete note', 'error');
    }
  };

  const showToast = (message, type = 'success') => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 3000);
  };

  const goToPage = (page) => {
    if (page >= 1 && page <= totalPages) {
      setCurrentPage(page);
    }
  };

  if (loading) {
    return (
      <div className="fixed inset-0 bg-gray-900 flex items-center justify-center z-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-primary-600 mx-auto mb-4"></div>
          <p className="text-white text-lg">Loading PDF...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="fixed inset-0 bg-gray-900 flex items-center justify-center z-50">
        <div className="bg-white rounded-xl p-8 max-w-md mx-4">
          <i className="ri-error-warning-line text-5xl text-red-500 mb-4"></i>
          <h3 className="text-xl font-bold mb-2">Error Loading PDF</h3>
          <p className="text-gray-600 mb-6">{error}</p>
          <button onClick={onClose} className="px-6 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700">
            Back to Library
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 bg-gray-900 z-50 flex flex-col">
      {/* Header */}
      <div className="bg-gray-800 text-white px-4 py-3 flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <button onClick={onClose} className="p-2 hover:bg-gray-700 rounded-lg">
            <i className="ri-arrow-left-line text-xl"></i>
          </button>
          <div>
            <h1 className="font-semibold">{bookInfo?.title}</h1>
            <p className="text-sm text-gray-400">Page {currentPage} of {totalPages}</p>
          </div>
        </div>

        <div className="flex items-center space-x-2">
          <button onClick={() => setShowAnnotations(!showAnnotations)} className="p-2 hover:bg-gray-700 rounded-lg">
            <i className="ri-sticky-note-line text-xl"></i>
          </button>
          <button onClick={() => setShowSettings(!showSettings)} className="p-2 hover:bg-gray-700 rounded-lg">
            <i className="ri-settings-3-line text-xl"></i>
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 overflow-auto flex items-center justify-center p-4">
        <canvas ref={canvasRef} className="shadow-2xl" />
      </div>

      {/* Controls */}
      <div className="bg-gray-800 text-white px-4 py-3 flex items-center justify-between">
        <div className="flex items-center space-x-2">
          <button onClick={() => goToPage(currentPage - 1)} disabled={currentPage === 1} className="p-2 hover:bg-gray-700 rounded-lg disabled:opacity-50">
            <i className="ri-arrow-left-s-line text-xl"></i>
          </button>
          <input
            type="number"
            value={currentPage}
            onChange={(e) => goToPage(parseInt(e.target.value))}
            className="w-16 px-2 py-1 bg-gray-700 rounded text-center"
            min="1"
            max={totalPages}
          />
          <button onClick={() => goToPage(currentPage + 1)} disabled={currentPage === totalPages} className="p-2 hover:bg-gray-700 rounded-lg disabled:opacity-50">
            <i className="ri-arrow-right-s-line text-xl"></i>
          </button>
        </div>

        <div className="flex items-center space-x-2">
          <button onClick={() => setScale(Math.max(0.5, scale - 0.25))} className="p-2 hover:bg-gray-700 rounded-lg">
            <i className="ri-zoom-out-line text-xl"></i>
          </button>
          <span className="text-sm">{Math.round(scale * 100)}%</span>
          <button onClick={() => setScale(Math.min(3, scale + 0.25))} className="p-2 hover:bg-gray-700 rounded-lg">
            <i className="ri-zoom-in-line text-xl"></i>
          </button>
        </div>

        <button onClick={() => setShowNoteInput(true)} className="px-4 py-2 bg-primary-600 hover:bg-primary-700 rounded-lg">
          <i className="ri-add-line mr-2"></i>Add Note
        </button>
      </div>

      {/* Settings Panel */}
      {showSettings && (
        <motion.div initial={{ x: 300 }} animate={{ x: 0 }} className="fixed right-0 top-0 h-full w-80 bg-white shadow-2xl z-50 overflow-y-auto">
          <div className="p-6">
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-xl font-bold">Settings</h3>
              <button onClick={() => setShowSettings(false)} className="p-2 hover:bg-gray-100 rounded-lg">
                <i className="ri-close-line text-xl"></i>
              </button>
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-2">Zoom Level</label>
                <input
                  type="range"
                  min="50"
                  max="300"
                  value={scale * 100}
                  onChange={(e) => setScale(e.target.value / 100)}
                  className="w-full"
                />
                <p className="text-sm text-gray-600 mt-1">{Math.round(scale * 100)}%</p>
              </div>
            </div>
          </div>
        </motion.div>
      )}

      {/* Annotations Panel */}
      {showAnnotations && (
        <motion.div initial={{ x: 300 }} animate={{ x: 0 }} className="fixed right-0 top-0 h-full w-80 bg-white shadow-2xl z-50 overflow-y-auto">
          <div className="p-6">
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-xl font-bold">Annotations</h3>
              <button onClick={() => setShowAnnotations(false)} className="p-2 hover:bg-gray-100 rounded-lg">
                <i className="ri-close-line text-xl"></i>
              </button>
            </div>

            <div className="flex space-x-2 mb-4">
              <button onClick={() => setActiveTab('notes')} className={`flex-1 py-2 rounded-lg ${activeTab === 'notes' ? 'bg-primary-600 text-white' : 'bg-gray-100'}`}>
                Notes ({notes.length})
              </button>
              <button onClick={() => setActiveTab('highlights')} className={`flex-1 py-2 rounded-lg ${activeTab === 'highlights' ? 'bg-primary-600 text-white' : 'bg-gray-100'}`}>
                Highlights ({highlights.length})
              </button>
            </div>

            {activeTab === 'notes' && (
              <div className="space-y-3">
                {notes.map(note => (
                  <div key={note.id} className="p-3 bg-gray-50 rounded-lg">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-xs text-gray-500">Page {note.page}</span>
                      <div className="flex space-x-1">
                        <button onClick={() => goToPage(note.page)} className="p-1 hover:bg-gray-200 rounded">
                          <i className="ri-eye-line text-sm"></i>
                        </button>
                        <button onClick={() => setEditingNote(note)} className="p-1 hover:bg-gray-200 rounded">
                          <i className="ri-edit-line text-sm"></i>
                        </button>
                        <button onClick={() => deleteNote(note.id)} className="p-1 hover:bg-gray-200 rounded text-red-600">
                          <i className="ri-delete-bin-line text-sm"></i>
                        </button>
                      </div>
                    </div>
                    {editingNote?.id === note.id ? (
                      <div>
                        <textarea
                          value={editingNote.content}
                          onChange={(e) => setEditingNote({ ...editingNote, content: e.target.value })}
                          className="w-full p-2 border rounded text-sm"
                          rows="3"
                        />
                        <div className="flex space-x-2 mt-2">
                          <button onClick={() => updateNote(note.id, editingNote.content)} className="px-3 py-1 bg-primary-600 text-white rounded text-sm">
                            Save
                          </button>
                          <button onClick={() => setEditingNote(null)} className="px-3 py-1 bg-gray-200 rounded text-sm">
                            Cancel
                          </button>
                        </div>
                      </div>
                    ) : (
                      <p className="text-sm">{note.content}</p>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>
        </motion.div>
      )}

      {/* Add Note Modal */}
      {showNoteInput && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 max-w-md w-full mx-4">
            <h3 className="text-xl font-bold mb-4">Add Note - Page {currentPage}</h3>
            <textarea
              value={noteContent}
              onChange={(e) => setNoteContent(e.target.value)}
              placeholder="Write your note..."
              className="w-full p-3 border rounded-lg resize-none"
              rows="4"
              autoFocus
            />
            <div className="flex space-x-3 mt-4">
              <button onClick={createNote} className="flex-1 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700">
                Save Note
              </button>
              <button onClick={() => { setShowNoteInput(false); setNoteContent(''); }} className="flex-1 py-2 bg-gray-200 rounded-lg hover:bg-gray-300">
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Toast */}
      {toast && (
        <motion.div initial={{ y: 50, opacity: 0 }} animate={{ y: 0, opacity: 1 }} className={`fixed bottom-4 right-4 px-6 py-3 rounded-lg shadow-lg ${toast.type === 'error' ? 'bg-red-500' : 'bg-green-500'} text-white`}>
          {toast.message}
        </motion.div>
      )}
    </div>
  );
}
