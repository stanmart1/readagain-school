export default function SettingsPanel({ 
  show, 
  onClose, 
  fontSize, 
  onFontSizeChange, 
  fontFamily, 
  onFontFamilyChange, 
  theme, 
  onThemeChange,
  scale,
  onScaleChange,
  readerType = 'epub' // 'epub' or 'pdf'
}) {
  if (!show) return null;

  return (
    <div className="absolute top-16 right-4 w-80 bg-white rounded-xl shadow-2xl z-10 p-6">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-lg font-bold">Reading Settings</h3>
        <button
          onClick={onClose}
          className="p-2 hover:bg-gray-100 rounded-lg"
        >
          <i className="ri-close-line text-xl"></i>
        </button>
      </div>

      {/* Font Size / Scale */}
      <div className="mb-6">
        <label className="block text-sm font-medium text-gray-700 mb-3">
          {readerType === 'pdf' ? 'Zoom Level' : 'Font Size'}
        </label>
        <div className="flex items-center space-x-4">
          <button
            onClick={() => readerType === 'pdf' ? onScaleChange(Math.max(0.5, scale - 0.25)) : onFontSizeChange(-10)}
            className="p-2 bg-gray-100 hover:bg-gray-200 rounded-lg"
          >
            <i className="ri-subtract-line"></i>
          </button>
          <span className="flex-1 text-center font-medium">
            {readerType === 'pdf' ? `${Math.round(scale * 100)}%` : `${fontSize}%`}
          </span>
          <button
            onClick={() => readerType === 'pdf' ? onScaleChange(Math.min(3, scale + 0.25)) : onFontSizeChange(10)}
            className="p-2 bg-gray-100 hover:bg-gray-200 rounded-lg"
          >
            <i className="ri-add-line"></i>
          </button>
        </div>
      </div>

      {/* Font Family - Only for EPUB */}
      {readerType === 'epub' && (
        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-3">Font Family</label>
          <div className="grid grid-cols-2 gap-2">
            {['Georgia', 'Arial', 'Times New Roman', 'Verdana'].map(font => (
              <button
                key={font}
                onClick={() => onFontFamilyChange(font)}
                className={`px-3 py-2 rounded-lg border-2 transition-all text-sm ${
                  fontFamily === font
                    ? 'border-primary-500 bg-primary-50 text-primary-700'
                    : 'border-gray-200 hover:border-gray-300'
                }`}
                style={{ fontFamily: font }}
              >
                {font}
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Theme - Only for EPUB */}
      {readerType === 'epub' && (
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-3">Theme</label>
          <div className="grid grid-cols-3 gap-2">
            <button
              onClick={() => onThemeChange('light')}
              className={`p-3 rounded-lg border-2 transition-all ${
                theme === 'light' ? 'border-primary-500 bg-primary-50' : 'border-gray-200'
              }`}
            >
              <div className="w-full h-8 bg-white border border-gray-300 rounded mb-2"></div>
              <p className="text-xs font-medium">Light</p>
            </button>
            <button
              onClick={() => onThemeChange('sepia')}
              className={`p-3 rounded-lg border-2 transition-all ${
                theme === 'sepia' ? 'border-primary-500 bg-primary-50' : 'border-gray-200'
              }`}
            >
              <div className="w-full h-8 bg-amber-50 border border-amber-200 rounded mb-2"></div>
              <p className="text-xs font-medium">Sepia</p>
            </button>
            <button
              onClick={() => onThemeChange('dark')}
              className={`p-3 rounded-lg border-2 transition-all ${
                theme === 'dark' ? 'border-primary-500 bg-primary-50' : 'border-gray-200'
              }`}
            >
              <div className="w-full h-8 bg-gray-900 border border-gray-700 rounded mb-2"></div>
              <p className="text-xs font-medium">Dark</p>
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
