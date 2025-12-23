import { useState, useEffect } from 'react';
import api from '../../lib/api';

const ReadingDetailsModal = ({ libraryId, onClose }) => {
  const [assignmentDetails, setAssignmentDetails] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadDetails = async () => {
      setLoading(true);
      try {
        const [detailsRes, analyticsRes] = await Promise.all([
          api.get(`/admin/library-assignment/${libraryId}/details`),
          api.get(`/admin/library-assignment/${libraryId}/analytics`)
        ]);
        setAssignmentDetails({
          ...detailsRes.data,
          analytics: analyticsRes.data
        });
      } catch (error) {
        console.error('Failed to load details:', error);
      } finally {
        setLoading(false);
      }
    };

    if (libraryId) {
      loadDetails();
    }
  }, [libraryId]);

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl shadow-2xl max-w-4xl w-full max-h-[90vh] overflow-hidden">
        {/* Header */}
        <div className="px-6 py-4 border-b border-gray-200 bg-gradient-to-r from-blue-50 to-indigo-50">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-gradient-to-r from-blue-600 to-indigo-600 rounded-full flex items-center justify-center">
                <i className="ri-book-open-line text-white text-lg"></i>
              </div>
              <div>
                <h3 className="text-xl font-bold text-gray-900">Reading Details</h3>
                <p className="text-sm text-gray-600">Notes, highlights, and progress</p>
              </div>
            </div>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-gray-600 transition-colors p-1 rounded-full hover:bg-white/50"
            >
              <i className="ri-close-line text-2xl"></i>
            </button>
          </div>
        </div>

        {/* Content */}
        <div className="p-6 overflow-y-auto max-h-[calc(90vh-140px)]">
          {loading ? (
            <div className="flex items-center justify-center py-12">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
            </div>
          ) : assignmentDetails ? (
            <div className="space-y-6">
              {/* Assignment Info */}
              <div className="bg-gray-50 rounded-xl p-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-sm text-gray-600">User</p>
                    <p className="font-semibold text-gray-900">{assignmentDetails.assignment.user_name}</p>
                    <p className="text-sm text-gray-500">{assignmentDetails.assignment.user_email}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Book</p>
                    <p className="font-semibold text-gray-900">{assignmentDetails.assignment.book_title}</p>
                    <p className="text-sm text-gray-500">{assignmentDetails.assignment.book_author}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Progress</p>
                    <p className="font-semibold text-gray-900">{assignmentDetails.assignment.progress.toFixed(2)}%</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Status</p>
                    <span className={`inline-block px-2 py-1 text-xs font-semibold rounded-full ${
                      assignmentDetails.assignment.status === 'reading' ? 'bg-green-100 text-green-800' :
                      assignmentDetails.assignment.status === 'completed' ? 'bg-blue-100 text-blue-800' :
                      'bg-gray-100 text-gray-800'
                    }`}>
                      {assignmentDetails.assignment.status}
                    </span>
                  </div>
                </div>
              </div>

              {/* Reading Analytics */}
              {assignmentDetails.analytics && (
                <div>
                  <h4 className="text-lg font-bold text-gray-900 mb-3 flex items-center gap-2">
                    <i className="ri-line-chart-line text-blue-600"></i>
                    Reading Analytics
                  </h4>
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                    <div className="bg-white rounded-lg shadow-md p-3">
                      <div className="flex items-center gap-2 mb-1">
                        <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-blue-600 rounded-lg flex items-center justify-center">
                          <i className="ri-time-line text-white text-sm"></i>
                        </div>
                        <div>
                          <p className="text-xs text-gray-600">Total Time</p>
                          <p className="text-sm font-bold text-gray-900">{assignmentDetails.analytics.total_reading_time || '0h'}</p>
                        </div>
                      </div>
                    </div>
                    <div className="bg-white rounded-lg shadow-md p-3">
                      <div className="flex items-center gap-2 mb-1">
                        <div className="w-8 h-8 bg-gradient-to-br from-green-500 to-green-600 rounded-lg flex items-center justify-center">
                          <i className="ri-calendar-line text-white text-sm"></i>
                        </div>
                        <div>
                          <p className="text-xs text-gray-600">Sessions</p>
                          <p className="text-sm font-bold text-gray-900">{assignmentDetails.analytics.total_sessions || 0}</p>
                        </div>
                      </div>
                    </div>
                    <div className="bg-white rounded-lg shadow-md p-3">
                      <div className="flex items-center gap-2 mb-1">
                        <div className="w-8 h-8 bg-gradient-to-br from-purple-500 to-purple-600 rounded-lg flex items-center justify-center">
                          <i className="ri-speed-line text-white text-sm"></i>
                        </div>
                        <div>
                          <p className="text-xs text-gray-600">Avg Session</p>
                          <p className="text-sm font-bold text-gray-900">{assignmentDetails.analytics.avg_session_time || '0m'}</p>
                        </div>
                      </div>
                    </div>
                    <div className="bg-white rounded-lg shadow-md p-3">
                      <div className="flex items-center gap-2 mb-1">
                        <div className="w-8 h-8 bg-gradient-to-br from-yellow-500 to-yellow-600 rounded-lg flex items-center justify-center">
                          <i className="ri-fire-line text-white text-sm"></i>
                        </div>
                        <div>
                          <p className="text-xs text-gray-600">Streak</p>
                          <p className="text-sm font-bold text-gray-900">{assignmentDetails.analytics.reading_streak || 0} days</p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {/* Reading Goals */}
              {assignmentDetails.analytics?.goals && assignmentDetails.analytics.goals.length > 0 && (
                <div>
                  <h4 className="text-lg font-bold text-gray-900 mb-3 flex items-center gap-2">
                    <i className="ri-flag-line text-green-600"></i>
                    Reading Goals ({assignmentDetails.analytics.goals.length})
                  </h4>
                  <div className="space-y-3">
                    {assignmentDetails.analytics.goals.map((goal) => (
                      <div key={goal.id} className="bg-white rounded-lg shadow-md p-4">
                        <div className="flex items-start justify-between mb-2">
                          <div className="flex-1">
                            <h5 className="font-semibold text-gray-900">{goal.title}</h5>
                            <p className="text-sm text-gray-600 mt-1">{goal.description}</p>
                          </div>
                          <span className={`px-2 py-1 text-xs font-semibold rounded-full ${
                            goal.status === 'completed' ? 'bg-green-100 text-green-800' :
                            goal.status === 'in_progress' ? 'bg-blue-100 text-blue-800' :
                            'bg-gray-100 text-gray-800'
                          }`}>
                            {goal.status === 'completed' ? 'Completed' : goal.status === 'in_progress' ? 'In Progress' : 'Not Started'}
                          </span>
                        </div>
                        <div className="space-y-2">
                          <div className="flex items-center justify-between text-sm">
                            <span className="text-gray-600">Progress</span>
                            <span className="font-semibold text-gray-900">{goal.current_value || 0} / {goal.target_value}</span>
                          </div>
                          <div className="w-full bg-gray-200 rounded-full h-2">
                            <div
                              className={`h-2 rounded-full ${
                                goal.status === 'completed' ? 'bg-green-500' : 'bg-blue-500'
                              }`}
                              style={{ width: `${Math.min((goal.current_value / goal.target_value) * 100, 100)}%` }}
                            ></div>
                          </div>
                          <div className="flex items-center justify-between text-xs text-gray-500">
                            <span>Started: {new Date(goal.start_date).toLocaleDateString()}</span>
                            <span>Target: {new Date(goal.target_date).toLocaleDateString()}</span>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Highlights Section */}
              <div>
                <div className="flex items-center justify-between mb-3">
                  <h4 className="text-lg font-bold text-gray-900 flex items-center gap-2">
                    <i className="ri-mark-pen-line text-yellow-600"></i>
                    Highlights ({assignmentDetails.highlights.length})
                  </h4>
                </div>
                {assignmentDetails.highlights.length === 0 ? (
                  <div className="text-center py-8 bg-gray-50 rounded-xl">
                    <i className="ri-mark-pen-line text-4xl text-gray-300 mb-2"></i>
                    <p className="text-gray-500">No highlights yet</p>
                  </div>
                ) : (
                  <div className="space-y-3">
                    {assignmentDetails.highlights.slice(0, 3).map((highlight) => (
                      <div key={highlight.id} className="bg-yellow-50 border-l-4 border-yellow-400 p-4 rounded-r-lg">
                        <p className="text-gray-800 italic mb-2">
                          "{highlight.text.length > 100 ? highlight.text.substring(0, 100) + '...' : highlight.text}"
                        </p>
                        <div className="flex items-center justify-between text-xs text-gray-500">
                          <span className="flex items-center gap-1">
                            <div className="w-3 h-3 rounded-full" style={{backgroundColor: highlight.color}}></div>
                            {highlight.color}
                          </span>
                          <span>{new Date(highlight.created_at).toLocaleString()}</span>
                        </div>
                      </div>
                    ))}
                    {assignmentDetails.highlights.length > 3 && (
                      <p className="text-sm text-gray-500 text-center py-2">
                        Showing 3 of {assignmentDetails.highlights.length} highlights
                      </p>
                    )}
                  </div>
                )}
              </div>

              {/* Notes Section */}
              <div>
                <div className="flex items-center justify-between mb-3">
                  <h4 className="text-lg font-bold text-gray-900 flex items-center gap-2">
                    <i className="ri-sticky-note-line text-blue-600"></i>
                    Notes ({assignmentDetails.notes?.length || 0})
                  </h4>
                </div>
                {!assignmentDetails.notes || assignmentDetails.notes.length === 0 ? (
                  <div className="text-center py-8 bg-gray-50 rounded-xl">
                    <i className="ri-sticky-note-line text-4xl text-gray-300 mb-2"></i>
                    <p className="text-gray-500">No notes yet</p>
                  </div>
                ) : (
                  <div className="space-y-3">
                    {assignmentDetails.notes.slice(0, 3).map((note) => (
                      <div key={note.id} className="bg-blue-50 border-l-4 border-blue-400 p-4 rounded-r-lg">
                        <p className="text-gray-800 mb-2">
                          {note.content.length > 150 ? note.content.substring(0, 150) + '...' : note.content}
                        </p>
                        <div className="flex items-center justify-between text-xs text-gray-500">
                          <span>Page {note.page || 0}</span>
                          <span>{new Date(note.created_at).toLocaleString()}</span>
                        </div>
                      </div>
                    ))}
                    {assignmentDetails.notes.length > 3 && (
                      <p className="text-sm text-gray-500 text-center py-2">
                        Showing 3 of {assignmentDetails.notes.length} notes
                      </p>
                    )}
                  </div>
                )}
              </div>
            </div>
          ) : (
            <div className="text-center py-12">
              <p className="text-gray-500">Failed to load details</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ReadingDetailsModal;
