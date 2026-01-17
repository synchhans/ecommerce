"use client";

import React, { useMemo, useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { motion, useScroll, useMotionValueEvent } from "framer-motion";
import { ShoppingCart, Search, Menu, User, X } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Sheet, SheetContent, SheetTrigger, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import { Separator } from "@/components/ui/separator";
import { useCartStore } from "@/store/cart";
import { cn } from "@/lib/utils";

const CATEGORIES = [
  { label: "All Products", href: "/catalog" },
  { label: "Electronics", href: "/catalog?cat=electronics" },
  { label: "Fashion", href: "/catalog?cat=fashion" },
  { label: "Beauty", href: "/catalog?cat=beauty" },
  { label: "Home & Living", href: "/catalog?cat=home" },
];

export function Header() {
  const cartCount = useCartStore((s) => s.count);
  const { scrollY } = useScroll();
  const [isScrolled, setIsScrolled] = useState(false);

  useMotionValueEvent(scrollY, "change", (latest) => {
    setIsScrolled(latest > 20);
  });

  return (
    <motion.header
      className={cn(
        "sticky top-0 z-50 w-full transition-all duration-300",
        isScrolled
          ? "bg-background/60 backdrop-blur-xl border-b shadow-sm py-2"
          : "bg-transparent py-4"
      )}
    >
      <div className="container mx-auto flex items-center justify-between gap-4 px-4 sm:px-6">
        {/* Mobile Menu */}
        <div className="md:hidden">
          <MobileNav />
        </div>

        {/* Logo */}
        <Link href="/" className="flex items-center gap-2 group">
          <motion.div
            whileHover={{ rotate: 10, scale: 1.1 }}
            className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-primary to-purple-600 text-white shadow-lg shadow-primary/20"
          >
            <span className="text-xl font-bold">E</span>
          </motion.div>
          <div className="hidden sm:block">
            <h1 className="text-lg font-bold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-foreground to-foreground/70 group-hover:to-primary transition-all">
              Ecommerce
            </h1>
          </div>
        </Link>

        {/* Desktop Navigation */}
        <nav className="hidden md:flex items-center gap-1 rounded-full bg-secondary/50 p-1 backdrop-blur-3xl border border-white/5">
          {CATEGORIES.slice(0, 5).map((c) => (
            <Link
              key={c.href}
              href={c.href}
              className="px-4 py-2 text-sm font-medium text-muted-foreground transition-colors hover:text-foreground hover:bg-background/80 rounded-full"
            >
              {c.label}
            </Link>
          ))}
        </nav>

        {/* Right Section: Search & Actions */}
        <div className="flex items-center gap-2 sm:gap-4">
          <div className="hidden lg:block w-64">
            <React.Suspense fallback={<div className="h-10 w-full rounded-full bg-secondary/50 animate-pulse" />}>
              <SearchBar />
            </React.Suspense>
          </div>

          <div className="flex items-center gap-2">
            <Link href="/search" className="lg:hidden">
              <Button variant="ghost" size="icon" className="rounded-full">
                <Search className="h-5 w-5" />
              </Button>
            </Link>

            <Link href="/cart">
              <Button variant="ghost" size="icon" className="relative rounded-full hover:bg-primary/10 hover:text-primary transition-colors">
                <ShoppingCart className="h-5 w-5" />
                {cartCount > 0 && (
                  <motion.span
                    initial={{ scale: 0 }}
                    animate={{ scale: 1 }}
                    className="absolute -right-1 -top-1 flex h-4 w-4 items-center justify-center rounded-full bg-red-500 text-[10px] font-bold text-white shadow-sm"
                  >
                    {cartCount}
                  </motion.span>
                )}
              </Button>
            </Link>

            <Separator orientation="vertical" className="h-6 hidden sm:block" />

            <div className="hidden sm:flex gap-2">
              <Link href="/login">
                <Button variant="ghost" size="sm" className="rounded-full">Sign in</Button>
              </Link>
              <Link href="/register">
                <Button size="sm" className="rounded-full bg-primary hover:bg-primary/90 shadow-lg shadow-primary/25">
                  Register
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </div>
    </motion.header>
  );
}

function MobileNav() {
  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button variant="ghost" size="icon" className="rounded-full">
          <Menu className="h-5 w-5" />
        </Button>
      </SheetTrigger>
      <SheetContent side="left" className="w-[300px] sm:w-[400px]">
        <SheetHeader>
          <SheetTitle className="text-left flex items-center gap-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-white">
              E
            </div>
            Menu
          </SheetTitle>
        </SheetHeader>
        <nav className="mt-8 flex flex-col gap-2">
          {CATEGORIES.map((c) => (
            <Link
              key={c.href}
              href={c.href}
              className="flex items-center justify-between rounded-lg px-4 py-3 text-sm font-medium transition-colors hover:bg-secondary"
            >
              {c.label}
            </Link>
          ))}
          <Separator className="my-4" />
          <Link href="/login" className="px-4 py-2 text-sm font-medium">Log in</Link>
          <Link href="/register" className="px-4 py-2 text-sm font-medium text-primary">Create Account</Link>
        </nav>
      </SheetContent>
    </Sheet>
  )
}

function SearchBar() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const initialQ = useMemo(() => searchParams.get("q") ?? "", [searchParams]);
  const [q, setQ] = useState(initialQ);
  const [isFocused, setIsFocused] = useState(false);

  function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    const query = q.trim();
    router.push(query ? `/catalog?q=${encodeURIComponent(query)}` : "/catalog");
  }

  return (
    <form onSubmit={onSubmit} className="w-full relative">
      <div className="relative group">
        <Search className={cn(
          "absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 transition-colors",
          isFocused ? "text-primary" : "text-muted-foreground"
        )} />
        <Input
          value={q}
          onChange={(e) => setQ(e.target.value)}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
          placeholder="Search..."
          className={cn(
            "h-10 w-full pl-9 rounded-full border-transparent bg-secondary/50 transition-all duration-300",
            "focus-visible:ring-primary focus-visible:bg-background focus-visible:shadow-lg focus-visible:shadow-primary/5",
            "hover:bg-secondary/70 placeholder:text-muted-foreground/50"
          )}
        />
      </div>
    </form>
  );
}
