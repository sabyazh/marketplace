'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useQuery } from '@tanstack/react-query';
import { type ColumnDef } from '@tanstack/react-table';
import { Eye, Star } from 'lucide-react';
import { apiGet } from '@/lib/api-client';
import { DataTable } from '@/components/ui/data-table';
import { StatusBadge } from '@/components/ui/status-badge';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Select } from '@/components/ui/select';
import { Button } from '@/components/ui/button';
import type { Vendor, PaginationMeta } from '@/types';

const columns: ColumnDef<Vendor, unknown>[] = [
  {
    accessorKey: 'business_name',
    header: 'Business Name',
    cell: ({ row }) => (
      <Link
        href={`/vendors/${row.original.id}`}
        className="font-medium text-indigo-600 hover:text-indigo-700"
      >
        {row.original.business_name}
      </Link>
    ),
  },
  {
    accessorKey: 'vendor_type',
    header: 'Type',
    cell: ({ row }) => (
      <Badge variant="indigo">
        {(row.original.vendor_type || 'product').replace(/_/g, ' ').replace(/\b\w/g, (c: string) => c.toUpperCase())}
      </Badge>
    ),
  },
  {
    accessorKey: 'status',
    header: 'Status',
    cell: ({ row }) => <StatusBadge status={row.original.status} />,
  },
  {
    accessorKey: 'city',
    header: 'City',
  },
  {
    accessorKey: 'avg_rating',
    header: 'Rating',
    cell: ({ row }) => (
      <div className="flex items-center gap-1">
        <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
        <span>{(row.original.avg_rating || 0).toFixed(1)}</span>
      </div>
    ),
  },
  {
    accessorKey: 'commission_pct',
    header: 'Commission %',
    cell: ({ row }) => `${row.original.commission_pct || 0}%`,
  },
  {
    id: 'actions',
    header: 'Actions',
    cell: ({ row }) => (
      <Link href={`/vendors/${row.original.id}`}>
        <Button variant="ghost" size="sm">
          <Eye className="h-4 w-4" />
          View
        </Button>
      </Link>
    ),
  },
];

export default function VendorsPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const pageSize = 10;

  const { data, isLoading } = useQuery({
    queryKey: ['vendors', page, search, statusFilter],
    queryFn: () =>
      apiGet<Vendor[]>(
        `/admin/vendors?page=${page}&per_page=${pageSize}${search ? `&search=${search}` : ''}${statusFilter ? `&status=${statusFilter}` : ''}`
      ),
  });

  const vendors = data?.data || [];
  const meta = (data?.meta || { total_pages: 1 }) as PaginationMeta;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Input
            placeholder="Search by business name..."
            value={search}
            onChange={(e) => {
              setSearch(e.target.value);
              setPage(1);
            }}
            className="w-64"
          />
          <Select
            value={statusFilter}
            onChange={(e) => {
              setStatusFilter(e.target.value);
              setPage(1);
            }}
            className="w-40"
          >
            <option value="">All Status</option>
            <option value="pending">Pending</option>
            <option value="approved">Approved</option>
            <option value="rejected">Rejected</option>
            <option value="suspended">Suspended</option>
          </Select>
        </div>
        <Link href="/vendors/pending">
          <Button variant="outline">Pending Approvals</Button>
        </Link>
      </div>

      <DataTable
        columns={columns}
        data={vendors}
        pageCount={meta.total_pages}
        page={page}
        onPageChange={setPage}
        pageSize={pageSize}
        isLoading={isLoading}
      />
    </div>
  );
}
