"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import Image from "next/image";
import { useParams } from "next/navigation";
import { motion } from "framer-motion";
import { ShoppingCart, Heart, Share2, Star, Truck, ShieldCheck, RefreshCw, ChevronRight, Minus, Plus } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { useCartStore } from "@/store/cart";
import { cn } from "@/lib/utils";

type Product = {
    id: string;
    name: string;
    slug: string;
    price: number;
    description: string;
    image_url: string;
    category?: string;
    stock?: number;
};

export default function ProductDetailPage() {
    const params = useParams();
    const slug = params.slug as string;
    const [product, setProduct] = useState<Product | null>(null);
    const [loading, setLoading] = useState(true);
    const [quantity, setQuantity] = useState(1);
    const [selectedImage, setSelectedImage] = useState(0);
    const addItem = useCartStore((s) => s.addItem);

    useEffect(() => {
        fetch(`http://localhost:8080/v1/products/${slug}`)
            .then((res) => res.json())
            .then((data) => setProduct(data))
            .catch(() => setProduct(null))
            .finally(() => setLoading(false));
    }, [slug]);

    const handleAddToCart = () => {
        if (!product) return;
        for (let i = 0; i < quantity; i++) {
            addItem({
                id: product.id,
                name: product.name,
                price: product.price,
                image_url: product.image_url,
            });
        }
    };

    if (loading) {
        return <ProductSkeleton />;
    }

    if (!product) {
        return (
            <div className="min-h-[60vh] flex flex-col items-center justify-center">
                <h2 className="text-2xl font-bold">Product not found</h2>
                <p className="text-muted-foreground mt-2">The product you're looking for doesn't exist.</p>
                <Link href="/catalog">
                    <Button className="mt-4">Browse Products</Button>
                </Link>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-background">
            {/* Breadcrumb */}
            <div className="container mx-auto px-6 py-4">
                <nav className="flex items-center gap-2 text-sm text-muted-foreground">
                    <Link href="/" className="hover:text-foreground">Home</Link>
                    <ChevronRight className="h-4 w-4" />
                    <Link href="/catalog" className="hover:text-foreground">Catalog</Link>
                    <ChevronRight className="h-4 w-4" />
                    <span className="text-foreground">{product.name}</span>
                </nav>
            </div>

            <div className="container mx-auto px-6 py-8">
                <div className="grid lg:grid-cols-2 gap-12">
                    {/* Image Gallery */}
                    <div className="space-y-4">
                        <motion.div
                            key={selectedImage}
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            className="aspect-square relative rounded-2xl overflow-hidden bg-secondary"
                        >
                            <Image
                                src={product.image_url || "/placeholder.png"}
                                alt={product.name}
                                fill
                                className="object-cover"
                            />
                        </motion.div>

                        {/* Thumbnails */}
                        <div className="flex gap-3">
                            {[0, 1, 2, 3].map((i) => (
                                <button
                                    key={i}
                                    onClick={() => setSelectedImage(i)}
                                    className={cn(
                                        "relative h-20 w-20 rounded-xl overflow-hidden bg-secondary border-2 transition-colors",
                                        selectedImage === i ? "border-primary" : "border-transparent hover:border-border"
                                    )}
                                >
                                    <Image
                                        src={product.image_url || "/placeholder.png"}
                                        alt={`${product.name} thumbnail ${i + 1}`}
                                        fill
                                        className="object-cover"
                                    />
                                </button>
                            ))}
                        </div>
                    </div>

                    {/* Product Info */}
                    <div className="space-y-6">
                        {product.category && (
                            <Badge variant="secondary" className="uppercase tracking-wider">
                                {product.category}
                            </Badge>
                        )}

                        <h1 className="text-3xl lg:text-4xl font-bold tracking-tight">
                            {product.name}
                        </h1>

                        {/* Rating */}
                        <div className="flex items-center gap-4">
                            <div className="flex items-center gap-1">
                                {[1, 2, 3, 4, 5].map((star) => (
                                    <Star
                                        key={star}
                                        className={cn(
                                            "h-5 w-5",
                                            star <= 4 ? "text-yellow-500 fill-yellow-500" : "text-muted-foreground"
                                        )}
                                    />
                                ))}
                            </div>
                            <span className="text-sm text-muted-foreground">
                                4.0 (128 reviews)
                            </span>
                        </div>

                        {/* Price */}
                        <div className="flex items-baseline gap-3">
                            <span className="text-4xl font-bold text-primary">
                                Rp {product.price?.toLocaleString("id-ID") || "0"}
                            </span>
                            <span className="text-lg text-muted-foreground line-through">
                                Rp {((product.price || 0) * 1.2).toLocaleString("id-ID")}
                            </span>
                            <Badge className="bg-red-500 hover:bg-red-500">-20%</Badge>
                        </div>

                        <Separator />

                        {/* Description */}
                        <div>
                            <h3 className="font-semibold mb-2">Description</h3>
                            <p className="text-muted-foreground leading-relaxed">
                                {product.description ||
                                    "Premium quality product with exceptional craftsmanship. Made with the finest materials to ensure durability and satisfaction. Perfect for everyday use or as a special gift."}
                            </p>
                        </div>

                        <Separator />

                        {/* Quantity & Add to Cart */}
                        <div className="space-y-4">
                            <div className="flex items-center gap-4">
                                <span className="font-medium">Quantity</span>
                                <div className="flex items-center gap-2">
                                    <button
                                        onClick={() => setQuantity(Math.max(1, quantity - 1))}
                                        className="h-10 w-10 rounded-lg border flex items-center justify-center hover:bg-secondary transition-colors"
                                    >
                                        <Minus className="h-4 w-4" />
                                    </button>
                                    <span className="w-12 text-center font-semibold text-lg">{quantity}</span>
                                    <button
                                        onClick={() => setQuantity(quantity + 1)}
                                        className="h-10 w-10 rounded-lg border flex items-center justify-center hover:bg-secondary transition-colors"
                                    >
                                        <Plus className="h-4 w-4" />
                                    </button>
                                </div>
                            </div>

                            <div className="flex gap-3">
                                <Button
                                    size="lg"
                                    className="flex-1 h-14 text-base gap-2"
                                    onClick={handleAddToCart}
                                >
                                    <ShoppingCart className="h-5 w-5" /> Add to Cart
                                </Button>
                                <Button size="lg" variant="outline" className="h-14 w-14">
                                    <Heart className="h-5 w-5" />
                                </Button>
                                <Button size="lg" variant="outline" className="h-14 w-14">
                                    <Share2 className="h-5 w-5" />
                                </Button>
                            </div>
                        </div>

                        {/* Benefits */}
                        <div className="grid grid-cols-3 gap-4 pt-4">
                            {[
                                { icon: Truck, label: "Free Shipping" },
                                { icon: ShieldCheck, label: "Secure Payment" },
                                { icon: RefreshCw, label: "Easy Returns" },
                            ].map(({ icon: Icon, label }) => (
                                <div key={label} className="text-center p-4 rounded-xl bg-secondary/50">
                                    <Icon className="h-6 w-6 mx-auto mb-2 text-primary" />
                                    <span className="text-xs font-medium">{label}</span>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

function ProductSkeleton() {
    return (
        <div className="container mx-auto px-6 py-12">
            <div className="grid lg:grid-cols-2 gap-12">
                <div className="aspect-square rounded-2xl bg-secondary animate-pulse" />
                <div className="space-y-6">
                    <div className="h-8 w-32 rounded bg-secondary animate-pulse" />
                    <div className="h-12 w-3/4 rounded bg-secondary animate-pulse" />
                    <div className="h-6 w-48 rounded bg-secondary animate-pulse" />
                    <div className="h-10 w-40 rounded bg-secondary animate-pulse" />
                    <div className="h-24 w-full rounded bg-secondary animate-pulse" />
                </div>
            </div>
        </div>
    );
}
