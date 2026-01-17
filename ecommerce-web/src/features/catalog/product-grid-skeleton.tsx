import { Skeleton } from "@/components/ui/skeleton";

export function ProductGridSkeleton() {
  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {Array.from({ length: 8 }).map((_, i) => (
        <div key={i} className="overflow-hidden rounded-xl border bg-background">
          <Skeleton className="aspect-[4/3] w-full" />
          <div className="grid gap-2 p-4">
            <Skeleton className="h-4 w-4/5" />
            <Skeleton className="h-4 w-3/5" />
            <div className="flex items-center justify-between pt-1">
              <Skeleton className="h-5 w-24" />
              <Skeleton className="h-4 w-10" />
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
