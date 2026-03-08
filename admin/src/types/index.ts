export type { ApiResponse, PaginationMeta, ErrorResponse } from './api';

// ---------- Status Union Types ----------

export type UserRole = 'customer' | 'vendor' | 'admin' | 'delivery_partner';

export type UserStatus = 'active' | 'inactive' | 'suspended' | 'banned';

export type VendorStatus = 'pending' | 'approved' | 'rejected' | 'suspended';

export type VendorType = 'restaurant' | 'grocery' | 'pharmacy' | 'electronics' | 'fashion' | 'home' | 'beauty' | 'service' | 'other';

export type DocumentStatus = 'pending' | 'approved' | 'rejected';

export type ProductStatus = 'active' | 'inactive' | 'out_of_stock' | 'draft';

export type CategoryType = 'product' | 'service';

export type OrderStatus =
  | 'pending'
  | 'confirmed'
  | 'preparing'
  | 'ready'
  | 'assigned'
  | 'out_for_delivery'
  | 'delivered'
  | 'cancelled'
  | 'refunded';

export type DeliveryPartnerStatus = 'pending' | 'approved' | 'rejected' | 'suspended';

export type DeliveryAssignmentStatus = 'assigned' | 'accepted' | 'rejected' | 'picked_up' | 'delivered' | 'cancelled';

export type VehicleType = 'bike' | 'scooter' | 'bicycle' | 'car';

export type BookingStatus =
  | 'pending'
  | 'confirmed'
  | 'in_progress'
  | 'completed'
  | 'cancelled'
  | 'no_show';

export type PaymentStatus = 'created' | 'authorized' | 'captured' | 'failed' | 'refunded';

export type PaymentMethod = 'razorpay' | 'wallet' | 'cod' | 'upi' | 'card' | 'net_banking';

export type TransactionType = 'credit' | 'debit' | 'payout' | 'refund' | 'commission';

export type NotificationType = 'promotion' | 'system' | 'order_update' | 'booking_update' | 'vendor_update';

// ---------- Entity Types ----------

export interface User {
  id: string;
  email: string;
  phone: string;
  first_name: string;
  last_name: string;
  role: UserRole;
  status: UserStatus;
  avatar_url?: string;
  email_verified: boolean;
  phone_verified: boolean;
  created_at: string;
  updated_at: string;
}

export interface Vendor {
  id: string;
  user_id: string;
  business_name: string;
  vendor_type: string;
  description?: string;
  logo_url?: string;
  banner_url?: string;
  address: string;
  city: string;
  state: string;
  pincode: string;
  latitude?: number;
  longitude?: number;
  avg_rating: number;
  total_reviews: number;
  commission_pct: number;
  service_radius_km: number;
  is_online: boolean;
  status: VendorStatus;
  created_at: string;
  updated_at: string;
  user?: User;
  documents?: VendorDocument[];
}

export interface VendorDocument {
  id: string;
  vendor_id: string;
  document_type: string;
  document_url: string;
  status: DocumentStatus;
  rejection_reason?: string;
  verified_at?: string;
  created_at: string;
  updated_at: string;
}

export interface Product {
  id: string;
  vendor_id: string;
  category_id: string;
  name: string;
  slug: string;
  description?: string;
  short_description?: string;
  price: number;
  compare_at_price?: number;
  cost_price?: number;
  sku?: string;
  stock_quantity: number;
  low_stock_threshold: number;
  images: string[];
  thumbnail_url?: string;
  weight?: number;
  unit?: string;
  status: ProductStatus;
  is_featured: boolean;
  rating: number;
  total_reviews: number;
  created_at: string;
  updated_at: string;
  vendor?: Vendor;
  category?: Category;
}

export interface Category {
  id: string;
  name: string;
  slug: string;
  description?: string;
  image_url?: string;
  category_type: CategoryType;
  parent_id?: string;
  sort_order: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  children?: Category[];
  parent?: Category;
}

export interface Order {
  id: string;
  order_number: string;
  customer_id: string;
  vendor_id: string;
  status: OrderStatus;
  subtotal: number;
  delivery_fee: number;
  tax_amount: number;
  discount_amount: number;
  total_amount: number;
  delivery_address_line1: string;
  delivery_address_line2?: string;
  delivery_city: string;
  delivery_state: string;
  delivery_pincode: string;
  delivery_latitude?: number;
  delivery_longitude?: number;
  notes?: string;
  estimated_delivery_at?: string;
  delivered_at?: string;
  cancelled_at?: string;
  cancellation_reason?: string;
  created_at: string;
  updated_at: string;
  customer?: User;
  vendor?: Vendor;
  items?: OrderItem[];
  payment?: Payment;
}

