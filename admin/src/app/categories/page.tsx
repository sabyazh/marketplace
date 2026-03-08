'use client';

import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { ChevronDown, ChevronRight, Plus, Pencil, Trash2 } from 'lucide-react';
import { apiGet, apiPost, apiPut, apiDelete } from '@/lib/api-client';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Dialog } from '@/components/ui/dialog';
import { InputField } from '@/components/ui/input';
import { Select } from '@/components/ui/select';
import { PageLoading } from '@/components/ui/loading';
import type { Category, CategoryType } from '@/types';

type TabType = 'product' | 'service';

export default function CategoriesPage() {
  const queryClient = useQueryClient();
  const [activeTab, setActiveTab] = useState<TabType>('product');
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingCategory, setEditingCategory] = useState<Category | null>(null);
  const [expandedIds, setExpandedIds] = useState<Set<string>>(new Set());

  // Form state
  const [formName, setFormName] = useState('');
  const [formSlug, setFormSlug] = useState('');
  const [formType, setFormType] = useState<CategoryType>('product');
  const [formParentId, setFormParentId] = useState('');

  const { data, isLoading } = useQuery({
    queryKey: ['categories', activeTab],
    queryFn: () => apiGet<Category[]>(`/categories?type=${activeTab}`),
  });

  const createMutation = useMutation({
    mutationFn: (payload: {
      name: string;
      slug: string;
      category_type: CategoryType;
      parent_id?: string;
    }) => apiPost('/admin/categories', payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] });
      closeDialog();
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({
      id,
      payload,
    }: {
      id: string;
      payload: { name: string; slug: string; parent_id?: string };
    }) => apiPut(`/admin/categories/${id}`, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] });
      closeDialog();
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => apiDelete(`/admin/categories/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] });
    },
  });

  const categories = data?.data || [];

  // Build tree structure
  const rootCategories = categories.filter((c) => !c.parent_id);
  const childMap = new Map<string, Category[]>();
  categories.forEach((c) => {
    if (c.parent_id) {
      const existing = childMap.get(c.parent_id) || [];
      existing.push(c);
      childMap.set(c.parent_id, existing);
    }
  });

  function toggleExpand(id: string) {
    setExpandedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  }

  function openCreate() {
    setEditingCategory(null);
    setFormName('');
    setFormSlug('');
    setFormType(activeTab);
    setFormParentId('');
    setDialogOpen(true);
  }

  function openEdit(cat: Category) {
    setEditingCategory(cat);
    setFormName(cat.name);
    setFormSlug(cat.slug);
    setFormType(cat.category_type);
    setFormParentId(cat.parent_id || '');
    setDialogOpen(true);
  }

  function closeDialog() {
    setDialogOpen(false);
    setEditingCategory(null);
    setFormName('');
    setFormSlug('');
    setFormParentId('');
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (editingCategory) {
      updateMutation.mutate({
        id: editingCategory.id,
        payload: {
          name: formName,
          slug: formSlug,
          parent_id: formParentId || undefined,
        },
      });
    } else {
      createMutation.mutate({
        name: formName,
        slug: formSlug,
        category_type: formType,
        parent_id: formParentId || undefined,
      });
    }
  }

  function handleNameChange(name: string) {
    setFormName(name);
    if (!editingCategory) {
      setFormSlug(
        name
          .toLowerCase()
          .replace(/[^a-z0-9]+/g, '-')
          .replace(/^-|-$/g, '')
      );
    }
  }

  if (isLoading) return <PageLoading />;

  return (
    <div className="space-y-4">
      {/* Tabs */}
      <div className="flex items-center justify-between">
        <div className="flex border-b border-gray-200">
          <button
            className={`border-b-2 px-4 py-2 text-sm font-medium ${
              activeTab === 'product'
                ? 'border-indigo-600 text-indigo-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
            onClick={() => setActiveTab('product')}
          >
            Product Categories
          </button>
          <button
            className={`border-b-2 px-4 py-2 text-sm font-medium ${
              activeTab === 'service'
                ? 'border-indigo-600 text-indigo-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
            onClick={() => setActiveTab('service')}
          >
            Service Categories
          </button>
        </div>
        <Button onClick={openCreate}>
          <Plus className="h-4 w-4" />
          Add Category
        </Button>
      </div>

      <Card>
        <CardContent className="p-0">
          {rootCategories.length === 0 ? (
            <div className="py-12 text-center text-sm text-gray-500">
              No categories found. Create one to get started.
            </div>
          ) : (
            <div className="divide-y divide-gray-100">
              {rootCategories.map((cat) => (
                <CategoryRow
                  key={cat.id}
                  category={cat}
                  childMap={childMap}
                  expandedIds={expandedIds}
                  onToggle={toggleExpand}
                  onEdit={openEdit}
                  onDelete={(id) => deleteMutation.mutate(id)}
                  depth={0}
                />
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Create/Edit Dialog */}
      <Dialog
        open={dialogOpen}
        onClose={closeDialog}
        title={editingCategory ? 'Edit Category' : 'Add Category'}
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <InputField
            label="Name"
            value={formName}
            onChange={(e) => handleNameChange(e.target.value)}
            required
          />
          <InputField
            label="Slug"
            value={formSlug}
            onChange={(e) => setFormSlug(e.target.value)}
            required
          />
          {!editingCategory && (
            <Select
              label="Type"
              value={formType}
              onChange={(e) => setFormType(e.target.value as CategoryType)}
            >
              <option value="product">Product</option>
              <option value="service">Service</option>
            </Select>
          )}
          <Select
            label="Parent Category"
            value={formParentId}
            onChange={(e) => setFormParentId(e.target.value)}
          >
            <option value="">None (Root Category)</option>
            {rootCategories
              .filter((c) => c.id !== editingCategory?.id)
              .map((c) => (
                <option key={c.id} value={c.id}>
                  {c.name}
                </option>
              ))}
          </Select>
          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={closeDialog}>
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={createMutation.isPending || updateMutation.isPending}
            >
              {editingCategory ? 'Update' : 'Create'}
            </Button>
          </div>
        </form>
      </Dialog>
    </div>
  );
}

function CategoryRow({
  category,
  childMap,
  expandedIds,
  onToggle,
  onEdit,
  onDelete,
  depth,
}: {
  category: Category;
  childMap: Map<string, Category[]>;
  expandedIds: Set<string>;
  onToggle: (id: string) => void;
  onEdit: (cat: Category) => void;
  onDelete: (id: string) => void;
  depth: number;
}) {
  const children = childMap.get(category.id) || [];
  const hasChildren = children.length > 0;
  const isExpanded = expandedIds.has(category.id);

  return (
    <>
      <div
        className="flex items-center justify-between px-4 py-3 hover:bg-gray-50"
        style={{ paddingLeft: `${16 + depth * 24}px` }}
      >
        <div className="flex items-center gap-2">
          {hasChildren ? (
            <button
              onClick={() => onToggle(category.id)}
              className="rounded p-0.5 hover:bg-gray-200"
            >
              {isExpanded ? (
                <ChevronDown className="h-4 w-4 text-gray-500" />
              ) : (
                <ChevronRight className="h-4 w-4 text-gray-500" />
              )}
            </button>
          ) : (
            <span className="w-5" />
          )}
          <span className="text-sm font-medium text-gray-900">{category.name}</span>
          <span className="text-xs text-gray-400">/{category.slug}</span>
          {!category.is_active && (
            <span className="rounded bg-gray-100 px-1.5 py-0.5 text-xs text-gray-500">
              Inactive
            </span>
          )}
        </div>
        <div className="flex items-center gap-1">
          <Button variant="ghost" size="sm" onClick={() => onEdit(category)}>
            <Pencil className="h-3.5 w-3.5" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onDelete(category.id)}
            className="text-red-500 hover:text-red-700"
          >
            <Trash2 className="h-3.5 w-3.5" />
          </Button>
        </div>
      </div>
      {isExpanded &&
        children.map((child) => (
          <CategoryRow
            key={child.id}
            category={child}
            childMap={childMap}
            expandedIds={expandedIds}
            onToggle={onToggle}
            onEdit={onEdit}
            onDelete={onDelete}
            depth={depth + 1}
          />
        ))}
    </>
  );
}
