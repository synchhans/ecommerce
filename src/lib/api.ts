const BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL!;

type ApiError = {
  error: string;
  message?: string;
};

export async function apiFetch<T>(
  path: string,
  options: RequestInit & { token?: string } = {}
): Promise<T> {
  const { token, headers, ...rest } = options;

  const res = await fetch(`${BASE_URL}${path}`, {
    ...rest,
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(headers || {}),
    },
    cache: "no-store",
  });

  const text = await res.text();
  const data = text ? JSON.parse(text) : null;

  if (!res.ok) {
    const err: ApiError = data || { error: "unknown_error" };
    throw Object.assign(new Error(err.message || err.error), {
      status: res.status,
      code: err.error,
      payload: err,
    });
  }

  return data as T;
}
