-- Drop indices --
DROP INDEX IF EXISTS idx_books_publisher_id;
DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_cart_user_id;
DROP INDEX IF EXISTS idx_order_items_book_id;
DROP INDEX IF EXISTS idx_cart_items_book_id;

-- Drop tables (children → parents) --
DROP TABLE IF EXISTS book_categories;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS publishers;
DROP TABLE IF EXISTS users;

-- Drop types --
DROP TYPE IF EXISTS USER_ROLE;
DROP TYPE IF EXISTS PAYMENT_STATUS;
DROP TYPE IF EXISTS PAYMENT_METHOD;
DROP TYPE IF EXISTS ORDER_STATUS;
