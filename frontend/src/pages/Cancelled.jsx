import { useNavigate } from 'react-router-dom'

/**
 * Cancelled Page
 * Displayed when user cancels Stripe checkout
 */
export default function Cancelled() {
  const navigate = useNavigate()

  return (
    <div className="min-h-screen bg-gradient-to-b from-orange-50 to-orange-100 flex items-center justify-center px-4">
      <div className="text-center max-w-md">
        {/* Cancelled Icon */}
        <div className="text-6xl mb-6">⚠️</div>

        <h1 className="text-4xl font-bold text-orange-900 mb-4">
          Payment Cancelled
        </h1>

        <p className="text-lg text-orange-800 mb-2">
          Your payment was cancelled.
        </p>

        <p className="text-gray-700 mb-8">
          No charges were made to your account. Your items are still in your
          cart and ready for checkout whenever you're ready.
        </p>

        {/* Info Box */}
        <div className="bg-white rounded-lg shadow p-6 mb-8 border-l-4 border-orange-500">
          <p className="text-gray-700 text-sm">
            If you encountered any issues during checkout, please try again or
            contact our support team for assistance.
          </p>
        </div>

        {/* Back to Cart Button */}
        <button
          onClick={() => navigate('/')}
          className="px-8 py-3 bg-orange-600 text-white font-semibold rounded-lg hover:bg-orange-700 transition-colors duration-200"
        >
          Return to Cart
        </button>
      </div>
    </div>
  )
}
