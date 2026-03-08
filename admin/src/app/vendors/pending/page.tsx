'use client';

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import Link from 'next/link';
import { Check, X, Eye, MapPin } from 'lucide-react';
import { apiGet, apiPut } from '@/lib/api-client';
import { formatDate } from '@/lib/utils';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { PageLoading } from '@/components/ui/loading';
import type { Vendor } from '@/types';

export default function PendingVendorsPage() {
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['pending-vendors'],
    queryFn: () => apiGet<Vendor[]>('/admin/vendors?status=pending&per_page=50'),
  });

  const statusMutation = useMutation({
    mutationFn: ({ vendorId, status }: { vendorId: string; status: string }) =>
      apiPut(`/admin/vendors/${vendorId}/status`, { status }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pending-vendors'] });
    },
  });

  if (isLoading) return <PageLoading />;

  const vendors = data?.data || [];

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-sm text-gray-500">
          {vendors.length} vendor{vendors.length !== 1 ? 's' : ''} pending approval
        </p>
      </div>

      {vendors.length === 0 ? (
        <Card>
          <CardContent className="py-12 text-center">
            <p className="text-gray-500">No pending vendor applications</p>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-3">
          {vendors.map((vendor) => (
            <Card key={vendor.id}>
              <CardContent className="p-4">
                <div className="flex items-center justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-3">
                      <h3 className="font-semibold text-gray-900">
                        {vendor.business_name}
                      </h3>
                      <Badge variant="indigo">
                        {(vendor.vendor_type || 'product').replace(/_/g, ' ').replace(/\b\w/g, (c: string) => c.toUpperCase())}
                      </Badge>
                    </div>
                    <div className="mt-1 flex items-center gap-4 text-sm text-gray-500">
                      <span className="flex items-center gap-1">
                        <MapPin className="h-3.5 w-3.5" />
                        {vendor.city}, {vendor.state}
                      </span>
                      <span>Applied: {formatDate(vendor.created_at)}</span>
                    </div>
                    {vendor.description && (
                      <p className="mt-1 text-sm text-gray-600 line-clamp-1">
                        {vendor.description}
                      </p>
                    )}
                  </div>
                  <div className="flex items-center gap-2">
                    <Link href={`/vendors/${vendor.id}`}>
                      <Button variant="ghost" size="sm">
                        <Eye className="h-4 w-4" />
                        View
                      </Button>
                    </Link>
                    <Button
                      variant="success"
                      size="sm"
                      onClick={() =>
                        statusMutation.mutate({ vendorId: vendor.id, status: 'approved' })
                      }
                      disabled={statusMutation.isPending}
                    >
                      <Check className="h-4 w-4" />
                      Approve
                    </Button>
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() =>
                        statusMutation.mutate({ vendorId: vendor.id, status: 'rejected' })
                      }
                      disabled={statusMutation.isPending}
                    >
                      <X className="h-4 w-4" />
                      Reject
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
