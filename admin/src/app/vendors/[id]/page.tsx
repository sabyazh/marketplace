'use client';

import { useState } from 'react';
import { useParams } from 'next/navigation';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Star,
  MapPin,
  Phone,
  Mail,
  Globe,
  Check,
  X,
  Ban,
  ShoppingCart,
  IndianRupee,
  FileText,
  Image as ImageIcon,
} from 'lucide-react';
import { apiGet, apiPut, apiPost } from '@/lib/api-client';
import { formatCurrency, formatDate } from '@/lib/utils';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { StatusBadge } from '@/components/ui/status-badge';
import { StatCard } from '@/components/ui/stat-card';
import { Input } from '@/components/ui/input';
import { PageLoading } from '@/components/ui/loading';
import type { Vendor, VendorDocument } from '@/types';

interface VendorStats {
  total_orders: number;
  total_revenue: number;
  avg_rating: number;
}

export default function VendorDetailPage() {
  const params = useParams();
  const vendorId = params.id as string;
  const queryClient = useQueryClient();
  const [commissionRate, setCommissionRate] = useState<string>('');

  const { data: vendorData, isLoading } = useQuery({
    queryKey: ['vendor', vendorId],
    queryFn: () => apiGet<Vendor>(`/admin/vendors/${vendorId}`),
  });

  const { data: statsData } = useQuery({
    queryKey: ['vendor-stats', vendorId],
    queryFn: () => apiGet<VendorStats>(`/admin/vendors/${vendorId}/stats`),
  });

  const statusMutation = useMutation({
    mutationFn: ({ status }: { status: string }) =>
      apiPut(`/admin/vendors/${vendorId}/status`, { status }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['vendor', vendorId] });
    },
  });

  const commissionMutation = useMutation({
    mutationFn: (rate: number) =>
      apiPut(`/admin/vendors/${vendorId}/commission`, { commission_pct: rate }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['vendor', vendorId] });
    },
  });

  const docMutation = useMutation({
    mutationFn: ({ docId, status }: { docId: string; status: string }) =>
      apiPost(`/admin/vendors/${vendorId}/documents/${docId}/verify`, { status }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['vendor', vendorId] });
    },
  });

  if (isLoading) return <PageLoading />;

  const vendor = vendorData?.data;
  const stats = statsData?.data;

  if (!vendor) {
    return (
      <div className="flex h-[60vh] items-center justify-center">
        <p className="text-gray-500">Vendor not found</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Vendor Info & Actions */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-2">
          <CardHeader>
            <div className="flex items-start justify-between">
              <div>
                <CardTitle className="text-xl">{vendor.business_name}</CardTitle>
                <div className="mt-2 flex items-center gap-2">
                  <StatusBadge status={vendor.status} />
                  <Badge variant={vendor.is_online ? 'success' : 'default'}>
                    {vendor.is_online ? 'Online' : 'Offline'}
                  </Badge>
                </div>
              </div>
              <div className="flex items-center gap-1">
                <Star className="h-5 w-5 fill-yellow-400 text-yellow-400" />
                <span className="text-lg font-semibold">{(vendor.avg_rating || 0).toFixed(1)}</span>
                <span className="text-sm text-gray-500">({vendor.total_reviews} reviews)</span>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div className="space-y-3">
                <div className="flex items-center gap-2 text-sm text-gray-600">
                  <Badge variant="indigo">{(vendor.vendor_type || 'product').replace(/_/g, ' ').replace(/\b\w/g, (c: string) => c.toUpperCase())}</Badge>
                </div>
                <div className="flex items-center gap-2 text-sm text-gray-600">
                  <MapPin className="h-4 w-4 text-gray-400" />
                  <span>{vendor.address}, {vendor.city}, {vendor.state} - {vendor.pincode}</span>
                </div>
                {vendor.user?.phone && (
                  <div className="flex items-center gap-2 text-sm text-gray-600">
                    <Phone className="h-4 w-4 text-gray-400" />
                    <span>{vendor.user.phone}</span>
                  </div>
                )}
                {vendor.user?.email && (
                  <div className="flex items-center gap-2 text-sm text-gray-600">
                    <Mail className="h-4 w-4 text-gray-400" />
                    <span>{vendor.user.email}</span>
                  </div>
                )}
              </div>
              <div className="space-y-3">
                {vendor.service_radius_km > 0 && (
                  <div className="text-sm">
                    <span className="font-medium text-gray-500">Service Radius: </span>
                    <span className="text-gray-700">{vendor.service_radius_km} km</span>
                  </div>
                )}
                <div className="text-sm">
                  <span className="font-medium text-gray-500">Joined: </span>
                  <span className="text-gray-700">{formatDate(vendor.created_at)}</span>
                </div>
              </div>
            </div>
            {vendor.description && (
              <p className="mt-4 text-sm text-gray-600">{vendor.description}</p>
            )}
          </CardContent>
        </Card>

        {/* Action buttons */}
        <Card>
          <CardHeader>
            <CardTitle>Actions</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {vendor.status === 'pending' && (
              <>
                <Button
                  variant="success"
                  className="w-full"
                  onClick={() => statusMutation.mutate({ status: 'approved' })}
                  disabled={statusMutation.isPending}
                >
                  <Check className="h-4 w-4" />
                  Approve Vendor
                </Button>
                <Button
                  variant="destructive"
                  className="w-full"
                  onClick={() => statusMutation.mutate({ status: 'rejected' })}
                  disabled={statusMutation.isPending}
                >
                  <X className="h-4 w-4" />
                  Reject Vendor
                </Button>
              </>
            )}
            {vendor.status === 'approved' && (
              <Button
                variant="warning"
                className="w-full"
                onClick={() => statusMutation.mutate({ status: 'suspended' })}
                disabled={statusMutation.isPending}
              >
                <Ban className="h-4 w-4" />
                Suspend Vendor
              </Button>
            )}
            {vendor.status === 'suspended' && (
              <Button
                variant="success"
                className="w-full"
                onClick={() => statusMutation.mutate({ status: 'approved' })}
                disabled={statusMutation.isPending}
              >
                <Check className="h-4 w-4" />
                Reactivate Vendor
              </Button>
            )}
            {vendor.status === 'rejected' && (
              <Button
                variant="success"
                className="w-full"
                onClick={() => statusMutation.mutate({ status: 'approved' })}
                disabled={statusMutation.isPending}
              >
                <Check className="h-4 w-4" />
                Approve Vendor
              </Button>
            )}

            <hr className="my-4" />

            <div className="space-y-2">
              <label className="text-sm font-medium text-gray-700">
                Commission Rate (%)
              </label>
              <div className="flex gap-2">
                <Input
                  type="number"
                  min="0"
                  max="100"
                  step="0.5"
                  placeholder={(vendor.commission_pct || 0).toString()}
                  value={commissionRate}
                  onChange={(e) => setCommissionRate(e.target.value)}
                />
                <Button
                  onClick={() => {
                    const rate = parseFloat(commissionRate);
                    if (!isNaN(rate) && rate >= 0 && rate <= 100) {
                      commissionMutation.mutate(rate);
                    }
                  }}
                  disabled={commissionMutation.isPending}
                  size="sm"
                >
                  Save
                </Button>
              </div>
              <p className="text-xs text-gray-500">
                Current: {vendor.commission_pct || 0}%
              </p>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Stats */}
      {stats && (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <StatCard
            icon={<ShoppingCart className="h-6 w-6" />}
            label="Total Orders"
            value={stats.total_orders}
          />
          <StatCard
            icon={<IndianRupee className="h-6 w-6" />}
            label="Total Revenue"
            value={formatCurrency(stats.total_revenue)}
          />
          <StatCard
            icon={<Star className="h-6 w-6" />}
            label="Avg Rating"
            value={stats.avg_rating.toFixed(1)}
          />
        </div>
      )}

      {/* Documents */}
      {vendor.documents && vendor.documents.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Documents</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="divide-y divide-gray-100">
              {vendor.documents.map((doc: VendorDocument) => (
                <div
                  key={doc.id}
                  className="flex items-center justify-between py-3"
                >
                  <div className="flex items-center gap-3">
                    {doc.document_url.match(/\.(jpg|jpeg|png|gif|webp)$/i) ? (
                      <ImageIcon className="h-5 w-5 text-gray-400" />
                    ) : (
                      <FileText className="h-5 w-5 text-gray-400" />
                    )}
                    <div>
                      <p className="text-sm font-medium text-gray-900">
                        {doc.document_type.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase())}
                      </p>
                      <a
                        href={doc.document_url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-xs text-indigo-600 hover:text-indigo-700"
                      >
                        View Document
                      </a>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <StatusBadge status={doc.status} />
                    {doc.status === 'pending' && (
                      <>
                        <Button
                          variant="success"
                          size="sm"
                          onClick={() =>
                            docMutation.mutate({ docId: doc.id, status: 'approved' })
                          }
                          disabled={docMutation.isPending}
                        >
                          <Check className="h-3 w-3" />
                        </Button>
                        <Button
                          variant="destructive"
                          size="sm"
                          onClick={() =>
                            docMutation.mutate({ docId: doc.id, status: 'rejected' })
                          }
                          disabled={docMutation.isPending}
                        >
                          <X className="h-3 w-3" />
                        </Button>
                      </>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
