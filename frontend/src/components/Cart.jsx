import { useState } from 'react'
import { useAuth } from '../hooks/useAuth'

/**
 * Cart Component
 * Displays cart items, allows quantity adjustment, and shows total.
 * Can be toggled open/closed via a button.
 */
export default function Cart({
  cart,
  onIncrement,
  onDecrement,
  onRemove,
  cartTotal,
  onCheckout,
}) {
  const [isOpen, setIsOpen] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState(null)
  const { isAuthenticated } = useAuth()

  const toggleCart = () => setIsOpen(!isOpen)

  const handleCheckoutClick = async () => {
    if (!isAuthenticated) {
      // Redirect to login
      window.location.href = '/login'
      return
    }

    setIsLoading(true)
    setError(null)

    try {
      await onCheckout(cart, cartTotal)
    } catch (err) {
      setError(err.message || 'Checkout failed. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <>
      {/* Cart Toggle Button (Floating) */}
      <button
        onClick={toggleCart}
        className="fixed bottom-6 right-6 w-16 h-16 bg-blue-600 text-white rounded-full shadow-lg hover:bg-blue-700 transition-colors duration-200 flex items-center justify-center text-2xl z-40"
        aria-label="Toggle cart"
      >
        🛒
        {cart.length > 0 && (
          <span className="absolute top-0 right-0 bg-red-500 text-white text-xs font-bold rounded-full w-6 h-6 flex items-center justify-center">
            {cart.length}
          </span>
        )}
      </button>

      {/* Cart Sidebar */}
      <div
        className={`fixed top-0 right-0 h-full w-full sm:w-96 bg-white shadow-2xl transform transition-transform duration-300 z-50 ${
          isOpen ? 'translate-x-0' : 'translate-x-full'
        }`}
      >
        {/* Header */}
        <div className="bg-gray-100 p-6 flex items-center justify-between border-b">
          <h2 className="text-2xl font-bold text-gray-900">Shopping Cart</h2>
          <button
            onClick={toggleCart}
            className="text-gray-500 hover:text-gray-700 text-2xl"
            aria-label="Close cart"
          >
            ✕
          </button>
        </div>

        {/* Cart Items */}
        <div className="flex-1 overflow-y-auto p-6" style={{ maxHeight: 'calc(100vh - 200px)' }}>
          {cart.length === 0 ? (
            <div className="text-center text-gray-500 py-12">
              <p className="text-lg">Your cart is empty</p>
            </div>
          ) : (
            <div className="space-y-4">
              {cart.map((item) => (
                <div
                  key={item.id}
                  className="border border-gray-200 rounded-lg p-4 flex gap-4"
                >
                  {/* Item Image */}
                  <img
                    src={item.image_url}
                    alt={item.name}
                    className="w-20 h-20 object-cover rounded"
                  />

                  {/* Item Details */}
                  <div className="flex-1">
                    <h3 className="font-semibold text-gray-900">
                      {item.name}
                    </h3>
                    <p className="text-blue-600 font-bold mt-1">
                      ${item.price.toFixed(2)}
                    </p>

                    {/* Quantity Controls */}
                    <div className="flex items-center gap-2 mt-3">
                      <button
                        onClick={() => onDecrement(item.id)}
                        className="px-2 py-1 bg-gray-200 hover:bg-gray-300 rounded text-sm"
                      >
                        −
                      </button>
                      <span className="px-3 py-1 bg-gray-100 rounded text-sm font-medium">
                        {item.quantity}
                      </span>
                      <button
                        onClick={() => onIncrement(item.id)}
                        className="px-2 py-1 bg-gray-200 hover:bg-gray-300 rounded text-sm"
                      >
                        +
                      </button>
                      <button
                        onClick={() => onRemove(item.id)}
                        className="ml-auto text-red-500 hover:text-red-700 text-sm"
                      >
                        Remove
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Footer with Checkout */}
        {cart.length > 0 && (
          <div className="border-t bg-gray-50 p-6 space-y-4 absolute bottom-0 w-full">
            {error && (
              <div className="p-3 bg-red-50 border border-red-200 rounded">
                <p className="text-red-700 text-sm">{error}</p>
              </div>
            )}
            <div className="flex items-center justify-between text-lg font-bold">
              <span>Total:</span>
              <span className="text-blue-600">${cartTotal.toFixed(2)}</span>
            </div>
            <button
              onClick={handleCheckoutClick}
              disabled={isLoading}
              className="w-full px-4 py-3 bg-green-600 text-white font-semibold rounded-lg hover:bg-green-700 transition-colors duration-200 disabled:bg-gray-400 disabled:cursor-not-allowed"
            >
              {isLoading ? 'Processing...' : 'Proceed to Checkout'}
            </button>
          </div>
        )}
      </div>

      {/* Overlay */}
      {isOpen && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 z-40"
          onClick={toggleCart}
        />
      )}
    </>
  )
}