const API_URL = "http://localhost:8080"

// apiFetch wraps fetch: prepends the base URL, attaches the JWT (if present),
// sets JSON headers, throws on non-2xx, and parses the JSON response.
// Generic <T> lets callers say what type they expect, e.g. apiFetch<Board[]>("/boards").
export async function apiFetch<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = localStorage.getItem("token")

  const res = await fetch(API_URL + path, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers,
    },
  })

  if (!res.ok) {
    // try to surface the server's error text; fall back to the status code
    const message = await res.text()
    throw new Error(message || `Request failed with status ${res.status}`)
  }

  if (res.status === 204) {
    return undefined as T // 204 No Content: nothing to parse
  }

  return res.json() as Promise<T>
}
