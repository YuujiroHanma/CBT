-- Updated Seed Products with Direct Product Image URLs
-- Execute this in Supabase SQL Editor to update all products

-- Delete existing products (optional, if you want to reset)
DELETE FROM products;

-- Insert new products with updated image URLs
INSERT INTO products (name, description, price, image_url, stock_quantity) VALUES
(
    'Wireless Headphones',
    'Premium noise-cancelling wireless headphones with 30-hour battery life.',
    149.99,
    'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?q=80&w=500&auto=format&fit=crop',
    50
),
(
    'USB-C Hub',
    '7-in-1 USB-C hub with HDMI, USB 3.0, and SD card reader.',
    79.99,
    'https://s13emagst.akamaized.net/products/54323/54322728/images/res_ab1fe51efaa866631379d88815bb65be.png?width=720&height=720&hash=43DC3A0A0FD51895AC136114DE21CC94',
    75
),
(
    'Mechanical Keyboard',
    'RGB mechanical keyboard with hot-swappable switches, ideal for gaming.',
    129.99,
    'https://cdn.thewirecutter.com/wp-content/media/2025/12/BEST-MECHANICAL-KEYBOARDS-2048px-EVOWORKS-80-926.jpg?width=2048&quality=60&crop=2048:1365&auto=webp',
    40
),
(
    '4K Webcam',
    'Crystal clear 4K USB webcam with auto-focus for streaming and calls.',
    199.99,
    'https://m.media-amazon.com/images/I/61b4sT7kQeL.jpg',
    25
),
(
    'Portable SSD - 1TB',
    'Ultra-fast portable SSD with 1050MB/s read speed, perfect for professionals.',
    99.99,
    'https://images-cdn.ubuy.co.in/6627b88a13046b5ecf751c4c-1-tb-ultra-high-speed-portable-ssd.jpg',
    60
),
(
    'Monitor Light Bar',
    'Premium monitor light bar that reduces eye strain for extended work sessions.',
    69.99,
    'https://www.deskup.com.au/media/catalog/product/cache/55e9079b027a86c564cfc6b679f550fc/d/e/deskup_rgb_monitor_light_-_hero.png',
    30
);
