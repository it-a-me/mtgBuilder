import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App'
import * as Comlink from "comlink"
import type CardQuery from './workers/cardQuery'

const w = new Worker(new URL("./workers/cardQuery", import.meta.url), { type: "module" })
const cardQuery = Comlink.wrap<CardQuery>(w);

(async () => {
  const cardsJson = await fetch(new URL("/cards.json", import.meta.url)).then(r => r.text())
  console.log("fetched cards")
  await cardQuery.feedCards(cardsJson)
  console.log("fed cards")
})()


createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App cardQuery={cardQuery} />
  </StrictMode>,
)
