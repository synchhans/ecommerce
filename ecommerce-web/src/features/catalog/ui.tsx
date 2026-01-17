"use client";

import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { useMemo } from "react";

import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { listProducts } from "@/features/catalog/api";
import { ProductCard } from "@/features/catalog/product-card";
import { ProductGridSkeleton } from "@/features/catalog/product-grid-skeleton";

export default function CatalogClient({
  q,
  cat,
  page,
}: {
  q?: string;
  cat?: string;
  page: number;
}) {
  const queryKey = useMemo(() => ["products", { q: q ?? "", cat: cat ?? "", page }], [q, cat, page]);

  const { data, isLoading, isError, error } = useQuery({
    queryKey,
    queryFn: () => listProducts({ q, cat, page, limit: 24 }),
  });

  const items = data?.items ?? [];

  return (
    <div className="grid gap-6">
      <div className="grid gap-2">
        <div className="text-2xl font-semibold tracking-tight">Catalog</div>
        <div className="text-sm text-muted-foreground">
          Browse products and discover deals.
        </div>
      </div>

      {/* Filters row */}
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex flex-1 items-center gap-2">
          <Input defaultValue={q ?? ""} placeholder="Search in catalog..." />
          <Link href="/catalog">
            <Button variant="outline">Reset</Button>
          </Link>
        </div>

        <div className="flex gap-2">
          <Link href="/catalog?cat=electronics">
            <Button variant={cat === "electronics" ? "default" : "outline"}>Electronics</Button>
          </Link>
          <Link href="/catalog?cat=fashion">
            <Button variant={cat === "fashion" ? "default" : "outline"}>Fashion</Button>
          </Link>
        </div>
      </div>

      {/* Content */}
      {isLoading ? (
        <ProductGridSkeleton />
      ) : isError ? (
        <div className="rounded-xl border bg-background p-6">
          <div className="font-medium">Failed to load products</div>
          <div className="text-sm text-muted-foreground">
            {(error as any)?.message ?? "Unknown error"}
          </div>
        </div>
      ) : items.length === 0 ? (
        <div className="rounded-xl border bg-background p-8 text-center">
          <div className="text-lg font-semibold">No products found</div>
          <div className="text-sm text-muted-foreground">Try a different search or category.</div>
        </div>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {items.map((p: any) => (
            <ProductCard
              key={p.id}
              id={p.id}
              slug={p.slug}
              name={p.name}
              price={p.price}
              image_url={p.image_url}
            />
          ))}
        </div>
      )}
    </div>
  );
}
