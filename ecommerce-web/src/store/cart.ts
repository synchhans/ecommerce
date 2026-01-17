// src/store/cart.ts
import { create } from "zustand";
import { persist } from "zustand/middleware";

export type CartItem = {
  id: string;
  name: string;
  price: number;
  quantity: number;
  image_url: string;
};

type CartState = {
  items: CartItem[];
  count: number;
  addItem: (item: Omit<CartItem, "quantity">) => void;
  removeItem: (id: string) => void;
  updateQuantity: (id: string, quantity: number) => void;
  clearCart: () => void;
  inc: () => void;
  setCount: (n: number) => void;
};

export const useCartStore = create<CartState>()(
  persist(
    (set, get) => ({
      items: [],
      count: 0,

      addItem: (item) => set((state) => {
        const existing = state.items.find((i) => i.id === item.id);
        if (existing) {
          return {
            items: state.items.map((i) =>
              i.id === item.id ? { ...i, quantity: i.quantity + 1 } : i
            ),
            count: state.count + 1,
          };
        }
        return {
          items: [...state.items, { ...item, quantity: 1 }],
          count: state.count + 1,
        };
      }),

      removeItem: (id) => set((state) => {
        const item = state.items.find((i) => i.id === id);
        return {
          items: state.items.filter((i) => i.id !== id),
          count: state.count - (item?.quantity || 0),
        };
      }),

      updateQuantity: (id, quantity) => set((state) => {
        const item = state.items.find((i) => i.id === id);
        const diff = quantity - (item?.quantity || 0);
        return {
          items: state.items.map((i) =>
            i.id === id ? { ...i, quantity: Math.max(1, quantity) } : i
          ),
          count: state.count + diff,
        };
      }),

      clearCart: () => set({ items: [], count: 0 }),

      inc: () => set((s) => ({ count: s.count + 1 })),
      setCount: (n) => set({ count: n }),
    }),
    {
      name: "cart-storage",
    }
  )
);
