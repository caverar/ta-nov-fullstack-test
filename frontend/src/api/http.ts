export async function http<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const res = await fetch(import.meta.env.VITE_API_URL + path, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers || {}),
    },
  });

  if (!res.ok) {
    throw new Error(`HTTP error: ${res.status}`);
  }

  return res.json() as Promise<T>;
}
