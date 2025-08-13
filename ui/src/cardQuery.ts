import * as Comlink from "comlink"
import type CardQuery from './workers/cardQuery'
import { createContext } from "react";

const w = new Worker(new URL("./workers/cardQuery", import.meta.url), { type: "module" })
const cardQuery = Comlink.wrap<CardQuery>(w);

export const CardsReady = (async () => {
  const start = performance.now()
  await cardQuery.feedCards(new URL("/cards.json", import.meta.url).toString())
  console.log(`fed cards in ${performance.now() - start}ms`)
})()

export const CardContext = createContext(cardQuery)
