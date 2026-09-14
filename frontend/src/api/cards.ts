import { apiFetch } from "./client"
import type { Card } from "../types"

export function listCards(columnId: string) {
  return apiFetch<Card[]>(`/columns/${columnId}/cards`)
}

export function createCard(columnId: string, title: string, description = "") {
  return apiFetch<Card>(`/columns/${columnId}/cards`, {
    method: "POST",
    body: JSON.stringify({ title, description }),
  })
}

// PATCH /cards/:id/move — returns 204 (no body), so apiFetch<void>
export function moveCard(cardId: string, columnId: string, position: number) {
  return apiFetch<void>(`/cards/${cardId}/move`, {
    method: "PATCH",
    body: JSON.stringify({ columnId, position }),
  })
}
