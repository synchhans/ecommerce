"use client";

import { useState } from "react";
import Link from "next/link";
import Image from "next/image";
import { motion } from "framer-motion";
import { Check, CreditCard, Truck, MapPin, ChevronRight, Lock } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { useCartStore } from "@/store/cart";
import { cn } from "@/lib/utils";

const STEPS = ["Shipping", "Payment", "Review"];

export default function CheckoutPage() {
    const [currentStep, setCurrentStep] = useState(0);
    const { items } = useCartStore();

    const subtotal = items.reduce((sum, item) => sum + item.price * item.quantity, 0);
    const shipping = subtotal > 500000 ? 0 : 25000;
    const total = subtotal + shipping;

    if (items.length === 0) {
        return (
            <div className="min-h-[60vh] flex flex-col items-center justify-center">
                <h2 className="text-2xl font-bold">Your cart is empty</h2>
                <p className="text-muted-foreground mt-2">Add items to your cart before checkout.</p>
                <Link href="/catalog">
                    <Button className="mt-4">Browse Products</Button>
                </Link>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-background">
            <div className="container mx-auto px-6 py-12">
                {/* Progress Steps */}
                <div className="flex justify-center mb-12">
                    <div className="flex items-center gap-4">
                        {STEPS.map((step, i) => (
                            <div key={step} className="flex items-center gap-4">
                                <button
                                    onClick={() => setCurrentStep(i)}
                                    className={cn(
                                        "flex items-center gap-3 transition-colors",
                                        i <= currentStep ? "text-foreground" : "text-muted-foreground"
                                    )}
                                >
                                    <div className={cn(
                                        "h-10 w-10 rounded-full flex items-center justify-center font-semibold transition-colors",
                                        i < currentStep ? "bg-primary text-primary-foreground" :
                                            i === currentStep ? "bg-primary text-primary-foreground" :
                                                "bg-secondary"
                                    )}>
                                        {i < currentStep ? <Check className="h-5 w-5" /> : i + 1}
                                    </div>
                                    <span className="font-medium hidden sm:inline">{step}</span>
                                </button>
                                {i < STEPS.length - 1 && (
                                    <div className={cn(
                                        "w-12 h-0.5 transition-colors",
                                        i < currentStep ? "bg-primary" : "bg-secondary"
                                    )} />
                                )}
                            </div>
                        ))}
                    </div>
                </div>

                <div className="grid lg:grid-cols-3 gap-8">
                    {/* Form Section */}
                    <div className="lg:col-span-2">
                        {currentStep === 0 && <ShippingForm onNext={() => setCurrentStep(1)} />}
                        {currentStep === 1 && <PaymentForm onNext={() => setCurrentStep(2)} onBack={() => setCurrentStep(0)} />}
                        {currentStep === 2 && <ReviewOrder onBack={() => setCurrentStep(1)} />}
                    </div>

                    {/* Order Summary */}
                    <div className="lg:col-span-1">
                        <div className="sticky top-24 rounded-2xl border bg-card p-6 space-y-6">
                            <h2 className="text-xl font-semibold">Order Summary</h2>

                            <div className="space-y-4 max-h-64 overflow-y-auto">
                                {items.map((item) => (
                                    <div key={item.id} className="flex gap-3">
                                        <div className="relative h-16 w-16 rounded-lg overflow-hidden bg-secondary flex-shrink-0">
                                            <Image
                                                src={item.image_url || "/placeholder.png"}
                                                alt={item.name}
                                                fill
                                                className="object-cover"
                                            />
                                        </div>
                                        <div className="flex-1 min-w-0">
                                            <p className="font-medium text-sm line-clamp-1">{item.name}</p>
                                            <p className="text-xs text-muted-foreground">Qty: {item.quantity}</p>
                                            <p className="text-sm font-medium text-primary">
                                                Rp {(item.price * item.quantity).toLocaleString("id-ID")}
                                            </p>
                                        </div>
                                    </div>
                                ))}
                            </div>

                            <Separator />

                            <div className="space-y-2 text-sm">
                                <div className="flex justify-between">
                                    <span className="text-muted-foreground">Subtotal</span>
                                    <span>Rp {subtotal.toLocaleString("id-ID")}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-muted-foreground">Shipping</span>
                                    <span>{shipping === 0 ? "FREE" : `Rp ${shipping.toLocaleString("id-ID")}`}</span>
                                </div>
                            </div>

                            <Separator />

                            <div className="flex justify-between text-lg font-semibold">
                                <span>Total</span>
                                <span className="text-primary">Rp {total.toLocaleString("id-ID")}</span>
                            </div>

                            <div className="flex items-center justify-center gap-2 text-xs text-muted-foreground">
                                <Lock className="h-3 w-3" /> Secure SSL Encryption
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

function ShippingForm({ onNext }: { onNext: () => void }) {
    return (
        <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            className="rounded-2xl border bg-card p-6 space-y-6"
        >
            <div className="flex items-center gap-3">
                <MapPin className="h-5 w-5 text-primary" />
                <h2 className="text-xl font-semibold">Shipping Address</h2>
            </div>

            <div className="grid sm:grid-cols-2 gap-4">
                <div className="space-y-2">
                    <label className="text-sm font-medium">First Name</label>
                    <Input placeholder="John" className="h-12" />
                </div>
                <div className="space-y-2">
                    <label className="text-sm font-medium">Last Name</label>
                    <Input placeholder="Doe" className="h-12" />
                </div>
            </div>

            <div className="space-y-2">
                <label className="text-sm font-medium">Address</label>
                <Input placeholder="Street address" className="h-12" />
            </div>

            <div className="grid sm:grid-cols-3 gap-4">
                <div className="space-y-2">
                    <label className="text-sm font-medium">City</label>
                    <Input placeholder="Jakarta" className="h-12" />
                </div>
                <div className="space-y-2">
                    <label className="text-sm font-medium">Province</label>
                    <Input placeholder="DKI Jakarta" className="h-12" />
                </div>
                <div className="space-y-2">
                    <label className="text-sm font-medium">Postal Code</label>
                    <Input placeholder="12345" className="h-12" />
                </div>
            </div>

            <div className="space-y-2">
                <label className="text-sm font-medium">Phone Number</label>
                <Input placeholder="+62 812 3456 7890" className="h-12" />
            </div>

            <Button onClick={onNext} className="w-full h-12 gap-2">
                Continue to Payment <ChevronRight className="h-4 w-4" />
            </Button>
        </motion.div>
    );
}

function PaymentForm({ onNext, onBack }: { onNext: () => void; onBack: () => void }) {
    return (
        <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            className="rounded-2xl border bg-card p-6 space-y-6"
        >
            <div className="flex items-center gap-3">
                <CreditCard className="h-5 w-5 text-primary" />
                <h2 className="text-xl font-semibold">Payment Method</h2>
            </div>

            <div className="space-y-3">
                {["Credit Card", "Bank Transfer", "E-Wallet"].map((method, i) => (
                    <label
                        key={method}
                        className={cn(
                            "flex items-center gap-4 p-4 rounded-xl border cursor-pointer transition-colors",
                            i === 0 ? "border-primary bg-primary/5" : "hover:bg-secondary/50"
                        )}
                    >
                        <input type="radio" name="payment" defaultChecked={i === 0} className="h-4 w-4" />
                        <span className="font-medium">{method}</span>
                    </label>
                ))}
            </div>

            <div className="space-y-4">
                <div className="space-y-2">
                    <label className="text-sm font-medium">Card Number</label>
                    <Input placeholder="1234 5678 9012 3456" className="h-12" />
                </div>

                <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                        <label className="text-sm font-medium">Expiry Date</label>
                        <Input placeholder="MM/YY" className="h-12" />
                    </div>
                    <div className="space-y-2">
                        <label className="text-sm font-medium">CVV</label>
                        <Input placeholder="123" className="h-12" />
                    </div>
                </div>
            </div>

            <div className="flex gap-3">
                <Button variant="outline" onClick={onBack} className="flex-1 h-12">
                    Back
                </Button>
                <Button onClick={onNext} className="flex-1 h-12 gap-2">
                    Review Order <ChevronRight className="h-4 w-4" />
                </Button>
            </div>
        </motion.div>
    );
}

function ReviewOrder({ onBack }: { onBack: () => void }) {
    const clearCart = useCartStore((s) => s.clearCart);

    const handlePlaceOrder = () => {
        // TODO: Implement actual order placement
        alert("Order placed successfully!");
        clearCart();
        window.location.href = "/";
    };

    return (
        <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            className="rounded-2xl border bg-card p-6 space-y-6"
        >
            <div className="flex items-center gap-3">
                <Check className="h-5 w-5 text-primary" />
                <h2 className="text-xl font-semibold">Review & Confirm</h2>
            </div>

            <div className="space-y-4">
                <div className="p-4 rounded-xl bg-secondary/50">
                    <h3 className="font-medium mb-2">Shipping Address</h3>
                    <p className="text-sm text-muted-foreground">
                        John Doe<br />
                        123 Main Street<br />
                        Jakarta, DKI Jakarta 12345<br />
                        +62 812 3456 7890
                    </p>
                </div>

                <div className="p-4 rounded-xl bg-secondary/50">
                    <h3 className="font-medium mb-2">Payment Method</h3>
                    <p className="text-sm text-muted-foreground">
                        Credit Card ending in •••• 3456
                    </p>
                </div>
            </div>

            <div className="flex gap-3">
                <Button variant="outline" onClick={onBack} className="flex-1 h-12">
                    Back
                </Button>
                <Button onClick={handlePlaceOrder} className="flex-1 h-12 gap-2">
                    Place Order <Check className="h-4 w-4" />
                </Button>
            </div>
        </motion.div>
    );
}
