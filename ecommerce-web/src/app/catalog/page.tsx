import CatalogClient from "./catalog-client";
import { Suspense } from "react";

export default function CatalogPage() {
  return (
    <Suspense fallback={<div className="min-h-screen bg-background" />}>
      <CatalogClient />
    </Suspense>
  );
}
