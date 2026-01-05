import { useState, useEffect, useRef } from 'react';
import * as pdfjsLib from 'pdfjs-dist';
import api from '../lib/api';
import { queueProgressUpdate, syncQueuedUpdates, isOnline } from '../lib/offlineSync';
import { getBookFileUrl } from '../lib/fileService';
import SettingsPanel from './e-reader/SettingsPanel';
import AnnotationsPanel from './e-reader/AnnotationsPanel';

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
  const [noteContent, setNoteContent] = useState('');
  const [activeTab, setActiveTab] = useState('notes');
  const [editingNote, setEditingNote] = useState(null);
  const [isSavingNote, setIsSavingNote] = useState(false);
  const [isDeletingNote, setIsDeletingNote] = useState(null);
  const [isDeletingHighlight, setIsDeletingHighlight] = useState(null);
  const [toast, setToast] = useState(null);

  const canvasRef = useRef(null);
  const saveTimeoutRef = useRef(null);
  const sessionIdRef = useRef(null);
  const sessionStartTimeRef = useRef(null);
  const sessionDurationRef = useRef(0);
  const sessionIntervalRef = useRef(null);

  useEffect(() => {
    loadPDF();
    loadAnnotations();
    startSession();

    return () => {
      if (sessionIntervalRef.current) {
        clearInterval(sessionIntervalRef.current);
      }
      
      if (sessionIdRef.current && sessionDurationRef.current > 0) {
        api.post('/reading/sessions/end', {
          session_id: sessionIdRef.current,
          duration: sessionDurationRef.current
        }).catch(err => console.error('Failed to end session:', err));
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
      sessionIdRef.current = response.data.session.id;
      sessionStartTimeRef.current = Date.now();
      
      sessionIntervalRef.current = setInterval(() => {
        sessionDurationRef.current = Math.floor((Date.now() - sessionStartTimeRef.current) / 1000);
      }, 1000);
      
      console.log('📊 Reading session started:', response.data.session.id);
    } catch (err) {
      console.error('Failed to start reading session:', err);
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

    setIsSavingNote(true);
    try {
      const response = await api.post(`/library/${bookId}/notes`, {
        page: currentPage,
        content: noteContent
      });
      setNotes([...notes, response.data.note]);
      setNoteContent('');
      showToast('Note added successfully');
    } catch (err) {
      console.error('Failed to add note:', err);
      showToast('Failed to add note', 'error');
    } finally {
      setIsSavingNote(false);
    }
  };

  const updateNote = async (noteId, content) => {
    try {
      const response = await api.put(`/library/notes/${noteId}`, { content });
      setNotes(notes.map(n => n.id === noteId ? response.data.note : n));
      setEditingNote(null);
      showToast('Note updated successfully');
    } catch (err) {
      console.error('Failed to update note:', err);
      showToast('Failed to update note', 'error');
    }
  };

  const deleteNote = async (noteId) => {
    setIsDeletingNote(noteId);
    try {
      await api.delete(`/library/notes/${noteId}`);
      setNotes(notes.filter(n => n.id !== noteId));
      showToast('Note deleted successfully');
    } catch (err) {
      console.error('Failed to delete note:', err);
      showToast('Failed to delete note', 'error');
    } finally {
      setIsDeletingNote(null);
    }
  };

  const deleteHighlight = async (highlightId) => {
    setIsDeletingHighlight(highlightId);
    try {
      await api.delete(`/library/highlights/${highlightId}`);
      setHighlights(highlights.filter(h => h.id !== highlightId));
      showToast('Highlight deleted successfully');
    } catch (err) {
      console.error('Failed to delete highlight:', err);
      showToast('Failed to delete highlight', 'error');
    } finally {
      setIsDeletingHighlight(null);
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
    <div className="fixed inset-0 bg-white z-50 flex flex-col">
      {/* Top Bar */}
      <div className="bg-white border-b border-gray-200 px-4 py-3 flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
            title="Back to Library"
          >
            <i className="ri-arrow-left-line text-xl"></i>
          </button>
          <div>
            <h1 className="font-semibold text-gray-900">{bookInfo?.title}</h1>
            <p className="text-sm text-gray-500">Page {currentPage} of {totalPages}</p>
          </div>
        </div>

        <div className="flex items-center space-x-2">
          <button
            onClick={() => setShowAnnotations(!showAnnotations)}
            className={`p-2 rounded-lg transition-colors ${
              showAnnotations ? 'bg-primary-100 text-primary-600' : 'hover:bg-gray-100'
            }`}
            title="Annotations"
          >
            <i className="ri-sticky-note-line text-xl"></i>
            {(notes.length + highlights.length) > 0 && (
              <span className="absolute -top-1 -right-1 bg-primary-600 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                {notes.length + highlights.length}
              </span>
            )}
          </button>
          <button
            onClick={() => setShowSettings(!showSettings)}
            className={`p-2 rounded-lg transition-colors ${
              showSettings ? 'bg-gray-200' : 'hover:bg-gray-100'
            }`}
            title="Settings"
          >
            <i className="ri-settings-3-line text-xl"></i>
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 overflow-auto bg-gray-100 flex items-center justify-center p-4 relative">
        <canvas ref={canvasRef} className="shadow-2xl bg-white" />

        {/* Settings Panel */}
        <SettingsPanel
          show={showSettings}
          onClose={() => setShowSettings(false)}
          scale={scale}
          onScaleChange={setScale}
          readerType="pdf"
        />

        {/* Annotations Panel */}
        <AnnotationsPanel
          show={showAnnotations}
          onClose={() => setShowAnnotations(false)}
          activeTab={activeTab}
          onTabChange={setActiveTab}
          notes={notes}
          highlights={highlights}
          noteContent={noteContent}
          onNoteContentChange={setNoteContent}
          onCreateNote={createNote}
          onUpdateNote={updateNote}
          onDeleteNote={deleteNote}
          onDeleteHighlight={deleteHighlight}
          onGoToPage={goToPage}
          isSavingNote={isSavingNote}
          editingNote={editingNote}
          onEditNote={setEditingNote}
          isDeletingNote={isDeletingNote}
          isDeletingHighlight={isDeletingHighlight}
        />
      </div>

      {/* Bottom Controls */}
      <div className="bg-white border-t border-gray-200 px-4 py-3 flex items-center justify-between">
        <div className="flex items-center space-x-2">
          <button
            onClick={() => goToPage(currentPage - 1)}
            disabled={currentPage === 1}
            className="p-2 hover:bg-gray-100 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            title="Previous Page"
          >
            <i className="ri-arrow-left-s-line text-xl"></i>
          </button>
          <div className="flex items-center space-x-2">
            <input
              type="number"
              value={currentPage}
              onChange={(e) => goToPage(parseInt(e.target.value))}
              className="w-16 px-2 py-1 border border-gray-300 rounded text-center focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              min="1"
              max={totalPages}
            />
            <span className="text-sm text-gray-600">/ {totalPages}</span>
          </div>
          <button
            onClick={() => goToPage(currentPage + 1)}
            disabled={currentPage === totalPages}
            className="p-2 hover:bg-gray-100 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            title="Next Page"
          >
            <i className="ri-arrow-right-s-line text-xl"></i>
          </button>
        </div>

        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-2">
            <button
              onClick={() => setScale(Math.max(0.5, scale - 0.25))}
              className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
              title="Zoom Out"
            >
              <i className="ri-zoom-out-line text-xl"></i>
            </button>
            <span className="text-sm text-gray-600 min-w-[50px] text-center">{Math.round(scale * 100)}%</span>
            <button
              onClick={() => setScale(Math.min(3, scale + 0.25))}
              className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
              title="Zoom In"
            >
              <i className="ri-zoom-in-line text-xl"></i>
            </button>
          </div>

          <div className="h-6 w-px bg-gray-300"></div>

          <div className="text-sm text-gray-600">
            {Math.round((currentPage / totalPages) * 100)}% Complete
          </div>
        </div>
      </div>

      {/* Toast Notification */}
      {toast && (
        <div className={`fixed bottom-4 right-4 px-6 py-3 rounded-lg shadow-lg ${
          toast.type === 'error' ? 'bg-red-500' : 'bg-green-500'
        } text-white animate-fade-in-up z-50`}>
          <div className="flex items-center gap-2">
            <i className={`${toast.type === 'error' ? 'ri-error-warning-line' : 'ri-check-line'} text-xl`}></i>
            {toast.message}
          </div>
        </div>
      )}
    </div>
  );
}
