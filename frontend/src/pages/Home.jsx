import React from 'react'
import ProductList from '../components/ProductList'
import Cart from '../components/Cart'
import { useCart } from '../hooks/useCart'
import apiClient from '../services/api'

/**
 * Home Page
 * Main page combining ProductList and Cart components.
 * Manages cart state and handles add-to-cart functionality.
 */
export default function Home() {
  const {
    cart,
    addToCart,
    incrementQuantity,
    decrementQuantity,
    removeFromCart,
    cartTotal,
  } = useCart()

  const handleCheckout = async (cartItems, total) => {
    try {
      // Prepare cart items for the API
      const items = cartItems.map((item) => ({
        product_id: item.id,
        quantity: item.quantity,
        price: item.price,
        name: item.name,
      }))

      // Call the /api/checkout endpoint
      const response = await apiClient.post('/checkout', {
        items,
        success_url: `${window.location.origin}/success`,
        cancel_url: `${window.location.origin}/cancelled`,
      })

      // Redirect to Stripe Checkout
      if (response.data.url) {
        window.location.href = response.data.url
      }
    } catch (error) {
      const errorMsg =
        error.response?.data?.error || 'Failed to initiate checkout'
      throw new Error(errorMsg)
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Main Content */}
      <main className="max-w-6xl mx-auto px-4 py-12">
        <ProductList onAddToCart={addToCart} />
      </main>

      {/* Cart Sidebar */}
      <Cart
        cart={cart}
        onIncrement={incrementQuantity}
        onDecrement={decrementQuantity}
        onRemove={removeFromCart}
        cartTotal={cartTotal}
        onCheckout={handleCheckout}
      />
    </div>
  )
}