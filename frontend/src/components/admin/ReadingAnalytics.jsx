import { useState, useEffect, useMemo, useRef } from 'react';
import { useReadingAnalytics } from '../../hooks/useReadingAnalytics';
import ReadingDetailsModal from './ReadingDetailsModal';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell
} from 'recharts';

const ReadingAnalytics = () => {
  const { analyticsData, loading, error, fetchReadingAnalytics } = useReadingAnalytics();
  const [selectedPeriod, setSelectedPeriod] = useState('month');
  const [showDetailsModal, setShowDetailsModal] = useState(false);
  const [selectedLibraryId, setSelectedLibraryId] = useState(null);
  const hasFetched = useRef(false);

  useEffect(() => {
    if (!hasFetched.current) {
      fetchReadingAnalytics(selectedPeriod);
      hasFetched.current = true;
    } else {
      fetchReadingAnalytics(selectedPeriod);
    }
  }, [selectedPeriod]);

  const handleRefresh = () => {
    fetchReadingAnalytics(selectedPeriod);
  };

  const safeData = useMemo(() => ({
    overview: analyticsData?.overview || {},
    classStats: analyticsData?.class_stats || [],
    strugglingReaders: analyticsData?.struggling_readers || [],
    mostReadBooks: analyticsData?.most_read_books || [],
    topReaders: analyticsData?.top_readers || [],
    activeReaders: analyticsData?.active_readers || []
  }), [analyticsData]);

  if (loading && !analyticsData) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="flex items-center justify-center h-64">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
            <p className="mt-4 text-gray-600">Loading reading analytics...</p>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="text-center text-red-600">
          <i className="ri-error-warning-line text-4xl mb-4"></i>
          <p>{error}</p>
          <button
            onClick={handleRefresh}
            className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-4 sm:space-y-6 p-3 sm:p-6">
      {/* Header */}
      <div className="bg-white rounded-lg shadow-md p-3 sm:p-6">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 sm:gap-6 mb-4 sm:mb-6">
          <div>
            <h2 className="text-xl sm:text-2xl font-bold text-gray-900">Reading Analytics</h2>
            <p className="text-xs sm:text-base text-gray-600 mt-1">Student reading patterns and performance insights</p>
          </div>
          <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-2 sm:gap-4 w-full sm:w-auto">
            <select
              value={selectedPeriod}
              onChange={(e) => setSelectedPeriod(e.target.value)}
              className="px-3 sm:px-4 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 w-full sm:w-auto"
            >
              <option value="week">Last Week</option>
              <option value="month">Last Month</option>
              <option value="quarter">Last Quarter</option>
              <option value="year">Last Year</option>
            </select>
            <button
              onClick={handleRefresh}
              className="px-3 sm:px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 flex items-center justify-center gap-2 text-sm font-medium whitespace-nowrap w-full sm:w-auto"
            >
              <i className="ri-refresh-line"></i>
              <span>Refresh</span>
            </button>
          </div>
        </div>

        {/* Key Metrics */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">
          <div className="bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-lg p-3 sm:p-4 md:p-6">
            <div className="flex items-center justify-between gap-2">
              <div className="min-w-0">
                <p className="text-blue-100 text-xs sm:text-sm">Active Readers</p>
                <p className="text-lg sm:text-2xl md:text-3xl font-bold truncate">{safeData.overview.total_active_readers || 0}</p>
              </div>
              <i className="ri-user-line text-2xl sm:text-3xl md:text-4xl opacity-50 flex-shrink-0"></i>
            </div>
          </div>

          <div className="bg-gradient-to-r from-green-500 to-green-600 text-white rounded-lg p-3 sm:p-4 md:p-6">
            <div className="flex items-center justify-between gap-2">
              <div className="min-w-0">
                <p className="text-green-100 text-xs sm:text-sm">Books Completed</p>
                <p className="text-lg sm:text-2xl md:text-3xl font-bold truncate">{safeData.overview.total_books_completed || 0}</p>
              </div>
              <i className="ri-book-line text-2xl sm:text-3xl md:text-4xl opacity-50 flex-shrink-0"></i>
            </div>
          </div>

          <div className="bg-gradient-to-r from-purple-500 to-purple-600 text-white rounded-lg p-3 sm:p-4 md:p-6">
            <div className="flex items-center justify-between gap-2">
              <div className="min-w-0">
                <p className="text-purple-100 text-xs sm:text-sm">Avg Reading Time</p>
                <p className="text-lg sm:text-2xl md:text-3xl font-bold truncate">{(safeData.overview.avg_reading_time || 0).toFixed(1)}h</p>
              </div>
              <i className="ri-time-line text-2xl sm:text-3xl md:text-4xl opacity-50 flex-shrink-0"></i>
            </div>
          </div>

          <div className="bg-gradient-to-r from-orange-500 to-orange-600 text-white rounded-lg p-3 sm:p-4 md:p-6">
            <div className="flex items-center justify-between gap-2">
              <div className="min-w-0">
                <p className="text-orange-100 text-xs sm:text-sm">Avg Completion</p>
                <p className="text-lg sm:text-2xl md:text-3xl font-bold truncate">{(safeData.overview.avg_completion_rate || 0).toFixed(0)}%</p>
              </div>
              <i className="ri-checkbox-circle-line text-2xl sm:text-3xl md:text-4xl opacity-50 flex-shrink-0"></i>
            </div>
          </div>
        </div>
      </div>

      {/* Active Readers */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Active Readers</h3>
        {safeData.activeReaders.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="min-w-full">
              <thead>
                <tr className="border-b border-gray-200">
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-700">Student</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-700">Class</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-700">Sessions</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-700">Total Time</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-700">Last Active</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-700">Actions</th>
                </tr>
              </thead>
              <tbody>
                {safeData.activeReaders.map((reader, idx) => (
                  <tr key={idx} className="border-b border-gray-100 hover:bg-gray-50">
                    <td className="py-3 px-4 text-sm text-gray-900">{reader.name}</td>
                    <td className="py-3 px-4 text-sm text-gray-600">{reader.class_level}</td>
                    <td className="py-3 px-4 text-sm text-gray-600">{reader.session_count}</td>
                    <td className="py-3 px-4 text-sm text-gray-600">{Math.floor(reader.total_time / 60)}m</td>
                    <td className="py-3 px-4 text-sm text-gray-600">
                      {new Date(reader.last_session).toLocaleDateString()}
                    </td>
                    <td className="py-3 px-4 text-sm">
                      <button
                        onClick={() => {
                          setSelectedLibraryId(reader.library_id);
                          setShowDetailsModal(true);
                        }}
                        className="text-blue-600 hover:text-blue-800"
                        title="View Details"
                      >
                        <i className="ri-eye-line"></i>
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="h-64 flex items-center justify-center text-gray-500">
            No active readers found
          </div>
        )}
      </div>

      {/* Charts Row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 sm:gap-6">
        {/* Class Performance */}
        <div className="bg-white rounded-lg shadow-md p-3 sm:p-6">
          <h3 className="text-base sm:text-lg font-semibold text-gray-900 mb-3 sm:mb-4">Performance by Class</h3>
          {safeData.classStats.length > 0 ? (
            <ResponsiveContainer width="100%" height={250}>
              <BarChart data={safeData.classStats} margin={{ left: -20, right: 0 }}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="class_level" tick={{ fontSize: 12 }} />
                <YAxis tick={{ fontSize: 12 }} />
                <Tooltip />
                <Bar dataKey="avg_completion" fill="#3B82F6" name="Avg Completion %" />
                <Bar dataKey="books_completed" fill="#10B981" name="Books Completed" />
              </BarChart>
            </ResponsiveContainer>
          ) : (
            <div className="h-48 sm:h-64 flex items-center justify-center text-gray-500">
              No data available
            </div>
          )}
        </div>

        {/* Reading Time by Class */}
        <div className="bg-white rounded-lg shadow-md p-3 sm:p-6">
          <h3 className="text-base sm:text-lg font-semibold text-gray-900 mb-3 sm:mb-4">Avg Reading Time by Class</h3>
          {safeData.classStats.length > 0 ? (
            <ResponsiveContainer width="100%" height={250}>
              <BarChart data={safeData.classStats} margin={{ left: -20, right: 0 }}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="class_level" tick={{ fontSize: 12 }} />
                <YAxis tick={{ fontSize: 12 }} />
                <Tooltip />
                <Bar dataKey="avg_reading_time" fill="#8B5CF6" name="Avg Time (min)" />
              </BarChart>
            </ResponsiveContainer>
          ) : (
            <div className="h-48 sm:h-64 flex items-center justify-center text-gray-500">
              No data available
            </div>
          )}
        </div>
      </div>

      {/* Tables Row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Top Readers with Streaks */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Top Readers</h3>
          {safeData.topReaders.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="min-w-full">
                <thead>
                  <tr className="border-b border-gray-200">
                    <th className="text-left py-3 px-4 text-sm font-medium text-gray-700">Student</th>
                    <th className="text-left py-3 px-4 text-sm font-medium text-gray-700">Class</th>
                    <th className="text-left py-3 px-4 text-sm font-medium text-gray-700">Completed</th>
                    <th className="text-left py-3 px-4 text-sm font-medium text-gray-700">Time (h)</th>
                  </tr>
                </thead>
                <tbody>
                  {safeData.topReaders.map((reader, idx) => (
                    <tr key={idx} className="border-b border-gray-100 hover:bg-gray-50">
                      <td className="py-3 px-4 text-sm text-gray-900">{reader.name}</td>
                      <td className="py-3 px-4 text-sm text-gray-600">{reader.class_level}</td>
                      <td className="py-3 px-4 text-sm text-gray-600">{reader.books_completed}</td>
                      <td className="py-3 px-4 text-sm text-gray-600">{(reader.reading_time / 3600).toFixed(1)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <div className="h-64 flex items-center justify-center text-gray-500">
              No data available
            </div>
          )}
        </div>

        {/* Struggling Readers */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Struggling Readers</h3>
          {safeData.strugglingReaders.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="min-w-full">
                <thead>
                  <tr className="border-b border-gray-200">
                    <th className="text-left py-3 px-4 text-sm font-medium text-gray-700">Student</th>
                    <th className="text-left py-3 px-4 text-sm font-medium text-gray-700">Class</th>
                    <th className="text-left py-3 px-4 text-sm font-medium text-gray-700">Completion</th>
                  </tr>
                </thead>
                <tbody>
                  {safeData.strugglingReaders.map((reader, idx) => (
                    <tr key={idx} className="border-b border-gray-100 hover:bg-gray-50">
                      <td className="py-3 px-4 text-sm text-gray-900">{reader.name}</td>
                      <td className="py-3 px-4 text-sm text-gray-600">{reader.class_level}</td>
                      <td className="py-3 px-4 text-sm text-red-600">{reader.avg_completion.toFixed(0)}%</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <div className="h-64 flex items-center justify-center text-gray-500">
              No struggling readers found
            </div>
          )}
        </div>
      </div>

      {/* Reading Details Modal */}
      {showDetailsModal && selectedLibraryId && (
        <ReadingDetailsModal
          libraryId={selectedLibraryId}
          onClose={() => {
            setShowDetailsModal(false);
            setSelectedLibraryId(null);
          }}
        />
      )}
    </div>
  );
};

export default ReadingAnalytics;
