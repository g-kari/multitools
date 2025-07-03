import { URLInput } from '@/components/URLInput';
import { OGPResult } from '@/components/OGPResult';
import { ErrorMessage } from '@/components/ErrorMessage';
import { useOGP } from '@/hooks/useOGP';

function App() {
  const { data, loading, error, verifyOGP, reset } = useOGP();

  const handleSubmit = async (url: string) => {
    await verifyOGP({ url });
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <header className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            OGP Verification Service
          </h1>
          <p className="text-gray-600">
            Analyze and validate Open Graph Protocol metadata for social media platforms
          </p>
        </header>

        <div className="space-y-8">
          <URLInput onSubmit={handleSubmit} loading={loading} />

          {error && (
            <ErrorMessage
              error={error}
              onRetry={reset}
            />
          )}

          {loading && (
            <div className="text-center py-8">
              <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
              <p className="mt-2 text-gray-600">Analyzing OGP data...</p>
            </div>
          )}

          {data && !loading && !error && (
            <OGPResult data={data} />
          )}
        </div>

        <footer className="mt-12 text-center text-gray-500 text-sm">
          <p>
            Supports Twitter/X, Facebook, and Discord preview formats
          </p>
        </footer>
      </div>
    </div>
  );
}

export default App;