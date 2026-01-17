import Image from "next/image";
import Link from "next/link";
import { ShoppingCart } from "lucide-react";
import { Button } from "@/components/ui/button";

export function ProductCard({ product }: { product: any }) {
  return (
    <Link href={`/product/${product.id}`}>
      <div className="group rounded-xl border bg-white p-3 transition hover:shadow-lg">
        <div className="relative aspect-square overflow-hidden rounded-lg bg-muted">
          <Image
            src={product.image_url || "/placeholder.png"}
            alt={product.name}
            fill
            className="object-cover transition group-hover:scale-105"
          />
        </div>

        <div className="mt-3 space-y-1">
          <h3 className="line-clamp-2 text-sm font-medium">
            {product.name}
          </h3>

          <div className="flex items-center justify-between">
            <span className="font-semibold text-primary">
              Rp {product.price.toLocaleString("id-ID")}
            </span>

            <Button size="icon" variant="ghost">
              <ShoppingCart className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>
    </Link>
  );
}
