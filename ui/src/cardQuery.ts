import * as Comlink from "comlink"
import type CardQuery from './workers/cardQuery'
import { createContext } from "react";

const w = new Worker(new URL("./workers/cardQuery", import.meta.url), { type: "module" })
const cardQuery = Comlink.wrap<CardQuery>(w);

export const CardsReady = (async () => {
  const cardsJson = await fetch(new URL("/cards.json", import.meta.url)).then(r => r.text())
  console.log("fetched cards")
  await cardQuery.feedCards(cardsJson)
  console.log("fed cards")
})()

export const CardContext = createContext(cardQuery)
