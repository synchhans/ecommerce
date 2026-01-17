import CatalogClient from "./ui";

export default function CatalogPage({
  searchParams,
}: {
  searchParams: { q?: string; cat?: string; page?: string };
}) {
  return (
    <CatalogClient
      q={searchParams.q}
      cat={searchParams.cat}
      page={searchParams.page ? Number(searchParams.page) : 1}
    />
  );
}
