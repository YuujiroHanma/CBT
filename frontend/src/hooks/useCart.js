import { useState, useCallback, useEffect } from 'react'

/**
 * Custom hook to manage cart state and operations.
 * Persists cart to localStorage for durability across page refreshes.
 */
export const useCart = () => {
  const [cart, setCart] = useState(() => {
    // Initialize cart from localStorage or empty array
    const savedCart = localStorage.getItem('cart')
    return savedCart ? JSON.parse(savedCart) : []
  })

  // Persist cart to localStorage whenever it changes
  useEffect(() => {
    localStorage.setItem('cart', JSON.stringify(cart))
  }, [cart])

  // Add item to cart or increment quantity if it already exists
  const addToCart = useCallback((product) => {
    setCart((prevCart) => {
      const existingItem = prevCart.find((item) => item.id === product.id)
      if (existingItem) {
        return prevCart.map((item) =>
          item.id === product.id
            ? { ...item, quantity: item.quantity + 1 }
            : item
        )
      }
      return [...prevCart, { ...product, quantity: 1 }]
    })
  }, [])

  // Increment quantity of an item
  const incrementQuantity = useCallback((productId) => {
    setCart((prevCart) =>
      prevCart.map((item) =>
        item.id === productId
          ? { ...item, quantity: item.quantity + 1 }
          : item
      )
    )
  }, [])

  // Decrement quantity of an item
  const decrementQuantity = useCallback((productId) => {
    setCart((prevCart) =>
      prevCart
        .map((item) =>
          item.id === productId
            ? { ...item, quantity: Math.max(0, item.quantity - 1) }
            : item
        )
        .filter((item) => item.quantity > 0)
    )
  }, [])

  // Remove item from cart
  const removeFromCart = useCallback((productId) => {
    setCart((prevCart) => prevCart.filter((item) => item.id !== productId))
  }, [])

  // Clear entire cart
  const clearCart = useCallback(() => {
    setCart([])
  }, [])

  // Calculate total price
  const cartTotal = cart.reduce(
    (total, item) => total + item.price * item.quantity,
    0
  )

  // Calculate total item count
  const cartItemCount = cart.reduce((count, item) => count + item.quantity, 0)

  return {
    cart,
    addToCart,
    incrementQuantity,
    decrementQuantity,
    removeFromCart,
    clearCart,
    cartTotal,
    cartItemCount,
  }
}
