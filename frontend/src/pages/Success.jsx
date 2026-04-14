import { useNavigate } from 'react-router-dom'
import { useCart } from '../hooks/useCart'
import { useEffect } from 'react'

/**
 * Success Page
 * Displayed after successful payment on Stripe
 */
export default function Success() {
  const navigate = useNavigate()
  const { clearCart } = useCart()

  // Clear the cart on successful payment
  useEffect(() => {
    clearCart()
  }, [clearCart])

  return (
    <div className="min-h-screen bg-gradient-to-b from-green-50 to-green-100 flex items-center justify-center px-4">
      <div className="text-center max-w-md">
        {/* Success Icon */}
        <div className="text-6xl mb-6">✅</div>

        <h1 className="text-4xl font-bold text-green-900 mb-4">
          Payment Successful!
        </h1>

        <p className="text-lg text-green-800 mb-2">
          Thank you for your order!
        </p>

        <p className="text-gray-700 mb-8">
          We've received your payment and your order is being processed. You'll
          receive a confirmation email shortly with tracking information.
        </p>

        {/* Order Details Placeholder */}
        <div className="bg-white rounded-lg shadow p-6 mb-8 text-left">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">
            Order Summary
          </h2>
          <p className="text-gray-600 text-sm mb-2">
            ✓ Payment Confirmed
          </p>
          <p className="text-gray-600 text-sm mb-2">
            ✓ Order Processing
          </p>
          <p className="text-gray-600 text-sm">
            ✓ Shipping Soon
          </p>
        </div>

        {/* Back to Home Button */}
        <button
          onClick={() => navigate('/')}
          className="px-8 py-3 bg-green-600 text-white font-semibold rounded-lg hover:bg-green-700 transition-colors duration-200"
        >
          Continue Shopping
        </button>
      </div>
    </div>
  )
}
