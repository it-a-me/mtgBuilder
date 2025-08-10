import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App'
import * as Comlink from "comlink"
import type CardQuery from './workers/cardQuery'

const w = new Worker(new URL("./workers/cardQuery", import.meta.url), { type: "module" })
const cardQuery = Comlink.wrap<CardQuery>(w);
// const cardQuery = new CardQuery()
(async () => {
  const cardsJson = await fetch("/cards.json").then(r => r.text())
  console.log("fetched cards")
  await cardQuery.feedCards(cardsJson)
  console.log("fed cards")
  const first = await cardQuery.getCard(0)
  console.log(`card 0 is ${first.name}`)
})()


createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
