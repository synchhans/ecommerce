"use client";

import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ShoppingCart, Truck, ShieldCheck, ArrowRight, Star, Zap } from "lucide-react";
import { motion } from "framer-motion";
import Image from "next/image";
import Link from "next/link";
import { cn } from "@/lib/utils";

const FADE_UP_VARIANT = {
  hidden: { opacity: 0, y: 30 },
  visible: { opacity: 1, y: 0, transition: { duration: 0.6, ease: "easeOut" } },
};

const STAGGER_CONTAINER = {
  visible: { transition: { staggerChildren: 0.1 } },
};

export default function HomePage() {
  return (
    <div className="min-h-screen flex flex-col">
      {/* Immersive Hero Section */}
      <section className="relative h-[85vh] flex items-center overflow-hidden bg-background">
        {/* Background Gradients */}
        <div className="absolute inset-0 z-0">
          <div className="absolute top-[-10%] right-[-5%] w-[500px] h-[500px] bg-primary/20 rounded-full blur-[120px] animate-pulse" />
          <div className="absolute bottom-[-10%] left-[-10%] w-[600px] h-[600px] bg-purple-500/10 rounded-full blur-[140px]" />
        </div>

        <div className="container relative z-10 px-4 sm:px-6 mx-auto grid lg:grid-cols-2 gap-12 items-center">
          {/* Text Content */}
          <motion.div
            initial="hidden"
            animate="visible"
            variants={STAGGER_CONTAINER}
            className="max-w-2xl"
          >
            <motion.div variants={FADE_UP_VARIANT} className="mb-6">
              <span className="inline-flex items-center gap-2 rounded-full border border-primary/20 bg-primary/5 px-4 py-1.5 text-sm font-medium text-primary">
                <Zap className="h-4 w-4 fill-primary" />
                New Collection 2026
              </span>
            </motion.div>

            <motion.h1 variants={FADE_UP_VARIANT} className="text-5xl sm:text-7xl font-bold tracking-tight leading-[1.1] mb-6">
              Elevate Your <br />
              <span className="text-transparent bg-clip-text bg-gradient-to-r from-primary to-purple-600">
                Digital Lifestyle
              </span>
            </motion.h1>

            <motion.p variants={FADE_UP_VARIANT} className="text-lg sm:text-xl text-muted-foreground mb-8 text-balance">
              Experience the future of shopping with our curated collection of premium tech and lifestyle essentials.
            </motion.p>

            <motion.div variants={FADE_UP_VARIANT} className="flex flex-wrap gap-4">
              <Link href="/catalog">
                <Button size="lg" className="h-12 px-8 rounded-full text-base shadow-lg shadow-primary/25 hover:shadow-primary/40 transition-shadow">
                  Start Shopping <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              </Link>
              <Link href="/catalog?cat=sale">
                <Button size="lg" variant="outline" className="h-12 px-8 rounded-full text-base bg-background/50 backdrop-blur border-primary/20 hover:bg-primary/5">
                  View Deals
                </Button>
              </Link>
            </motion.div>

            <motion.div variants={FADE_UP_VARIANT} className="mt-12 flex items-center gap-8 text-muted-foreground">
              <div className="flex -space-x-3">
                {[1, 2, 3, 4].map((i) => (
                  <div key={i} className="h-10 w-10 rounded-full bg-secondary border-2 border-background flex items-center justify-center text-xs font-bold">
                    U{i}
                  </div>
                ))}
              </div>
              <div className="flex flex-col">
                <div className="flex items-center gap-1 text-yellow-500">
                  {[1, 2, 3, 4, 5].map((i) => <Star key={i} className="h-4 w-4 fill-current" />)}
                </div>
                <span className="text-sm font-medium">Trusted by 10k+ Customers</span>
              </div>
            </motion.div>
          </motion.div>

          {/* Visual/Image Area (Abstract Representation) */}
          <motion.div
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.8, delay: 0.2 }}
            className="hidden lg:block relative h-[600px] w-full"
          >
            <div className="absolute inset-0 bg-gradient-to-tr from-primary/10 to-transparent rounded-3xl backdrop-blur-3xl border border-white/5" />
            {/* Placeholder for Hero Image - using CSS Composition to look good without asset */}
            <div className="absolute inset-4 rounded-2xl bg-gradient-to-br from-gray-900 to-black overflow-hidden flex items-center justify-center border border-white/10 shadow-2xl">
              <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_120%,rgba(120,119,198,0.3),rgba(255,255,255,0))]" />
              <div className="text-center p-12">
                <div className="inline-block p-4 rounded-3xl bg-white/5 backdrop-blur-md border border-white/10 mb-8 animate-float">
                  <ShoppingBagIcon className="h-32 w-32 text-primary/80" />
                </div>
                <p className="text-white/50 text-sm tracking-widest uppercase">Premium Goods</p>
              </div>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Featured Categories */}
      <section className="py-24 bg-secondary/30 relative">
        <div className="container mx-auto px-6">
          <div className="flex items-end justify-between mb-12">
            <div>
              <h2 className="text-3xl font-bold tracking-tight mb-2">Shop by Category</h2>
              <p className="text-muted-foreground">Explore our most popular collections</p>
            </div>
            <Link href="/catalog" className="hidden sm:flex items-center text-primary font-medium hover:underline">
              View All <ArrowRight className="ml-1 h-4 w-4" />
            </Link>
          </div>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 sm:gap-6">
            {['Electronics', 'Fashion', 'Home', 'Beauty'].map((cat, i) => (
              <Link key={cat} href={`/catalog?cat=${cat.toLowerCase()}`}>
                <motion.div
                  whileHover={{ y: -5 }}
                  className="group relative aspect-square overflow-hidden rounded-2xl bg-card border border-border/50 shadow-sm"
                >
                  <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent z-10" />
                  <div className="absolute bottom-0 left-0 p-4 z-20">
                    <h3 className="text-white text-xl font-bold">{cat}</h3>
                    <p className="text-white/80 text-sm opacity-0 group-hover:opacity-100 transition-opacity">Browse Collection</p>
                  </div>
                  {/* Colorful placeholder background */}
                  <div className={cn("absolute inset-0 transition-transform duration-500 group-hover:scale-110",
                    i === 0 ? "bg-blue-500/20" :
                      i === 1 ? "bg-purple-500/20" :
                        i === 2 ? "bg-emerald-500/20" : "bg-pink-500/20"
                  )} />
                </motion.div>
              </Link>
            ))}
          </div>
        </div>
      </section>

      {/* Value Props */}
      <section className="py-24 container mx-auto px-6">
        <div className="grid md:grid-cols-3 gap-8">
          <FeatureCard
            icon={<Zap className="h-6 w-6 text-yellow-500" />}
            title="Super Fast Delivery"
            desc="Get your orders delivered in 24 hours or less."
          />
          <FeatureCard
            icon={<ShieldCheck className="h-6 w-6 text-emerald-500" />}
            title="Secure Payment"
            desc="100% secure payment with advanced encryption."
          />
          <FeatureCard
            icon={<Star className="h-6 w-6 text-purple-500" />}
            title="Premium Quality"
            desc="Handpicked products from top tier brands."
          />
        </div>
      </section>
    </div>
  );
}

function FeatureCard({ icon, title, desc }: { icon: React.ReactNode, title: string, desc: string }) {
  return (
    <div className="p-6 rounded-2xl bg-background border border-border/50 hover:border-primary/50 transition-colors shadow-sm">
      <div className="h-12 w-12 rounded-xl bg-secondary flex items-center justify-center mb-4">
        {icon}
      </div>
      <h3 className="font-bold text-lg mb-2">{title}</h3>
      <p className="text-muted-foreground">{desc}</p>
    </div>
  )
}

function ShoppingBagIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1" className={className}>
      <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z" />
      <line x1="3" y1="6" x2="21" y2="6" />
      <path d="M16 10a4 4 0 0 1-8 0" />
    </svg>
  )
}
