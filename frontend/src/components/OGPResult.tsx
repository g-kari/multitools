import type { OGPResponse } from '@/types/ogp';

interface OGPResultProps {
  data: OGPResponse;
}

export const OGPResult: React.FC<OGPResultProps> = ({ data }) => {
  const { ogp_data, validation, previews } = data;

  return (
    <div className="w-full max-w-4xl mx-auto space-y-6">
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-semibold mb-4 text-gray-800">OGP Data</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <h3 className="font-medium text-gray-700 mb-2">Title</h3>
            <p className="text-gray-600 bg-gray-50 p-2 rounded">
              {ogp_data.title || 'Not found'}
            </p>
          </div>
          
          <div>
            <h3 className="font-medium text-gray-700 mb-2">Site Name</h3>
            <p className="text-gray-600 bg-gray-50 p-2 rounded">
              {ogp_data.site_name || 'Not found'}
            </p>
          </div>
          
          <div className="md:col-span-2">
            <h3 className="font-medium text-gray-700 mb-2">Description</h3>
            <p className="text-gray-600 bg-gray-50 p-2 rounded">
              {ogp_data.description || 'Not found'}
            </p>
          </div>
          
          <div className="md:col-span-2">
            <h3 className="font-medium text-gray-700 mb-2">Image</h3>
            {ogp_data.image ? (
              <div className="space-y-2">
                <p className="text-gray-600 bg-gray-50 p-2 rounded break-all">
                  {ogp_data.image}
                </p>
                <img
                  src={ogp_data.image}
                  alt={ogp_data.image_alt || 'OGP Image'}
                  className="max-w-full h-auto max-h-64 rounded border"
                  onError={(e) => {
                    e.currentTarget.style.display = 'none';
                  }}
                />
              </div>
            ) : (
              <p className="text-gray-600 bg-gray-50 p-2 rounded">Not found</p>
            )}
          </div>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-semibold mb-4 text-gray-800">Validation Results</h2>
        
        <div className="space-y-4">
          <div className={`p-4 rounded-lg ${validation.is_valid ? 'bg-green-50 border border-green-200' : 'bg-red-50 border border-red-200'}`}>
            <h3 className={`font-medium ${validation.is_valid ? 'text-green-800' : 'text-red-800'}`}>
              {validation.is_valid ? '✓ Valid' : '✗ Invalid'}
            </h3>
          </div>
          
          {validation.warnings.length > 0 && (
            <div className="bg-yellow-50 border border-yellow-200 p-4 rounded-lg">
              <h3 className="font-medium text-yellow-800 mb-2">Warnings</h3>
              <ul className="list-disc list-inside text-yellow-700 space-y-1">
                {validation.warnings.map((warning, index) => (
                  <li key={index}>{warning}</li>
                ))}
              </ul>
            </div>
          )}
          
          {validation.errors.length > 0 && (
            <div className="bg-red-50 border border-red-200 p-4 rounded-lg">
              <h3 className="font-medium text-red-800 mb-2">Errors</h3>
              <ul className="list-disc list-inside text-red-700 space-y-1">
                {validation.errors.map((error, index) => (
                  <li key={index}>{error}</li>
                ))}
              </ul>
            </div>
          )}
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-semibold mb-4 text-gray-800">Platform Previews</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {Object.entries(previews).map(([platform, preview]) => (
            <div key={platform} className="border rounded-lg p-4">
              <h3 className="font-medium text-gray-800 mb-3 capitalize flex items-center">
                {platform}
                {preview.platform === 'twitter' && ' (X)'}
              </h3>
              
              <div className="space-y-2">
                <div>
                  <p className="text-sm text-gray-600">
                    Title ({preview.title_length}/{preview.max_title_len})
                  </p>
                  <p className="text-sm bg-gray-50 p-2 rounded">
                    {preview.title || 'No title'}
                  </p>
                </div>
                
                <div>
                  <p className="text-sm text-gray-600">
                    Description ({preview.desc_length}/{preview.max_desc_len})
                  </p>
                  <p className="text-sm bg-gray-50 p-2 rounded">
                    {preview.description || 'No description'}
                  </p>
                </div>
                
                {preview.image && (
                  <div>
                    <p className="text-sm text-gray-600 mb-1">Image</p>
                    <img
                      src={preview.image}
                      alt="Preview"
                      className="w-full h-24 object-cover rounded border"
                      onError={(e) => {
                        e.currentTarget.style.display = 'none';
                      }}
                    />
                  </div>
                )}
                
                {preview.warnings.length > 0 && (
                  <div className="text-xs text-yellow-700 bg-yellow-50 p-2 rounded">
                    {preview.warnings.join(', ')}
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};