export default function AnnotationsPanel({
  show,
  onClose,
  activeTab,
  onTabChange,
  notes,
  highlights,
  noteContent,
  onNoteContentChange,
  onCreateNote,
  onUpdateNote,
  onDeleteNote,
  onDeleteHighlight,
  onGoToPage,
  onGoToHighlight,
  isSavingNote,
  editingNote,
  onEditNote,
  isDeletingNote,
  isDeletingHighlight
}) {
  if (!show) return null;

  return (
    <div className="absolute top-0 right-0 bottom-0 w-96 bg-white shadow-2xl z-10 flex flex-col">
      {/* Header */}
      <div className="p-4 border-b border-gray-200 flex items-center justify-between bg-white">
        <h2 className="text-lg font-bold">Annotations</h2>
        <button
          onClick={onClose}
          className="p-2 hover:bg-gray-100 rounded-lg"
        >
          <i className="ri-close-line text-xl"></i>
        </button>
      </div>

      {/* Tabs */}
      <div className="flex border-b border-gray-200 bg-white">
        <button
          onClick={() => onTabChange('notes')}
          className={`flex-1 px-4 py-3 font-semibold transition-colors relative ${
            activeTab === 'notes'
              ? 'text-primary-600 bg-primary-50'
              : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
          }`}
        >
          <div className="flex items-center justify-center gap-2">
            <i className="ri-sticky-note-line"></i>
            Notes
            {notes.length > 0 && (
              <span className={`px-2 py-0.5 rounded-full text-xs ${
                activeTab === 'notes' ? 'bg-primary-600 text-white' : 'bg-gray-200 text-gray-600'
              }`}>
                {notes.length}
              </span>
            )}
          </div>
          {activeTab === 'notes' && (
            <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-primary-600"></div>
          )}
        </button>
        <button
          onClick={() => onTabChange('highlights')}
          className={`flex-1 px-4 py-3 font-semibold transition-colors relative ${
            activeTab === 'highlights'
              ? 'text-yellow-600 bg-yellow-50'
              : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
          }`}
        >
          <div className="flex items-center justify-center gap-2">
            <i className="ri-mark-pen-line"></i>
            Highlights
            {highlights.length > 0 && (
              <span className={`px-2 py-0.5 rounded-full text-xs ${
                activeTab === 'highlights' ? 'bg-yellow-600 text-white' : 'bg-gray-200 text-gray-600'
              }`}>
                {highlights.length}
              </span>
            )}
          </div>
          {activeTab === 'highlights' && (
            <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-yellow-600"></div>
          )}
        </button>
      </div>

      {/* Tab Content */}
      <div className="flex-1 overflow-y-auto">
        {/* Notes Tab */}
        {activeTab === 'notes' && (
          <div>
            {/* Add Note Section */}
            <div className="p-4 border-b border-gray-200 bg-primary-50">
              <h3 className="font-semibold mb-2 text-sm text-gray-700">Add New Note</h3>
              <textarea
                value={noteContent}
                onChange={(e) => onNoteContentChange(e.target.value)}
                placeholder="Write your note here..."
                className="w-full p-3 border border-gray-300 rounded-lg resize-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                rows="3"
              />
              <button
                onClick={onCreateNote}
                disabled={!noteContent.trim() || isSavingNote}
                className="mt-2 w-full px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:bg-gray-300 disabled:cursor-not-allowed flex items-center justify-center gap-2"
              >
                {isSavingNote ? (
                  <>
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                    Saving...
                  </>
                ) : (
                  <>
                    <i className="ri-save-line"></i>
                    Save Note
                  </>
                )}
              </button>
            </div>

            {/* Notes List */}
            <div className="p-4 space-y-3">
              {notes.length === 0 ? (
                <div className="text-center py-8 text-gray-500">
                  <i className="ri-sticky-note-line text-4xl mb-2"></i>
                  <p>No notes yet</p>
                </div>
              ) : (
                notes.map(note => (
                  <div key={note.id} className="p-3 bg-gray-50 rounded-lg border border-gray-200 hover:border-primary-300 transition-colors">
                    {editingNote?.id === note.id ? (
                      <div>
                        <textarea
                          value={editingNote.content}
                          onChange={(e) => onEditNote({ ...editingNote, content: e.target.value })}
                          className="w-full p-2 border border-gray-300 rounded-lg text-sm resize-none focus:ring-2 focus:ring-primary-500"
                          rows="3"
                        />
                        <div className="flex gap-2 mt-2">
                          <button
                            onClick={() => onUpdateNote(note.id, editingNote.content)}
                            className="flex-1 px-3 py-1.5 bg-primary-600 text-white rounded-lg text-sm hover:bg-primary-700"
                          >
                            Save
                          </button>
                          <button
                            onClick={() => onEditNote(null)}
                            className="flex-1 px-3 py-1.5 bg-gray-200 text-gray-700 rounded-lg text-sm hover:bg-gray-300"
                          >
                            Cancel
                          </button>
                        </div>
                      </div>
                    ) : (
                      <>
                        <div className="flex items-center justify-between mb-2">
                          <span className="text-xs text-gray-500 font-medium">Page {note.page}</span>
                          <div className="flex gap-1">
                            {onGoToPage && (
                              <button
                                onClick={() => onGoToPage(note.page)}
                                className="p-1.5 hover:bg-gray-200 rounded transition-colors"
                                title="Go to page"
                              >
                                <i className="ri-eye-line text-sm text-gray-600"></i>
                              </button>
                            )}
                            <button
                              onClick={() => onEditNote(note)}
                              className="p-1.5 hover:bg-gray-200 rounded transition-colors"
                              title="Edit note"
                            >
                              <i className="ri-edit-line text-sm text-gray-600"></i>
                            </button>
                            <button
                              onClick={() => onDeleteNote(note.id)}
                              disabled={isDeletingNote === note.id}
                              className="p-1.5 hover:bg-red-50 rounded transition-colors disabled:opacity-50"
                              title="Delete note"
                            >
                              {isDeletingNote === note.id ? (
                                <div className="animate-spin rounded-full h-3 w-3 border-b-2 border-red-600"></div>
                              ) : (
                                <i className="ri-delete-bin-line text-sm text-red-600"></i>
                              )}
                            </button>
                          </div>
                        </div>
                        <p className="text-sm text-gray-700 whitespace-pre-wrap">{note.content}</p>
                      </>
                    )}
                  </div>
                ))
              )}
            </div>
          </div>
        )}

        {/* Highlights Tab */}
        {activeTab === 'highlights' && (
          <div className="p-4 space-y-3">
            {highlights.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <i className="ri-mark-pen-line text-4xl mb-2"></i>
                <p>No highlights yet</p>
              </div>
            ) : (
              highlights.map(highlight => (
                <div key={highlight.id} className="p-3 bg-yellow-50 rounded-lg border-l-4 border-yellow-400">
                  <div className="flex items-start justify-between mb-2">
                    <div className={`w-4 h-4 rounded-full`} style={{ backgroundColor: highlight.color || '#fbbf24' }}></div>
                    <div className="flex gap-1">
                      {onGoToHighlight && (
                        <button
                          onClick={() => onGoToHighlight(highlight)}
                          className="p-1.5 hover:bg-yellow-100 rounded transition-colors"
                          title="Go to highlight"
                        >
                          <i className="ri-eye-line text-sm text-gray-600"></i>
                        </button>
                      )}
                      <button
                        onClick={() => onDeleteHighlight(highlight.id)}
                        disabled={isDeletingHighlight === highlight.id}
                        className="p-1.5 hover:bg-red-50 rounded transition-colors disabled:opacity-50"
                        title="Delete highlight"
                      >
                        {isDeletingHighlight === highlight.id ? (
                          <div className="animate-spin rounded-full h-3 w-3 border-b-2 border-red-600"></div>
                        ) : (
                          <i className="ri-delete-bin-line text-sm text-red-600"></i>
                        )}
                      </button>
                    </div>
                  </div>
                  <p className="text-sm text-gray-700 italic">"{highlight.text}"</p>
                </div>
              ))
            )}
          </div>
        )}
      </div>
    </div>
  );
}
