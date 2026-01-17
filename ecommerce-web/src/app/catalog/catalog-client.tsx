"use client";

import React, { useEffect, useState } from "react";
import { motion } from "framer-motion";
import Image from "next/image";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { ShoppingCart, Filter, ChevronDown, Grid3X3, LayoutGrid, Star } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { cn } from "@/lib/utils";

type Product = {
  id: string;
  name: string;
  slug: string;
  price: number;
  image_url: string;
  category?: string;
  rating?: number;
};

const CATEGORIES = [
  { name: "All", value: "" },
  { name: "Electronics", value: "electronics" },
  { name: "Fashion", value: "fashion" },
  { name: "Beauty", value: "beauty" },
  { name: "Home & Living", value: "home" },
];

const STAGGER = {
  visible: { transition: { staggerChildren: 0.05 } },
};

const FADE_UP = {
  hidden: { opacity: 0, y: 20 },
  visible: { opacity: 1, y: 0, transition: { duration: 0.4, ease: "easeOut" } },
};

export default function CatalogClient() {
  const searchParams = useSearchParams();
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeCategory, setActiveCategory] = useState(searchParams.get("cat") || "");
  const [viewMode, setViewMode] = useState<"grid" | "large">("grid");

  useEffect(() => {
    setLoading(true);
    const url = new URL("http://localhost:8080/v1/products");
    if (activeCategory) url.searchParams.set("category", activeCategory);

    fetch(url.toString())
      .then((res) => res.json())
      .then((data) => setProducts(data.items ?? data.data ?? []))
      .catch(() => setProducts([]))
      .finally(() => setLoading(false));
  }, [activeCategory]);

  return (
    <div className="min-h-screen bg-background">
      {/* Hero Banner */}
      <section className="relative bg-gradient-to-br from-primary/10 via-background to-background py-16">
        <div className="container mx-auto px-6">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="max-w-2xl"
          >
            <Badge variant="secondary" className="mb-4">Shop Collection</Badge>
            <h1 className="text-4xl font-bold tracking-tight mb-4">
              {activeCategory ? CATEGORIES.find(c => c.value === activeCategory)?.name : "All Products"}
            </h1>
            <p className="text-muted-foreground text-lg">
              Discover our curated selection of premium products.
            </p>
          </motion.div>
        </div>
      </section>

      <div className="container mx-auto px-6 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-[280px_1fr] gap-8">
          {/* Sidebar Filters */}
          <aside className="hidden lg:block space-y-6">
            <div className="sticky top-24">
              <div className="rounded-2xl border bg-card p-6 shadow-sm">
                <h3 className="font-semibold mb-4 flex items-center gap-2">
                  <Filter className="h-4 w-4" /> Filters
                </h3>

                <div className="space-y-4">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground mb-3">Category</p>
                    <div className="space-y-2">
                      {CATEGORIES.map((cat) => (
                        <button
                          key={cat.value}
                          onClick={() => setActiveCategory(cat.value)}
                          className={cn(
                            "w-full text-left px-3 py-2 rounded-lg text-sm transition-colors",
                            activeCategory === cat.value
                              ? "bg-primary text-primary-foreground font-medium"
                              : "hover:bg-secondary"
                          )}
                        >
                          {cat.name}
                        </button>
                      ))}
                    </div>
                  </div>

                  <Separator />

                  <div>
                    <p className="text-sm font-medium text-muted-foreground mb-3">Price Range</p>
                    <div className="flex gap-2">
                      <Input placeholder="Min" className="h-9" type="number" />
                      <Input placeholder="Max" className="h-9" type="number" />
                    </div>
                  </div>

                  <Separator />

                  <div>
                    <p className="text-sm font-medium text-muted-foreground mb-3">Rating</p>
                    <div className="space-y-2">
                      {[4, 3, 2].map((r) => (
                        <button key={r} className="flex items-center gap-2 w-full text-left px-3 py-2 rounded-lg text-sm hover:bg-secondary transition-colors">
                          <div className="flex text-yellow-500">
                            {Array.from({ length: r }).map((_, i) => <Star key={i} className="h-3.5 w-3.5 fill-current" />)}
                          </div>
                          <span className="text-muted-foreground">& Up</span>
                        </button>
                      ))}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </aside>

          {/* Main Content */}
          <div className="space-y-6">
            {/* Toolbar */}
            <div className="flex items-center justify-between gap-4 p-4 rounded-xl bg-card border">
              <p className="text-sm text-muted-foreground">
                Showing <span className="font-medium text-foreground">{products.length}</span> products
              </p>

              <div className="flex items-center gap-4">
                <div className="hidden sm:flex items-center gap-1 border rounded-lg p-1">
                  <button
                    onClick={() => setViewMode("grid")}
                    className={cn("p-1.5 rounded", viewMode === "grid" ? "bg-secondary" : "hover:bg-secondary/50")}
                  >
                    <Grid3X3 className="h-4 w-4" />
                  </button>
                  <button
                    onClick={() => setViewMode("large")}
                    className={cn("p-1.5 rounded", viewMode === "large" ? "bg-secondary" : "hover:bg-secondary/50")}
                  >
                    <LayoutGrid className="h-4 w-4" />
                  </button>
                </div>

                <select className="h-9 rounded-lg border bg-background px-3 text-sm focus:ring-2 focus:ring-primary/20 outline-none">
                  <option>Sort by: Popular</option>
                  <option>Newest</option>
                  <option>Price: Low → High</option>
                  <option>Price: High → Low</option>
                </select>
              </div>
            </div>

            {/* Product Grid */}
            {loading ? (
              <div className={cn(
                "grid gap-4",
                viewMode === "grid" ? "grid-cols-2 md:grid-cols-3 lg:grid-cols-4" : "grid-cols-1 md:grid-cols-2"
              )}>
                {Array.from({ length: 8 }).map((_, i) => (
                  <div key={i} className="aspect-[3/4] rounded-2xl bg-secondary/50 animate-pulse" />
                ))}
              </div>
            ) : products.length === 0 ? (
              <div className="py-20 text-center">
                <p className="text-muted-foreground">No products found.</p>
              </div>
            ) : (
              <motion.div
                initial="hidden"
                animate="visible"
                variants={STAGGER}
                className={cn(
                  "grid gap-4",
                  viewMode === "grid" ? "grid-cols-2 md:grid-cols-3 lg:grid-cols-4" : "grid-cols-1 md:grid-cols-2"
                )}
              >
                {products.map((product) => (
                  <ProductCard key={product.id} product={product} viewMode={viewMode} />
                ))}
              </motion.div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

function ProductCard({ product, viewMode }: { product: Product; viewMode: string }) {
  return (
    <motion.div variants={FADE_UP}>
      <Link href={`/product/${product.slug || product.id}`}>
        <div className={cn(
          "group relative rounded-2xl border bg-card overflow-hidden transition-all duration-300 hover:shadow-xl hover:border-primary/20",
          viewMode === "large" ? "flex gap-4" : ""
        )}>
          {/* Image */}
          <div className={cn(
            "relative overflow-hidden bg-secondary/50",
            viewMode === "large" ? "w-48 h-48 flex-shrink-0" : "aspect-square"
          )}>
            <Image
              src={product.image_url || "/placeholder.png"}
              alt={product.name}
              fill
              className="object-cover transition-transform duration-500 group-hover:scale-110"
            />

            {/* Quick Add Button */}
            <motion.div
              initial={{ opacity: 0, y: 10 }}
              whileHover={{ opacity: 1, y: 0 }}
              className="absolute bottom-3 left-3 right-3 opacity-0 group-hover:opacity-100 transition-opacity"
            >
              <Button size="sm" className="w-full rounded-lg gap-2 shadow-lg">
                <ShoppingCart className="h-4 w-4" /> Add to Cart
              </Button>
            </motion.div>
          </div>

          {/* Content */}
          <div className={cn("p-4", viewMode === "large" ? "flex-1 flex flex-col justify-center" : "")}>
            {product.category && (
              <span className="text-xs text-muted-foreground uppercase tracking-wider">{product.category}</span>
            )}
            <h3 className="font-medium line-clamp-2 mt-1 group-hover:text-primary transition-colors">
              {product.name}
            </h3>

            <div className="flex items-center gap-2 mt-2">
              <div className="flex text-yellow-500">
                {Array.from({ length: 5 }).map((_, i) => (
                  <Star key={i} className={cn("h-3 w-3", i < (product.rating || 4) ? "fill-current" : "fill-muted stroke-muted-foreground/30")} />
                ))}
              </div>
              <span className="text-xs text-muted-foreground">({Math.floor(Math.random() * 100) + 10})</span>
            </div>

            <p className="mt-3 text-lg font-bold text-primary">
              Rp {product.price?.toLocaleString("id-ID") || "0"}
            </p>
          </div>
        </div>
      </Link>
    </motion.div>
  );
}