export interface OrderItem {
  id: string;
  order_id: string;
  product_id: string;
  product_name: string;
  product_image?: string;
  quantity: number;
  unit_price: number;
  total_price: number;
  notes?: string;
  created_at: string;
  product?: Product;
}

export interface Booking {
  id: string;
  booking_number: string;
  customer_id: string;
  vendor_id: string;
  service_id: string;
  service_slot_id?: string;
  status: BookingStatus;
  booking_date: string;
  start_time: string;
  end_time: string;
  duration_minutes: number;
  subtotal: number;
  tax_amount: number;
  discount_amount: number;
  total_amount: number;
  notes?: string;
  cancellation_reason?: string;
  completed_at?: string;
  cancelled_at?: string;
  created_at: string;
  updated_at: string;
  customer?: User;
  vendor?: Vendor;
  service?: Service;
  payment?: Payment;
}

export interface Service {
  id: string;
  vendor_id: string;
  category_id: string;
  name: string;
  slug: string;
  description?: string;
  price: number;
  compare_at_price?: number;
  duration_minutes: number;
  images: string[];
  thumbnail_url?: string;
  is_active: boolean;
  is_featured: boolean;
  rating: number;
  total_reviews: number;
  created_at: string;
  updated_at: string;
  vendor?: Vendor;
  category?: Category;
  slots?: ServiceSlot[];
}

export interface ServiceSlot {
  id: string;
  service_id: string;
  day_of_week: number;
  start_time: string;
  end_time: string;
  max_bookings: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Payment {
  id: string;
  order_id?: string;
  booking_id?: string;
  customer_id: string;
  vendor_id: string;
  amount: number;
  currency: string;
  method: PaymentMethod;
  status: PaymentStatus;
  gateway_payment_id?: string;
  gateway_order_id?: string;
  gateway_signature?: string;
  refund_amount?: number;
  refund_reason?: string;
  refunded_at?: string;
  paid_at?: string;
  created_at: string;
  updated_at: string;
  customer?: User;
  vendor?: Vendor;
}

export interface Wallet {
  id: string;
  user_id: string;
  vendor_id?: string;
  delivery_partner_id?: string;
  balance: number;
  currency: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  user?: User;
}

export interface WalletTransaction {
  id: string;
  wallet_id: string;
  type: TransactionType;
  amount: number;
  balance_after: number;
  description?: string;
  reference_type?: string;
  reference_id?: string;
  status: 'pending' | 'completed' | 'failed';
  created_at: string;
  updated_at: string;
  wallet?: Wallet;
}

export interface Review {
  id: string;
  customer_id: string;
  vendor_id: string;
  product_id?: string;
  service_id?: string;
  order_id?: string;
  booking_id?: string;
  rating: number;
  title?: string;
  comment?: string;
  images?: string[];
  is_verified: boolean;
  created_at: string;
  updated_at: string;
  customer?: User;
  vendor?: Vendor;
  product?: Product;
  service?: Service;
}

export interface Notification {
  id: string;
  user_id?: string;
  title: string;
  body: string;
  type: NotificationType;
  data?: Record<string, unknown>;
  is_read: boolean;
  read_at?: string;
  sent_at?: string;
  created_at: string;
  updated_at: string;
}

export interface DeliveryPartner {
  id: string;
  user_id: string;
  vehicle_type: VehicleType;
  vehicle_number?: string;
  license_number?: string;
  status: DeliveryPartnerStatus;
  is_available: boolean;
  is_on_shift: boolean;
  current_latitude?: number;
  current_longitude?: number;
  current_order_id?: string;
  zone_preference?: string;
  avg_rating: number;
  total_deliveries: number;
  total_earnings: number;
  commission_pct: number;
  created_at: string;
  updated_at: string;
  user?: User;
}

export interface DeliveryAssignment {
  id: string;
  order_id: string;
  delivery_partner_id: string;
  status: DeliveryAssignmentStatus;
  assigned_at?: string;
  accepted_at?: string;
  picked_up_at?: string;
  delivered_at?: string;
  delivery_proof_url?: string;
  delivery_otp?: string;
  distance_km?: number;
  earnings?: number;
  rejection_reason?: string;
  created_at: string;
  order?: Order;
  delivery_partner?: DeliveryPartner;
}
