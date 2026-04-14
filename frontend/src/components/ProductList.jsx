import { useEffect, useState } from 'react'
import { fetchProducts } from '../services/api'

/**
 * ProductList Component
 * Fetches and displays products from the backend API in a responsive grid.
 * Includes loading and error state handling.
 */
export default function ProductList({ onAddToCart }) {
  const [products, setProducts] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  // Fetch products on component mount
  useEffect(() => {
    const getProducts = async () => {
      try {
        setLoading(true)
        setError(null)
        const response = await fetchProducts()
        setProducts(response.data || [])
      } catch (err) {
        console.error('Error fetching products:', err)
        setError('Failed to load products. Please try again later.')
      } finally {
        setLoading(false)
      }
    }

    getProducts()
  }, [])

  // Loading state
  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-lg text-gray-500">Loading products...</div>
      </div>
    )
  }

  // Error state
  if (error) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-lg text-red-500">{error}</div>
      </div>
    )
  }

  // No products
  if (products.length === 0) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-lg text-gray-500">No products available.</div>
      </div>
    )
  }

  return (
    <div className="py-8">
      <h2 className="text-3xl font-bold text-gray-900 mb-8">Our Products</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {products.map((product) => (
          <div
            key={product.id}
            className="border border-gray-200 rounded-lg overflow-hidden hover:shadow-lg transition-shadow duration-300"
          >
            {/* Product Image */}
            <div className="w-full h-64 bg-gray-100 overflow-hidden">
              <img
                src={product.image_url}
                alt={product.name}
                className="w-full h-full object-cover hover:scale-105 transition-transform duration-300"
              />
            </div>

            {/* Product Info */}
            <div className="p-4">
              <h3 className="text-lg font-semibold text-gray-900 mb-2">
                {product.name}
              </h3>
              <p className="text-gray-600 text-sm mb-4 line-clamp-2">
                {product.description}
              </p>

              {/* Stock Info */}
              <div className="flex items-center justify-between mb-4">
                <span
                  className={`text-sm font-medium ${
                    product.stock_quantity > 0
                      ? 'text-green-600'
                      : 'text-red-600'
                  }`}
                >
                  {product.stock_quantity > 0
                    ? `${product.stock_quantity} in stock`
                    : 'Out of stock'}
                </span>
              </div>

              {/* Price and Button */}
              <div className="flex items-center justify-between">
                <span className="text-2xl font-bold text-gray-900">
                  ${product.price.toFixed(2)}
                </span>
                <button
                  onClick={() => onAddToCart(product)}
                  disabled={product.stock_quantity === 0}
                  className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors duration-200 disabled:bg-gray-400 disabled:cursor-not-allowed"
                >
                  Add to Cart
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
