import { apiFetch } from "./client"
import type { Column } from "../types"

export function listColumns(boardId: string) {
  return apiFetch<Column[]>(`/boards/${boardId}/columns`)
}

export function createColumn(boardId: string, name: string) {
  return apiFetch<Column>(`/boards/${boardId}/columns`, {
    method: "POST",
    body: JSON.stringify({ name }),
  })
}
