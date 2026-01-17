"use client";

import Link from "next/link";
import Image from "next/image";
import { motion, AnimatePresence } from "framer-motion";
import { Trash2, Plus, Minus, ShoppingBag, ArrowRight, Tag } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { useCartStore, CartItem } from "@/store/cart";
import { cn } from "@/lib/utils";

export default function CartPage() {
    const { items, removeItem, updateQuantity, clearCart } = useCartStore();

    const subtotal = items.reduce((sum, item) => sum + item.price * item.quantity, 0);
    const shipping = subtotal > 500000 ? 0 : 25000;
    const total = subtotal + shipping;

    if (items.length === 0) {
        return <EmptyCart />;
    }

    return (
        <div className="min-h-screen bg-background">
            <div className="container mx-auto px-6 py-12">
                <div className="mb-8">
                    <h1 className="text-3xl font-bold tracking-tight">Shopping Cart</h1>
                    <p className="text-muted-foreground mt-1">
                        {items.length} item{items.length > 1 ? "s" : ""} in your cart
                    </p>
                </div>

                <div className="grid lg:grid-cols-3 gap-8">
                    {/* Cart Items */}
                    <div className="lg:col-span-2 space-y-4">
                        <AnimatePresence mode="popLayout">
                            {items.map((item) => (
                                <CartItemCard
                                    key={item.id}
                                    item={item}
                                    onRemove={() => removeItem(item.id)}
                                    onUpdateQuantity={(qty) => updateQuantity(item.id, qty)}
                                />
                            ))}
                        </AnimatePresence>

                        <div className="flex justify-between pt-4">
                            <Button variant="outline" onClick={clearCart} className="text-destructive hover:text-destructive">
                                Clear Cart
                            </Button>
                            <Link href="/catalog">
                                <Button variant="ghost">Continue Shopping</Button>
                            </Link>
                        </div>
                    </div>

                    {/* Order Summary */}
                    <div className="lg:col-span-1">
                        <div className="sticky top-24 rounded-2xl border bg-card p-6 space-y-6">
                            <h2 className="text-xl font-semibold">Order Summary</h2>

                            <div className="space-y-3 text-sm">
                                <div className="flex justify-between">
                                    <span className="text-muted-foreground">Subtotal</span>
                                    <span>Rp {subtotal.toLocaleString("id-ID")}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-muted-foreground">Shipping</span>
                                    <span className={shipping === 0 ? "text-green-600" : ""}>
                                        {shipping === 0 ? "FREE" : `Rp ${shipping.toLocaleString("id-ID")}`}
                                    </span>
                                </div>
                                {shipping > 0 && (
                                    <p className="text-xs text-muted-foreground">
                                        Free shipping on orders over Rp 500,000
                                    </p>
                                )}
                            </div>

                            <Separator />

                            {/* Promo Code */}
                            <div className="space-y-2">
                                <label className="text-sm font-medium flex items-center gap-2">
                                    <Tag className="h-4 w-4" /> Promo Code
                                </label>
                                <div className="flex gap-2">
                                    <Input placeholder="Enter code" className="flex-1" />
                                    <Button variant="secondary">Apply</Button>
                                </div>
                            </div>

                            <Separator />

                            <div className="flex justify-between text-lg font-semibold">
                                <span>Total</span>
                                <span className="text-primary">Rp {total.toLocaleString("id-ID")}</span>
                            </div>

                            <Link href="/checkout" className="block">
                                <Button className="w-full h-12 text-base gap-2">
                                    Proceed to Checkout <ArrowRight className="h-4 w-4" />
                                </Button>
                            </Link>

                            <div className="text-center text-xs text-muted-foreground">
                                Taxes calculated at checkout
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

function CartItemCard({
    item,
    onRemove,
    onUpdateQuantity,
}: {
    item: CartItem;
    onRemove: () => void;
    onUpdateQuantity: (qty: number) => void;
}) {
    return (
        <motion.div
            layout
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, x: -100 }}
            className="flex gap-4 p-4 rounded-xl border bg-card"
        >
            {/* Image */}
            <div className="relative h-24 w-24 rounded-lg overflow-hidden bg-secondary flex-shrink-0">
                <Image
                    src={item.image_url || "/placeholder.png"}
                    alt={item.name}
                    fill
                    className="object-cover"
                />
            </div>

            {/* Info */}
            <div className="flex-1 flex flex-col justify-between">
                <div>
                    <h3 className="font-medium line-clamp-1">{item.name}</h3>
                    <p className="text-sm text-muted-foreground">SKU: {item.id.slice(0, 8)}</p>
                </div>

                <div className="flex items-center justify-between mt-2">
                    {/* Quantity */}
                    <div className="flex items-center gap-2">
                        <button
                            onClick={() => onUpdateQuantity(item.quantity - 1)}
                            className="h-8 w-8 rounded-lg border flex items-center justify-center hover:bg-secondary transition-colors"
                            disabled={item.quantity <= 1}
                        >
                            <Minus className="h-3 w-3" />
                        </button>
                        <span className="w-8 text-center font-medium">{item.quantity}</span>
                        <button
                            onClick={() => onUpdateQuantity(item.quantity + 1)}
                            className="h-8 w-8 rounded-lg border flex items-center justify-center hover:bg-secondary transition-colors"
                        >
                            <Plus className="h-3 w-3" />
                        </button>
                    </div>

                    {/* Price */}
                    <p className="font-semibold text-primary">
                        Rp {(item.price * item.quantity).toLocaleString("id-ID")}
                    </p>
                </div>
            </div>

            {/* Remove */}
            <button
                onClick={onRemove}
                className="self-start p-2 text-muted-foreground hover:text-destructive transition-colors"
            >
                <Trash2 className="h-4 w-4" />
            </button>
        </motion.div>
    );
}

function EmptyCart() {
    return (
        <div className="min-h-[70vh] flex flex-col items-center justify-center text-center px-6">
            <motion.div
                initial={{ opacity: 0, scale: 0.9 }}
                animate={{ opacity: 1, scale: 1 }}
                className="space-y-6"
            >
                <div className="h-24 w-24 mx-auto rounded-full bg-secondary flex items-center justify-center">
                    <ShoppingBag className="h-10 w-10 text-muted-foreground" />
                </div>
                <div>
                    <h2 className="text-2xl font-bold">Your cart is empty</h2>
                    <p className="text-muted-foreground mt-2">
                        Looks like you haven't added anything yet.
                    </p>
                </div>
                <Link href="/catalog">
                    <Button size="lg" className="gap-2">
                        Start Shopping <ArrowRight className="h-4 w-4" />
                    </Button>
                </Link>
            </motion.div>
        </div>
    );
}
