export function CatalogFilters() {
  return (
    <aside className="hidden md:block rounded-xl border bg-white p-4">
      <h2 className="mb-4 font-semibold">Filter</h2>

      <div className="space-y-4 text-sm">
        <div>
          <p className="mb-2 font-medium">Category</p>
          <div className="space-y-1">
            <label><input type="checkbox" /> Electronics</label>
            <label><input type="checkbox" /> Fashion</label>
            <label><input type="checkbox" /> Beauty</label>
          </div>
        </div>

        <div>
          <p className="mb-2 font-medium">Price</p>
          <input type="range" className="w-full" />
        </div>
      </div>
    </aside>
  );
}
