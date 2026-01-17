import Link from "next/link";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

type Props = {
  id: string;
  slug: string;
  name: string;
  price?: number;
  image_url?: string;
};

export function ProductCard({ slug, name, price, image_url }: Props) {
  return (
    <Link href={`/p/${slug}`} className="group">
      <Card className="overflow-hidden transition-shadow hover:shadow-md">
        <div className="aspect-[4/3] w-full bg-muted">
          {/* placeholder image */}
          {image_url ? (
            // eslint-disable-next-line @next/next/no-img-element
            <img
              src={image_url}
              alt={name}
              className="h-full w-full object-cover transition-transform group-hover:scale-[1.02]"
            />
          ) : (
            <div className="flex h-full w-full items-center justify-center text-xs text-muted-foreground">
              No image
            </div>
          )}
        </div>

        <CardContent className="grid gap-2 p-4">
          <div className="line-clamp-2 text-sm font-medium">{name}</div>

          <div className="flex items-center justify-between">
            {typeof price === "number" ? (
              <div className="text-base font-semibold">
                Rp {price.toLocaleString("id-ID")}
              </div>
            ) : (
              <Badge variant="secondary">View details</Badge>
            )}

            <div className="text-xs text-muted-foreground">★ 4.8</div>
          </div>
        </CardContent>
      </Card>
    </Link>
  );
}
