'use client';

import { useQuery } from '@tanstack/react-query';
import Link from 'next/link';
import {
  IndianRupee,
  ShoppingCart,
  Store,
  Users,
  ArrowRight,
  Bike,
  Clock,
  UserCheck,
} from 'lucide-react';
import { apiGet } from '@/lib/api-client';
import { formatCurrency, formatDate } from '@/lib/utils';
import { StatCard } from '@/components/ui/stat-card';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { StatusBadge } from '@/components/ui/status-badge';
import { RevenueChart } from '@/components/charts/revenue-chart';
import { OrderChart } from '@/components/charts/order-chart';
import { PageLoading } from '@/components/ui/loading';
import type { Order, Vendor, DeliveryPartner } from '@/types';

interface DashboardStats {
  total_revenue: number;
  total_orders: number;
  active_vendors: number;
  total_users: number;
  revenue_change: string;
  orders_change: string;
  vendors_change: string;
  users_change: string;
  active_delivery_partners: number;
  pending_delivery_partners: number;
  avg_delivery_time_mins: number;
  total_deliveries_today: number;
}

interface RevenueDataPoint {
  date: string;
  revenue: number;
}

interface OrderStatusCount {
  status: string;
  count: number;
}

const fallbackStats: DashboardStats = {
  total_revenue: 1245600,
  total_orders: 856,
  active_vendors: 42,
  total_users: 3240,
  revenue_change: '+12.5%',
  orders_change: '+8.2%',
  vendors_change: '+3',
  users_change: '+156',
  active_delivery_partners: 18,
  pending_delivery_partners: 5,
  avg_delivery_time_mins: 32,
  total_deliveries_today: 47,
};

const fallbackRevenueData: RevenueDataPoint[] = [
  { date: 'Mon', revenue: 45000 },
  { date: 'Tue', revenue: 52000 },
  { date: 'Wed', revenue: 48000 },
  { date: 'Thu', revenue: 61000 },
  { date: 'Fri', revenue: 55000 },
  { date: 'Sat', revenue: 67000 },
  { date: 'Sun', revenue: 43000 },
];

const fallbackOrderData: OrderStatusCount[] = [
  { status: 'pending', count: 24 },
  { status: 'confirmed', count: 18 },
  { status: 'preparing', count: 12 },
  { status: 'out_for_delivery', count: 8 },
  { status: 'delivered', count: 156 },
  { status: 'cancelled', count: 6 },
];

