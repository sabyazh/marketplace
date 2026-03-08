'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useQuery } from '@tanstack/react-query';
import { type ColumnDef } from '@tanstack/react-table';
import { Eye, Package } from 'lucide-react';
import { apiGet } from '@/lib/api-client';
import { formatCurrency } from '@/lib/utils';
import { DataTable } from '@/components/ui/data-table';
import { StatusBadge } from '@/components/ui/status-badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import type { Product, PaginationMeta } from '@/types';

const columns: ColumnDef<Product, unknown>[] = [
  {
    accessorKey: 'thumbnail_url',
    header: 'Image',
    enableSorting: false,
    cell: ({ row }) => (
      <div className="h-10 w-10 overflow-hidden rounded-lg border border-gray-200 bg-gray-50">
        {row.original.thumbnail_url ? (
          <img
            src={row.original.thumbnail_url}
            alt={row.original.name}
            className="h-full w-full object-cover"
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center">
            <Package className="h-5 w-5 text-gray-300" />
          </div>
        )}
      </div>
    ),
  },
  {
    accessorKey: 'name',
    header: 'Name',
    cell: ({ row }) => (
      <span className="font-medium text-gray-900">{row.original.name}</span>
    ),
  },
  {
    accessorKey: 'vendor.business_name',
    header: 'Vendor',
    cell: ({ row }) => (
      <span className="text-gray-600">
        {row.original.vendor?.business_name || '-'}
      </span>
    ),
  },
  {
    accessorKey: 'category.name',
    header: 'Category',
    cell: ({ row }) => (
      <span className="text-gray-600">
        {row.original.category?.name || '-'}
      </span>
    ),
  },
  {
    accessorKey: 'price',
    header: 'Price',
    cell: ({ row }) => formatCurrency(row.original.price),
  },
  {
    accessorKey: 'stock_quantity',
    header: 'Stock',
    cell: ({ row }) => {
      const stock = row.original.stock_quantity;
      const low = row.original.low_stock_threshold;
      return (
        <span
          className={
            stock <= 0
              ? 'font-medium text-red-600'
              : stock <= low
                ? 'font-medium text-yellow-600'
                : 'text-gray-700'
          }
        >
          {stock}
        </span>
      );
    },
  },
  {
    accessorKey: 'status',
    header: 'Status',
    cell: ({ row }) => {
      const p = row.original;
      // Backend returns is_active (boolean), derive status string
      const isActive = (p as unknown as Record<string, unknown>).is_active;
      const status = isActive === false
        ? 'inactive'
        : p.stock_quantity <= 0
          ? 'out_of_stock'
          : p.status || 'active';
      return <StatusBadge status={status} />;
    },
  },
  {
    id: 'actions',
    header: 'Actions',
    enableSorting: false,
    cell: ({ row }) => (
      <Link href={`/products/${row.original.id}`}>
        <Button variant="ghost" size="sm">
          <Eye className="h-4 w-4" />
        </Button>
      </Link>
    ),
  },
];

export default function ProductsPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const pageSize = 10;

  const { data, isLoading } = useQuery({
    queryKey: ['products', page, search],
    queryFn: () =>
      apiGet<Product[]>(
        `/admin/products?page=${page}&per_page=${pageSize}${search ? `&search=${search}` : ''}`
      ),
  });

  const products = data?.data || [];
  const meta = (data?.meta || { total_pages: 1 }) as PaginationMeta;

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-3">
        <Input
          placeholder="Search products..."
          value={search}
          onChange={(e) => {
            setSearch(e.target.value);
            setPage(1);
          }}
          className="w-64"
        />
      </div>

      <DataTable
        columns={columns}
        data={products}
        pageCount={meta.total_pages}
        page={page}
        onPageChange={setPage}
        pageSize={pageSize}
        isLoading={isLoading}
      />
    </div>
  );
}
