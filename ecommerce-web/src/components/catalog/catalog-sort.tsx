export function CatalogSort() {
  return (
    <div className="flex items-center justify-between rounded-xl border bg-white p-3">
      <span className="text-sm text-muted-foreground">
        Showing products
      </span>

      <select className="rounded-md border px-2 py-1 text-sm">
        <option>Popular</option>
        <option>Newest</option>
        <option>Price: Low → High</option>
        <option>Price: High → Low</option>
      </select>
    </div>
  );
}
