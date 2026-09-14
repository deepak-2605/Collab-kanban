import { apiFetch } from "./client"
import type { Board } from "../types"

export function listBoards() {
  return apiFetch<Board[]>("/boards")
}

export function createBoard(name: string) {
  return apiFetch<Board>("/boards", {
    method: "POST",
    body: JSON.stringify({ name }),
  })
}
