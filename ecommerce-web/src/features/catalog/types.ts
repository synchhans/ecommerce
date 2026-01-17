export type Product = {
  id: string;
  slug: string;
  name: string;
  description?: string;
  price?: number; // optional: kalau backend mengirim
  image_url?: string; // optional
};
