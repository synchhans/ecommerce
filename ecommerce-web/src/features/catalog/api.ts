import { apiFetch } from "@/lib/api";
import type { Product } from "./types";

export type ListProductsParams = {
  q?: string;
  cat?: string;
  page?: number;
  limit?: number;
};

type ListProductsResponse =
  | { items: Product[]; page?: number; limit?: number; total?: number }
  | Product[]; // fallback kalau backend return array langsung

export async function listProducts(params: ListProductsParams) {
  const sp = new URLSearchParams();
  if (params.q) sp.set("q", params.q);
  if (params.cat) sp.set("cat", params.cat);
  if (params.page) sp.set("page", String(params.page));
  if (params.limit) sp.set("limit", String(params.limit));

  const qs = sp.toString() ? `?${sp.toString()}` : "";
  const res = await apiFetch<ListProductsResponse>(`/v1/products${qs}`);

  // normalize
  if (Array.isArray(res)) return { items: res };
  return res;
}