export default function DashboardPage() {
  const { data: statsData, isLoading: statsLoading } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: () => apiGet<DashboardStats>('/admin/dashboard/stats'),
  });

  const { data: revenueData } = useQuery({
    queryKey: ['dashboard-revenue'],
    queryFn: () => apiGet<RevenueDataPoint[]>('/admin/dashboard/revenue'),
  });

  const { data: orderStatusData } = useQuery({
    queryKey: ['dashboard-order-status'],
    queryFn: () => apiGet<OrderStatusCount[]>('/admin/dashboard/order-status'),
  });

  const { data: recentOrdersData } = useQuery({
    queryKey: ['dashboard-recent-orders'],
    queryFn: () => apiGet<Order[]>('/admin/orders?per_page=5&sort=-created_at'),
  });

  const { data: pendingVendorsData } = useQuery({
    queryKey: ['dashboard-pending-vendors'],
    queryFn: () => apiGet<Vendor[]>('/admin/vendors?status=pending&per_page=5'),
  });

  const { data: pendingPartnersData } = useQuery({
    queryKey: ['dashboard-pending-partners'],
    queryFn: () => apiGet<DeliveryPartner[]>('/admin/delivery-partners?status=pending&per_page=5'),
  });

  if (statsLoading) return <PageLoading />;

  const stats = statsData?.data || fallbackStats;
  const revenue = revenueData?.data || fallbackRevenueData;
  const orderStatus = orderStatusData?.data || fallbackOrderData;
  const recentOrders = recentOrdersData?.data || [];
  const pendingVendors = pendingVendorsData?.data || [];
  const pendingPartners = pendingPartnersData?.data || [];

  return (
    <div className="space-y-6">
      {/* Stat Cards */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard
          icon={<IndianRupee className="h-6 w-6" />}
          label="Total Revenue"
          value={formatCurrency(stats.total_revenue)}
          change={stats.revenue_change}
          changeType="positive"
        />
        <StatCard
          icon={<ShoppingCart className="h-6 w-6" />}
          label="Total Orders"
          value={stats.total_orders.toLocaleString()}
          change={stats.orders_change}
          changeType="positive"
        />
        <StatCard
          icon={<Store className="h-6 w-6" />}
          label="Active Vendors"
          value={stats.active_vendors}
          change={stats.vendors_change}
          changeType="positive"
        />
        <StatCard
          icon={<Users className="h-6 w-6" />}
          label="Total Users"
          value={stats.total_users.toLocaleString()}
          change={stats.users_change}
          changeType="positive"
        />
      </div>

      {/* Delivery KPIs */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard
          icon={<Bike className="h-6 w-6" />}
          label="Active Delivery Partners"
          value={stats.active_delivery_partners}
          changeType="neutral"
        />
        <StatCard
          icon={<UserCheck className="h-6 w-6" />}
          label="Pending Approvals"
          value={stats.pending_delivery_partners}
          changeType={stats.pending_delivery_partners > 0 ? 'negative' : 'neutral'}
        />
        <StatCard
          icon={<Clock className="h-6 w-6" />}
          label="Avg Delivery Time"
          value={`${stats.avg_delivery_time_mins} min`}
          changeType="neutral"
        />
        <StatCard
          icon={<ShoppingCart className="h-6 w-6" />}
          label="Deliveries Today"
          value={stats.total_deliveries_today}
          changeType="positive"
        />
      </div>

      {/* Charts */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Revenue Trend</CardTitle>
          </CardHeader>
          <CardContent>
            <RevenueChart data={revenue} />
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Orders by Status</CardTitle>
          </CardHeader>
          <CardContent>
            <OrderChart data={orderStatus} />
          </CardContent>
        </Card>
      </div>

      {/* Recent Orders & Pending Vendors & Pending Partners */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-2">
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle>Recent Orders</CardTitle>
            <Link
              href="/orders"
              className="flex items-center gap-1 text-sm text-indigo-600 hover:text-indigo-700"
            >
              View all <ArrowRight className="h-4 w-4" />
            </Link>
          </CardHeader>
          <CardContent>
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead>
                  <tr>
                    <th className="px-3 py-2 text-left text-xs font-medium uppercase text-gray-500">
                      Order
                    </th>
                    <th className="px-3 py-2 text-left text-xs font-medium uppercase text-gray-500">
                      Customer
                    </th>
                    <th className="px-3 py-2 text-left text-xs font-medium uppercase text-gray-500">
                      Total
                    </th>
                    <th className="px-3 py-2 text-left text-xs font-medium uppercase text-gray-500">
                      Status
                    </th>
                    <th className="px-3 py-2 text-left text-xs font-medium uppercase text-gray-500">
                      Date
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  {recentOrders.length === 0 ? (
                    <tr>
                      <td colSpan={5} className="px-3 py-8 text-center text-sm text-gray-500">
                        No recent orders
                      </td>
                    </tr>
                  ) : (
                    recentOrders.map((order) => (
                      <tr key={order.id} className="hover:bg-gray-50">
                        <td className="px-3 py-2 text-sm font-medium text-indigo-600">
                          <Link href={`/orders/${order.id}`}>{order.order_number}</Link>
                        </td>
                        <td className="px-3 py-2 text-sm text-gray-700">
                          {order.customer
                            ? `${order.customer.first_name} ${order.customer.last_name}`
                            : '-'}
                        </td>
                        <td className="px-3 py-2 text-sm text-gray-700">
                          {formatCurrency(order.total_amount)}
                        </td>
                        <td className="px-3 py-2">
                          <StatusBadge status={order.status} />
                        </td>
                        <td className="px-3 py-2 text-sm text-gray-500">
                          {formatDate(order.created_at)}
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle>Pending Vendors</CardTitle>
            <Link
              href="/vendors/pending"
              className="flex items-center gap-1 text-sm text-indigo-600 hover:text-indigo-700"
            >
              View all <ArrowRight className="h-4 w-4" />
            </Link>
          </CardHeader>
          <CardContent>
            {pendingVendors.length === 0 ? (
              <p className="py-8 text-center text-sm text-gray-500">
                No pending vendors
              </p>
            ) : (
              <div className="space-y-3">
                {pendingVendors.map((vendor) => (
                  <Link
                    key={vendor.id}
                    href={`/vendors/${vendor.id}`}
                    className="block rounded-lg border border-gray-100 p-3 hover:bg-gray-50"
                  >
                    <p className="font-medium text-gray-900">
                      {vendor.business_name}
                    </p>
                    <p className="text-xs text-gray-500">
                      {vendor.city} &middot; {vendor.vendor_type || 'product'}
                    </p>
                  </Link>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Pending Delivery Partners */}
      {pendingPartners.length > 0 && (
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle>Pending Delivery Partners</CardTitle>
            <Link
              href="/delivery-partners?status=pending"
              className="flex items-center gap-1 text-sm text-indigo-600 hover:text-indigo-700"
            >
              View all <ArrowRight className="h-4 w-4" />
            </Link>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {pendingPartners.map((partner) => (
                <Link
                  key={partner.id}
                  href={`/delivery-partners/${partner.id}`}
                  className="flex items-center justify-between rounded-lg border border-gray-100 p-3 hover:bg-gray-50"
                >
                  <div>
                    <p className="font-medium text-gray-900">
                      {partner.user?.first_name} {partner.user?.last_name}
                    </p>
                    <p className="text-xs text-gray-500">
                      {partner.vehicle_type} &middot; {partner.vehicle_number || 'No vehicle'}
                    </p>
                  </div>
                  <StatusBadge status="pending" />
                </Link>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
